# AcademAide - AI-Powered Academic Management System

## Project Overview

AcademAide is an intelligent academic management platform that leverages AI and Retrieval-Augmented Generation (RAG) to provide personalized academic assistance to both students and faculty members. The system combines traditional student information management with advanced AI capabilities to create an adaptive learning environment.

### Key Highlights

- **Dual-Role System**: Separate interfaces and features for students and teachers
- **AI-Powered Chatbot**: Context-aware conversations with multiple specialized agent personalities
- **RAG Integration**: Course material embeddings for intelligent content retrieval
- **Smart Quiz Generation**: AI-generated quizzes based on course materials and syllabus
- **Academic Analytics**: Real-time insights, risk assessment, and predictive analytics
- **Blockchain Wallet Integration**: Ethereum wallet generation and management for students
- **Study Group Matching**: Intelligent peer matching and collaboration features

## Core Features

### 1. **Multi-Agent AI Chat System**
Students and teachers can interact with specialized AI agents tailored to different needs:
- **General Agent**: Standard academic assistant for course-related queries
- **Socratic Tutor**: Guides learning through probing questions
- **Code Reviewer**: Analyzes code for bugs, complexity, and best practices
- **Research Assistant**: Provides academic-level, comprehensive answers
- **Exam Coach**: Focuses on high-yield topics and exam strategies
- **Motivational Coach**: Provides encouragement and study planning assistance
- **Teacher Assistant**: Helps faculty with course planning and student analysis

The chat system uses RAG (Retrieval-Augmented Generation) to:
- Retrieve relevant course materials from vectorized PDF embeddings
- Provide context-aware responses based on student's enrollment, schedule, and academic history
- Filter content based on enrolled courses and specific units
- Maintain conversation history for contextual continuity

### 2. **AI Quiz Generator**
- Generate customized quizzes based on course syllabus and materials
- Filter by specific units (1-5) or generate comprehensive quizzes
- Configurable question count (5, 10, 15, or 20 questions)
- Automatic quiz analysis identifying weak areas and study priorities
- Reference-based questions linked to source materials

### 3. **Student Dashboard Features**
- **Profile Management**: View personal information, CGPA, and enrollment status
- **Dynamic Timetable**: Real-time schedule with current/upcoming class notifications
- **Course Resources**: Access course materials organized by units
- **Announcements**: View course-specific announcements from faculty
- **AI Insights**: Risk assessment based on attendance and grade trends
- **What-If Scenarios**: Simulate impact of missed classes on attendance

### 4. **Teacher Dashboard**
- **Class Health Analytics**: Grade distribution and performance metrics
- **Student Management**: View enrolled students with detailed profiles
- **At-Risk Student Detection**: Early warning system for struggling students
- **Announcement System**: Broadcast course announcements to students
- **Student Performance Tracking**: Individual student progress monitoring
- **Grade Analysis**: Performance heatmaps and distribution charts


### 5. **Authentication & Security**
- **Dual-Role Login**: Separate authentication for students and teachers
- **Google OAuth Integration**: Sign in/up with Google accounts
- **JWT-based Authentication**: Secure API access with token validation
- **Onboarding Flow**: New user registration with profile completion


### 6. **RAG-Powered Material Ingestion**
- PDF course material processing and chunking
- Semantic embedding generation using Ollama (nomic-embed-text model)
- Vector storage in PostgreSQL with cosine similarity search
- Unit-based material organization (1-5 per course)
- Course-filtered content retrieval

### 7. **Intelligent Context Management**
The system maintains rich contextual awareness including:
- **For Students**: 
  - Real-time schedule (current/next class detection)
  - Academic standing (CGPA, grades, year)
  - Course enrollment status
  - Historical performance trends
- **For Teachers**:
  - Courses taught and enrolled student lists
  - Class performance statistics
  - Alert generation for at-risk students

