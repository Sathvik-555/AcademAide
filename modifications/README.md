# AcademAide Modifications - README

## Overview
This directory contains detailed implementation guides for potential modifications to the AcademAide DBMS project. Each modification includes complete SQL migrations, backend code changes, and frontend updates with actual code diffs.

## Purpose
These modifications are designed for:
- **DBMS Project Presentations**: Demonstrate advanced SQL knowledge
- **Viva Questions**: Be prepared to explain and implement changes on the spot
- **Feature Extensions**: Add real functionality to the project

---

## Available Modifications

### 📊 Modification 1: Attendance Tracking System
**File:** `01_attendance_tracking.md`

**What it adds:**
- Complete attendance table with status tracking (Present/Absent/Late)
- Course-wise and overall attendance percentage calculations
- Warning system for students below 75% (exam eligibility threshold)
- Backend API endpoints for fetching attendance stats
- Beautiful frontend page with progress bars and risk indicators

**Key SQL Features:**
- Materialized view pattern for performance (`ATTENDANCE_SUMMARY`)
- Aggregate functions with `CASE` statements
- Composite indexes for optimization
- CHECK constraints for data integrity

**When to use:**
- Interviewer asks: "How would you track student attendance?"
- Need to demonstrate aggregation queries
- Want to show view optimization techniques

---

### 🎯 Modification 2: Intent Classification for RAG
**File:** `02_intent_classification.md`

**What it adds:**
- Classifies user queries into intents (code, concept, exam_prep, advice, etc.)
- Routes queries to specialized prompts for better AI responses
- Filters course materials by type (Code, Theory, Example, Definition)
- Enhanced metadata on chunks (keywords, difficulty, chunk_type)

**Key SQL Features:**
- Array columns (`keywords TEXT[]`)
- GIN indexes for array searching
- Metadata filtering in vector search queries
- Custom filtering logic with `ANY()` operator

**When to use:**
- Interviewer asks: "How can you improve RAG relevance?"
- Need to demonstrate advanced data types (arrays)
- Want to show query optimization with metadata

---

### 🔗 Modification 3: Course Prerequisites System
**File:** `03_prerequisites_system.md`

**What it adds:**
- Prerequisite table with mandatory/optional requirements
- Minimum grade requirements for prerequisites
- Database trigger to prevent invalid enrollments
- Student eligibility view for available courses
- Frontend course catalog with prerequisite warnings

**Key SQL Features:**
- Recursive CTEs for prerequisite chains
- Database triggers (`BEFORE INSERT`)
- CHECK constraints to prevent self-reference
- Complex eligibility logic with `NOT EXISTS` subqueries
- PL/pgSQL functions

**When to use:**
- Interviewer asks: "How do you ensure data integrity?"
- Need to demonstrate triggers and constraints
- Want to show recursive queries

---

### 📈 Modification 4: Advanced Analytics with Window Functions
**File:** `04_advanced_analytics.md`

**What it adds:**
- Semester-wise performance trends with cumulative GPA
- Department ranking and percentile calculation
- Course difficulty metrics (avg grade, std deviation, failure rate)
- At-risk student identification with risk scoring
- Recursive prerequisite tree traversal
- Time-series enrollment trend analysis

**Key SQL Features:**
- Window functions (`RANK()`, `LAG()`, `AVG() OVER()`)
- Recursive CTEs for prerequisite graphs
- CTEs for complex multi-step queries
- Moving averages for trend analysis
- `PERCENT_RANK()` for percentile calculation
- Materialized views for performance

**When to use:**
- Interviewer asks: "Demonstrate window functions"
- Need to show advanced SQL mastery
- Want impressive analytics queries
- **THIS IS THE BEST ONE FOR IMPRESSING INTERVIEWERS**

---

## How to Use This Directory

### During Preparation
1. Read each modification thoroughly
2. Understand the SQL queries and logic
3. Practice explaining window functions, CTEs, triggers
4. Try implementing one modification completely

### During Presentation
1. **If asked about a missing feature:**
   - "Yes, we can add that. Let me show you the implementation."
   - Open the relevant modification file
   - Explain the SQL schema changes
   - Show the code diffs

