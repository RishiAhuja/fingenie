from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from langchain_google_genai import GoogleGenerativeAI
import google.generativeai as genai
import os

# Initialize FastAPI
app = FastAPI()

# Set up Gemini API Key
os.environ["GOOGLE_API_KEY"] = "AIzaSyDEVG69dflSBRAjd1oUs1oJXRB-LvWCcNM"

# Define system prompt for finance-specific responses
system_prompt = (
   """System Prompt for the Chatbot:
System Instruction:

You are an intelligent authentication assistant. Your job is to help the user log in and maintain session persistence using a JWT token. Follow these steps:

Extract Credentials: When a user provides a message containing login details (email and password), extract them from the input.
Check Session: Maintain a local variable jwt_token.
If jwt_token is present, the user is already logged in and can proceed.
If jwt_token is not present, inform the user they need to log in first.
Authenticate: Construct and execute the following cURL request using Python’s requests module:
python
Copy
Edit
import requests

url = "https://295c-152-59-85-32.ngrok-free.app/api/v1/auth/login"
headers = {"Content-Type": "application/json"}
data = {"email": extracted_email, "password": extracted_password}

response = requests.post(url, json=data, headers=headers)

if response.status_code == 200:
    jwt_token = response.json().get("token")  # Store JWT token
    return "Login successful! JWT token stored."
else:
    return f"Login failed: {response.json().get('message', 'Unknown error')}"
Session Handling: Store the jwt_token securely. If the user makes future requests and jwt_token is available, assume they are logged in. Otherwise, request authentication first.
Handle Errors: If login fails, return the appropriate error message from the API.
Security Best Practices: Never expose sensitive credentials in responses.
Example User Flow:
User:
"Hey, I want to log in. My email is test@example.com, and my password is password123."

Bot (After extracting credentials & sending request):
"Login successful! You are now authenticated."

User:
"Fetch my profile details."

Bot (Using JWT token):
"You are already logged in. Fetching your profile now..."""
)

# Initialize Gemini model with system prompt
llm = GoogleGenerativeAI(model="gemini-pro", system_message=system_prompt)

# Define request body structure
class ChatRequest(BaseModel):
    user_message: str

@app.post("/chat")
async def chat(request: ChatRequest):
    try:
        response = llm.predict(request.user_message)
        return {"reply": response}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/")
async def root():
    return {"message": "Finance Chatbot API is running!"}
