-- AcademAide PostgreSQL Schema
-- Enable pgvector extension for embeddings
CREATE EXTENSION IF NOT EXISTS vector;

-- 3NF Justification:
-- 1NF: All columns are atomic. No repeating groups.
-- 2NF: No partial dependencies. All non-key attributes depend on the full primary key.
-- 3NF: No transitive dependencies. Non-key attributes depend ONLY on the primary key, not on other non-key attributes.
--      Example: 'dept_name' is in DEPARTMENT, not repeated in STUDENT or COURSE.

-- Enable UUID extension if needed, though schema uses standard types (assuming serial or text/int)
-- Using SERIAL/INTEGER for IDs for simplicity as per common academic projects, or VARCHAR for specific ID formats.
-- Based on typical academic IDs, VARCHAR is safer for 'student_id' (e.g., 'S2023001').

-- Corrected Drop Order to handle FK dependencies
-- Child tables must be dropped before Parents


DROP TABLE IF EXISTS SCHEDULE;
DROP TABLE IF EXISTS ENROLLS_IN;
DROP TABLE IF EXISTS RESOURCE;
DROP TABLE IF EXISTS SYLLABUS_UNIT;
DROP TABLE IF EXISTS TEACHES;

-- Dependents of Department must go first
DROP TABLE IF EXISTS STUDENT;
DROP TABLE IF EXISTS COURSE;
DROP TABLE IF EXISTS SECTION;

-- DEPARTMENT depends on FACULTY (via hod_id), so drop Dept first OR use CASCADE
DROP TABLE IF EXISTS DEPARTMENT CASCADE; 
DROP TABLE IF EXISTS FACULTY CASCADE;

-- 1. DEPARTMENT
CREATE TABLE DEPARTMENT (
    dept_id VARCHAR(10) PRIMARY KEY,
    dept_name VARCHAR(100) NOT NULL,
    hod_id VARCHAR(20) -- FK added later to avoid circular dependency or defer
);

-- 2. FACULTY
CREATE TABLE FACULTY (
    faculty_id VARCHAR(20) PRIMARY KEY,
    f_first_name VARCHAR(50) NOT NULL,
    f_last_name VARCHAR(50) NOT NULL,
    f_email VARCHAR(100) UNIQUE NOT NULL,
    f_phone_no VARCHAR(15)
);

-- Add Circular FK for Department HOD
ALTER TABLE DEPARTMENT
ADD CONSTRAINT fk_dept_hod
FOREIGN KEY (hod_id) REFERENCES FACULTY(faculty_id);

-- 3. STUDENT
CREATE TABLE STUDENT (
    student_id VARCHAR(20) PRIMARY KEY,
    s_first_name VARCHAR(50) NOT NULL,
    s_last_name VARCHAR(50) NOT NULL,
    s_email VARCHAR(100) UNIQUE NOT NULL,
    s_phone_no VARCHAR(15),
    semester INTEGER NOT NULL,
    year_of_joining INTEGER NOT NULL,
    dept_id VARCHAR(10) NOT NULL,
    CONSTRAINT fk_student_dept FOREIGN KEY (dept_id) REFERENCES DEPARTMENT(dept_id)
);

-- 4. SECTION
CREATE TABLE SECTION (
    section_name VARCHAR(10) PRIMARY KEY, -- e.g., 'A', 'B', 'CS-A'
    dept_id VARCHAR(10) NOT NULL,
    CONSTRAINT fk_section_dept FOREIGN KEY (dept_id) REFERENCES DEPARTMENT(dept_id)
);

-- 5. COURSE
CREATE TABLE COURSE (
    course_id VARCHAR(10) PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    credits INTEGER NOT NULL CHECK (credits > 0),
    dept_id VARCHAR(10) NOT NULL,
    description TEXT,
    embedding vector(768),
    CONSTRAINT fk_course_dept FOREIGN KEY (dept_id) REFERENCES DEPARTMENT(dept_id)
);

-- 6. TEACHES
CREATE TABLE TEACHES (
    faculty_id VARCHAR(20),
    course_id VARCHAR(10),
    section_name VARCHAR(10),
    PRIMARY KEY (faculty_id, course_id, section_name),
    CONSTRAINT fk_teaches_faculty FOREIGN KEY (faculty_id) REFERENCES FACULTY(faculty_id),
    CONSTRAINT fk_teaches_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id),
    CONSTRAINT fk_teaches_section FOREIGN KEY (section_name) REFERENCES SECTION(section_name)
);

