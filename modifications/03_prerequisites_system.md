# Modification 3: Add Course Prerequisites System

## Overview
Implements a prerequisite tracking system to prevent students from enrolling in courses without completing required prerequisites.

## Files Changed
- `database/06_add_prerequisites.sql` (NEW)
- `internal/handlers/student.go` (MODIFIED)
- `frontend/AcademAide/app/(dashboard)/courses/page.tsx` (NEW)

---

## File 1: database/06_add_prerequisites.sql (NEW FILE)

```sql
-- Migration: Add Course Prerequisites
-- Ensures academic progression integrity

-- Create Prerequisites Table
CREATE TABLE PREREQUISITE (
    course_id VARCHAR(10) NOT NULL,
    prerequisite_course_id VARCHAR(10) NOT NULL,
    is_mandatory BOOLEAN DEFAULT TRUE,
    minimum_grade VARCHAR(2),  -- Minimum grade required in prerequisite
    PRIMARY KEY (course_id, prerequisite_course_id),
    FOREIGN KEY (course_id) REFERENCES COURSE(course_id) ON DELETE CASCADE,
    FOREIGN KEY (prerequisite_course_id) REFERENCES COURSE(course_id) ON DELETE CASCADE,
    CHECK (course_id != prerequisite_course_id)  -- Prevent self-reference
);

-- Create index for faster prerequisite lookups
CREATE INDEX idx_prerequisite_course ON PREREQUISITE(course_id);
CREATE INDEX idx_prerequisite_prereq ON PREREQUISITE(prerequisite_course_id);

-- Sample Data: Define realistic prerequisites
INSERT INTO PREREQUISITE (course_id, prerequisite_course_id, is_mandatory, minimum_grade) VALUES
-- Data Structures requires Intro to Programming
('CS102', 'CS101', TRUE, 'C'),
-- Database Systems requires Data Structures
('CS103', 'CS102', TRUE, 'C+'),
-- Advanced algorithms requires both
('CS201', 'CS101', TRUE, 'B'),
('CS201', 'CS102', TRUE, 'B+');

-- View: Student Eligibility for Courses
CREATE OR REPLACE VIEW STUDENT_COURSE_ELIGIBILITY AS
SELECT 
    s.student_id,
    s.s_first_name,
    s.s_last_name,
    c.course_id,
    c.title as course_title,
    c.credits,
    CASE 
        WHEN NOT EXISTS (
            SELECT 1 FROM PREREQUISITE p
            WHERE p.course_id = c.course_id
        ) THEN TRUE
        WHEN NOT EXISTS (
            SELECT 1 FROM PREREQUISITE p
            WHERE p.course_id = c.course_id
            AND NOT EXISTS (
                SELECT 1 FROM ENROLLS_IN e
                WHERE e.student_id = s.student_id
                AND e.course_id = p.prerequisite_course_id
                AND e.status = 'Completed'
                AND (p.minimum_grade IS NULL 
                     OR e.grade >= p.minimum_grade)
            )
        ) THEN TRUE
        ELSE FALSE
    END as is_eligible,
    ARRAY(
        SELECT pc.title 
        FROM PREREQUISITE p
        JOIN COURSE pc ON p.prerequisite_course_id = pc.course_id
        WHERE p.course_id = c.course_id
    ) as prerequisite_courses
FROM STUDENT s
CROSS JOIN COURSE c
WHERE c.dept_id = s.dept_id  -- Only show courses from student's department
  AND NOT EXISTS (
      SELECT 1 FROM ENROLLS_IN e
      WHERE e.student_id = s.student_id
      AND e.course_id = c.course_id
      AND e.status IN ('Enrolled', 'Completed')
  );

-- Function: Check if student can enroll in course
CREATE OR REPLACE FUNCTION can_enroll(
    p_student_id VARCHAR(20),
    p_course_id VARCHAR(10)
) RETURNS TABLE(
    can_enroll BOOLEAN,
    missing_prerequisites TEXT[]
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        CASE 
            WHEN COUNT(*) = 0 THEN TRUE
            ELSE FALSE
        END as can_enroll,
        ARRAY_AGG(pc.title)::TEXT[] as missing_prerequisites
    FROM PREREQUISITE p
    LEFT JOIN COURSE pc ON p.prerequisite_course_id = pc.course_id
    WHERE p.course_id = p_course_id
    AND NOT EXISTS (
        SELECT 1 FROM ENROLLS_IN e
        WHERE e.student_id = p_student_id
        AND e.course_id = p.prerequisite_course_id
        AND e.status = 'Completed'
        AND (p.minimum_grade IS NULL 
             OR e.grade >= p.minimum_grade)
    );
END;
$$ LANGUAGE plpgsql;

-- Trigger: Prevent enrollment without prerequisites
CREATE OR REPLACE FUNCTION check_prerequisites_before_enrollment()
RETURNS TRIGGER AS $$
DECLARE
    v_missing_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO v_missing_count
    FROM PREREQUISITE p
    WHERE p.course_id = NEW.course_id
    AND NOT EXISTS (
        SELECT 1 FROM ENROLLS_IN e
        WHERE e.student_id = NEW.student_id
        AND e.course_id = p.prerequisite_course_id
        AND e.status = 'Completed'
        AND (p.minimum_grade IS NULL 
             OR e.grade >= p.minimum_grade)
    );
    
    IF v_missing_count > 0 THEN
        RAISE EXCEPTION 'Cannot enroll: Missing % prerequisite(s) for course %', 
            v_missing_count, NEW.course_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_prerequisites
BEFORE INSERT ON ENROLLS_IN
FOR EACH ROW
EXECUTE FUNCTION check_prerequisites_before_enrollment();

-- Query: Find all courses a student is eligible for
-- Example usage:
-- SELECT * FROM STUDENT_COURSE_ELIGIBILITY 
-- WHERE student_id = 'S1001' AND is_eligible = TRUE;
```

