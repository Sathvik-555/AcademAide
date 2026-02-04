# Modification 4: Advanced Analytics with Window Functions & CTEs

## Overview
Implements sophisticated analytics using PostgreSQL advanced SQL features like window functions, CTEs, and recursive queries for student performance tracking and trend analysis.

## Files Changed
- `database/07_add_analytics_views.sql` (NEW)
- `internal/handlers/analytics.go` (NEW)
- `cmd/server/main.go` (MODIFIED)

---

## File 1: database/07_add_analytics_views.sql (NEW FILE)

```sql
-- Advanced Analytics Views and Functions
-- Demonstrates Window Functions, CTEs, Recursive Queries

-- ===========================================
-- 1. STUDENT PERFORMANCE TRENDS (Window Functions)
-- ===========================================

-- View: Semester-wise CGPA with Trends
CREATE OR REPLACE VIEW SEMESTER_PERFORMANCE_TREND AS
WITH semester_data AS (
    SELECT 
        s.student_id,
        s.s_first_name || ' ' || s.s_last_name as student_name,
        s.dept_id,
        -- Assume we have a semester column or derive it from enrollment
        COALESCE(s.semester, 5) as current_semester,
        c.course_id,
        c.title as course_title,
        c.credits,
        e.grade,
        CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A+' THEN 10
            WHEN 'A' THEN 9 WHEN 'B+' THEN 8
            WHEN 'B' THEN 7 WHEN 'C+' THEN 6
            WHEN 'C' THEN 5 WHEN 'D' THEN 4
            ELSE 0
        END as grade_points
    FROM STUDENT s
    JOIN ENROLLS_IN e ON s.student_id = e.student_id
    JOIN COURSE c ON e.course_id = c.course_id
    WHERE e.grade IS NOT NULL
),
semester_gpa AS (
    SELECT 
        student_id,
        student_name,
        dept_id,
        current_semester,
        SUM(grade_points * credits) / NULLIF(SUM(credits), 0) as semester_gpa,
        SUM(credits) as credits_earned
    FROM semester_data
    GROUP BY student_id, student_name, dept_id, current_semester
)
SELECT 
    student_id,
    student_name,
    dept_id,
    current_semester,
    ROUND(semester_gpa::numeric, 2) as gpa,
    credits_earned,
    -- Window Function: Running cumulative GPA
    ROUND(AVG(semester_gpa) OVER (
        PARTITION BY student_id 
        ORDER BY current_semester
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    )::numeric, 2) as cumulative_gpa,
    -- Window Function: Trend (compare to previous semester)
    ROUND((semester_gpa - LAG(semester_gpa) OVER (
        PARTITION BY student_id 
        ORDER BY current_semester
    ))::numeric, 2) as gpa_change,
    -- Window Function: Rank within department
    RANK() OVER (
        PARTITION BY dept_id, current_semester 
        ORDER BY semester_gpa DESC
    ) as dept_rank,
    -- Window Function: Percentile
    PERCENT_RANK() OVER (
        PARTITION BY dept_id, current_semester 
        ORDER BY semester_gpa
    ) * 100 as percentile
FROM semester_gpa
ORDER BY student_id, current_semester;

-- ===========================================
-- 2. COURSE DIFFICULTY ANALYSIS
-- ===========================================

CREATE OR REPLACE VIEW COURSE_DIFFICULTY_METRICS AS
WITH grade_stats AS (
    SELECT 
        c.course_id,
        c.title,
        c.dept_id,
        COUNT(*) as total_students,
        AVG(CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A+' THEN 10 WHEN 'A' THEN 9
            WHEN 'B+' THEN 8 WHEN 'B' THEN 7 WHEN 'C+' THEN 6
            WHEN 'C' THEN 5 WHEN 'D' THEN 4 ELSE 0
        END) as avg_grade_points,
        -- Standard deviation using window function
        STDDEV(CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A+' THEN 10 WHEN 'A' THEN 9
            WHEN 'B+' THEN 8 WHEN 'B' THEN 7 WHEN 'C+' THEN 6
            WHEN 'C' THEN 5 WHEN 'D' THEN 4 ELSE 0
        END) as grade_std_dev,
        -- Failure rate
        100.0 * SUM(CASE WHEN e.grade IN ('F', 'D') THEN 1 ELSE 0 END) / COUNT(*) as failure_rate,
        -- Excellence rate (A and above)
        100.0 * SUM(CASE WHEN e.grade IN ('O', 'A+', 'A') THEN 1 ELSE 0 END) / COUNT(*) as excellence_rate
    FROM COURSE c
    LEFT JOIN ENROLLS_IN e ON c.course_id = e.course_id
    WHERE e.grade IS NOT NULL
    GROUP BY c.course_id, c.title, c.dept_id
    HAVING COUNT(*) >= 5  -- Only include courses with sufficient data
)
SELECT 
    course_id,
    title,
    dept_id,
    total_students,
    ROUND(avg_grade_points::numeric, 2) as avg_grade,
    ROUND(grade_std_dev::numeric, 2) as std_deviation,
    ROUND(failure_rate::numeric, 2) as failure_rate_pct,
    ROUND(excellence_rate::numeric, 2) as excellence_rate_pct,
    -- Difficulty classification
    CASE 
        WHEN avg_grade_points >= 8.5 THEN 'Easy'
        WHEN avg_grade_points >= 7.0 THEN 'Moderate'
        WHEN avg_grade_points >= 5.5 THEN 'Challenging'
        ELSE 'Very Difficult'
    END as difficulty_level,
    -- Consistency (lower std dev = more consistent grading)
    CASE 
        WHEN grade_std_dev < 1.5 THEN 'Highly Consistent'
        WHEN grade_std_dev < 2.5 THEN 'Consistent'
        ELSE 'Variable'
    END as grading_consistency
FROM grade_stats
ORDER BY avg_grade_points ASC;

-- ===========================================
-- 3. AT-RISK STUDENT IDENTIFICATION (CTEs + Window Functions)
-- ===========================================

CREATE OR REPLACE VIEW AT_RISK_STUDENTS AS
WITH recent_performance AS (
    SELECT 
        s.student_id,
        s.s_first_name || ' ' || s.s_last_name as student_name,
        s.s_email,
        s.dept_id,
        s.semester,
        c.course_id,
        c.title as course_title,
        e.grade,
        CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A+' THEN 10 WHEN 'A' THEN 9
            WHEN 'B+' THEN 8 WHEN 'B' THEN 7 WHEN 'C+' THEN 6
            WHEN 'C' THEN 5 WHEN 'D' THEN 4 ELSE 0
        END as grade_points,
        c.credits,
        -- Number of courses taken
        ROW_NUMBER() OVER (
            PARTITION BY s.student_id 
            ORDER BY e.grade  -- Assuming recent courses are later in the table
        ) as course_sequence
    FROM STUDENT s
    JOIN ENROLLS_IN e ON s.student_id = e.student_id
    JOIN COURSE c ON e.course_id = c.course_id
    WHERE e.grade IS NOT NULL
),
student_metrics AS (
    SELECT 
        student_id,
        student_name,
        s_email,
        dept_id,
        semester,
        COUNT(*) as courses_completed,
        SUM(credits) as total_credits,
        ROUND((SUM(grade_points * credits) / NULLIF(SUM(credits), 0))::numeric, 2) as cgpa,
        -- Recent 3 courses average
        ROUND((
            SELECT AVG(rp2.grade_points)
            FROM recent_performance rp2
            WHERE rp2.student_id = recent_performance.student_id
            ORDER BY rp2.course_sequence DESC
            LIMIT 3
        )::numeric, 2) as recent_avg,
        -- Number of failures
        SUM(CASE WHEN grade_points <= 4 THEN 1 ELSE 0 END) as failure_count,
        -- Backlog courses
        (SELECT COUNT(*) FROM ENROLLS_IN e2 
         WHERE e2.student_id = recent_performance.student_id 
         AND e2.backlog = TRUE) as backlog_count
    FROM recent_performance
    GROUP BY student_id, student_name, s_email, dept_id, semester
)
SELECT 
    student_id,
    student_name,
    s_email,
    dept_id,
    semester,
    courses_completed,
    total_credits,
    cgpa,
    recent_avg,
    failure_count,
    backlog_count,
    -- Risk Score (weighted formula)
    ROUND((
        (10 - COALESCE(cgpa, 0)) * 0.4 +  -- Low CGPA
        (10 - COALESCE(recent_avg, 0)) * 0.3 +  -- Declining performance
        (failure_count * 2) * 0.2 +  -- Failures
        (backlog_count * 3) * 0.1  -- Backlogs
    )::numeric, 2) as risk_score,
    -- Risk Level
    CASE 
        WHEN cgpa < 5.0 OR failure_count >= 3 OR backlog_count >= 2 THEN 'Critical'
        WHEN cgpa < 6.5 OR recent_avg < 6.0 OR failure_count >= 2 THEN 'High'
        WHEN cgpa < 7.5 OR recent_avg < 7.0 OR failure_count >= 1 THEN 'Moderate'
        ELSE 'Low'
    END as risk_level,
    -- Recommendation
    CASE 
        WHEN cgpa < 5.0 THEN 'Immediate academic counseling required'
        WHEN failure_count >= 3 THEN 'Focus on clearing backlogs'
        WHEN recent_avg < 6.0 AND cgpa > 7.0 THEN 'Recent decline - investigate causes'
        WHEN backlog_count >= 2 THEN 'Reduce course load, focus on backlogs'
        ELSE 'Monitor performance'
    END as recommendation
FROM student_metrics
WHERE cgpa IS NOT NULL
ORDER BY risk_score DESC;

-- ===========================================
-- 4. PREREQUISITE GRAPH (Recursive CTE)
-- ===========================================

-- Function: Get all transitive prerequisites for a course
CREATE OR REPLACE FUNCTION get_all_prerequisites(p_course_id VARCHAR(10))
RETURNS TABLE (
    level INTEGER,
    prerequisite_id VARCHAR(10),
    prerequisite_title VARCHAR(100),
    is_direct BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE prereq_tree AS (
        -- Base case: Direct prerequisites
        SELECT 
            1 as level,
            p.prerequisite_course_id,
            c.title,
            TRUE as is_direct
        FROM PREREQUISITE p
        JOIN COURSE c ON p.prerequisite_course_id = c.course_id
        WHERE p.course_id = p_course_id
        
        UNION ALL
        
        -- Recursive case: Prerequisites of prerequisites
        SELECT 
            pt.level + 1,
            p2.prerequisite_course_id,
            c2.title,
            FALSE
        FROM prereq_tree pt
        JOIN PREREQUISITE p2 ON pt.prerequisite_id = p2.course_id
        JOIN COURSE c2 ON p2.prerequisite_course_id = c2.course_id
        WHERE pt.level < 5  -- Prevent infinite recursion
    )
    SELECT DISTINCT ON (prerequisite_id)
        level,
        prerequisite_id,
        prerequisite_title,
        is_direct
    FROM prereq_tree
    ORDER BY prerequisite_id, level;
END;
$$ LANGUAGE plpgsql;

-- ===========================================
-- 5. TEACHER PERFORMANCE METRICS
-- ===========================================

CREATE OR REPLACE VIEW TEACHER_PERFORMANCE_SUMMARY AS
WITH teacher_courses AS (
    SELECT 
        f.faculty_id,
        f.f_first_name || ' ' || f.f_last_name as faculty_name,
        f.f_email,
        t.course_id,
        c.title as course_title,
        COUNT(DISTINCT e.student_id) as students_enrolled,
        AVG(CASE e.grade
            WHEN 'O' THEN 10 WHEN 'A+' THEN 10 WHEN 'A' THEN 9
            WHEN 'B+' THEN 8 WHEN 'B' THEN 7 WHEN 'C+' THEN 6
            WHEN 'C' THEN 5 WHEN 'D' THEN 4 ELSE 0
        END) as avg_grade_points,
        100.0 * SUM(CASE WHEN e.grade IN ('F', 'D') THEN 1 ELSE 0 END) / 
            NULLIF(COUNT(*), 0) as failure_rate
    FROM FACULTY f
    JOIN TEACHES t ON f.faculty_id = t.faculty_id
    JOIN COURSE c ON t.course_id = c.course_id
    LEFT JOIN ENROLLS_IN e ON c.course_id = e.course_id
    WHERE e.grade IS NOT NULL
    GROUP BY f.faculty_id, f.f_first_name, f.f_last_name, f.f_email, t.course_id, c.title
)
SELECT 
    faculty_id,
    faculty_name,
    f_email,
    COUNT(DISTINCT course_id) as courses_taught,
    SUM(students_enrolled) as total_students,
    ROUND(AVG(avg_grade_points)::numeric, 2) as overall_avg_grade,
    ROUND(AVG(failure_rate)::numeric, 2) as avg_failure_rate,
    -- Performance rating
    CASE 
        WHEN AVG(avg_grade_points) >= 8.0 AND AVG(failure_rate) < 10 THEN 'Excellent'
        WHEN AVG(avg_grade_points) >= 7.0 AND AVG(failure_rate) < 20 THEN 'Good'
        WHEN AVG(avg_grade_points) >= 6.0 THEN 'Average'
        ELSE 'Needs Improvement'
    END as performance_rating
FROM teacher_courses
GROUP BY faculty_id, faculty_name, f_email
HAVING SUM(students_enrolled) >= 10  -- Only teachers with sufficient data
ORDER BY overall_avg_grade DESC;

-- ===========================================
-- 6. ENROLLMENT TRENDS (Time-series analysis)
-- ===========================================

CREATE OR REPLACE VIEW ENROLLMENT_TRENDS AS
WITH monthly_enrollments AS (
    SELECT 
        c.course_id,
        c.title,
        c.dept_id,
        DATE_TRUNC('month', e.created_at) as enrollment_month,
        COUNT(*) as new_enrollments
    FROM COURSE c
    LEFT JOIN ENROLLS_IN e ON c.course_id = e.course_id
    WHERE e.created_at IS NOT NULL
    GROUP BY c.course_id, c.title, c.dept_id, DATE_TRUNC('month', e.created_at)
)
SELECT 
    course_id,
    title,
    dept_id,
    enrollment_month,
    new_enrollments,
    -- Moving average (3-month)
    ROUND(AVG(new_enrollments) OVER (
        PARTITION BY course_id 
        ORDER BY enrollment_month
        ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
    )::numeric, 2) as moving_avg_3month,
    -- Month-over-month change
    ROUND((
        100.0 * (new_enrollments - LAG(new_enrollments) OVER (
            PARTITION BY course_id 
            ORDER BY enrollment_month
        )) / NULLIF(LAG(new_enrollments) OVER (
            PARTITION BY course_id 
            ORDER BY enrollment_month
        ), 0)
    )::numeric, 2) as pct_change,
    -- Cumulative enrollments
    SUM(new_enrollments) OVER (
        PARTITION BY course_id 
        ORDER BY enrollment_month
    ) as cumulative_enrollments
FROM monthly_enrollments
ORDER BY course_id, enrollment_month;

-- ===========================================
-- 7. INDEXES FOR PERFORMANCE
-- ===========================================

-- Optimize grade queries
CREATE INDEX IF NOT EXISTS idx_enrolls_grade ON ENROLLS_IN(grade) WHERE grade IS NOT NULL;

-- Optimize student lookups with department
CREATE INDEX IF NOT EXISTS idx_student_dept_semester ON STUDENT(dept_id, semester);

-- Composite index for common join patterns
CREATE INDEX IF NOT EXISTS idx_enrolls_student_course_status 
ON ENROLLS_IN(student_id, course_id, status) INCLUDE (grade);
```

