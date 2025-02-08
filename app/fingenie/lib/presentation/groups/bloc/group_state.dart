import 'package:fingenie/domain/models/group_model.dart';

class GroupState {
  final bool isLoading;
  final String? errorMessage;
  final List<GroupModel> groups;
  final GroupModel? selectedGroup;

  GroupState({
    this.isLoading = false,
    this.errorMessage,
    this.groups = const [],
    this.selectedGroup,
  });

  GroupState copyWith({
    bool? isLoading,
    String? errorMessage,
    List<GroupModel>? groups,
    GroupModel? selectedGroup,
  }) =>
      GroupState(
        isLoading: isLoading ?? this.isLoading,
        errorMessage: errorMessage,
        groups: groups ?? this.groups,
        selectedGroup: selectedGroup ?? this.selectedGroup,
      );
}