-- 7. SYLLABUS_UNIT
CREATE TABLE SYLLABUS_UNIT (
    unit_id SERIAL PRIMARY KEY,
    unit_no INTEGER NOT NULL,
    topic VARCHAR(255) NOT NULL,
    course_id VARCHAR(10) NOT NULL,
    CONSTRAINT fk_syllabus_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);

-- 8. RESOURCE
CREATE TABLE RESOURCE (
    resource_id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(50), -- e.g., 'PDF', 'Video', 'Slide'
    course_id VARCHAR(10) NOT NULL,
    CONSTRAINT fk_resource_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);

-- 9. ENROLLS_IN
CREATE TABLE ENROLLS_IN (
    student_id VARCHAR(20),
    course_id VARCHAR(10),
    grade VARCHAR(2), -- e.g., 'A', 'B+'
    status VARCHAR(20) DEFAULT 'Enrolled', -- 'Completed', 'Dropped'
    backlog BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (student_id, course_id),
    CONSTRAINT fk_enrolls_student FOREIGN KEY (student_id) REFERENCES STUDENT(student_id),
    CONSTRAINT fk_enrolls_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id)
);

-- 10. SCHEDULE
CREATE TABLE SCHEDULE (
    schedule_id SERIAL PRIMARY KEY,
    course_id VARCHAR(10) NOT NULL,
    section_name VARCHAR(10) NOT NULL,
    day_of_week VARCHAR(10) NOT NULL CHECK (day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday')),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    room_number VARCHAR(20),
    CONSTRAINT fk_schedule_course FOREIGN KEY (course_id) REFERENCES COURSE(course_id),
    CONSTRAINT fk_schedule_section FOREIGN KEY (section_name) REFERENCES SECTION(section_name)
);



-- ==========================================
-- SEED DATA
-- ==========================================

-- Insert Faculty first (needed for Dept HOD)
INSERT INTO FACULTY (faculty_id, f_first_name, f_last_name, f_email, f_phone_no) VALUES
('F001', 'Alice', 'Smith', 'alice.smith@univ.edu', '555-0101'),
('F002', 'Bob', 'Jones', 'bob.jones@univ.edu', '555-0102');

-- Insert Departments
INSERT INTO DEPARTMENT (dept_id, dept_name, hod_id) VALUES
('CS', 'Computer Science', 'F001'),
('EC', 'Electronics', 'F002');

-- Insert Sections
INSERT INTO SECTION (section_name, dept_id) VALUES
('CS-A', 'CS'),
('CS-B', 'CS'),
('EC-A', 'EC');

-- Insert Students
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) VALUES
('S1001', 'John', 'Doe', 'john.doe@student.univ.edu', '555-1001', 5, 2023, 'CS'),
('S1002', 'Jane', 'Roe', 'jane.roe@student.univ.edu', '555-1002', 3, 2024, 'CS');

-- Insert Courses
INSERT INTO COURSE (course_id, title, credits, dept_id, description) VALUES
('CS101', 'Intro to Programming', 4, 'CS', 'Fundamental concepts of programming using C++. Covers loops, logic, functions, and basic data structures.'),
('CS102', 'Data Structures', 4, 'CS', 'Advanced storage and retrieval of data. Topics include arrays, linked lists, trees, graphs, and sorting algorithms.'),
('CS103', 'Database Systems', 3, 'CS', 'Design and implementation of relational databases. SQL, normalization, indexing, and transaction management.');

-- Insert Teaches
INSERT INTO TEACHES (faculty_id, course_id, section_name) VALUES
('F001', 'CS101', 'CS-A'),
('F002', 'CS102', 'CS-B');

-- Insert Syllabus
INSERT INTO SYLLABUS_UNIT (unit_no, topic, course_id) VALUES
(1, 'Introduction to C++', 'CS101'),
(2, 'Loops and Logic', 'CS101'),
(1, 'Arrays and Linked Lists', 'CS102');

-- Insert Resources
INSERT INTO RESOURCE (title, description, type, course_id) VALUES
('C++ Basics', 'Introductory slides', 'PDF', 'CS101'),
('Sorting Algos', 'Video lecture on sorting', 'Video', 'CS102');

-- Insert Enrollment
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('S1001', 'CS101', 'A', 'Completed'),
('S1001', 'CS102', NULL, 'Enrolled'), -- Ongoing
('S1002', 'CS101', NULL, 'Enrolled');

-- Insert Schedule
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CS101', 'CS-A', 'Monday', '09:00:00', '10:00:00', 'Room-101'),
('CS102', 'CS-B', 'Tuesday', '11:00:00', '12:30:00', 'Room-102');