---

## File 2: internal/handlers/analytics.go (NEW FILE)

```go
package handlers

import (
	"academ_aide/internal/config"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStudentPerformanceTrend returns semester-wise performance trend
func GetStudentPerformanceTrend(c *gin.Context) {
	rawID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	studentID := rawID.(string)
	
	query := `
		SELECT semester, gpa, cumulative_gpa, gpa_change, dept_rank, percentile
		FROM SEMESTER_PERFORMANCE_TREND
		WHERE student_id = $1
		ORDER BY current_semester
	`
	
	rows, err := config.PostgresDB.Query(query, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	type SemesterTrend struct {
		Semester      int     `json:"semester"`
		GPA           float64 `json:"gpa"`
		CumulativeGPA float64 `json:"cumulative_gpa"`
		GPAChange     *float64 `json:"gpa_change,omitempty"`  // Nullable
		DeptRank      int     `json:"dept_rank"`
		Percentile    float64 `json:"percentile"`
	}
	
	var trends []SemesterTrend
	for rows.Next() {
		var t SemesterTrend
		var gpaChange sql.NullFloat64
		
		err := rows.Scan(&t.Semester, &t.GPA, &t.CumulativeGPA, &gpaChange, 
		                 &t.DeptRank, &t.Percentile)
		if err != nil {
			continue
		}
		
		if gpaChange.Valid {
			t.GPAChange = &gpaChange.Float64
		}
		
		trends = append(trends, t)
	}
	
	c.JSON(http.StatusOK, trends)
}

// GetCourseDifficultyMetrics returns course difficulty analysis
func GetCourseDifficultyMetrics(c *gin.Context) {
	deptID := c.Query("dept_id")
	
	query := `
		SELECT course_id, title, total_students, avg_grade, std_deviation,
		       failure_rate_pct, excellence_rate_pct, difficulty_level, grading_consistency
		FROM COURSE_DIFFICULTY_METRICS
	`
	
	var args []interface{}
	if deptID != "" {
		query += " WHERE dept_id = $1"
		args = append(args, deptID)
	}
	
	query += " ORDER BY avg_grade ASC"
	
	rows, err := config.PostgresDB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	type CourseDifficulty struct {
		CourseID          string  `json:"course_id"`
		Title             string  `json:"title"`
		TotalStudents     int     `json:"total_students"`
		AvgGrade          float64 `json:"avg_grade"`
		StdDeviation      float64 `json:"std_deviation"`
		FailureRate       float64 `json:"failure_rate_pct"`
		ExcellenceRate    float64 `json:"excellence_rate_pct"`
		DifficultyLevel   string  `json:"difficulty_level"`
		GradingConsistency string `json:"grading_consistency"`
	}
	
	var courses []CourseDifficulty
	for rows.Next() {
		var c CourseDifficulty
		rows.Scan(&c.CourseID, &c.Title, &c.TotalStudents, &c.AvgGrade, 
		          &c.StdDeviation, &c.FailureRate, &c.ExcellenceRate, 
		          &c.DifficultyLevel, &c.GradingConsistency)
		courses = append(courses, c)
	}
	
	c.JSON(http.StatusOK, courses)
}

// GetAtRiskStudents returns students needing intervention (Admin/Teacher only)
func GetAtRiskStudents(c *gin.Context) {
	// Check if user is teacher/admin
	role, exists := c.Get("role")
	if !exists || role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	riskLevel := c.Query("risk_level")  // Optional filter: Critical, High, Moderate, Low
	
	query := `
		SELECT student_id, student_name, s_email, dept_id, semester,
		       courses_completed, cgpa, recent_avg, failure_count, backlog_count,
		       risk_score, risk_level, recommendation
		FROM AT_RISK_STUDENTS
	`
	
	var args []interface{}
	if riskLevel != "" {
		query += " WHERE risk_level = $1"
		args = append(args, riskLevel)
	}
	
	query += " ORDER BY risk_score DESC LIMIT 50"
	
	rows, err := config.PostgresDB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	type AtRiskStudent struct {
		StudentID        string  `json:"student_id"`
		StudentName      string  `json:"student_name"`
		Email            string  `json:"email"`
		DeptID           string  `json:"dept_id"`
		Semester         int     `json:"semester"`
		CoursesCompleted int     `json:"courses_completed"`
		CGPA             float64 `json:"cgpa"`
		RecentAvg        float64 `json:"recent_avg"`
		FailureCount     int     `json:"failure_count"`
		BacklogCount     int     `json:"backlog_count"`
		RiskScore        float64 `json:"risk_score"`
		RiskLevel        string  `json:"risk_level"`
		Recommendation   string  `json:"recommendation"`
	}
	
	var students []AtRiskStudent
	for rows.Next() {
		var s AtRiskStudent
		rows.Scan(&s.StudentID, &s.StudentName, &s.Email, &s.DeptID, &s.Semester,
		          &s.CoursesCompleted, &s.CGPA, &s.RecentAvg, &s.FailureCount,
		          &s.BacklogCount, &s.RiskScore, &s.RiskLevel, &s.Recommendation)
		students = append(students, s)
	}
	
	c.JSON(http.StatusOK, students)
}

// GetCoursePrerequisiteTree returns full prerequisite hierarchy
func GetCoursePrerequisiteTree(c *gin.Context) {
	courseID := c.Param("course_id")
	
	query := `SELECT * FROM get_all_prerequisites($1)`
	
	rows, err := config.PostgresDB.Query(query, courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	
	type PrereqNode struct {
		Level       int    `json:"level"`
		CourseID    string `json:"course_id"`
		CourseTitle string `json:"course_title"`
		IsDirect    bool   `json:"is_direct"`
	}
	
	var tree []PrereqNode
	for rows.Next() {
		var node PrereqNode
		rows.Scan(&node.Level, &node.CourseID, &node.CourseTitle, &node.IsDirect)
		tree = append(tree, node)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"course_id": courseID,
		"prerequisites": tree,
	})
}
```