---

## File 2: internal/handlers/student.go (MODIFICATION)

### Diff:
```diff
--- a/internal/handlers/student.go
+++ b/internal/handlers/student.go
@@ -393,3 +393,120 @@ func GetStudentResources(c *gin.Context) {
 
 	c.JSON(http.StatusOK, resources)
 }
+
+// GetEligibleCourses returns courses a student can enroll in
+func GetEligibleCourses(c *gin.Context) {
+	rawID, exists := c.Get("user_id")
+	if !exists {
+		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
+		return
+	}
+	studentID := rawID.(string)
+	
+	query := `
+		SELECT course_id, course_title, credits, is_eligible, prerequisite_courses
+		FROM STUDENT_COURSE_ELIGIBILITY
+		WHERE student_id = $1
+		ORDER BY is_eligible DESC, course_id
+	`
+	
+	rows, err := config.PostgresDB.Query(query, studentID)
+	if err != nil {
+		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
+		return
+	}
+	defer rows.Close()
+	
+	type EligibleCourse struct {
+		CourseID      string   `json:"course_id"`
+		CourseTitle   string   `json:"course_title"`
+		Credits       int      `json:"credits"`
+		IsEligible    bool     `json:"is_eligible"`
+		Prerequisites []string `json:"prerequisites"`
+	}
+	
+	var courses []EligibleCourse
+	for rows.Next() {
+		var course EligibleCourse
+		var prereqArray sql.NullString  // PostgreSQL array as string
+		
+		err := rows.Scan(&course.CourseID, &course.CourseTitle, &course.Credits,
+		                 &course.IsEligible, &prereqArray)
+		if err != nil {
+			continue
+		}
+		
+		// Parse PostgreSQL array {item1,item2} to Go slice
+		if prereqArray.Valid && prereqArray.String != "{}" {
+			prereqStr := strings.Trim(prereqArray.String, "{}")
+			course.Prerequisites = strings.Split(prereqStr, ",")
+		}
+		
+		courses = append(courses, course)
+	}
+	
+	c.JSON(http.StatusOK, courses)
+}
+
+// CheckEnrollmentEligibility checks if student can enroll in a specific course
+func CheckEnrollmentEligibility(c *gin.Context) {
+	rawID, exists := c.Get("user_id")
+	if !exists {
+		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
+		return
+	}
+	studentID := rawID.(string)
+	courseID := c.Param("course_id")
+	
+	if courseID == "" {
+		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id required"})
+		return
+	}
+	
+	// Use the function we created in SQL
+	query := `SELECT * FROM can_enroll($1, $2)`
+	
+	var canEnroll bool
+	var missingPrereqs sql.NullString
+	
+	err := config.PostgresDB.QueryRow(query, studentID, courseID).Scan(&canEnroll, &missingPrereqs)
+	if err != nil {
+		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
+		return
+	}
+	
+	result := gin.H{
+		"can_enroll": canEnroll,
+		"course_id":  courseID,
+	}
+	
+	if !canEnroll && missingPrereqs.Valid {
+		prereqStr := strings.Trim(missingPrereqs.String, "{}")
+		result["missing_prerequisites"] = strings.Split(prereqStr, ",")
+		result["message"] = "You must complete the prerequisite courses first"
+	} else if canEnroll {
+		result["message"] = "You are eligible to enroll in this course"
+	}
+	
+	c.JSON(http.StatusOK, result)
+}
+
+// GetCoursePrerequisites returns prerequisites for a given course
+func GetCoursePrerequisites(c *gin.Context) {
+	courseID := c.Param("course_id")
+	
+	query := `
+		SELECT pc.course_id, pc.title, p.minimum_grade, p.is_mandatory
+		FROM PREREQUISITE p
+		JOIN COURSE pc ON p.prerequisite_course_id = pc.course_id
+		WHERE p.course_id = $1
+		ORDER BY p.is_mandatory DESC, pc.course_id
+	`
+	
+	rows, err := config.PostgresDB.Query(query, courseID)
+	if err != nil {
+		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
+		return
+	}
+	defer rows.Close()
+	
+	// ... Parse and return prerequisites
+	c.JSON(http.StatusOK, gin.H{"prerequisites": []gin.H{}})
+}
```

