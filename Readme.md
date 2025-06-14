# Go Gin To-Do List API

A secure and feature-rich RESTful API for a To-Do list application, built with Go and the Gin framework. This project includes user authentication, email verification, password reset functionality, and full CRUD operations for tasks.

---

## 1. Setup & Installation

### Prerequisites
- Go (version 1.18+)
- PostgreSQL
- An email account (like Gmail) for sending verification emails.

### Configuration
1.  **Create `.env` file**: Create a file named `.env` in the project root and add your configuration details.

    ```env
    # PostgreSQL Database Configuration
    DB_HOST=localhost
    DB_USER=your_postgres_user
    DB_PASSWORD=your_postgres_password
    DB_NAME=todo_db
    DB_PORT=5432
    
    # Application & Email Configuration
    # The base URL of your server (NO trailing slash). This is crucial for email links.
    # For local development: APP_URL=http://localhost:8080
    # For production on bob.com: APP_URL=[http://bob.com](http://bob.com)
    APP_URL=http://localhost:8080
    
    # SMTP settings for sending emails. See the "Email Setup" section below.
    SMTP_HOST=smtp.gmail.com
    SMTP_PORT=587
    SMTP_USER=your_email@gmail.com
    SMTP_PASSWORD=your_gmail_app_password
    
    # JWT Configuration
    # LEAVE THIS BLANK! The application will generate a secure secret for you.
    JWT_SECRET=
    ```

2.  **Gmail App Password**: To use Gmail, you must enable 2-Step Verification and create a 16-character [App Password](https://myaccount.google.com/apppasswords). Use this App Password for `SMTP_PASSWORD`, not your regular password.

### Running the Server
1.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```
2.  **Run the Server**:
    ```bash
    go run main.go
    ```
    The API server will start on `http://localhost:8080`.

---

## 2. How to Use the API

All endpoints are available under the base URL: `http://localhost:8080/api/v1`

This guide uses `curl` for examples. You can use any API client like [Postman](https://www.postman.com/) or [Insomnia](https://insomnia.rest/).

### Step 1: Register a New User

Send a `POST` request with your desired username, email, and password.

```bash
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{
    "username": "bob",
    "email": "bob@example.com",
    "password": "a-very-secure-password"
}'
```
**Expected Response:**
```json
{
    "message": "User registered successfully. Please check your email to verify your account."
}
```

### Step 2: Verify Your Email

Check the server logs or your email inbox for the verification link. Click the link or use `curl` to send a `GET` request to it.

```bash
# The token will be long and unique
curl http://localhost:8080/api/v1/verify-email?token=...your_unique_token...
```
**Expected Response:**
```json
{
    "message": "Email verified successfully. You can now log in."
}
```

### Step 3: Log In to Get Your JWT

Now that your account is verified, you can log in to get a JSON Web Token (JWT). This token is required to access the protected task endpoints.

```bash
curl -X POST http://localhost:8080/api/v1/login \
-H "Content-Type: application/json" \
-d '{
    "email": "bob@example.com",
    "password": "a-very-secure-password"
}'
```
**Expected Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJvYkBleGFtcGxlLmNvbSIsImV4cCI6MTc0OTg2ODAwOSwidXNlcl9pZCI6MX0.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```
**Important:** Copy the `token` value. You will need it for the next steps.

### Step 4: Access Protected Routes (Manage Tasks)

To make using the token easier in the terminal, you can save it to an environment variable.

```bash
# Replace the token with the one you received from the login step
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImJvYkBleGFtcGxlLmNvbSIsImV4cCI6MTc0OTg2ODAwOSwidXNlcl9pZCI6MX0.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

Now you can use `$TOKEN` in your `curl` commands.

#### Create a New Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "title": "Learn Go",
    "description": "Finish the advanced concurrency chapter."
}'
```

#### Get All Your Tasks
```bash
curl -X GET http://localhost:8080/api/v1/tasks \
-H "Authorization: Bearer $TOKEN"
```

#### Get a Single Task by ID
(Assuming a task with ID `1` was created)
```bash
curl -X GET http://localhost:8080/api/v1/tasks/1 \
-H "Authorization: Bearer $TOKEN"
```

#### Update a Task
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
    "title": "Master Go Concurrency",
    "status": "in-progress"
}'
```

#### Delete a Task
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/1 \
-H "Authorization: Bearer $TOKEN"
```
**Expected Response:**
```json
{
    "message": "Task deleted successfully"
}
```