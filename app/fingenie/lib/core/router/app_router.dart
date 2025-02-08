import 'package:dio/dio.dart';
import 'package:fingenie/data/auth/auth_repository.dart';
import 'package:fingenie/domain/models/user_model.dart';
import 'package:fingenie/presentation/auth/bloc/signup_bloc/signup_bloc.dart';
import 'package:fingenie/presentation/auth/screens/login.dart';
import 'package:fingenie/presentation/auth/screens/signup.dart';
import 'package:fingenie/presentation/onboarding/screens/intro_screen.dart';
import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fingenie/presentation/auth/bloc/login_bloc/login_bloc.dart';

class AppRouter {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    // Get dependencies
    final dio = Dio();
    final userBox = Hive.box<UserModel>('userBox');

    switch (settings.name) {
      case '/':
        return MaterialPageRoute(
          builder: (_) => FutureBuilder<SharedPreferences>(
            future: SharedPreferences.getInstance(),
            builder: (context, snapshot) {
              if (!snapshot.hasData) {
                return const CircularProgressIndicator();
              }

              final hasSeenIntro =
                  snapshot.data!.getBool('hasSeenIntro') ?? false;

              if (hasSeenIntro) {
                return const LoginScreen();
              } else {
                return const IntroScreen();
              }
            },
          ),
        );
      case '/login':
        return MaterialPageRoute(
          builder: (_) => BlocProvider(
            create: (context) => LoginBloc(
              authRepository: AuthRepository(
                userBox: userBox,
                dio: dio,
              ),
              dio: dio,
            ),
            child: const LoginScreen(),
          ),
        );
      case '/signup':
        return MaterialPageRoute(
          builder: (_) => BlocProvider(
            create: (context) => SignUpBloc(
              authRepository: AuthRepository(
                userBox: userBox,
                dio: dio,
              ),
            ),
            child: const SignUpScreen(),
          ),
        );
      default:
        return MaterialPageRoute(
          builder: (_) => const Scaffold(
            body: Center(child: Text('Route not found')),
          ),
        );
    }
  }
}
