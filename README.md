# AcademAide Setup Guide

## Prerequisites
- Go 1.21+
- PostgreSQL (Port 5432)
- MongoDB (Port 27017)
- Redis (Port 6379)
- Ollama (Running `llama3.2` on Port 11434)
- Node.js (for Frontend)

## Installation & Running

### Backend

1.  **Initialize Module** (if not done):
    ```powershell
    go mod init academ_aide
    go mod tidy
    ```

2.  **Setup Database**:
    - Run the `database/schema.sql` script in PostgreSQL to create tables and seed data.
    - MongoDB and Redis do not require schema setup.

3.  **Environment Variables**:
    - Create a `.env` file in the root (optional, defaults provided):
      ```
      POSTGRES_DSN=host=localhost user=postgres password=postgres dbname=academ_aide port=5432 sslmode=disable
      MONGO_URI=mongodb://localhost:27017
      REDIS_ADDR=localhost:6379
      ```

4.  **Run Server**:
    ```powershell
    go run cmd/server/main.go
    ```
    Server will start on `http://localhost:8080`.

### Frontend

The frontend is located in `frontend/AcademAide`.

1.  **Navigate to Frontend**:
    ```powershell
    cd frontend/AcademAide
    ```

2.  **Install Dependencies**:
    ```bash
    npm install
    # or
    yarn install
    ```

3.  **Run Development Server**:
    ```bash
    npm run dev
    ```
    Open [http://localhost:3000](http://localhost:3000) with your browser.

## API Endpoints

- **POST** `/login`
  - Body: `{"student_id": "S1001"}`
- **GET** `/student/profile?student_id=S1001`
- **GET** `/student/timetable?student_id=S1001`
- **POST** `/chat/message`
  - Body: `{"student_id": "S1001", "message": "When is my next class?"}`