### 8. **Caching & Performance**
- Redis caching for frequently accessed data (profiles, timetables)
- Response caching for common chat queries
- TTL-based cache invalidation (5 minutes for profiles, 1 hour for timetables)


## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Databases**: 
  - PostgreSQL (relational data, vector embeddings)
  - MongoDB (chat logs, quizzes, context)
  - Redis (caching)
- **AI/ML**: Ollama (llama3.2 for chat, nomic-embed-text for embeddings)
- **Blockchain**: Ethereum wallet generation (go-ethereum)
- **Authentication**: JWT, Google OAuth 2.0

### Frontend
- **Framework**: Next.js 14 (React)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: Custom component library
- **State Management**: React Hooks
- **HTTP Client**: Fetch API

### Infrastructure
- Docker support for PostgreSQL
- Environment-based configuration
- CORS-enabled API

## Prerequisites

### Required Software
- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher (for frontend)
- **PostgreSQL**: Version 14+ (Port 5432 or 5435 for Docker)
- **MongoDB**: Version 5+ (Port 27017)
- **Redis**: Version 6+ (Port 6379)
- **Ollama**: Local AI runtime (Port 11434)
  - Models required: `llama3.2` (chat) and `nomic-embed-text` (embeddings)
- **Docker** (Optional): For containerized PostgreSQL

### Environment Setup

Before starting, ensure all required services are running:

```powershell
# Check PostgreSQL
psql -U postgres -c "SELECT version();"

# Check MongoDB
mongosh --eval "db.version()"

# Check Redis
redis-cli ping

# Check Ollama and pull required models
ollama list
ollama pull llama3.2
ollama pull nomic-embed-text
```

## Installation & Setup

### Step 1: Clone and Navigate

```powershell
cd C:\Users\Sathvik\OneDrive\Desktop\DBMS_Lab
```

### Step 2: Backend Setup

#### 2.1 Install Go Dependencies

```powershell
go mod download
go mod tidy
```

#### 2.2 Database Initialization

**Option A: Standard PostgreSQL (Port 5432)**

```powershell
# Create database
psql -U postgres -c "CREATE DATABASE academ_aide;"

# Run schema setup
psql -U postgres -d academ_aide -f database/schema.sql

# Insert sample data (optional)
psql -U postgres -d academ_aide -f database/insert_real_data.sql
```

**Option B: Docker PostgreSQL (Port 5435)**

```powershell
# Start PostgreSQL in Docker
docker run --name academ-postgres -e POSTGRES_PASSWORD=postgres -p 5435:5432 -d postgres:14

# Initialize database
python database/init_docker_db.py

# Run schema
docker exec -i academ-postgres psql -U postgres -d academ_aide < database/schema.sql
```

MongoDB and Redis require no schema setup.

#### 2.3 Environment Configuration

Create a `.env` file in the project root:

```env
# Database Configuration
POSTGRES_DSN=host=localhost user=postgres password=postgres dbname=academ_aide port=5432 sslmode=disable
MONGO_URI=mongodb://localhost:27017
REDIS_ADDR=localhost:6379

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Google OAuth (Optional - for OAuth login)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# Ollama Configuration (default values)
OLLAMA_URL=http://localhost:11434
```

**Note**: If using Docker PostgreSQL, change port to `5435` in POSTGRES_DSN.

#### 2.4 Material Ingestion (Optional but Recommended)

To enable RAG features, ingest course materials:

```powershell
# Ensure materials are in materials/<CourseID>/<UnitNumber>/*.pdf
# Example: materials/CD252IA/1/introduction.pdf

# Update database config in ingest_materials.py
# Set DB_PORT = "5432" or "5435" based on your setup

# Run ingestion script
python ingest_materials.py
```

This process:
- Extracts text from PDFs
- Generates embeddings using Ollama
- Stores chunks in PostgreSQL with vector search capability

#### 2.5 Start Backend Server

