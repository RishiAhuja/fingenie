import 'package:fingenie/domain/models/group_model.dart';
import 'package:fingenie/presentation/contacts/screens/contact_selection_screen.dart';
import 'package:fingenie/presentation/groups/bloc/group_bloc.dart';
import 'package:fingenie/presentation/groups/bloc/group_events.dart';
import 'package:fingenie/presentation/groups/bloc/group_state.dart';
import 'package:fingenie/utils/app_logger.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class GroupDetailScreen extends StatelessWidget {
  final GroupModel group;
  final String apiUrl; // Add this parameter

  const GroupDetailScreen({
    required this.group,
    required this.apiUrl, // Add this parameter
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    final parentGroupBloc = context.read<GroupBloc>();

    AppLogger.info(
        'GroupBloc in GroupDetailScreen: ${context.read<GroupBloc>()}');

    return BlocProvider.value(
      value: parentGroupBloc, // Use the existing GroupBloc
      child: Builder(
        builder: (context) => Scaffold(
          appBar: AppBar(
            title: Text(group.name),
            actions: [
              IconButton(
                icon: const Icon(Icons.person_add),
                onPressed: () {
                  final currentGroupBloc = context.read<GroupBloc>();

                  AppLogger.warning(
                      'GroupBloc before navigation: $currentGroupBloc');

                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => BlocProvider<GroupBloc>.value(
                        value: currentGroupBloc,
                        child: ContactSelectionScreen(
                          onContactsSelected: (contacts) {
                            AppLogger.info(
                                'Contacts selected: ${contacts.length}');
                            currentGroupBloc.add(AddGroupMembers(
                              groupId: group.id,
                              memberIds: contacts.map((e) => e.id).toList(),
                            ));
                          },
                        ),
                      ),
                    ),
                  );
                },
              ),
            ],
          ),
          body: BlocBuilder<GroupBloc, GroupState>(
            builder: (context, state) {
              if (state.isLoading) {
                return const Center(child: CircularProgressIndicator());
              }

              final selectedGroup = state.selectedGroup ?? group;
              return ListView.builder(
                itemCount: selectedGroup.memberIds.length,
                itemBuilder: (context, index) {
                  final memberId = selectedGroup.memberIds[index];
                  return ListTile(
                    leading: CircleAvatar(
                      child: Text(memberId[0].toUpperCase()),
                    ),
                    title: Text(memberId), // Replace with actual contact name
                  );
                },
              );
            },
          ),
        ),
      ),
    );
  }
}
