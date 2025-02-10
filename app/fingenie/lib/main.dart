import 'package:dio/dio.dart';
import 'package:fingenie/core/config/theme/app_themes.dart';
import 'package:fingenie/core/router/app_router.dart';
import 'package:fingenie/core/services/hive/service/contact_service.dart';
import 'package:fingenie/data/auth/auth_repository.dart';
import 'package:fingenie/data/groups/group_repository.dart';
import 'package:fingenie/domain/models/user_model.dart';
import 'package:fingenie/presentation/groups/bloc/group_bloc.dart';
import 'package:fingenie/presentation/home/bloc/expense_bloc.dart';
import 'package:fingenie/presentation/home/screens/home.dart';
import 'package:fingenie/presentation/onboarding/screens/intro/intro_screen.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:hive_flutter/hive_flutter.dart';

void main() async {
  try {
    WidgetsFlutterBinding.ensureInitialized();
    await Hive.initFlutter();
    await dotenv.load(fileName: ".env");

    if (!Hive.isAdapterRegistered(3)) {
      Hive.registerAdapter(UserModelAdapter());
    }

    await ContactsService.initializeHive();
    await Hive.openBox<UserModel>('userBox');

    await AuthRepository.init();

    // Open user box
    final userBox = await Hive.openBox<UserModel>('userBox');
    final apiUrl = dotenv.env['API_URL'] ?? '';

    // Check for API key
    final apiKey = dotenv.env['GEMINI_API_KEY'];
    if (apiKey == null || apiKey.isEmpty) {
      AppLogger.error('GEMINI_API_KEY not found in environment variables');
      throw Exception('GEMINI_API_KEY not found in environment variables');
    }

    // Check for current user
    final currentUser = userBox.get('current_user');
    if (currentUser != null) {
      AppLogger.debug('''
Current user found:
ID: ${currentUser.id}
Name: ${currentUser.name}
Email: ${currentUser.email}
IsLoggedIn: ${currentUser.isLoggedIn ?? false}
''');
    } else {
      AppLogger.debug('No current user found in the box');
    }

    // Determine initial screen
    Widget initialScreen;
    if (currentUser != null && currentUser.isLoggedIn == true) {
      AppLogger.debug('Logged in user found, redirecting to HomeScreen');
      initialScreen = MultiBlocProvider(
        providers: [
          BlocProvider(
            create: (context) => ExpenseBloc(),
          ),
          BlocProvider(
            create: (context) => GroupBloc(
              repository: GroupRepository(
                dio: Dio(),
                apiUrl: apiUrl,
              ),
              apiUrl: apiUrl,
            ),
          ),
        ],
        child: const HomeScreen(),
      );
    } else {
      AppLogger.debug('No logged in user found, showing IntroScreen');
      initialScreen = const IntroScreen();
    }

    runApp(MyApp(initialScreen: initialScreen));
  } catch (e, stackTrace) {
    AppLogger.error('Initialization error: $e');
    AppLogger.error('Stack trace: $stackTrace');
    rethrow;
  }
}

class MyApp extends StatelessWidget {
  final Widget initialScreen;

  const MyApp({
    super.key,
    required this.initialScreen,
  });

  @override
  Widget build(BuildContext context) {
    final apiUrl = dotenv.env['API_URL'] ?? '';

    // Create repositories
    final groupRepository = GroupRepository(
      dio: Dio(),
      apiUrl: apiUrl,
    );

    return MaterialApp(
      title: 'FinGenie',
      debugShowCheckedModeBanner: false,
      theme: AppThemes.lightTheme,
      darkTheme: AppThemes.lightTheme,
      onGenerateRoute: AppRouter.generateRoute,
      home: RepositoryProvider(
        create: (context) => groupRepository,
        child: MultiBlocProvider(
          providers: [
            BlocProvider<ExpenseBloc>(
              create: (context) => ExpenseBloc(),
            ),
            BlocProvider<GroupBloc>(
              create: (context) => GroupBloc(
                apiUrl: apiUrl,
                repository: groupRepository,
              ),
            ),
          ],
          child: initialScreen,
        ),
      ),
    );
  }
}
