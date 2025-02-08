import 'package:fingenie/core/config/theme/app_themes.dart';
import 'package:fingenie/core/router/app_router.dart';
import 'package:fingenie/core/services/hive/service/contact_service.dart';
import 'package:fingenie/data/auth/auth_repository.dart';
import 'package:fingenie/domain/models/user_model.dart';
import 'package:fingenie/presentation/home/bloc/expense_bloc.dart';
import 'package:fingenie/presentation/onboarding/screens/intro_screen.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:hive_flutter/hive_flutter.dart';

void main() async {
  try {
    WidgetsFlutterBinding.ensureInitialized();
    await Hive.initFlutter();
    await AuthRepository.init();
    await dotenv.load(fileName: ".env");
    await ContactsService.initializeHive();
    if (!Hive.isAdapterRegistered(0)) {
      Hive.registerAdapter(UserModelAdapter());
    }
    final userBox = await Hive.openBox<UserModel>('userBox');

    final currentUser = userBox.get('current_user');
    if (currentUser != null) {
      AppLogger.debug('Current user found:');
      AppLogger.debug('ID: ${currentUser.id}');
      AppLogger.debug('Name: ${currentUser.name}');
      AppLogger.debug('Email: ${currentUser.email}');
    } else {
      AppLogger.debug('No current user found in the box');
    }
    runApp(const MyApp());
  } catch (e, stackTrace) {
    debugPrint('Initialization error: $e');
    debugPrint(stackTrace.toString());
    rethrow;
  }
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'FinGenie',
      debugShowCheckedModeBanner: false,
      theme: AppThemes.lightTheme,
      darkTheme: AppThemes.lightTheme,
      onGenerateRoute: AppRouter.generateRoute,
      home: MultiBlocProvider(
        providers: [
          BlocProvider(
            create: (context) => ExpenseBloc(),
          ),
        ],
        // child: const IntroScreen(),
        child: const IntroScreen(),
      ),
    );
  }
}
