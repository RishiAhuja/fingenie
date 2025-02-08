abstract class GroupEvent {}

class LoadGroups extends GroupEvent {}

class CreateGroup extends GroupEvent {
  final String name;
  final String tag;
  final List<String> initialMembers;

  CreateGroup({
    required this.name,
    required this.tag,
    this.initialMembers = const [],
  });
}

class AddGroupMembers extends GroupEvent {
  final String groupId;
  final List<String> memberIds;

  AddGroupMembers({
    required this.groupId,
    required this.memberIds,
  });
}