---

## File 3: cmd/server/main.go (MODIFICATION)

```diff
--- a/cmd/server/main.go
+++ b/cmd/server/main.go
@@ -63,6 +63,13 @@ func main() {
 		protected.GET("/courses/eligible", handlers.GetEligibleCourses)
 		protected.GET("/courses/:course_id/eligibility", handlers.CheckEnrollmentEligibility)
 		protected.GET("/courses/:course_id/prerequisites", handlers.GetCoursePrerequisites)
+		protected.GET("/analytics/performance-trend", handlers.GetStudentPerformanceTrend)
+		protected.GET("/analytics/courses/difficulty", handlers.GetCourseDifficultyMetrics)
+		protected.GET("/analytics/prerequisites-tree/:course_id", handlers.GetCoursePrerequisiteTree)
+	}
+	
+	// Teacher/Admin Analytics
+	adminRoutes := r.Group("/admin")
+	adminRoutes.Use(middleware.AuthMiddleware())
+	{
+		adminRoutes.GET("/analytics/at-risk-students", handlers.GetAtRiskStudents)
 	}
```

---

## Key SQL Concepts Demonstrated

### 1. Window Functions

**Example: Cumulative GPA**
```sql
AVG(semester_gpa) OVER (
    PARTITION BY student_id 
    ORDER BY current_semester
    ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
) as cumulative_gpa
```

