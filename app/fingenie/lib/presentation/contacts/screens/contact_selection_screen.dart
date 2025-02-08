import 'package:fingenie/presentation/contacts/bloc/contact_bloc.dart';
import 'package:fingenie/presentation/contacts/bloc/contact_events.dart';
import 'package:fingenie/presentation/contacts/bloc/contact_states.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:fingenie/core/config/theme/app_colors.dart';
import 'package:flutter_contacts/contact.dart';

class ContactSelectionScreen extends StatefulWidget {
  final Function(List<Contact>) onContactsSelected;

  const ContactSelectionScreen({
    super.key,
    required this.onContactsSelected,
  });

  @override
  // ignore: library_private_types_in_public_api
  _ContactSelectionScreenState createState() => _ContactSelectionScreenState();
}

class _ContactSelectionScreenState extends State<ContactSelectionScreen> {
  final ValueNotifier<Set<Contact>> _selectedContacts =
      ValueNotifier<Set<Contact>>({});
  String _searchQuery = '';
  final TextEditingController _searchController = TextEditingController();

  @override
  void dispose() {
    _selectedContacts.dispose();
    _searchController.dispose();
    super.dispose();
  }

  void _toggleContact(Contact contact) {
    final updatedSet = Set<Contact>.from(_selectedContacts.value);
    if (updatedSet.contains(contact)) {
      updatedSet.remove(contact);
    } else {
      updatedSet.add(contact);
    }
    _selectedContacts.value = updatedSet;
  }

  @override
  Widget build(BuildContext context) {
    return BlocProvider(
      create: (context) => ContactsBloc()..add(FetchContactsEvent()),
      child: Scaffold(
        appBar: AppBar(
          title: Text(
            'Select Members',
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.w600,
                ),
          ),
          actions: [
            BlocBuilder<ContactsBloc, ContactsState>(
              builder: (context, state) {
                if (state.isLoading) {
                  return const Padding(
                    padding: EdgeInsets.all(16.0),
                    child: SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                      ),
                    ),
                  );
                }
                return IconButton(
                  icon: const Icon(Icons.refresh),
                  onPressed: () =>
                      context.read<ContactsBloc>().add(RefreshContactsEvent()),
                );
              },
            ),
          ],
          bottom: PreferredSize(
            preferredSize: const Size.fromHeight(70),
            child: Padding(
              padding: const EdgeInsets.fromLTRB(16, 0, 16, 8),
              child: SearchBar(
                controller: _searchController,
                onChanged: (value) {
                  setState(() {
                    _searchQuery = value.toLowerCase();
                  });
                },
                leading: const Icon(Icons.search),
                trailing: [
                  if (_searchQuery.isNotEmpty)
                    IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        _searchController.clear();
                        setState(() {
                          _searchQuery = '';
                        });
                      },
                    ),
                ],
                hintText: 'Search contacts...',
              ),
            ),
          ),
        ),
        body: BlocBuilder<ContactsBloc, ContactsState>(
          builder: (context, state) {
            if (state.isLoading && state.contacts.isEmpty) {
              return const Center(child: CircularProgressIndicator());
            }

            if (state.errorMessage != null) {
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.error_outline,
                      size: 48,
                      color: Theme.of(context).colorScheme.error,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      state.errorMessage!,
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.error,
                      ),
                    ),
                    const SizedBox(height: 16),
                    ElevatedButton.icon(
                      onPressed: () => context
                          .read<ContactsBloc>()
                          .add(FetchContactsEvent()),
                      icon: const Icon(Icons.refresh),
                      label: const Text('Retry'),
                    ),
                  ],
                ),
              );
            }

            if (state.contacts.isEmpty) {
              return const Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.contacts_outlined,
                      size: 48,
                      color: Colors.grey,
                    ),
                    SizedBox(height: 16),
                    Text(
                      'No contacts found',
                      style: TextStyle(
                        color: Colors.grey,
                        fontSize: 16,
                      ),
                    ),
                  ],
                ),
              );
            }

            final filteredContacts = state.contacts.where((contact) {
              return contact.displayName.toLowerCase().contains(_searchQuery);
            }).toList();

            // Group contacts by first letter
            final groupedContacts = <String, List<Contact>>{};
            for (var contact in filteredContacts) {
              if (contact.displayName.isEmpty) {
                continue; // Skip contacts with empty names
              }
              final firstLetter = contact.displayName[0].toUpperCase();
              groupedContacts.putIfAbsent(firstLetter, () => []);
              groupedContacts[firstLetter]!.add(contact);
            }

            final sortedKeys = groupedContacts.keys.toList()..sort();

            return CustomScrollView(
              slivers: [
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: ValueListenableBuilder<Set<Contact>>(
                      valueListenable: _selectedContacts,
                      builder: (context, selectedContacts, _) {
                        return Text(
                          'Selected: ${selectedContacts.length}',
                          style:
                              Theme.of(context).textTheme.bodyLarge?.copyWith(
                                    color: AppColors.primary,
                                    fontWeight: FontWeight.w600,
                                  ),
                        );
                      },
                    ),
                  ),
                ),
                for (final key in sortedKeys) ...[
                  SliverPersistentHeader(
                    pinned: true,
                    delegate: _SectionHeaderDelegate(
                      key,
                      Theme.of(context).colorScheme.surface,
                    ),
                  ),
                  SliverList(
                    delegate: SliverChildBuilderDelegate(
                      (context, index) {
                        final contact = groupedContacts[key]![index];
                        return ContactTile(
                          contact: contact,
                          isSelected: _selectedContacts.value.contains(contact),
                          onTap: () => _toggleContact(contact),
                        );
                      },
                      childCount: groupedContacts[key]!.length,
                    ),
                  ),
                ],
              ],
            );
          },
        ),
        floatingActionButton: ValueListenableBuilder<Set<Contact>>(
          valueListenable: _selectedContacts,
          builder: (context, selectedContacts, _) {
            return AnimatedOpacity(
              opacity: selectedContacts.isNotEmpty ? 1.0 : 0.0,
              duration: const Duration(milliseconds: 200),
              child: FloatingActionButton.extended(
                onPressed: selectedContacts.isEmpty
                    ? null
                    : () {
                        widget.onContactsSelected(selectedContacts.toList());
                        Navigator.pop(context);
                      },
                icon: const Icon(Icons.check),
                label: Text('Add ${selectedContacts.length} Members'),
              ),
            );
          },
        ),
      ),
    );
  }
}

