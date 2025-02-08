import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:fingenie/domain/models/user_model.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen> {
  late Box<UserModel> userBox;
  UserModel? currentUser;

  @override
  void initState() {
    super.initState();
    _loadUserData();
  }

  Future<void> _loadUserData() async {
    try {
      userBox = await Hive.openBox<UserModel>('userBox');

      // Print all users in the box
      AppLogger.debug('Number of users in box: ${userBox.length}');

      for (var i = 0; i < userBox.length; i++) {
        final user = userBox.getAt(i);
        if (user != null) {
          AppLogger.debug('User $i data:');
          AppLogger.debug('ID: ${user.id}');
          AppLogger.debug('Name: ${user.name}');
          AppLogger.debug('Email: ${user.email}');
          AppLogger.debug('Phone: ${user.phoneNumber}');
          AppLogger.debug('Created At: ${user.createdAt}');
          AppLogger.debug('Is Logged In: ${user.isLoggedIn}');
          AppLogger.debug('-------------------');
        }
      }

      // Get the first user (assuming single user storage)
      if (userBox.isNotEmpty) {
        setState(() {
          currentUser = userBox.getAt(0);
        });
      }
    } catch (e) {
      AppLogger.error('Error loading user data: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: currentUser != null
          ? Text('Welcome ${currentUser!.name}')
          : const Text('No user data found'),
    );
  }

  @override
  void dispose() {
    userBox.close();
    super.dispose();
  }
}
