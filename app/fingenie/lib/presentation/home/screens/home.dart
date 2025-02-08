// home_screen.dart
import 'package:fingenie/presentation/contacts/screens/contacts.dart';
import 'package:fingenie/presentation/groups/bloc/group_bloc.dart';
import 'package:fingenie/presentation/groups/screens/creat_group_modal.dart';
import 'package:fingenie/presentation/ocr/screens/ocr.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:intl/intl.dart';

import '../bloc/expense_bloc.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  void initState() {
    super.initState();
    context.read<ExpenseBloc>().add(LoadExpenses());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        actions: [
          IconButton(
            icon: const Icon(Icons.group_add),
            onPressed: () {
              showModalBottomSheet(
                context: context,
                isScrollControlled: true,
                backgroundColor: Colors.transparent,
                builder: (context) => Padding(
                  padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom,
                  ),
                  child: BlocProvider(
                    create: (context) =>
                        GroupBloc(apiUrl: dotenv.env['API_URL'] ?? ''),
                    child: const CreateGroupModal(),
                  ),
                ),
              );
            },
          ),
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () {},
          ),
        ],
        title: const Text(
          'Fin Genie.',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
          ),
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: BlocBuilder<ExpenseBloc, ExpenseState>(
            builder: (context, state) {
              if (state is ExpenseLoading) {
                return const Center(child: CircularProgressIndicator());
              }
              if (state is ExpenseLoaded) {
                return Column(
                  children: [
                    const SizedBox(height: 20),
                    GradientCard(
                      balance: 20.0,
                      // balance: state.totalBalance,
                      // onAddExpense: () => _showAddExpenseDialog(context),
                      onAddExpense: () {},
                    ),
                    const SizedBox(height: 20),
                    _buildBalanceCards(state),
                    const SizedBox(height: 20),
                    _buildExpensesList(state.expenses),
                    const SizedBox(
                      height: 20,
                    ),
                    TextButton(
                        onPressed: () {
                          Navigator.push(
                              context,
                              MaterialPageRoute(
                                  builder: (context) => const OcrScreen()));
                        },
                        child: const Text('OCR')),
                    TextButton(
                        onPressed: () {
                          Navigator.push(
                              context,
                              MaterialPageRoute(
                                  builder: (context) =>
                                      const ContactsScreen()));
                        },
                        child: const Text('Contacts'))
                  ],
                );
              }
              return const Center(child: Text('Something went wrong'));
            },
          ),
        ),
      ),
    );
  }
}

Widget _buildBalanceCards(ExpenseState state) {
  return const Row(
    mainAxisAlignment: MainAxisAlignment.spaceBetween,
    children: [
      BalanceCard(
        title: 'YOU OWE',
        amount: 100.0,
        // amount: state.,
        subtitle: 'You should Pay to others',
        icon: Icons.arrow_upward,
        iconColor: Colors.red,
      ),
      BalanceCard(
        title: 'YOU OWED',
        // amount: state.youOwed,
        amount: 50.0,
        subtitle: 'Others should Pay to you',
        icon: Icons.arrow_downward,
        iconColor: Colors.green,
      ),
    ],
  );
}

Widget _buildExpensesList(List<int> expenses) {
  return Expanded(
    child: ListView.separated(
      itemCount: expenses.length,
      separatorBuilder: (context, index) => const Divider(height: 1),
      itemBuilder: (context, index) {
        final expense = expenses[index];
        return ExpenseItem(
          icon: 'expense.icon',
          title: 'expense.title',
          amount: 100.0,
          date: DateTime(10, 10, 10),
          onTap: () => Navigator.pushNamed(
            context,
            '/expense-detail',
            arguments: expense,
          ),
        );
      },
    ),
  );
}

class BalanceCard extends StatelessWidget {
  final String title;
  final double amount;
  final String subtitle;
  final IconData icon;
  final Color iconColor;

  const BalanceCard({
    super.key,
    required this.title,
    required this.amount,
    required this.subtitle,
    required this.icon,
    required this.iconColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: MediaQuery.of(context).size.width * 0.43,
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: Offset(0, 5),
          ),
        ],
      ),
      child: Column(
        children: [
          CircleAvatar(
            backgroundColor: iconColor.withOpacity(0.1),
            child: Icon(icon, color: iconColor),
          ),
          SizedBox(height: 8),
          Text(
            title,
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '\$${amount.toStringAsFixed(2)}',
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            subtitle,
            style: TextStyle(
              fontSize: 10,
              color: Colors.grey[500],
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}

class ExpenseItem extends StatelessWidget {
  final String icon;
  final String title;
  final double amount;
  final DateTime date;
  final VoidCallback onTap;

  const ExpenseItem({
    super.key,
    required this.icon,
    required this.title,
    required this.amount,
    required this.date,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: EdgeInsets.symmetric(vertical: 12, horizontal: 16),
        child: Row(
          children: [
            Container(
              padding: EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.grey[100],
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text(
                icon,
                style: TextStyle(fontSize: 24),
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  Text(
                    DateFormat('dd MMM - HH:mm').format(date),
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey[600],
                    ),
                  ),
                ],
              ),
            ),
            Text(
              amount > 0
                  ? '+\$${amount.toStringAsFixed(2)}'
                  : '-\$${amount.abs().toStringAsFixed(2)}',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w500,
                color: amount > 0 ? Colors.green : Colors.black,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// gradient_card.dart
class GradientCard extends StatelessWidget {
  final double balance;
  final VoidCallback onAddExpense;

  const GradientCard({
    Key? key,
    required this.balance,
    required this.onAddExpense,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            Color(0xFF2DD4BF),
            Color(0xFFFCD34D),
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(24),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Total Balance',
            style: TextStyle(
              color: Colors.white.withOpacity(0.8),
              fontSize: 14,
            ),
          ),
          SizedBox(height: 8),
          Text(
            '\$${balance.toStringAsFixed(2)}',
            style: TextStyle(
              color: Colors.white,
              fontSize: 32,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: 20),
          ElevatedButton(
            onPressed: onAddExpense,
            child: Text('+ ADD EXPENSE'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.black,
              foregroundColor: Colors.white,
              padding: EdgeInsets.symmetric(horizontal: 24, vertical: 12),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