2. **If asked to write a query on the spot:**
   - Reference similar queries from the modifications
   - Adapt the logic to the specific question
   - Example: "This is similar to the attendance percentage calculation I designed..."

3. **If asked about advanced SQL:**
   - Use Modification 4 (Advanced Analytics)
   - Explain window functions with the ranking example
   - Show recursive CTE with prerequisites
   - Demonstrate CTEs with at-risk student detection

---

## Common Viva Questions & Answers

### Q1: "How would you track student attendance?"
**Answer:** Reference `01_attendance_tracking.md`

**Quick Explanation:**
- Create `ATTENDANCE` table with student_id, schedule_id, date, status
- Use `UNIQUE(student_id, schedule_id, date)` to prevent duplicates
- Create a view to aggregate: `100.0 * SUM(CASE WHEN status='Present' THEN 1 ELSE 0 END) / COUNT(*)`
- Add indexes on (student_id, date) for fast lookups

**SQL to write:**
```sql
SELECT c.title, 
       COUNT(*) as total,
       SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END) as attended,
       ROUND(100.0 * SUM(...) / COUNT(*), 2) as percentage
FROM ATTENDANCE a
JOIN SCHEDULE sch ON a.schedule_id = sch.schedule_id
JOIN COURSE c ON sch.course_id = c.course_id
WHERE a.student_id = 'S1001'
GROUP BY c.title;
```

---

### Q2: "Explain window functions with an example"
**Answer:** Reference `04_advanced_analytics.md`

**Quick Explanation:**
- Window functions perform calculations across a set of rows related to the current row
- Unlike GROUP BY, they don't collapse rows
- Example: Calculate cumulative GPA

**SQL to write:**
```sql
SELECT student_id, semester, gpa,
       AVG(gpa) OVER (
           PARTITION BY student_id 
           ORDER BY semester
           ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
       ) as cumulative_gpa,
       RANK() OVER (
           PARTITION BY dept_id 
           ORDER BY gpa DESC
       ) as dept_rank
FROM semester_grades;
```

**Explanation:**
- `PARTITION BY`: Like GROUP BY, defines the window
- `ORDER BY`: Defines the order within the window
- `ROWS BETWEEN`: Defines the frame (which rows to include)
- `RANK()`: Assigns ranks, ties get the same rank

---

### Q3: "How do you prevent invalid data insertion?"
**Answer:** Reference `03_prerequisites_system.md`

**Quick Explanation:**
- Use CHECK constraints for simple validations
- Use triggers for complex business logic
- Example: Prevent enrolling without prerequisites

**SQL to write:**
```sql
CREATE OR REPLACE FUNCTION check_prerequisites()
RETURNS TRIGGER AS $$
DECLARE
    missing_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO missing_count
    FROM PREREQUISITE p
    WHERE p.course_id = NEW.course_id
    AND NOT EXISTS (
        SELECT 1 FROM ENROLLS_IN e
        WHERE e.student_id = NEW.student_id
        AND e.course_id = p.prerequisite_course_id
        AND e.status = 'Completed'
    );
    
    IF missing_count > 0 THEN
        RAISE EXCEPTION 'Missing prerequisites';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_prerequisites
BEFORE INSERT ON ENROLLS_IN
FOR EACH ROW EXECUTE FUNCTION check_prerequisites();
```

---

### Q4: "Show me a recursive query"
**Answer:** Reference `04_advanced_analytics.md`

**Quick Explanation:**
- Recursive CTEs are used for hierarchical or graph data
- Example: Find all prerequisites (including transitive)

**SQL to write:**
```sql
WITH RECURSIVE prereq_tree AS (
    -- Base case
    SELECT prerequisite_course_id, 1 as level
    FROM PREREQUISITE
    WHERE course_id = 'CS103'
    
    UNION ALL
    
    -- Recursive case
    SELECT p.prerequisite_course_id, pt.level + 1
    FROM prereq_tree pt
    JOIN PREREQUISITE p ON pt.prerequisite_course_id = p.course_id
    WHERE pt.level < 5
)
SELECT * FROM prereq_tree;
```

**Result:**
```
prerequisite_course_id | level
-----------------------+-------
CS102                  | 1
CS101                  | 2
```

---

### Q5: "How do you optimize this query?"
**Answer:** General optimization techniques

