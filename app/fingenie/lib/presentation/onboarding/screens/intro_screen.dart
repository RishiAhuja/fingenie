import 'package:dio/dio.dart';
import 'package:fingenie/data/groups/group_repository.dart';
import 'package:fingenie/presentation/groups/bloc/group_bloc.dart';
import 'package:fingenie/presentation/home/bloc/expense_bloc.dart';
import 'package:fingenie/presentation/home/screens/home.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:introduction_screen/introduction_screen.dart';
import 'package:fingenie/core/config/theme/app_colors.dart';
import 'package:lottie/lottie.dart';
import 'package:shared_preferences/shared_preferences.dart';

class IntroScreen extends StatefulWidget {
  const IntroScreen({super.key});

  @override
  State<IntroScreen> createState() => _IntroScreenState();
}

class _IntroScreenState extends State<IntroScreen> {
  @override
  void initState() {
    super.initState();
    _checkIntroStatus();
  }

  Future<void> _checkIntroStatus() async {
    final prefs = await SharedPreferences.getInstance();
    final hasSeenIntro = prefs.getBool('hasSeenIntro') ?? false;

    if (hasSeenIntro && mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (context) => MultiBlocProvider(
            providers: [
              BlocProvider(
                create: (context) => ExpenseBloc(),
              ),
              BlocProvider(
                create: (context) => GroupBloc(
                    repository: GroupRepository(
                        dio: Dio(), apiUrl: dotenv.env['API_URL'] ?? ''),
                    apiUrl: dotenv.env['API_URL'] ?? ''),
              ),
            ],
            child: const HomeScreen(),
          ),
        ),
      );
    }
  }

  Future<void> _onIntroEnd(BuildContext context) async {
    // Save that user has seen the intro
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('hasSeenIntro', true);
    if (context.mounted) {
      Navigator.of(context).pushReplacementNamed('/login');
    }
  }

  @override
  Widget build(BuildContext context) {
    final titleStyle = Theme.of(context).textTheme.headlineMedium?.copyWith(
          color: AppColors.textPrimary,
          fontWeight: FontWeight.bold,
        );

    final bodyStyle = Theme.of(context).textTheme.bodyLarge?.copyWith(
          color: AppColors.textSecondary,
        );

    Widget _buildImage(String assetName) {
      return Container(
        height: 300,
        margin: const EdgeInsets.only(top: 40),
        // decoration: BoxDecoration(
        //   borderRadius: BorderRadius.circular(20),
        //   boxShadow: [
        //     BoxShadow(
        //       // ignore: deprecated_member_use
        //       color: Colors.black.withOpacity(0.1),
        //       blurRadius: 20,
        //       offset: const Offset(0, 10),
        //     ),
        //   ],
        // ),
        child: Lottie.asset(
          'assets/intro/$assetName',
          fit: BoxFit.contain,
        ),
      );
    }

    return IntroductionScreen(
      pages: [
        PageViewModel(
          title: "Track Your Expenses",
          body:
              "Keep track of your spending habits and stay on top of your finances with easy expense tracking.",
          image: _buildImage('expense_tracking.json'),
          decoration: PageDecoration(
            titleTextStyle: titleStyle ?? const TextStyle(),
            bodyTextStyle: bodyStyle ?? const TextStyle(),
            bodyPadding: const EdgeInsets.symmetric(horizontal: 16),
            imagePadding: const EdgeInsets.only(top: 40),
          ),
        ),
        PageViewModel(
          title: "Smart Analytics",
          body:
              "Get insights into your spending patterns with detailed analytics and visualization tools.",
          image: _buildImage('analytics.json'),
          decoration: PageDecoration(
            titleTextStyle: titleStyle ?? const TextStyle(),
            bodyTextStyle: bodyStyle ?? const TextStyle(),
            bodyPadding: const EdgeInsets.symmetric(horizontal: 16),
            imagePadding: const EdgeInsets.only(top: 40),
          ),
        ),
        PageViewModel(
          title: "Split Expenses",
          body:
              "Easily split bills and expenses with friends and family. Keep track of who owes what.",
          image: _buildImage('split_bills.json'),
          decoration: PageDecoration(
            titleTextStyle: titleStyle ?? const TextStyle(),
            bodyTextStyle: bodyStyle ?? const TextStyle(),
            bodyPadding: const EdgeInsets.symmetric(horizontal: 16),
            imagePadding: const EdgeInsets.only(top: 40),
          ),
        ),
        PageViewModel(
          title: "Secure & Private",
          body:
              "Your financial data is encrypted and secure. We prioritize your privacy and data protection.",
          image: _buildImage('security.json'),
          decoration: PageDecoration(
            titleTextStyle: titleStyle ?? const TextStyle(),
            bodyTextStyle: bodyStyle ?? const TextStyle(),
            bodyPadding: const EdgeInsets.symmetric(horizontal: 16),
            imagePadding: const EdgeInsets.only(top: 40),
          ),
        ),
      ],
      onDone: () => _onIntroEnd(context),
      showSkipButton: true,
      skipOrBackFlex: 0,
      nextFlex: 0,
      skip: Text(
        'Skip',
        style: Theme.of(context)
            .textTheme
            .bodySmall
            ?.copyWith(color: AppColors.primary),
      ),
      next: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: AppColors.primary,
          borderRadius: BorderRadius.circular(30),
        ),
        child: const Icon(
          Icons.arrow_forward,
          color: Colors.white,
        ),
      ),
      done: Container(
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
        decoration: BoxDecoration(
          color: AppColors.primary,
          borderRadius: BorderRadius.circular(30),
        ),
        child: Text(
          'Get Started',
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              ),
        ),
      ),
      curve: Curves.fastLinearToSlowEaseIn,
      controlsMargin: const EdgeInsets.all(16),
      controlsPadding: const EdgeInsets.fromLTRB(8.0, 4.0, 8.0, 4.0),
      dotsDecorator: DotsDecorator(
        size: const Size(10.0, 10.0),
        color: Colors.grey,
        activeSize: const Size(22.0, 10.0),
        activeColor: AppColors.primary,
        activeShape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(25.0),
        ),
      ),
    );
  }
}
