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
      AppLogger.debug('Profile: Box opened with ${userBox.length} items');

      final user = userBox.get('current_user');
      if (user != null) {
        setState(() {
          currentUser = user;
        });
        AppLogger.debug('''
Profile loaded user:
ID: ${user.id}
Name: ${user.name}
Email: ${user.email}
IsLoggedIn: ${user.isLoggedIn}
''');
      } else {
        AppLogger.warning('Profile: No user found in box');
      }
    } catch (e) {
      AppLogger.error('Profile: Error loading user data: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: currentUser != null
          ? Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(currentUser!.name,
                    style: Theme.of(context).textTheme.displayMedium),
                const SizedBox(
                  height: 5,
                ),
                Text('email: ${currentUser!.email}',
                    style: Theme.of(context).textTheme.bodyMedium),
                Text('phoneNumber: ${currentUser!.phoneNumber}',
                    style: Theme.of(context).textTheme.bodyMedium),
                Text('id: ${currentUser!.id}',
                    style: Theme.of(context).textTheme.bodyMedium),
                Text('isLoggedIn: ${currentUser!.isLoggedIn}',
                    style: Theme.of(context).textTheme.bodyMedium),
                Text('createdAt: ${currentUser!.createdAt}',
                    style: Theme.of(context).textTheme.bodyMedium),
              ],
            )
          : const Text('No user data found'),
    );
  }

  @override
  void dispose() {
    userBox.close();
    super.dispose();
  }
}
