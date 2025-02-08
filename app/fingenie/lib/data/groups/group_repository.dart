import 'package:dio/dio.dart';
import 'package:fingenie/domain/models/group_model.dart';
// import 'package:flutter_dotenv/flutter_dotenv.dart';

class GroupRepository {
  final Dio _dio;
  final String apiUrl;

  GroupRepository({
    required Dio dio,
    required this.apiUrl,
  }) : _dio = dio;

  Future<List<GroupModel>> fetchUserGroups() async {
    try {
      final response = await _dio.get('$apiUrl/groups');

      if (response.statusCode == 200 || response.statusCode == 201) {
        final List<dynamic> data = response.data['groups'];
        return data.map((json) => GroupModel.fromJson(json)).toList();
      }

      throw Exception('Failed to fetch groups');
    } on DioException catch (e) {
      throw Exception('Network error: ${e.message}');
    }
  }

  Future<GroupModel> createGroup({
    required String name,
    required List<String> memberIds,
    required String icon,
    required String color,
  }) async {
    try {
      final response = await _dio.post(
        '$apiUrl/groups',
        data: {
          'name': name,
          'memberIds': memberIds,
          'icon': icon,
          'color': color,
        },
      );

      if (response.statusCode == 201) {
        return GroupModel.fromJson(response.data['group']);
      }

      throw Exception('Failed to create group');
    } on DioException catch (e) {
      throw Exception('Network error: ${e.message}');
    }
  }

  Future<void> addGroupMembers({
    required String groupId,
    required List<String> memberIds,
  }) async {
    try {
      await _dio.post(
        '$apiUrl/groups/$groupId/members',
        data: {
          'member_ids': memberIds,
        },
      );
    } catch (e) {
      throw Exception('Failed to add members: ${e.toString()}');
    }
  }
}