---

## File 3: cmd/server/main.go (MODIFICATION)

### Diff:
```diff
--- a/cmd/server/main.go
+++ b/cmd/server/main.go
@@ -60,6 +60,9 @@ func main() {
 		protected.GET("/announcements", handlers.GetAnnouncements)
 		protected.GET("/attendance", handlers.GetAttendanceStats)
 		protected.GET("/attendance/detailed", handlers.GetDetailedAttendance)
+		protected.GET("/courses/eligible", handlers.GetEligibleCourses)
+		protected.GET("/courses/:course_id/eligibility", handlers.CheckEnrollmentEligibility)
+		protected.GET("/courses/:course_id/prerequisites", handlers.GetCoursePrerequisites)
 	}
```

---

## File 4: frontend/AcademAide/app/(dashboard)/courses/page.tsx (NEW FILE)

```tsx
"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { CheckCircle, XCircle, AlertCircle, BookOpen } from "lucide-react"
import Cookies from "js-cookie"

interface Course {
    course_id: string
    course_title: string
    credits: number
    is_eligible: boolean
    prerequisites: string[]
}

export default function CoursesPage() {
    const [courses, setCourses] = useState<Course[]>([])
    const [loading, setLoading] = useState(true)
    const [filter, setFilter] = useState<"all" | "eligible" | "blocked">("all")

    useEffect(() => {
        const fetchCourses = async () => {
            const token = Cookies.get("token")
            if (!token) return

            try {
                const res = await fetch("http://localhost:8080/student/courses/eligible", {
                    headers: { "Authorization": `Bearer ${token}` }
                })
                if (res.ok) {
                    const data = await res.json()
                    setCourses(data || [])
                }
            } catch (error) {
                console.error("Failed to fetch courses", error)
            } finally {
                setLoading(false)
            }
        }
        fetchCourses()
    }, [])

    const filteredCourses = courses.filter(course => {
        if (filter === "eligible") return course.is_eligible
        if (filter === "blocked") return !course.is_eligible
        return true
    })

    if (loading) {
        return <div className="flex items-center justify-center h-64">Loading courses...</div>
    }

    return (
        <div className="flex flex-col gap-6 max-w-6xl mx-auto pb-10">
            {/* Header */}
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight gradient-text flex items-center gap-2">
                    <BookOpen className="h-8 w-8 text-primary" />
                    Available Courses
                </h1>
                <p className="text-muted-foreground">
                    View courses you can enroll in based on your academic progress
                </p>
            </div>

            {/* Filter Tabs */}
            <div className="flex gap-2 border-b">
                <button
                    onClick={() => setFilter("all")}
                    className={`px-4 py-2 font-medium transition-colors ${
                        filter === "all"
                            ? "border-b-2 border-primary text-primary"
                            : "text-muted-foreground hover:text-foreground"
                    }`}
                >
                    All Courses ({courses.length})
                </button>
                <button
                    onClick={() => setFilter("eligible")}
                    className={`px-4 py-2 font-medium transition-colors ${
                        filter === "eligible"
                            ? "border-b-2 border-green-600 text-green-600"
                            : "text-muted-foreground hover:text-foreground"
                    }`}
                >
                    Eligible ({courses.filter(c => c.is_eligible).length})
                </button>
                <button
                    onClick={() => setFilter("blocked")}
                    className={`px-4 py-2 font-medium transition-colors ${
                        filter === "blocked"
                            ? "border-b-2 border-red-600 text-red-600"
                            : "text-muted-foreground hover:text-foreground"
                    }`}
                >
                    Prerequisites Required ({courses.filter(c => !c.is_eligible).length})
                </button>
            </div>

            {/* Course Cards */}
            <div className="grid gap-4">
                {filteredCourses.length === 0 ? (
                    <Card className="glass">
                        <CardContent className="pt-6 text-center text-muted-foreground">
                            No courses found in this category
                        </CardContent>
                    </Card>
                ) : (
                    filteredCourses.map((course) => (
                        <Card 
                            key={course.course_id} 
                            className={`glass border-none hover:shadow-md transition-shadow ${
                                !course.is_eligible ? "opacity-75" : ""
                            }`}
                        >
                            <CardHeader>
                                <div className="flex items-start justify-between">
                                    <div className="flex-1">
                                        <CardTitle className="flex items-center gap-2">
                                            {course.course_title}
                                            {course.is_eligible ? (
                                                <CheckCircle className="h-5 w-5 text-green-600" />
                                            ) : (
                                                <XCircle className="h-5 w-5 text-red-600" />
                                            )}
                                        </CardTitle>
                                        <CardDescription className="mt-1">
                                            {course.course_id} • {course.credits} Credits
                                        </CardDescription>
                                    </div>
                                    {course.is_eligible ? (
                                        <Button className="bg-green-600 hover:bg-green-700">
                                            Enroll Now
                                        </Button>
                                    ) : (
                                        <Button variant="outline" disabled>
                                            Not Eligible
                                        </Button>
                                    )}
                                </div>
                            </CardHeader>

                            {!course.is_eligible && course.prerequisites.length > 0 && (
                                <CardContent>
                                    <div className="p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
                                        <div className="flex items-start gap-2">
                                            <AlertCircle className="h-4 w-4 text-amber-600 mt-0.5" />
                                            <div className="flex-1">
                                                <div className="font-medium text-amber-900 dark:text-amber-100 mb-2">
                                                    Prerequisites Required:
                                                </div>
                                                <ul className="space-y-1 text-sm text-amber-800 dark:text-amber-200">
                                                    {course.prerequisites.map((prereq, idx) => (
                                                        <li key={idx} className="flex items-center gap-2">
                                                            <span className="w-1.5 h-1.5 bg-amber-600 rounded-full" />
                                                            {prereq}
                                                        </li>
                                                    ))}
                                                </ul>
                                                <p className="text-xs text-amber-700 dark:text-amber-300 mt-3">
                                                    Complete these courses with a passing grade before enrolling
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            )}
                        </Card>
                    ))
                )}
            </div>
        </div>
    )
}
```

