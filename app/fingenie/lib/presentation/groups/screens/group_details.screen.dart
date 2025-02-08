import 'package:fingenie/domain/models/group_model.dart';
import 'package:fingenie/presentation/contacts/bloc/contact_bloc.dart';
import 'package:fingenie/presentation/contacts/bloc/contact_events.dart';
import 'package:fingenie/presentation/contacts/screens/contact_selection_screen.dart';
import 'package:fingenie/presentation/groups/bloc/group_bloc.dart';
import 'package:fingenie/presentation/groups/bloc/group_events.dart';
import 'package:fingenie/presentation/groups/bloc/group_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class GroupDetailScreen extends StatelessWidget {
  final GroupModel group;

  const GroupDetailScreen({
    required this.group,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(group.name),
        actions: [
          IconButton(
            icon: Icon(Icons.person_add),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => BlocProvider(
                    create: (context) =>
                        ContactsBloc()..add(FetchContactsEvent()),
                    child: ContactSelectionScreen(
                      onContactsSelected: (selectedContacts) {
                        context.read<GroupBloc>().add(
                              AddGroupMembers(
                                groupId: group.id,
                                memberIds:
                                    selectedContacts.map((c) => c.id).toList(),
                              ),
                            );
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
              // You'll need to fetch contact details based on memberId
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
    );
  }
}