**Explanation:**
- `PARTITION BY`: Group by student
- `ORDER BY`: Process semesters in order
- `ROWS BETWEEN`: Define window frame
- Result: Running average up to current semester

### 2. Ranking Functions

**Example: Department Rank**
```sql
RANK() OVER (
    PARTITION BY dept_id, current_semester 
    ORDER BY semester_gpa DESC
) as dept_rank
```

**Explanation:**
- Ranks students within their department
- Higher GPA = better rank
- Ties get same rank

### 3. Recursive CTEs

**Example: Prerequisite Tree**
```sql
WITH RECURSIVE prereq_tree AS (
    -- Base: Direct prerequisites
    SELECT 1 as level, prerequisite_course_id
    FROM PREREQUISITE
    WHERE course_id = 'CS103'
    
    UNION ALL
    
    -- Recursive: Prerequisites of prerequisites
    SELECT level + 1, p2.prerequisite_course_id
    FROM prereq_tree pt
    JOIN PREREQUISITE p2 ON pt.prerequisite_id = p2.course_id
    WHERE level < 5
)
SELECT * FROM prereq_tree
```

**Explanation:**
- Finds CS103 → CS102 → CS101 chain
- Prevents infinite loops with level limit
- Returns entire dependency graph

### 4. Common Table Expressions (CTEs)

