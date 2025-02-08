import 'package:dio/dio.dart';
import 'package:fingenie/domain/entities/auth/create_login_request.dart';
import 'package:fingenie/domain/entities/auth/create_signup_request.dart';
import 'package:fingenie/domain/models/user_model.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:hive_flutter/hive_flutter.dart';

class AuthRepository {
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
      final response = await _dio.post(
        '$apiUrl/auth/signup',
        data: request.toJson(),
      );

      AppLogger.success('signUp: Signup request successful');
      AppLogger.info('signUp: Saving user to local storage');
      final user = UserModel.fromJson(response.data);
      await _userBox.put('current_user', user);
      AppLogger.success('signUp: User saved to local storage');
      return user;
    } on DioException catch (e) {
      AppLogger.error(
          'signUp: DioException: ${e.message} at status code: ${e.response?.statusCode}');
      throw _handleDioError(e);
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
    return _userBox.get('current_user');
  }
}