```powershell
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

**Verify Backend**:
```powershell
curl http://localhost:8080/student/profile
# Should return 401 Unauthorized (authentication required)
```

### Step 3: Frontend Setup

#### 3.1 Navigate to Frontend Directory

```powershell
cd frontend/AcademAide
```

#### 3.2 Install Dependencies

```powershell
npm install
# or
yarn install
# or
pnpm install
```

#### 3.3 Configure Frontend (if needed)

API endpoints are pre-configured to `http://localhost:8080`. No changes needed for local development.

#### 3.4 Start Development Server

```powershell
npm run dev
# or
yarn dev
# or
pnpm dev
```

Frontend will start on `http://localhost:3000`

### Step 4: Access the Application

1. **Open Browser**: Navigate to `http://localhost:3000`
2. **Login**:
   - **Student**: Use student ID (e.g., `S1001`) with any password
   - **Teacher**: Use faculty ID (e.g., `F1001`) with any password
   - **Google OAuth**: Click "Sign in with Google" (requires OAuth setup)

### Default Test Accounts

**Students**:
- ID: `S1001` - `S1100` (Password: any)

**Teachers**:
- ID: `F1001` - `F1020` (Password: any)

**Courses**:
- `CD252IA` - Database Management Systems
- `CS354TA` - Theory of Computation
- `IS353IA` - Artificial Intelligence & ML
- `XX355TBX` - Cloud Computing
- `HS251TA` - Economics & Management

## API Endpoints

- **POST** `/login`
  - Body: `{"student_id": "S1001"}`
- **GET** `/student/profile?student_id=S1001`
- **GET** `/student/timetable?student_id=S1001`
- **POST** `/chat/message`
  - Body: `{"student_id": "S1001", "message": "When is my next class?"}`


## Extended API Documentation

### Authentication Endpoints
- **POST** `/login` - Dual-role login (students/teachers)
- **GET** `/auth/google/login` - Google OAuth initiation
- **GET** `/auth/google/callback` - OAuth callback
- **POST** `/auth/complete-registration` - Complete user onboarding

### Chat Endpoints (All require authentication)
- **POST** `/chat/message` - AI conversation with context
- **DELETE** `/chat/history` - Clear user chat history

**Available Agents**: general, socratic, code_reviewer, research, exam, motivational, teacher

### Quiz & Learning
- **POST** `/quiz/generate` - Generate customized quiz (unit-specific or comprehensive)
- **POST** `/ai/quiz-analysis` - Get personalized study recommendations

### Study Groups
- **GET** `/groups/peers` - Find classmates in same course
- **POST** `/groups/create` - Create study group
- **GET** `/groups/list` - List available groups

### Teacher Dashboard
- **GET** `/teacher/courses` - Get taught courses
- **GET** `/teacher/students` - View enrolled students
- **GET** `/teacher/class-health` - Performance analytics
- **GET** `/teacher/at-risk` - Identify struggling students
- **POST** `/teacher/announce` - Broadcast announcements

### AI Insights
- **GET** `/ai/insights` - Academic risk assessment
- **POST** `/ai/what-if` - Attendance scenario simulation

## Troubleshooting Guide

**Database Connection Issues**:
- Verify PostgreSQL is running on correct port (5432 or 5435)
- Check credentials in .env file
- Test connection: `psql -U postgres -d academ_aide`

**Ollama/AI Issues**:
- Ensure models are downloaded: `ollama list`
- Pull required models: `ollama pull llama3.2` and `ollama pull nomic-embed-text`
- Verify Ollama is running: `curl http://localhost:11434/api/tags`

**Frontend Build Errors**:
- Clear Next.js cache: `rm -rf .next`
- Reinstall dependencies: `npm install`
- Check Node version: `node --version` (should be 18+)

**RAG Not Working**:
- Verify materials ingested: `SELECT COUNT(*) FROM COURSE_MATERIAL_CHUNK;`
- Re-run ingestion script if needed
- Check vector embeddings exist

---

**Project Version**: 1.0.0
**Last Updated**: January 30, 2026
**Status**: Active Development
