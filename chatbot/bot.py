from telegram import Update
from telegram.ext import Updater, CommandHandler, MessageHandler, filters, CallbackContext, Application
import requests
import os
from datetime import datetime

# Telegram Bot Token (Replace with your bot token)
BOT_TOKEN = "7662156020:AAHzqezB046bG8oOA12T4FghdZjXdVPHx5M"
BACKEND_URL = "http://127.0.0.1:5000"

# Start command
async def start(update: Update, context: CallbackContext):
    await update.message.reply_text("Welcome! You can add expenses like this:\n\n"
                                    "/add 500 Food Lunch at cafe\n"
                                    "Or type /summary to get your expense summary.")

# Add expense command
async def add_expense(update: Update, context: CallbackContext):
    try:
        args = context.args
        if len(args) < 2:
            await update.message.reply_text("Usage: /add <amount> <category> [description]")
            return

        user_id = str(update.message.chat_id)
        amount = float(args[0])
        category = args[1]
        description = " ".join(args[2:]) if len(args) > 2 else ""

        # Send data to backend
        data = {
            "user_id": user_id,
            "amount": amount,
            "category": category,
            "description": description,
            "date": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        }
        response = requests.post(f"{BACKEND_URL}/add_expense", json=data)

        if response.status_code == 200:
            await update.message.reply_text("Expense added successfully! ✅")
        else:
            await update.message.reply_text("Error adding expense. ❌")

    except ValueError:
        await update.message.reply_text("Invalid amount! Please enter a number.")

# Expense summary command
async def summary(update: Update, context: CallbackContext):
    user_id = str(update.message.chat_id)
    response = requests.get(f"{BACKEND_URL}/get_expenses/{user_id}")

    if response.status_code == 200:
        expenses = response.json()
        if not expenses:
            await update.message.reply_text("No expenses recorded yet.")
        else:
            summary_text = "📊 Expense Summary:\n"
            for category, amount in expenses.items():
                summary_text += f"{category}: ₹{amount}\n"
            await update.message.reply_text(summary_text)
    else:
        await update.message.reply_text("Error fetching summary. ❌")

# Main function to run bot
def main():
    app = Application.builder().token(BOT_TOKEN).build()

    app.add_handler(CommandHandler("start", start))
    app.add_handler(CommandHandler("add", add_expense))
    app.add_handler(CommandHandler("summary", summary))

    print("🤖 Bot is running...")
    app.run_polling()

if __name__ == '__main__':
    main()