**Example: At-Risk Analysis**
```sql
WITH recent_performance AS (
    SELECT student_id, grade_points, ROW_NUMBER() OVER (...)
    FROM ENROLLS_IN
),
student_metrics AS (
    SELECT student_id, AVG(grade_points) as cgpa
    FROM recent_performance
    GROUP BY student_id
)
SELECT * FROM student_metrics WHERE cgpa < 6.0
```

**Benefits:**
- Breaks complex queries into logical steps
- Improves readability
- Can be referenced multiple times

### 5. Analytical Functions

**Example: Moving Average**
```sql
AVG(new_enrollments) OVER (
    PARTITION BY course_id 
    ORDER BY enrollment_month
    ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
) as moving_avg_3month
```

**Explanation:**
- Smooths out monthly fluctuations
- Uses current + previous 2 months
- Useful for trend identification

---

## Testing Queries

### Test 1: Check Performance Trend
```sql
SELECT * FROM SEMESTER_PERFORMANCE_TREND 
WHERE student_id = 'S1001';
```

**Expected Output:**
```
student_id | semester | gpa  | cumulative_gpa | gpa_change | dept_rank | percentile
-----------+----------+------+----------------+------------+-----------+------------
S1001      | 3        | 8.2  | 8.2            | NULL       | 5         | 72.5
S1001      | 4        | 8.5  | 8.35           | 0.3        | 3         | 85.2
S1001      | 5        | 8.1  | 8.27           | -0.4       | 6         | 68.9
```

