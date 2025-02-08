import 'package:dio/dio.dart';
import 'package:fingenie/domain/models/group_model.dart';
import 'package:fingenie/presentation/groups/bloc/group_events.dart';
import 'package:fingenie/presentation/groups/bloc/group_state.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class GroupBloc extends Bloc<GroupEvent, GroupState> {
  final Dio _dio;
  final String apiUrl;

  GroupBloc({
    required this.apiUrl,
    Dio? dio,
  })  : _dio = dio ?? Dio(),
        super(GroupState()) {
    on<CreateGroup>(_onCreateGroup);
    on<AddGroupMembers>(_onAddGroupMembers);
  }

  Future<void> _onCreateGroup(
    CreateGroup event,
    Emitter<GroupState> emit,
  ) async {
    emit(state.copyWith(isLoading: true, errorMessage: null));

    try {
      // final response = await _dio.post(
      //   '$apiUrl/groups',
      //   data: {
      //     'name': event.name,
      //     'tag': event.tag,
      //     'member_ids': event.initialMembers,
      //   },
      // );

      // final newGroup = GroupModel.fromJson(response.data);
      final newGroup = GroupModel(
        id: 'new-group-id',
        name: event.name,
        tag: event.tag,
        memberIds: event.initialMembers,
        createdAt: DateTime.now(),
      );
      emit(state.copyWith(
        isLoading: false,
        groups: [...state.groups, newGroup],
        selectedGroup: newGroup,
      ));
    } catch (e) {
      emit(state.copyWith(
        isLoading: false,
        errorMessage: 'Failed to create group: ${e.toString()}',
      ));
    }
  }

  Future<void> _onAddGroupMembers(
    AddGroupMembers event,
    Emitter<GroupState> emit,
  ) async {
    emit(state.copyWith(isLoading: true, errorMessage: null));

    try {
      await _dio.post(
        '$apiUrl/groups/${event.groupId}/members',
        data: {
          'member_ids': event.memberIds,
        },
      );

      // Update the selected group with new members
      final updatedGroup = state.selectedGroup?.copyWith(
        memberIds: [...state.selectedGroup!.memberIds, ...event.memberIds],
      );

      emit(state.copyWith(
        isLoading: false,
        selectedGroup: updatedGroup,
      ));
    } catch (e) {
      emit(state.copyWith(
        isLoading: false,
        errorMessage: 'Failed to add members: ${e.toString()}',
      ));
    }
  }
}
