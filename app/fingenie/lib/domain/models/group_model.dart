class GroupModel {
  final String id;
  final String name;
  final String tag;
  final List<String> memberIds;
  final DateTime createdAt;

  GroupModel({
    required this.id,
    required this.name,
    required this.tag,
    required this.memberIds,
    required this.createdAt,
  });

  Map<String, dynamic> toJson() => {
        'id': id,
        'name': name,
        'tag': tag,
        'member_ids': memberIds,
        'created_at': createdAt.toIso8601String(),
      };

  factory GroupModel.fromJson(Map<String, dynamic> json) => GroupModel(
        id: json['id'],
        name: json['name'],
        tag: json['tag'],
        memberIds: List<String>.from(json['member_ids']),
        createdAt: DateTime.parse(json['created_at']),
      );

  copyWith({required List<String> memberIds}) {}
}
