import 'package:fingenie/domain/models/group_model.dart';

class GroupState {
  final List<GroupModel> groups;
  final GroupModel? selectedGroup;
  final bool isLoading;
  final String? errorMessage;

  GroupState({
    this.groups = const [],
    this.selectedGroup,
    this.isLoading = false,
    this.errorMessage,
  });

  GroupState copyWith({
    List<GroupModel>? groups,
    GroupModel? selectedGroup,
    bool? isLoading,
    String? errorMessage,
  }) {
    return GroupState(
      groups: groups ?? this.groups,
      selectedGroup: selectedGroup ?? this.selectedGroup,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
    );
  }
}