**Strategies:**
1. **Add indexes:**
   ```sql
   CREATE INDEX idx_student_id ON ENROLLS_IN(student_id);
   ```

2. **Use views for common queries:**
   ```sql
   CREATE MATERIALIZED VIEW student_summary AS
   SELECT student_id, COUNT(*) as courses FROM ENROLLS_IN GROUP BY student_id;
   ```

3. **Avoid SELECT *:** Only fetch needed columns

4. **Use EXPLAIN ANALYZE:**
   ```sql
   EXPLAIN ANALYZE SELECT * FROM ENROLLS_IN WHERE student_id = 'S1001';
   ```

5. **Partition large tables** (if millions of rows)

6. **Use connection pooling** in application

---

## Implementation Checklist

When implementing a modification:

- [ ] Run SQL migration file
- [ ] Update Go models in `internal/models/models.go`
- [ ] Add handler functions in `internal/handlers/`
- [ ] Register routes in `cmd/server/main.go`
- [ ] Create frontend page in `frontend/AcademAide/app/`
- [ ] Test with sample data
- [ ] Verify SQL query performance with `EXPLAIN`

---

## SQL Concepts Coverage

| Concept | Modification File |
|---------|------------------|
| Aggregate Functions | All |
| Window Functions | 04_advanced_analytics |
| Recursive CTEs | 03_prerequisites, 04_advanced_analytics |
| Non-recursive CTEs | 04_advanced_analytics |
| Triggers | 03_prerequisites |
| Views | 01_attendance, 04_advanced_analytics |
| Materialized Views | 04_advanced_analytics |
| Indexes | All |
| CHECK Constraints | 01_attendance, 03_prerequisites |
| Foreign Keys | All |
| Array Columns | 02_intent_classification |
| PL/pgSQL Functions | 03_prerequisites, 04_advanced_analytics |

---

## Tips for Presentation

1. **Start with the main documentation:** Show `PROJECT_DOCUMENTATION.md` first
2. **Highlight 3NF normalization:** Explain STUDENT/DEPARTMENT separation
3. **Demonstrate RAG:** Show the vector search query with <=> operator
4. **Prepare one modification:** Practice Modification 4 (most impressive)
5. **Know your metrics:** Mention 5,000 chunks, 768 dimensions, 3NF normalized
6. **Be honest:** If asked about something you haven't implemented, show the modification file and explain how you would add it

---

## Quick Reference: SQL Query Templates

### Aggregation with CASE
```sql
SELECT course_id,
       SUM(CASE WHEN grade = 'A' THEN 1 ELSE 0 END) as a_count,
       AVG(CASE WHEN grade = 'A' THEN 10 WHEN grade = 'B' THEN 8 ELSE 0 END) as avg_points
FROM ENROLLS_IN
GROUP BY course_id;
```

### Window Function (Ranking)
```sql
SELECT student_id, gpa,
       RANK() OVER (ORDER BY gpa DESC) as rank,
       DENSE_RANK() OVER (ORDER BY gpa DESC) as dense_rank,
       ROW_NUMBER() OVER (ORDER BY gpa DESC) as row_num
FROM students;
```

### Recursive CTE
```sql
WITH RECURSIVE tree AS (
    SELECT id, parent_id, 1 as level FROM nodes WHERE parent_id IS NULL
    UNION ALL
    SELECT n.id, n.parent_id, t.level + 1
    FROM nodes n JOIN tree t ON n.parent_id = t.id
)
SELECT * FROM tree;
```

### Trigger
```sql
CREATE TRIGGER trigger_name
BEFORE INSERT ON table_name
FOR EACH ROW EXECUTE FUNCTION function_name();
```

### View with Joins
```sql
CREATE VIEW view_name AS
SELECT s.id, s.name, c.title
FROM STUDENT s
JOIN ENROLLS_IN e ON s.student_id = e.student_id
JOIN COURSE c ON e.course_id = c.course_id;
```

---

## Need Help?

- Review main documentation: `../PROJECT_DOCUMENTATION.md`
- Practice SQL queries on actual database
- Understand the "Why" behind each modification
- Be able to explain trade-offs (e.g., views vs materialized views)

**Good luck with your presentation! 🚀**
