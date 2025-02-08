import 'package:json_annotation/json_annotation.dart';

part 'group_model.g.dart';

@JsonSerializable()
class GroupModel {
  final String id;
  final String name;
  final double balance;
  final String icon;
  final String color;
  final List<String> memberIds;
  final DateTime createdAt;
  final String createdBy;
  final String tag;

  GroupModel({
    required this.id,
    required this.name,
    required this.balance,
    required this.icon,
    required this.color,
    required this.memberIds,
    required this.createdAt,
    required this.createdBy,
    required this.tag,
  });

  factory GroupModel.fromJson(Map<String, dynamic> json) =>
      _$GroupModelFromJson(json);

  Map<String, dynamic> toJson() => _$GroupModelToJson(this);

  copyWith({required List<String> memberIds}) {}
}