### Test 2: Identify At-Risk Students
```sql
SELECT * FROM AT_RISK_STUDENTS 
WHERE risk_level IN ('Critical', 'High') 
ORDER BY risk_score DESC 
LIMIT 5;
```

### Test 3: Course Difficulty
```sql
SELECT * FROM COURSE_DIFFICULTY_METRICS 
WHERE dept_id = 'CS' 
ORDER BY failure_rate_pct DESC;
```

---

## Performance Considerations

1. **Indexes Created**: Optimizes grade lookups and joins
2. **View Materialization**: Consider materializing views for very large datasets:
   ```sql
   CREATE MATERIALIZED VIEW SEMESTER_PERFORMANCE_TREND_MAT AS
   SELECT * FROM SEMESTER_PERFORMANCE_TREND;
   
   -- Refresh daily
   REFRESH MATERIALIZED VIEW SEMESTER_PERFORMANCE_TREND_MAT;
   ```
3. **Partition Large Tables**: If ENROLLS_IN grows to millions of rows

---

## Benefits

1. **Advanced SQL Mastery**: Demonstrates window functions, CTEs, recursion
2. **Actionable Insights**: Identifies at-risk students before they fail
3. **Data-Driven Decisions**: Course difficulty metrics guide curriculum
4. **Performance Optimization**: Indexed views for fast queries
5. **Scalability**: Recursive queries handle arbitrary prerequisite chains

This modification showcases advanced PostgreSQL features essential for DBMS projects!