class _SectionHeaderDelegate extends SliverPersistentHeaderDelegate {
  final String title;
  final Color backgroundColor;

  _SectionHeaderDelegate(this.title, this.backgroundColor);

  @override
  Widget build(
      BuildContext context, double shrinkOffset, bool overlapsContent) {
    return Container(
      color: backgroundColor,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
        child: Text(
          title,
          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                color: AppColors.textSecondary,
                fontWeight: FontWeight.w600,
              ),
        ),
      ),
    );
  }

  @override
  double get maxExtent => 40;

  @override
  double get minExtent => 40;

  @override
  bool shouldRebuild(covariant _SectionHeaderDelegate oldDelegate) {
    return title != oldDelegate.title ||
        backgroundColor != oldDelegate.backgroundColor;
  }
}

class ContactTile extends StatelessWidget {
  final Contact contact;
  final bool isSelected;
  final VoidCallback onTap;

  const ContactTile({
    Key? key,
    required this.contact,
    required this.isSelected,
    required this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: isSelected ? AppColors.primary : Colors.grey[200],
                shape: BoxShape.circle,
              ),
              child: Center(
                child: isSelected
                    ? const Icon(
                        Icons.check,
                        color: Colors.white,
                        size: 20,
                      )
                    : Text(
                        contact.displayName.isNotEmpty
                            ? contact.displayName[0].toUpperCase()
                            : '#', // Fallback for empty names
                        style: TextStyle(
                          color: isSelected ? Colors.white : Colors.grey[800],
                          fontSize: 16,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    contact.displayName,
                    style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          fontWeight: FontWeight.w500,
                        ),
                  ),
                  if (contact.phones.isNotEmpty)
                    Text(
                      contact.phones.first.number,
                      style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: AppColors.textSecondary,
                          ),
                    ),
                ],
              ),
            ),
            AnimatedContainer(
              duration: const Duration(milliseconds: 200),
              width: 24,
              height: 24,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(
                  color: isSelected ? AppColors.primary : Colors.grey[400]!,
                  width: 2,
                ),
                color: isSelected ? AppColors.primary : Colors.transparent,
              ),
              child: isSelected
                  ? const Icon(
                      Icons.check,
                      size: 16,
                      color: Colors.white,
                    )
                  : null,
            ),
          ],
        ),
      ),
    );
  }
}