---

## Key SQL Queries Explained

### Query 1: Check Eligibility for All Courses
```sql
SELECT course_id, course_title, credits, is_eligible, prerequisite_courses
FROM STUDENT_COURSE_ELIGIBILITY
WHERE student_id = 'S1001'
ORDER BY is_eligible DESC, course_id
```

**How it works:**
- View pre-computes eligibility using `NOT EXISTS` subquery
- Checks if all prerequisites are completed with required grade
- Returns array of prerequisite course titles
- Orders eligible courses first

### Query 2: Check Specific Course
```sql
SELECT * FROM can_enroll('S1001', 'CS103')
```

**Returns:**
```
can_enroll | missing_prerequisites
-----------+-----------------------
false      | {Data Structures}
```

**Logic:**
1. Finds all prerequisites for CS103
2. Checks if student completed each with required grade
3. Returns FALSE if any are missing
4. Lists titles of missing prerequisites

### Query 3: Prevent Invalid Enrollment (Trigger)
```sql
-- Automatically fires BEFORE INSERT on ENROLLS_IN
-- Example: Try to enroll S1002 in CS103 without CS102

INSERT INTO ENROLLS_IN (student_id, course_id, status) 
VALUES ('S1002', 'CS103', 'Enrolled');

-- ERROR: Cannot enroll: Missing 1 prerequisite(s) for course CS103
```

