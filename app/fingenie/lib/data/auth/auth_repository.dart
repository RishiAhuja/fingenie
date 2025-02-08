import 'package:dio/dio.dart';
import 'package:fingenie/domain/entities/auth/create_login_request.dart';
import 'package:fingenie/domain/entities/auth/create_signup_request.dart';
import 'package:fingenie/domain/models/user_model.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:hive_flutter/hive_flutter.dart';

class AuthRepository {
  static Future<void> init() async {
    try {
      if (!Hive.isAdapterRegistered(0)) {
        Hive.registerAdapter(UserModelAdapter());
      }

      // Always close the box if it's open before reopening
      if (Hive.isBoxOpen(userBoxName)) {
        await Hive.box<UserModel>(userBoxName).close();
      }

      // Open the box
      final box = await Hive.openBox<UserModel>(userBoxName);
      AppLogger.debug('AuthRepository initialized with ${box.length} items');

      // Debug box contents
      for (var key in box.keys) {
        final user = box.get(key);
        AppLogger.debug('Box contains key: $key with user: ${user?.name}');
      }
    } catch (e) {
      AppLogger.error('Failed to initialize AuthRepository: $e');
      rethrow;
    }
  }

  final Dio _dio;
  final Box<UserModel> _userBox;

  static const String userBoxName = 'userBox';
  final String apiUrl = dotenv.env['API_URL'] ?? '';

  AuthRepository({
    Dio? dio,
    Box<UserModel>? userBox,
  })  : _dio = dio ?? Dio(),
        _userBox = userBox ?? Hive.box<UserModel>(userBoxName);

  Future<UserModel> signUp(SignUpRequest request) async {
    try {
      AppLogger.info('signUp: Making signup request');

      final user = UserModel(
          id: 'id123',
          name: request.name,
          email: request.email,
          phoneNumber: request.phoneNumber,
          createdAt: DateTime.now(),
          isLoggedIn: true,
          token: 'jwt');

      // Close and reopen box to ensure fresh state
      if (Hive.isBoxOpen(userBoxName)) {
        await _userBox.close();
      }
      final box = await Hive.openBox<UserModel>(userBoxName);

      // Save user and verify immediately
      await box.put('current_user', user);
      AppLogger.success('signUp: User saved to local storage');

      // Immediate verification
      final savedUser = box.get('current_user');
      if (savedUser != null) {
        AppLogger.debug('''
Immediate verification after save:
Box is open: ${Hive.isBoxOpen(userBoxName)}
Box length: ${box.length}
Box keys: ${box.keys.toList()}
User ID: ${savedUser.id}
User Name: ${savedUser.name}
User Email: ${savedUser.email}
''');
      } else {
        AppLogger.error('Failed to verify saved user immediately after save');
      }

      return user;
    } catch (e) {
      AppLogger.error('signUp: Error: $e');
      throw Exception('Failed to sign up. Please try again.');
    }
  }

  Future<UserModel> login(LoginRequest request) async {
    try {
      AppLogger.info('login: Making login request');
      final response = await _dio.post(
        '$apiUrl/auth/login',
        data: request.toJson(),
      );

      AppLogger.success('login: Login request successful');
      AppLogger.info('login: Saving user to local storage');

      final user = UserModel.fromJson(response.data);
      await _userBox.put('current_user', user);

      AppLogger.success('login: User saved to local storage');
      return user;
    } on DioException catch (e) {
      AppLogger.error(
        'login: DioException: ${e.message} at status code: ${e.response?.statusCode}',
      );
      throw _handleLoginDioError(e);
    }
  }

  Exception _handleLoginDioError(DioException e) {
    if (e.response?.statusCode == 401) {
      return Exception('Invalid email or password');
    }
    if (e.response?.statusCode == 400) {
      return Exception('Invalid input data');
    }
    return Exception('Failed to log in. Please try again.');
  }

  Exception _handleDioError(DioException e) {
    if (e.response?.statusCode == 409) {
      return Exception('Email already exists');
    }
    if (e.response?.statusCode == 400) {
      return Exception('Invalid input data');
    }
    return Exception('Failed to sign up. Please try again.');
  }

  Future<void> logout() async {
    await _userBox.delete('current_user');
  }

  UserModel? getCurrentUser() {
    try {
      if (!Hive.isBoxOpen(userBoxName)) {
        Hive.openBox<UserModel>(userBoxName);
      }
      return _userBox.get('current_user');
    } catch (e) {
      AppLogger.error('Error getting current user: $e');
      return null;
    }
  }

  // Add a method to verify user persistence
  Future<void> verifyUserPersistence() async {
    try {
      if (!Hive.isBoxOpen(userBoxName)) {
        await Hive.openBox<UserModel>(userBoxName);
      }

      final user = _userBox.get('current_user');
      if (user != null) {
        AppLogger.debug('''
Verified persisted user:
ID: ${user.id}
Name: ${user.name}
Email: ${user.email}
IsLoggedIn: ${user.isLoggedIn}
''');
      } else {
        AppLogger.warning('No user found in persistence');
      }
    } catch (e) {
      AppLogger.error('Error verifying user persistence: $e');
    }
  }

  // Add this method to help with debugging
  static Future<void> debugBoxContents() async {
    try {
      if (!Hive.isBoxOpen(userBoxName)) {
        await Hive.openBox<UserModel>(userBoxName);
      }

      final box = Hive.box<UserModel>(userBoxName);
      AppLogger.debug('''
Box Debug Information:
Box is open: ${Hive.isBoxOpen(userBoxName)}
Box length: ${box.length}
Box keys: ${box.keys.toList()}
Current user exists: ${box.get('current_user') != null}
''');
    } catch (e) {
      AppLogger.error('Error debugging box contents: $e');
    }
  }
}
