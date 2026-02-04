# AcademAide - Comprehensive Project Documentation
## DBMS Project Presentation Guide

---

## Table of Contents
1. [Project Overview](#1-project-overview)
2. [System Architecture](#2-system-architecture)
3. [Database Design & Schema](#3-database-design--schema)
4. [Chatbot Personalization Implementation](#4-chatbot-personalization-implementation)
5. [Quiz System & RAG Implementation](#5-quiz-system--rag-implementation)
6. [Key SQL Queries Explained](#6-key-sql-queries-explained)
7. [Data Flow & Integration](#7-data-flow--integration)
8. [AI/ML Components](#8-aiml-components)
9. [Caching & Performance Optimization](#9-caching--performance-optimization)
10. [Potential Modifications with Code Diffs](#10-potential-modifications-with-code-diffs)

---

## 1. Project Overview

### What is AcademAide?
AcademAide is an **AI-powered academic assistant platform** that provides personalized educational support to students and teachers through:
- **Context-aware chatbot** with multiple AI agent personalities
- **Intelligent quiz generation** using RAG (Retrieval-Augmented Generation)
- **Academic performance insights** with predictive analytics
- **Real-time schedule management** and course tracking

### Technology Stack
- **Frontend**: Next.js 14, TypeScript, Tailwind CSS, shadcn/ui
- **Backend**: Go (Gin framework)
- **Databases**: 
  - PostgreSQL (with pgvector extension) - Primary relational data
  - MongoDB - Chat logs and quiz storage
  - Redis - Caching layer
- **AI/ML**: 
  - Ollama (llama3.2 for generation)
  - nomic-embed-text (768-dimensional embeddings)
  - pgvector for similarity search

---

## 2. System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Next.js)                      │
│  - React Components                                          │
│  - JWT Auth (Cookies)                                        │
│  - REST API Calls                                            │
└────────────────┬────────────────────────────────────────────┘
                 │ HTTP/JSON
                 ▼
┌─────────────────────────────────────────────────────────────┐
│                   Backend (Go/Gin)                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Handlers   │  │   Services   │  │ Repository   │      │
│  │  (API Layer) │  │  (Business)  │  │ (Data Access)│      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────┬─────────┬─────────┬─────────────────────────┬────────┘
      │         │         │                         │
      ▼         ▼         ▼                         ▼
┌──────────┐ ┌────────┐ ┌──────────┐        ┌──────────────┐
│PostgreSQL│ │ MongoDB│ │  Redis   │        │Ollama (Local)│
│ (pgvector│ │(Logs & │ │ (Cache)  │        │- llama3.2    │
│ enabled) │ │ Quizzes│ │          │        │- nomic-embed │
└──────────┘ └────────┘ └──────────┘        └──────────────┘
```

### Request Flow Example (Chat Message):
1. User types message in frontend `/chat` page
2. Frontend sends POST to `/chat/message` with JWT token
3. Middleware validates JWT and extracts `user_id` and `role`
4. `ChatHandler` calls `RAGService.ProcessChat()`
5. RAG Service:
   - Checks Redis cache for similar query
   - Queries PostgreSQL for student profile, grades, schedule
   - Generates embedding of user query
   - Searches `COURSE_MATERIAL_CHUNK` table using vector similarity
   - Fetches last 5 messages from MongoDB
   - Constructs personalized prompt with context
   - Calls Ollama API for LLM response
   - Stores messages in MongoDB
   - Caches response in Redis (5 min TTL)
6. Returns AI response to frontend

---

## 3. Database Design & Schema

### 3.1 PostgreSQL Schema (3NF Normalized)

#### Core Tables

**DEPARTMENT**
```sql
CREATE TABLE DEPARTMENT (
    dept_id VARCHAR(10) PRIMARY KEY,
    dept_name VARCHAR(100) NOT NULL,
    hod_id VARCHAR(20),
    FOREIGN KEY (hod_id) REFERENCES FACULTY(faculty_id)
);
```

**FACULTY**
```sql
CREATE TABLE FACULTY (
    faculty_id VARCHAR(20) PRIMARY KEY,
    f_first_name VARCHAR(50) NOT NULL,
    f_last_name VARCHAR(50) NOT NULL,
    f_email VARCHAR(100) UNIQUE NOT NULL,
    f_phone_no VARCHAR(15)
);
```

**STUDENT**
```sql
CREATE TABLE STUDENT (
    student_id VARCHAR(20) PRIMARY KEY,
    s_first_name VARCHAR(50) NOT NULL,
    s_last_name VARCHAR(50) NOT NULL,
    s_email VARCHAR(100) UNIQUE NOT NULL,
    s_phone_no VARCHAR(15),
    semester INTEGER NOT NULL,
    year_of_joining INTEGER NOT NULL,
    dept_id VARCHAR(10) NOT NULL,
    FOREIGN KEY (dept_id) REFERENCES DEPARTMENT(dept_id)
);
```

**COURSE** (with vector embedding)
```sql
CREATE TABLE COURSE (
    course_id VARCHAR(10) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    credits INTEGER NOT NULL CHECK (credits > 0),
    dept_id VARCHAR(10) NOT NULL,
    description TEXT,
    embedding vector(768),  -- pgvector type for semantic search
    FOREIGN KEY (dept_id) REFERENCES DEPARTMENT(dept_id)
);
```

#### Relationship Tables

**ENROLLS_IN** (Student-Course enrollment)
```sql
CREATE TABLE ENROLLS_IN (
    student_id VARCHAR(20),
    course_id VARCHAR(10),
    grade VARCHAR(2),
    status VARCHAR(20) DEFAULT 'Enrolled',
    backlog BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES STUDENT(student_id),
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);
```

**TEACHES** (Faculty-Course-Section assignment)
```sql
CREATE TABLE TEACHES (
    faculty_id VARCHAR(20),
    course_id VARCHAR(10),
    section_name VARCHAR(10),
    PRIMARY KEY (faculty_id, course_id, section_name),
    FOREIGN KEY (faculty_id) REFERENCES FACULTY(faculty_id),
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id),
    FOREIGN KEY (section_name) REFERENCES SECTION(section_name)
);
```

**SCHEDULE** (Timetable management)
```sql
CREATE TABLE SCHEDULE (
    schedule_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    section_name VARCHAR(10) NOT NULL,
    day_of_week VARCHAR(10) NOT NULL CHECK (day_of_week IN 
        ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday')),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_number VARCHAR(20),
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id),
    FOREIGN KEY (section_name) REFERENCES SECTION(section_name)
);
```

#### RAG/AI Table (Critical for Chatbot)

**COURSE_MATERIAL_CHUNK** (Vector embeddings of course materials)
```sql
CREATE TABLE COURSE_MATERIAL_CHUNK (
    chunk_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    unit_no INTEGER NOT NULL,           -- 1-5 for units
    content_text TEXT NOT NULL,         -- Actual paragraph/note content
    embedding vector(768),              -- 768-dim vector from nomic-embed-text
    source_file VARCHAR(255),           -- e.g., 'Unit1_Introduction.pdf'
    chunk_index INTEGER,                -- Order within document
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);
```

### 3.2 MongoDB Collections

**ChatLogs** (Conversation history)
```javascript
{
  _id: ObjectId,
  student_id: "S1001",
  message: "Explain arrays in C++",
  intent: "chat",
  sentiment: "neutral",
  timestamp: ISODate,
  is_bot: false  // false = user message, true = bot response
}
```

**quizzes** (Generated quiz data)
```javascript
{
  _id: ObjectId,
  course_id: "CS101",
  topic: "Generated from Course Materials",
  questions: [
    {
      id: 1,
      text: "What is the time complexity of binary search?",
      options: ["O(n)", "O(log n)", "O(n²)", "O(1)"],
      correct_option: 1,
      reference: "Unit2_Algorithms.pdf"
    }
  ],
  created_at: ISODate
}
```

### 3.3 Entity-Relationship Diagram

```
DEPARTMENT (1) ──< (M) STUDENT
    │                    │
    │                    │
    └──< (M) COURSE     (M)
              │          │
              │          └──> (M,M) ENROLLS_IN
              │                     [grade, status, backlog]
              │
              └──> (M) COURSE_MATERIAL_CHUNK
                        [embedding vector(768)]
                        
FACULTY (1) ──< (M,M,M) TEACHES
                        │
                        └──> COURSE + SECTION

COURSE (1) ──< (M) SCHEDULE
                   [day_of_week, start_time, end_time]
```

### 3.4 Normalization (3NF Compliance)

**1NF (First Normal Form):**
- All attributes are atomic (no repeating groups)
- Each cell contains single value
- Example: `s_first_name` and `s_last_name` are separate, not combined

**2NF (Second Normal Form):**
- No partial dependencies
- All non-key attributes depend on the *entire* primary key
- Example: In `ENROLLS_IN(student_id, course_id, grade)`, `grade` depends on both student AND course, not just one

**3NF (Third Normal Form):**
- No transitive dependencies
- Non-key attributes depend ONLY on primary key
- Example: `dept_name` is in `DEPARTMENT`, not duplicated in `STUDENT` or `COURSE`
- Student has `dept_id` (FK), not `dept_name` directly

---

## 4. Chatbot Personalization Implementation

### 4.1 Multi-Agent System

The chatbot supports **7 different AI personalities**:

1. **General**: Basic academic advisor with grade-aware recommendations
2. **Socratic**: Asks guiding questions instead of giving direct answers
3. **Code Reviewer**: Analyzes code for bugs, complexity, and best practices
4. **Research**: PhD-level detailed explanations with citations
5. **Exam Coach**: Focus on high-yield topics, mnemonics, time management
6. **Motivational**: Supportive mentor for stress management and study planning
7. **Teacher**: Assists faculty with course planning and student analysis

**Agent selection logic (in frontend):**
```typescript
// components/AgentSelector.tsx
const agents = [
  { id: "general", name: "General Advisor", icon: "👨‍🎓" },
  { id: "socratic", name: "Socratic Tutor", icon: "🤔" },
  { id: "code_reviewer", name: "Code Reviewer", icon: "💻" },
  // ... etc
]
```

### 4.2 Context Building for Students

**SQL Query 1: Student Profile & Identity**
```sql
SELECT s_first_name, dept_id, year_of_joining 
FROM STUDENT 
WHERE student_id = $1
```
- Calculates current year: `currentYear - year_of_joining`
- Used in prompt: *"You are talking to John, a 3rd Year CS student"*

**SQL Query 2: Academic History & CGPA Calculation**
```sql
SELECT c.title, e.grade, c.credits, e.course_id
FROM ENROLLS_IN e
JOIN COURSE c ON e.course_id = c.course_id
WHERE e.student_id = $1 AND e.grade IS NOT NULL
```

**Grade Point Mapping:**
```go
switch grade {
case "O", "A+": points = 10.0
case "A":       points = 9.0
case "B+":      points = 8.0
case "B":       points = 7.0
case "C+":      points = 6.0
case "C":       points = 5.0
case "D":       points = 4.0
default:        points = 0.0
}
totalPoints += points * float64(credits)
totalCredits += credits
cgpa = totalPoints / totalCredits
```

**Context in Prompt:**
```
Academic Standing: CGPA 8.45 (Academic History: [Data Structures: A, Intro to Programming: B+, ...])
```

**SQL Query 3: Enrolled Courses (for RAG filtering)**
```sql
SELECT c.course_id, c.title 
FROM ENROLLS_IN e 
JOIN COURSE c ON e.course_id = c.course_id 
WHERE e.student_id = $1 AND e.status = 'Enrolled'
```
- Returns: `["CS101", "CS102", "CS103"]`
- Used to filter vector search to only relevant courses

**SQL Query 4: Next Class Detection**
```sql
-- Check for ongoing class
SELECT c.title, sch.room_number, sch.end_time
FROM SCHEDULE sch 
JOIN ENROLLS_IN e ON sch.course_id = e.course_id
JOIN COURSE c ON sch.course_id = c.course_id
WHERE e.student_id = $1 
  AND sch.day_of_week = $2       -- e.g., 'Monday'
  AND sch.start_time <= $3       -- Current time: '10:30:00'
  AND sch.end_time > $3
LIMIT 1
```

If no ongoing class:
```sql
-- Find next upcoming class today
SELECT c.title, sch.start_time, sch.room_number 
FROM SCHEDULE sch 
JOIN ENROLLS_IN e ON sch.course_id = e.course_id
JOIN COURSE c ON sch.course_id = c.course_id
WHERE e.student_id = $1 
  AND sch.day_of_week = $2
  AND sch.start_time > $3        -- Future time only
ORDER BY sch.start_time ASC
LIMIT 1
```

**Context in Prompt:**
```
Current Schedule Status: HAPPENING NOW: You should be in Data Structures (Room 102) until 11:30
OR
Current Schedule Status: Next Class: Database Systems at 14:00 in Room 205 (Today)
```

### 4.3 Context Building for Teachers

**SQL Query: Taught Courses**
```sql
SELECT c.title 
FROM TEACHES t 
JOIN COURSE c ON t.course_id = c.course_id 
WHERE t.faculty_id = $1
```

**SQL Query: Enrolled Students in Taught Courses**
```sql
SELECT DISTINCT c.title, s.s_first_name, s.s_last_name, s.student_id
FROM TEACHES t
JOIN ENROLLS_IN e ON t.course_id = e.course_id
JOIN STUDENT s ON e.student_id = s.student_id
JOIN COURSE c ON t.course_id = c.course_id
WHERE t.faculty_id = $1 AND e.status = 'Enrolled'
ORDER BY c.title, s.s_last_name
```

**Teacher Context in Prompt:**
```
You are talking to Dr. Alice Smith, a Faculty Member.
- Email: alice.smith@univ.edu
- Courses Taught: [Intro to Programming, Data Structures]

[ENROLLED STUDENTS]:
Course: Intro to Programming
- John Doe (S1001)
- Jane Roe (S1002)
```

### 4.4 RAG (Retrieval-Augmented Generation)

**Step 1: Generate Query Embedding**
```go
embedding, err := s.Embedder.GenerateEmbedding(message)
// Calls Ollama API: POST http://localhost:11434/api/embeddings
// Model: nomic-embed-text
// Returns: []float32 of 768 dimensions
```

**Step 2: Vector Similarity Search**
```sql
SELECT content_text, course_id, unit_no, 
       1 - (embedding <=> $1::vector) as score, 
       source_file
FROM COURSE_MATERIAL_CHUNK
WHERE course_id = ANY($3::text[])  -- Filter by enrolled courses
  AND unit_no = $4                 -- Optional unit filter
ORDER BY embedding <=> $1::vector ASC  -- Cosine distance
LIMIT $2
```

**Operator Explanation:**
- `<=>` : Cosine distance operator (pgvector)
- `1 - distance` = similarity score (0-1 range)
- Lower distance = higher similarity
- `ANY($3::text[])` : Match any course_id in array

**Example Retrieved Context:**
```
Relevant Materials:
- [Course CS101 Unit 2] Arrays are contiguous memory blocks. Access time is O(1) for index-based retrieval. They are fixed-size in C++, unlike vectors which are dynamic.
- [Course CS101 Unit 2] Pointer arithmetic allows traversing arrays efficiently. ptr++ moves to next element based on type size.
```

### 4.5 Final Prompt Construction

```go
prompt := fmt.Sprintf(`
[SYSTEM PREAMBLE]
You are a Socratic Tutor. Never provide direct answers immediately.
Ask guiding questions to check understanding...

[CONTEXTUAL AWARENESS]
You are talking to John, a 3 Year CS student.
- Current Time: Monday, 10:45 AM
- Academic Standing: CGPA 8.45 (Academic History: [Data Structures: A, ...])
- Current Schedule Status: Next Class: Database Systems at 14:00 in Room 205
- Enrolled Courses: [Intro to Programming, Data Structures, Database Systems]

[RELEVANT STUDY MATERIALS (RAG)]
- [Course CS101 Unit 2] Arrays are contiguous memory blocks...

[USER SENTIMENT]
User seems neutral.

[CONVERSATION HISTORY]
User: What is a pointer?
Bot: A pointer stores memory address of another variable...
User: How do I use arrays with pointers?

User: %s
Assistant:`, message)
```

**This prompt goes to Ollama (llama3.2) for generation.**

### 4.6 Caching & Storage

**Redis Cache:**
```go
hash := sha256.Sum256([]byte(message))
cacheKey := "response:" + hex.EncodeToString(hash[:])
cached, err := config.RedisClient.Get(ctx, cacheKey).Result()
if err == nil {
    return cached, nil  // Cache hit
}
// ... generate response ...
config.RedisClient.Set(ctx, cacheKey, response, 5*time.Minute)
```

**MongoDB Storage:**
```go
// User message
userLog := models.ChatLog{
    StudentID: userID,
    Message:   message,
    Sentiment: "neutral",
    IsBot:     false,
}
coll.InsertOne(ctx, userLog)

// Bot response
botLog := models.ChatLog{
    StudentID: userID,
    Message:   response,
    IsBot:     true,
}
coll.InsertOne(ctx, botLog)
```

---

## 5. Quiz System & RAG Implementation

### 5.1 Quiz Generation Flow

**Frontend Request:**
```typescript
const res = await fetch("http://localhost:8080/quiz/generate", {
    method: "POST",
    body: JSON.stringify({
        course_id: "CS101",
        unit: 2,           // 0 = All units, 1-5 = specific unit
        num_questions: 10
    })
})
```

**Backend Processing:**

**Step 1: Fetch Syllabus Topics**
```sql
-- If specific unit:
SELECT topic FROM SYLLABUS_UNIT 
WHERE course_id = $1 AND unit_no = $2

-- If all units:
SELECT topic FROM SYLLABUS_UNIT 
WHERE course_id = $1
```

**Step 2: RAG Material Retrieval**
```go
query := fmt.Sprintf("Important concepts in %s Unit %d", courseID, unit)
embedding, _ := s.Embedder.GenerateEmbedding(query)

materials, _ := s.Repo.SearchMaterials(ctx, embedding, 5, []string{courseID}, unit)
```

**SQL Query:**
```sql
SELECT content_text, course_id, unit_no, source_file
FROM COURSE_MATERIAL_CHUNK
WHERE course_id = $1 AND unit_no = $2
ORDER BY embedding <=> $3::vector ASC
LIMIT 5
```

**Step 3: Generate Questions with Ollama**
```go
prompt := fmt.Sprintf(`
You are a professor. Generate a quiz with %d multiple-choice questions 
for course %s, focusing STRICTLY on Unit %d.

CRITICAL: Use ONLY the content provided below. Do NOT use outside knowledge.
For each question, "reference" MUST be exact Source filename.

Context Materials:
---
Source: Unit2_Pointers.pdf
Unit: 2
Content: Pointers store memory addresses. Dereferencing with * operator...
---
Source: Unit2_Arrays.pdf
Unit: 2  
Content: Arrays are fixed-size contiguous blocks...

Topics: Pointers, Arrays, Pointer Arithmetic

Return ONLY valid JSON:
{
  "questions": [
    {
      "id": 1,
      "text": "What does the * operator do to a pointer?",
      "options": ["Increments it", "Dereferences it", "Deletes it", "None"],
      "correct_option": 1,
      "reference": "Unit2_Pointers.pdf"
    }
  ]
}
`, numQuestions, courseID, unit)

jsonResp, _ := s.callOllamaJSON(prompt)
```

**Ollama API Call:**
```go
reqBody := map[string]interface{}{
    "model":  "llama3.2",
    "prompt": prompt,
    "stream": false,
    "format": "json",  // Force JSON output mode
}
resp, _ := http.Post("http://localhost:11434/api/generate", 
                     "application/json", 
                     bytes.NewBuffer(jsonData))
```

**Step 4: Save Quiz to MongoDB**
```go
quiz := &models.Quiz{
    CourseID:  courseID,
    Topic:     "Generated from Course Materials",
    Questions: parsedQuestions,
    CreatedAt: time.Now(),
}
coll := config.MongoDB.Collection("quizzes")
coll.InsertOne(context.Background(), quiz)
```

### 5.2 Quiz Submission & Analysis

**Frontend Submission:**
```typescript
const wrongQuestions = quiz.questions.filter((q, idx) => 
    answers[q.id] !== q.correct_option
).map(q => ({
    question_text: q.text,
    correct_answer: q.options[q.correct_option],
    user_answer: q.options[answers[q.id]] || "Skipped",
    reference: q.reference
}))

// Send for AI analysis
await fetch("http://localhost:8080/ai/quiz-analysis", {
    method: "POST",
    body: JSON.stringify({
        course_id: "CS101",
        wrong_questions: wrongQuestions,
        total_questions: 10,
        score: 7
    })
})
```

**Backend Analysis (ai_service.go):**
```go
prompt := fmt.Sprintf(`
You are an academic tutor. A student took a quiz for %s and got these wrong:

1. Question: What does * do to a pointer?
   Correct: Dereferences it
   User Answered: Increments it
   Source: Unit2_Pointers.pdf

Analyze mistakes to identify weak topics. Return JSON:
{
  "weak_areas": ["Pointer Operations", "Memory Management"],
  "study_priorities": [
    {
      "topic": "Pointer Dereferencing",
      "priority": "High",
      "reason": "Confused dereferencing with arithmetic - core concept"
    }
  ]
}
`, courseID)

jsonResp, _ := callOllamaJSON(prompt)
```

**Response to Frontend:**
```json
{
  "weak_areas": ["Pointer Operations", "Array Indexing"],
  "study_priorities": [
    {
      "topic": "Pointer Dereferencing vs Arithmetic",
      "priority": "High",
      "reason": "Fundamental misunderstanding of * operator usage"
    },
    {
      "topic": "Array Bounds Checking",
      "priority": "Medium",
      "reason": "Off-by-one errors indicate need for practice"
    }
  ]
}
```

---

## 6. Key SQL Queries Explained

### 6.1 Student Dashboard Queries

**Query 1: Student Profile with Course Count**
```sql
-- Main profile
SELECT student_id, s_first_name, s_last_name, s_email, 
       s_phone_no, semester, year_of_joining, dept_id 
FROM STUDENT 
WHERE student_id = 'S1001';

-- Course count
SELECT COUNT(*) 
FROM ENROLLS_IN 
WHERE student_id = 'S1001' AND status = 'Enrolled';
```

**Query 2: CGPA Calculation**
```sql
SELECT e.grade, c.credits
FROM ENROLLS_IN e
JOIN COURSE c ON e.course_id = c.course_id
WHERE e.student_id = 'S1001' AND e.grade IS NOT NULL;

-- Application logic calculates:
-- CGPA = SUM(grade_points * credits) / SUM(credits)
```

**Query 3: Timetable Retrieval**
```sql
SELECT c.course_id, c.title, sch.section_name, sch.day_of_week, 
       TO_CHAR(sch.start_time, 'HH24:MI') as start_time,
       TO_CHAR(sch.end_time, 'HH24:MI') as end_time,
       sch.room_number
FROM ENROLLS_IN e
JOIN COURSE c ON e.course_id = c.course_id
JOIN SCHEDULE sch ON c.course_id = sch.course_id 
WHERE e.student_id = 'S1001'
ORDER BY sch.day_of_week, sch.start_time;
```

**Query 4: Resources for Enrolled Courses**
```sql
SELECT r.resource_id, r.title, r.description, r.type, r.course_id
FROM RESOURCE r
JOIN ENROLLS_IN e ON r.course_id = e.course_id
WHERE e.student_id = 'S1001';
```

### 6.2 RAG Vector Search Queries

**Query 5: Course Material Semantic Search**
```sql
-- With course and unit filters
SELECT content_text, course_id, unit_no, 
       1 - (embedding <=> '[0.123, 0.456, ...]'::vector) as similarity_score,
       source_file
FROM COURSE_MATERIAL_CHUNK
WHERE course_id = ANY(ARRAY['CS101', 'CS102', 'CS103']::text[])
  AND unit_no = 2
ORDER BY embedding <=> '[0.123, 0.456, ...]'::vector ASC
LIMIT 3;
```

**Explanation:**
- `embedding <=> $1::vector`: Calculates cosine distance
- `1 - distance`: Converts to similarity (higher = better match)
- `ANY(ARRAY[...])`: Filters by multiple course IDs (student's enrolled courses)
- `ORDER BY ... ASC`: Closest matches first
- Vector literal format: `'[0.1, 0.2, ..., 0.768]'`

**Query 6: Course Recommendation (if implemented)**
```sql
SELECT course_id, title, description, credits,
       1 - (embedding <=> $1::vector) as relevance_score
FROM COURSE
ORDER BY embedding <=> $1::vector ASC
LIMIT 5;
```

### 6.3 Teacher-Specific Queries

**Query 7: Taught Courses & Sections**
```sql
SELECT c.course_id, c.title, t.section_name
FROM TEACHES t
JOIN COURSE c ON t.course_id = c.course_id
WHERE t.faculty_id = 'F001';
```

**Query 8: Students in Taught Courses**
```sql
SELECT c.title as course_title, 
       s.student_id, s.s_first_name, s.s_last_name,
       e.grade, e.status
FROM TEACHES t
JOIN ENROLLS_IN e ON t.course_id = e.course_id
JOIN STUDENT s ON e.student_id = s.student_id
JOIN COURSE c ON t.course_id = c.course_id
WHERE t.faculty_id = 'F001'
ORDER BY c.title, s.s_last_name;
```

### 6.4 Performance Queries

**Query 9: Student Risk Analysis (Backlog Detection)**
```sql
SELECT c.title, e.grade, e.status, e.backlog
FROM ENROLLS_IN e
JOIN COURSE c ON e.course_id = c.course_id
WHERE e.student_id = 'S1001'
  AND (e.grade IN ('D', 'F') OR e.backlog = TRUE);
```

**Query 10: Course-wise Performance Stats**
```sql
SELECT c.course_id, c.title,
       COUNT(*) as total_students,
       AVG(CASE e.grade
           WHEN 'O' THEN 10
           WHEN 'A' THEN 9
           WHEN 'B+' THEN 8
           WHEN 'B' THEN 7
           WHEN 'C+' THEN 6
           WHEN 'C' THEN 5
           WHEN 'D' THEN 4
           ELSE 0
       END) as average_grade_points
FROM COURSE c
LEFT JOIN ENROLLS_IN e ON c.course_id = e.course_id
WHERE e.grade IS NOT NULL
GROUP BY c.course_id, c.title
ORDER BY average_grade_points DESC;
```

---

## 7. Data Flow & Integration

### 7.1 Complete Chat Request Flow

```
[Frontend: chat/page.tsx]
  ↓ User types: "Explain pointers"
  ↓ POST /chat/message
  {
    student_id: "S1001",
    role: "student", 
    message: "Explain pointers",
    agent_id: "socratic"
  }
  
[Backend: handlers/chat.go]
  ↓ JWT validation (middleware/auth.go)
  ↓ Extract user_id from token
  ↓ Call RAGService.ProcessChat()
  
[services/rag_service.go]
  ↓
  ├─→ [Redis] Check cache: "response:<sha256_hash>"
  │   └─→ Cache hit? Return immediately
  │
  ├─→ [PostgreSQL] Build student context:
  │   ├─→ SELECT from STUDENT (profile)
  │   ├─→ SELECT from ENROLLS_IN + COURSE (grades, CGPA)
  │   └─→ SELECT from SCHEDULE (next class)
  │
  ├─→ [ai/embedder.go] Generate query embedding
  │   └─→ POST to Ollama: /api/embeddings
  │       Model: nomic-embed-text
  │       Returns: [0.123, 0.456, ..., 0.768]
  │
  ├─→ [repository/course_repo.go] SearchMaterials()
  │   └─→ [PostgreSQL] Vector search:
  │       SELECT ... FROM COURSE_MATERIAL_CHUNK
  │       WHERE course_id IN (enrolled courses)
  │       ORDER BY embedding <=> query_vector
  │       LIMIT 3
  │
  ├─→ [MongoDB] Fetch last 5 messages from ChatLogs
  │
  ├─→ Construct prompt with:
  │   - Agent personality instructions
  │   - Student context (CGPA, schedule, courses)
  │   - RAG materials (top 3 relevant chunks)
  │   - Conversation history
  │   - User sentiment
  │
  ├─→ POST to Ollama: /api/generate
  │   Model: llama3.2
  │   Prompt: <constructed_prompt>
  │   Returns: AI response
  │
  ├─→ [MongoDB] Save user message + bot response to ChatLogs
  │
  └─→ [Redis] Cache response (5 min TTL)

[Response to Frontend]
  {
    response: "Great question! Before I explain pointers...",
    user_id: "S1001",
    role: "student"
  }
  
[Frontend: chat/page.tsx]
  ↓ Display response with ReactMarkdown
  ↓ Update messages state
  ↓ Scroll to bottom
```

### 7.2 Quiz Generation Data Flow

```
[Frontend: quizzes/page.tsx]
  ↓ Select: course_id="CS101", unit=2, num_questions=10
  ↓ POST /quiz/generate
  
[Backend: handlers/quiz.go]
  ↓ Call QuizService.GenerateQuiz()
  
[services/quiz_service.go]
  ↓
  ├─→ [PostgreSQL] Fetch syllabus:
  │   SELECT topic FROM SYLLABUS_UNIT 
  │   WHERE course_id='CS101' AND unit_no=2
  │   Returns: ["Pointers", "Arrays", "Pointer Arithmetic"]
  │
  ├─→ [ai/embedder.go] Generate embedding:
  │   Query: "Important concepts in CS101 Unit 2"
  │   POST to Ollama: /api/embeddings
  │
  ├─→ [PostgreSQL] RAG search:
  │   SELECT content_text, source_file
  │   FROM COURSE_MATERIAL_CHUNK
  │   WHERE course_id='CS101' AND unit_no=2
  │   ORDER BY embedding <=> query_vector
  │   LIMIT 5
  │   Returns: Top 5 relevant paragraphs
  │
  ├─→ Construct quiz generation prompt:
  │   - Course: CS101
  │   - Unit: 2
  │   - Topics: Pointers, Arrays, Pointer Arithmetic
  │   - Materials: [5 paragraphs from PDFs]
  │   - Format: JSON with questions array
  │
  ├─→ POST to Ollama: /api/generate
  │   Model: llama3.2
  │   Format: "json" (enforces JSON mode)
  │   Returns: {questions: [...]}
  │
  ├─→ Parse JSON response
  │
  └─→ [MongoDB] Save quiz to quizzes collection

[Response to Frontend]
  {
    id: "...",
    course_id: "CS101",
    topic: "Generated from Course Materials",
    questions: [
      {
        id: 1,
        text: "What does the * operator do?",
        options: [...],
        correct_option: 1,
        reference: "Unit2_Pointers.pdf"
      }
    ]
  }

[Frontend: quizzes/page.tsx]
  ↓ Render quiz questions
  ↓ User answers questions
  ↓ On submit: Calculate score
  ↓ POST /ai/quiz-analysis (if wrong answers exist)
  ↓ Display analysis with weak areas & study priorities
```

---

## 8. AI/ML Components

### 8.1 Ollama Setup

**Local Installation:**
```bash
# Install Ollama
curl https://ollama.ai/install.sh | sh

# Pull models
ollama pull llama3.2
ollama pull nomic-embed-text

# Run server (default port 11434)
ollama serve
```

**Model Specifications:**
- **llama3.2**: 3B parameter LLM for text generation
- **nomic-embed-text**: 768-dimensional text embeddings

### 8.2 Embedding Generation (embedder.go)

```go
type EmbeddingRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}

func (e *Embedder) GenerateEmbedding(text string) ([]float32, error) {
    reqBody := EmbeddingRequest{
        Model:  "nomic-embed-text",
        Prompt: text,
    }
    
    resp, err := http.Post(
        "http://localhost:11434/api/embeddings",
        "application/json",
        bytes.NewBuffer(jsonData)
    )
    
    var embeddingResp EmbeddingResponse
    json.NewDecoder(resp.Body).Decode(&embeddingResp)
    
    // Convert float64 to float32 for pgvector
    float32Embedding := make([]float32, 768)
    for i, v := range embeddingResp.Embedding {
        float32Embedding[i] = float32(v)
    }
    
    return float32Embedding, nil
}
```

### 8.3 pgvector Extension

**Installation:**
```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

**Vector Operations:**
```sql
-- Cosine Distance (0 = identical, 2 = opposite)
embedding <=> '[0.1, 0.2, ...]'::vector

-- Euclidean Distance (L2)
embedding <-> '[0.1, 0.2, ...]'::vector

-- Inner Product (negative)
embedding <#> '[0.1, 0.2, ...]'::vector
```

**Indexing (for large datasets):**
```sql
-- IVFFlat index (approximate nearest neighbor)
CREATE INDEX ON COURSE_MATERIAL_CHUNK 
USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);
```

**Why no index in current schema?**
- Dataset is small (~5,000 chunks for 25 units)
- PostgreSQL can scan 5,000 rows in < 50ms
- Index overhead not justified for this scale
- Add index when > 100,000 chunks

### 8.4 Sentiment Analysis (Simple Heuristic)

```go
func (s *RAGService) AnalyzeSentiment(message string) string {
    lower := strings.ToLower(message)
    if strings.Contains(lower, "bad") || 
       strings.Contains(lower, "hate") || 
       strings.Contains(lower, "fail") {
        return "negative"
    } else if strings.Contains(lower, "good") || 
              strings.Contains(lower, "love") || 
              strings.Contains(lower, "thanks") {
        return "positive"
    }
    return "neutral"
}
```

**Used in prompt:**
```
[USER SENTIMENT]
User seems negative.
```
(Helps AI respond more empathetically)

---

## 9. Caching & Performance Optimization

### 9.1 Redis Caching Strategy

**Cache Layers:**

1. **Chat Response Cache** (5 min TTL)
```go
hash := sha256.Sum256([]byte(message))
cacheKey := "response:" + hex.EncodeToString(hash[:])
config.RedisClient.Set(ctx, cacheKey, response, 5*time.Minute)
```

2. **Student Profile Cache** (5 min TTL)
```go
cacheKey := "student_profile_v2:" + studentID
config.RedisClient.Set(ctx, cacheKey, jsonBytes, 5*time.Minute)
```

3. **Timetable Cache** (1 hour TTL)
```go
cacheKey := "timetable:" + studentID
config.RedisClient.Set(ctx, cacheKey, jsonBytes, time.Hour)
```

4. **Session Storage** (24 hour TTL)
```go
config.RedisClient.Set(ctx, "session:"+userID, jwtToken, 24*time.Hour)
```

### 9.2 Database Query Optimization

**Indexes (implicit from PRIMARY KEY / UNIQUE):**
- `STUDENT(student_id)` - PRIMARY KEY
- `COURSE(course_id)` - PRIMARY KEY
- `ENROLLS_IN(student_id, course_id)` - PRIMARY KEY (composite)
- `STUDENT(s_email)` - UNIQUE constraint creates index

**Additional Recommended Indexes:**
```sql
-- For schedule queries (day + time filtering)
CREATE INDEX idx_schedule_day_time 
ON SCHEDULE(day_of_week, start_time);

-- For enrollment status queries
CREATE INDEX idx_enrolls_status 
ON ENROLLS_IN(status) WHERE status = 'Enrolled';

-- For course material searches (if not using vector index)
CREATE INDEX idx_material_course_unit 
ON COURSE_MATERIAL_CHUNK(course_id, unit_no);
```

### 9.3 Connection Pooling

**PostgreSQL (database/config/db.go):**
```go
db, _ := sql.Open("postgres", connStr)
db.SetMaxOpenConns(25)      // Max concurrent connections
db.SetMaxIdleConns(10)      // Keep 10 idle for reuse
db.SetConnMaxLifetime(5 * time.Minute)
```

**MongoDB:**
```go
clientOptions := options.Client().ApplyURI(mongoURI).
    SetMaxPoolSize(50).
    SetMinPoolSize(10)
```

### 9.4 N+1 Query Prevention

**Example: Fetching Timetable**
```sql
-- ❌ BAD (N+1): Loop over courses, query schedule for each
SELECT course_id FROM ENROLLS_IN WHERE student_id = 'S1001';
-- Then for each course:
SELECT * FROM SCHEDULE WHERE course_id = 'CS101';

-- ✅ GOOD: Single JOIN query
SELECT c.course_id, c.title, sch.*
FROM ENROLLS_IN e
JOIN COURSE c ON e.course_id = c.course_id
JOIN SCHEDULE sch ON c.course_id = sch.course_id
WHERE e.student_id = 'S1001';
```

---

## 10. Potential Modifications with Code Diffs

This section covers common changes that might be requested during a DBMS project presentation or viva.

---

### Modification 1: Add Attendance Tracking Table

**Requirement:** "Add a table to track student attendance per class session."

**SQL Migration:**
```sql
-- File: database/04_add_attendance.sql

CREATE TABLE ATTENDANCE (
    attendance_id SERIAL PRIMARY KEY,
    student_id VARCHAR(20) NOT NULL,
    schedule_id INTEGER NOT NULL,
    date DATE NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('Present', 'Absent', 'Late')),
    marked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, schedule_id, date),
    FOREIGN KEY (student_id) REFERENCES STUDENT(student_id),
    FOREIGN KEY (schedule_id) REFERENCES SCHEDULE(schedule_id)
);

-- Index for quick lookups
CREATE INDEX idx_attendance_student_date 
ON ATTENDANCE(student_id, date);
```

**Go Model (internal/models/models.go):**
```diff
+type Attendance struct {
+    AttendanceID int       `json:"attendance_id"`
+    StudentID    string    `json:"student_id"`
+    ScheduleID   int       `json:"schedule_id"`
+    Date         time.Time `json:"date"`
+    Status       string    `json:"status"` // Present, Absent, Late
+    MarkedAt     time.Time `json:"marked_at"`
+}
```

**Handler (internal/handlers/student.go):**
```diff
+func GetAttendanceStats(c *gin.Context) {
+    rawID, exists := c.Get("user_id")
+    if !exists {
+        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
+        return
+    }
+    studentID := rawID.(string)
+    
+    query := `
+        SELECT c.course_id, c.title,
+               COUNT(*) as total_classes,
+               SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) as attended,
+               ROUND(100.0 * SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) / COUNT(*), 2) as percentage
+        FROM ATTENDANCE a
+        JOIN SCHEDULE sch ON a.schedule_id = sch.schedule_id
+        JOIN COURSE c ON sch.course_id = c.course_id
+        WHERE a.student_id = $1
+        GROUP BY c.course_id, c.title
+    `
+    
+    rows, err := config.PostgresDB.Query(query, studentID)
+    if err != nil {
+        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
+        return
+    }
+    defer rows.Close()
+    
+    type AttendanceStat struct {
+        CourseID   string  `json:"course_id"`
+        Title      string  `json:"title"`
+        Total      int     `json:"total_classes"`
+        Attended   int     `json:"attended"`
+        Percentage float64 `json:"percentage"`
+    }
+    
+    var stats []AttendanceStat
+    for rows.Next() {
+        var s AttendanceStat
+        rows.Scan(&s.CourseID, &s.Title, &s.Total, &s.Attended, &s.Percentage)
+        stats = append(stats, s)
+    }
+    
+    c.JSON(http.StatusOK, stats)
+}
```

**Router Registration (cmd/server/main.go):**
```diff
protected := r.Group("/student")
protected.Use(middleware.AuthMiddleware())
{
    protected.GET("/profile", handlers.GetStudentProfile)
    protected.GET("/timetable", handlers.GetStudentTimetable)
+   protected.GET("/attendance", handlers.GetAttendanceStats)
}
```

**Frontend (app/(dashboard)/attendance/page.tsx - NEW FILE):**
```tsx
"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import Cookies from "js-cookie"

interface AttendanceStat {
    course_id: string
    title: string
    total_classes: number
    attended: number
    percentage: number
}

export default function AttendancePage() {
    const [stats, setStats] = useState<AttendanceStat[]>([])

    useEffect(() => {
        const fetchStats = async () => {
            const token = Cookies.get("token")
            const res = await fetch("http://localhost:8080/student/attendance", {
                headers: { "Authorization": `Bearer ${token}` }
            })
            if (res.ok) {
                const data = await res.json()
                setStats(data)
            }
        }
        fetchStats()
    }, [])

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold">Attendance Records</h1>
            <div className="grid gap-4">
                {stats.map(stat => (
                    <Card key={stat.course_id}>
                        <CardHeader>
                            <CardTitle>{stat.title}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex justify-between">
                                <span>Attended: {stat.attended} / {stat.total_classes}</span>
                                <span className={stat.percentage < 75 ? "text-red-600" : "text-green-600"}>
                                    {stat.percentage}%
                                </span>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    )
}
```

---

### Modification 2: Add Assignment Submission Tracking

**SQL Migration:**
```sql
-- File: database/05_add_assignments.sql

CREATE TABLE ASSIGNMENT (
    assignment_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    due_date TIMESTAMP NOT NULL,
    max_marks INTEGER DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);

CREATE TABLE SUBMISSION (
    submission_id SERIAL PRIMARY KEY,
    assignment_id INTEGER NOT NULL,
    student_id VARCHAR(20) NOT NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    file_url TEXT,
    marks_obtained INTEGER,
    feedback TEXT,
    UNIQUE(assignment_id, student_id),
    FOREIGN KEY (assignment_id) REFERENCES ASSIGNMENT(assignment_id),
    FOREIGN KEY (student_id) REFERENCES STUDENT(student_id)
);

CREATE INDEX idx_assignment_course_due 
ON ASSIGNMENT(course_id, due_date);
```

**Query to Integrate in Chatbot Context:**
```diff
// In rag_service.go, ProcessChat() function:

+// Fetch Pending Assignments
+var pendingAssignments []string
+assQuery := `
+    SELECT a.title, TO_CHAR(a.due_date, 'YYYY-MM-DD HH24:MI')
+    FROM ASSIGNMENT a
+    JOIN ENROLLS_IN e ON a.course_id = e.course_id
+    LEFT JOIN SUBMISSION s ON a.assignment_id = s.assignment_id 
+        AND s.student_id = $1
+    WHERE e.student_id = $1 
+      AND a.due_date > NOW()
+      AND s.submission_id IS NULL
+    ORDER BY a.due_date ASC
+    LIMIT 3
+`
+assRows, err := config.PostgresDB.Query(assQuery, studentID)
+if err == nil {
+    defer assRows.Close()
+    for assRows.Next() {
+        var title, due string
+        if err := assRows.Scan(&title, &due); err == nil {
+            pendingAssignments = append(pendingAssignments, 
+                fmt.Sprintf("%s (Due: %s)", title, due))
+        }
+    }
+}
+
+var assignmentContext string
+if len(pendingAssignments) > 0 {
+    assignmentContext = fmt.Sprintf("\n[PENDING ASSIGNMENTS]:\n- %s", 
+        strings.Join(pendingAssignments, "\n- "))
+}

contextString = fmt.Sprintf(`
[CONTEXTUAL AWARENESS]
You are talking to %s, a %d Year %s student.
...
+%s
`, studentName, studentYear, deptID, assignmentContext)
```

---

### Modification 3: Filter Quiz by Difficulty Level

**SQL Migration:**
```sql
ALTER TABLE COURSE_MATERIAL_CHUNK 
ADD COLUMN difficulty VARCHAR(10) CHECK (difficulty IN ('Easy', 'Medium', 'Hard'));

UPDATE COURSE_MATERIAL_CHUNK 
SET difficulty = 'Medium' 
WHERE difficulty IS NULL;
```

**Frontend Update (quizzes/page.tsx):**
```diff
const [difficulty, setDifficulty] = useState<string>("Medium")

// In JSX:
+<select
+    value={difficulty}
+    onChange={(e) => setDifficulty(e.target.value)}
+    className="flex h-10 w-full rounded-md border"
+>
+    <option value="Easy">Easy</option>
+    <option value="Medium">Medium</option>
+    <option value="Hard">Hard</option>
+</select>

// In generateQuiz():
body: JSON.stringify({
    course_id: courseId,
    unit: unitId,
    num_questions: numQuestions,
+   difficulty: difficulty
})
```

**Backend Update (quiz_service.go):**
```diff
-func (s *QuizService) GenerateQuiz(courseID string, unit int, numQuestions int) (*models.Quiz, error) {
+func (s *QuizService) GenerateQuiz(courseID string, unit int, numQuestions int, difficulty string) (*models.Quiz, error) {

    // ... existing code ...
    
-   materials, err := s.Repo.SearchMaterials(ctx, embedding, 5, []string{courseID}, unit)
+   materials, err := s.Repo.SearchMaterialsWithDifficulty(ctx, embedding, 5, []string{courseID}, unit, difficulty)
```

**Repository Update (course_repo.go):**
```diff
+func (r *CourseRepository) SearchMaterialsWithDifficulty(
+    ctx context.Context, 
+    embedding []float32, 
+    limit int, 
+    courseIDFilter []string, 
+    unitFilter int,
+    difficulty string,
+) ([]CourseMaterial, error) {
+    vectorStr := vecToString(embedding)
+    
+    query := `
+        SELECT content_text, course_id, unit_no, source_file
+        FROM COURSE_MATERIAL_CHUNK
+        WHERE course_id = ANY($3::text[])
+          AND unit_no = $4
+          AND difficulty = $5
+        ORDER BY embedding <=> $1::vector ASC
+        LIMIT $2
+    `
+    
+    rows, err := r.DB.QueryContext(ctx, query, vectorStr, limit, 
+                                   courseIDFilter, unitFilter, difficulty)
+    // ... rest of function
+}
```

---

### Modification 4: Add Course Prerequisites

**SQL Migration:**
```sql
CREATE TABLE PREREQUISITE (
    course_id VARCHAR(10),
    prerequisite_course_id VARCHAR(10),
    PRIMARY KEY (course_id, prerequisite_course_id),
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id),
    FOREIGN KEY (prerequisite_course_id) REFERENCES COURSE(course_id)
);

-- Example data
INSERT INTO PREREQUISITE VALUES 
('CS102', 'CS101'),  -- Data Structures requires Intro to Programming
('CS103', 'CS102');  -- DBMS requires Data Structures
```

**Query to Check Eligibility:**
```sql
-- Check if student has completed prerequisites for CS103
SELECT p.prerequisite_course_id, c.title, 
       COALESCE(e.status, 'Not Taken') as completion_status
FROM PREREQUISITE p
JOIN COURSE c ON p.prerequisite_course_id = c.course_id
LEFT JOIN ENROLLS_IN e ON p.prerequisite_course_id = e.course_id 
    AND e.student_id = 'S1001'
WHERE p.course_id = 'CS103';
```

**Frontend: Course Recommendations with Prerequisites:**
```diff
// In recommendations/page.tsx

+// Fetch eligible courses (prerequisites met)
+const eligibleQuery = `
+    SELECT c.course_id, c.title, c.description
+    FROM COURSE c
+    WHERE c.course_id NOT IN (
+        SELECT course_id FROM ENROLLS_IN WHERE student_id = $1
+    )
+    AND NOT EXISTS (
+        SELECT 1 FROM PREREQUISITE p
+        WHERE p.course_id = c.course_id
+        AND NOT EXISTS (
+            SELECT 1 FROM ENROLLS_IN e
+            WHERE e.course_id = p.prerequisite_course_id
+            AND e.student_id = $1
+            AND e.status = 'Completed'
+        )
+    )
+`
```

---

### Modification 5: Enhanced RAG with Metadata Filtering

**SQL Migration:**
```sql
ALTER TABLE COURSE_MATERIAL_CHUNK 
ADD COLUMN keywords TEXT[],
ADD COLUMN chunk_type VARCHAR(20) CHECK (chunk_type IN ('Definition', 'Example', 'Theory', 'Code'));

-- Example update
UPDATE COURSE_MATERIAL_CHUNK 
SET keywords = ARRAY['pointer', 'memory', 'address'],
    chunk_type = 'Definition'
WHERE content_text LIKE '%pointer is a variable%';
```

**Enhanced Search Query:**
```diff
func (r *CourseRepository) SearchMaterialsAdvanced(
    ctx context.Context, 
    embedding []float32, 
    limit int, 
    filters map[string]interface{},
) ([]CourseMaterial, error) {
    vectorStr := vecToString(embedding)
    
    query := `
        SELECT content_text, course_id, unit_no, source_file, keywords
        FROM COURSE_MATERIAL_CHUNK
        WHERE 1=1
    `
    args := []interface{}{vectorStr, limit}
    argIdx := 3
    
+   if courseIDs, ok := filters["course_ids"].([]string); ok {
+       query += fmt.Sprintf(" AND course_id = ANY($%d::text[])", argIdx)
+       args = append(args, courseIDs)
+       argIdx++
+   }
+   
+   if chunkType, ok := filters["chunk_type"].(string); ok {
+       query += fmt.Sprintf(" AND chunk_type = $%d", argIdx)
+       args = append(args, chunkType)
+       argIdx++
+   }
+   
+   if keywords, ok := filters["keywords"].([]string); ok {
+       query += fmt.Sprintf(" AND keywords && $%d::text[]", argIdx)
+       args = append(args, pq.Array(keywords))
+       argIdx++
+   }
    
    query += " ORDER BY embedding <=> $1::vector ASC LIMIT $2"
    
    rows, err := r.DB.QueryContext(ctx, query, args...)
    // ...
}
```

---

### Modification 6: Add Real-Time Notifications (WebSocket)

**Backend: WebSocket Handler:**
```go
// File: internal/handlers/websocket.go

package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "sync"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins (restrict in production)
    },
}

type Hub struct {
    clients    map[string]*websocket.Conn
    mu         sync.RWMutex
}

var hub = &Hub{
    clients: make(map[string]*websocket.Conn),
}

func WebSocketHandler(c *gin.Context) {
    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(400, gin.H{"error": "user_id required"})
        return
    }
    
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }
    
    hub.mu.Lock()
    hub.clients[userID] = conn
    hub.mu.Unlock()
    
    defer func() {
        hub.mu.Lock()
        delete(hub.clients, userID)
        hub.mu.Unlock()
        conn.Close()
    }()
    
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func BroadcastNotification(userID string, message map[string]interface{}) {
    hub.mu.RLock()
    conn, exists := hub.clients[userID]
    hub.mu.RUnlock()
    
    if exists {
        conn.WriteJSON(message)
    }
}
```

**Router:**
```diff
+r.GET("/ws", handlers.WebSocketHandler)
```

**Frontend Hook (lib/useNotifications.ts - NEW FILE):**
```typescript
import { useEffect, useState } from 'react'

export function useNotifications(userId: string) {
    const [notifications, setNotifications] = useState<any[]>([])
    
    useEffect(() => {
        const ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}`)
        
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data)
            setNotifications(prev => [data, ...prev])
        }
        
        return () => ws.close()
    }, [userId])
    
    return notifications
}
```

---

### Modification 7: Add Grade History Visualization Query

**SQL Query for Semester-wise Performance:**
```sql
-- Get semester-wise CGPA trend
WITH semester_grades AS (
    SELECT 
        s.semester,
        c.credits,
        CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A' THEN 9 WHEN 'B+' THEN 8
            WHEN 'B' THEN 7 WHEN 'C+' THEN 6 WHEN 'C' THEN 5
            WHEN 'D' THEN 4 ELSE 0
        END as grade_points
    FROM ENROLLS_IN e
    JOIN COURSE c ON e.course_id = c.course_id
    JOIN (
        SELECT course_id, MIN(semester) as semester
        FROM enrollment_history  -- Assuming you track when courses were taken
        WHERE student_id = 'S1001'
        GROUP BY course_id
    ) s ON e.course_id = s.course_id
    WHERE e.student_id = 'S1001' AND e.grade IS NOT NULL
)
SELECT 
    semester,
    ROUND(SUM(grade_points * credits)::numeric / SUM(credits), 2) as semester_gpa
FROM semester_grades
GROUP BY semester
ORDER BY semester;
```

**Handler:**
```go
func GetGradeHistory(c *gin.Context) {
    studentID := c.GetString("user_id")
    
    query := `
        WITH semester_grades AS (
            -- ... (query above)
        )
        SELECT semester, semester_gpa FROM semester_grades
    `
    
    rows, _ := config.PostgresDB.Query(query, studentID)
    defer rows.Close()
    
    type SemesterGPA struct {
        Semester int     `json:"semester"`
        GPA      float64 `json:"gpa"`
    }
    
    var history []SemesterGPA
    for rows.Next() {
        var s SemesterGPA
        rows.Scan(&s.Semester, &s.GPA)
        history = append(history, s)
    }
    
    c.JSON(200, history)
}
```

---

## Summary of SQL Query Categories

### **Read Queries (SELECT)**
1. Student profile retrieval
2. CGPA calculation with joins
3. Timetable with schedule joins
4. Vector similarity search (RAG)
5. Chat history from MongoDB
6. Quiz generation material search
7. Teacher's course and student lists
8. Attendance statistics
9. Assignment deadlines
10. Grade history trends

### **Write Queries (INSERT/UPDATE)**
1. Chat log insertion (MongoDB)
2. Quiz storage (MongoDB)
3. Attendance marking
4. Assignment submission
5. Grade updates
6. Cache invalidation triggers

### **Complex Queries**
1. Recursive CTEs for prerequisites
2. Window functions for ranking
3. Aggregations for statistics
4. Joins across 4+ tables
5. Vector operations with filtering

---

## Key Takeaways for Presentation

1. **3NF Normalization**: Explain with examples from STUDENT/DEPARTMENT
2. **Vector Search**: Demonstrate <=> operator and similarity scoring
3. **Personalization**: Show how context is built from multiple tables
4. **RAG Pipeline**: Explain embedding → search → prompt → generate
5. **Performance**: Discuss caching strategy and index usage
6. **Scalability**: MongoDB for logs, PostgreSQL for structured data
7. **Real-world SQL**: Show actual queries used in production code

---

## Questions to Prepare For

**Q: Why use pgvector instead of a dedicated vector database?**
A: For our scale (5,000 vectors), PostgreSQL with pgvector is sufficient and reduces operational complexity. It also allows ACID transactions combining vector search with relational joins in a single query.

**Q: How do you prevent SQL injection?**
A: We use parameterized queries throughout. Example: `config.PostgresDB.Query("SELECT * FROM STUDENT WHERE student_id=$1", userID)` - the `$1` placeholder ensures safe parameter binding.

**Q: How is the chatbot personalized differently for students vs teachers?**
A: Students get context about their grades, CGPA, next class, and enrolled courses. Teachers get lists of taught courses and enrolled students. The RAG filter also changes - students see only their course materials, teachers see materials from courses they teach.

**Q: What happens if Ollama is down?**
A: The system has fallback logic. In `rag_service.go`, if the HTTP call fails, it returns a simulated response. In production, you'd implement a circuit breaker pattern or queue system.

**Q: How do you ensure quiz questions come from actual course material?**
A: The prompt explicitly instructs the LLM to use ONLY provided context. Each question includes a `reference` field showing the source PDF. We retrieve top 5 most relevant chunks using vector search before generating questions.

**Q: Can you explain the CGPA calculation logic?**
A: We use a 10-point scale: O/A+=10, A=9, B+=8, etc. For each course: `points = grade_value * credits`. Then `CGPA = SUM(points) / SUM(credits)`. The calculation happens in Go after fetching grades from the database.

---

**End of Documentation**

This document covers the complete implementation of AcademAide with focus on database design, SQL queries, personalization mechanisms, and RAG integration. Use this as a reference for your DBMS project presentation.
