\# Smart Budget Setup Guide



\## Prerequisites

\- Go 1.23 or higher

\- Google Cloud Console account (for OAuth)



\## Installation



\### 1. Clone the Repository

```bash

git clone https://github.com/YourUsername/SmartBudget.git

cd SmartBudget

```



\### 2. Install Dependencies

```bash

go mod download

```



\### 3. Setup Google OAuth (Optional)

1\. Go to \[Google Cloud Console](https://console.cloud.google.com/)

2\. Create a new project

3\. Enable OAuth 2.0

4\. Create credentials (OAuth 2.0 Client ID)

5\. Add authorized redirect URI: `http://localhost:8080/auth/google/callback`

6\. Copy your Client ID and Client Secret



\### 4. Configure Environment Variables



\*\*Windows PowerShell:\*\*

```powershell

$env:GOOGLE\_CLIENT\_ID = "your-client-id-here"

$env:GOOGLE\_CLIENT\_SECRET = "your-client-secret-here"

$env:GOOGLE\_REDIRECT\_URL = "http://localhost:8080/auth/google/callback"

$env:SESSION\_KEY = "your-random-32-character-secret-key"

```



\*\*Linux/Mac:\*\*

```bash

export GOOGLE\_CLIENT\_ID="your-client-id-here"

export GOOGLE\_CLIENT\_SECRET="your-client-secret-here"

export GOOGLE\_REDIRECT\_URL="http://localhost:8080/auth/google/callback"

export SESSION\_KEY="your-random-32-character-secret-key"

```



\### 5. Build and Run



\*\*Build:\*\*

```bash

go build -o smartbudget ./cmd/expenseowl

```



\*\*Run:\*\*

```bash

./smartbudget

```



Or run directly:

```bash

go run ./cmd/expenseowl

```



\### 6. Access the Application

Open your browser and go to: `http://localhost:8080`



\## Features

\- ✅ Google OAuth Login

\- ✅ Email/Password Authentication

\- ✅ Expense Tracking

\- ✅ Budget Limits \& Warnings

\- ✅ Monthly Pie Chart Visualization

\- ✅ Cashflow Tracking

\- ✅ Recurring Transactions

\- ✅ CSV Import/Export

\- ✅ PWA Support



\## Login Options

1\. \*\*Google Sign-In\*\*: Click "Sign in with Google"

2\. \*\*Email/Password\*\*: Register a new account or login with existing credentials



\## Notes

\- Data is stored locally in JSON format by default

\- For production use, consider using PostgreSQL backend

\- The app runs on port 8080 by default