**Trigger prevents:**
- Enrolling in advanced courses without basics
- Bypassing prerequisite requirements
- Data integrity violations

---

## Testing Procedure

### Step 1: Run Migration
```bash
psql -U your_user -d your_database -f database/06_add_prerequisites.sql
```

### Step 2: Test Eligibility View
```sql
-- Should show S1001 can't enroll in CS103 if CS102 not completed
SELECT * FROM STUDENT_COURSE_ELIGIBILITY 
WHERE student_id = 'S1001' AND course_id = 'CS103';
```

### Step 3: Test Trigger
```sql
-- This should fail if prerequisites aren't met
INSERT INTO ENROLLS_IN (student_id, course_id, status) 
VALUES ('S1001', 'CS103', 'Enrolled');
```

### Step 4: Complete Prerequisite
```sql
-- Mark CS102 as completed
UPDATE ENROLLS_IN 
SET status = 'Completed', grade = 'B+' 
WHERE student_id = 'S1001' AND course_id = 'CS102';

-- Now enrollment should succeed
INSERT INTO ENROLLS_IN (student_id, course_id, status) 
VALUES ('S1001', 'CS103', 'Enrolled');
```

---

## Benefits

1. **Academic Integrity**: Ensures students take courses in proper order
2. **Database Enforcement**: Prevents invalid enrollments at DB level
3. **Better Student Experience**: Clear visibility of what courses are available
4. **Advisor Support**: Faculty can see why students are blocked from courses
5. **Performance**: Pre-computed view makes queries fast

---

## Integration with RAG Chatbot

Update `rag_service.go` to include prerequisite info in student context:

```go
// Fetch courses student can't take yet
prereqQuery := `
    SELECT c.course_id, c.title, 
           ARRAY_TO_STRING(prerequisite_courses, ', ') as missing
    FROM STUDENT_COURSE_ELIGIBILITY
    WHERE student_id = $1 AND is_eligible = FALSE
    LIMIT 3
`
// Add to context:
// "You cannot enroll in CS103 (DBMS) yet. Complete: Data Structures first."
```

This enables the chatbot to give better course selection advice!
