import 'package:fingenie/data/groups/group_repository.dart';
import 'package:fingenie/domain/models/group_model.dart';
import 'package:fingenie/presentation/groups/bloc/group_events.dart';
import 'package:fingenie/presentation/groups/bloc/group_state.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class GroupBloc extends Bloc<GroupEvent, GroupState> {
  final GroupRepository repository;
  final String apiUrl;

  GroupBloc({
    required this.repository,
    required this.apiUrl,
  }) : super(GroupState()) {
    on<LoadGroups>(_onLoadGroups);
    on<CreateGroup>(_onCreateGroup);
    on<AddGroupMembers>(_onAddGroupMembers);
  }

  Future<void> _onLoadGroups(
    LoadGroups event,
    Emitter<GroupState> emit,
  ) async {
    emit(state.copyWith(isLoading: true, errorMessage: null));

    try {
      // final response = await _dio.get('$apiUrl/groups');
      // final List<GroupModel> groups = (response.data as List)
      //     .map((json) => GroupModel.fromJson(json))
      //     .toList();

      // emit(state.copyWith(
      //   isLoading: false,
      //   groups: groups,
      // ));
    } catch (e) {
      emit(state.copyWith(
        isLoading: false,
        errorMessage: 'Failed to load groups: ${e.toString()}',
      ));
    }
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
        balance: 0,
        icon: '',
        color: '',
        createdBy: 'current-user-id',
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
      await repository.addGroupMembers(
        groupId: event.groupId,
        memberIds: event.memberIds,
      );

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
