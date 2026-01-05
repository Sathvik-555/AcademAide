-- ============================================================
-- 1. CLEAN UP (Delete existing mock data)
-- ============================================================
TRUNCATE TABLE ENROLLS_IN, SCHEDULE, SYLLABUS_UNIT, TEACHES, SECTION, COURSE, STUDENT, FACULTY, DEPARTMENT RESTART IDENTITY CASCADE;

-- ============================================================
-- 2. INSERT REAL DATA
-- ============================================================

-- Insert Faculty first (needed for Dept HOD)
INSERT INTO FACULTY (faculty_id, f_first_name, f_last_name, f_email, f_phone_no) VALUES
('F101', 'Dr. Pratiba ', 'D', 'pratibad@rvce.edu.in', '9876543210'),
('F102', 'Dr. Smriti', 'Srivastava', 'smritis@rvce.edu.in', '9876543211'),
('F103', 'Dr. Nagaraja', 'G S', 'nagarajags@rvce.edu.in', '9876543212'),
('F104', 'Dr. Sindhu ', 'D V', 'sindhudv@rvce.edu.in', '9876543213'),
('F105', 'Dr. Hemavathy', 'R.', 'hemavathyr@rvce.edu.in', '9876543214'),
('F106', 'Dr. Vikram N ', 'Bahadurdesai', 'vikramnb@rvce.edu.in', '9876543218'),
('F107', 'Dr. Shanta', 'Rangaswamy', 'shantharangaswamy@rvce.edu.in', '9876543219');

-- Insert Departments
INSERT INTO DEPARTMENT (dept_id, dept_name, hod_id) VALUES
('CSE', 'Computer Science & Engineering', 'F107'),
('IEM', 'Industrial Management', NULL),
('CD', 'Data Science', NULL);

-- Insert Sections
INSERT INTO SECTION (section_name, dept_id) VALUES
('CSE-A', 'CSE'), 
('CSE-B', 'CSE'),
('CSE-C', 'CSE'),
('CSE-D', 'CSE'),
('CSE-E', 'CSE'),
('CD-A', 'CD');


-- Insert Students (Update with your details)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) VALUES
('1RV23CS221', 'Sathvik', 'Vasudeva', 'sathvikvasudeva.cs23@rvce.edu.in', '7019865562', 5, 2023, 'CSE'),
('1RV23CS225', 'Sharanya', 'Narendran', 'sharanyan.cs23@rvce.edu.in', '9980706884', 5, 2023, 'CSE'),
('1RV23CS234', 'Shrinivas', 'Deshpande', 'shrinivasdeshpande.cs23@rvce.edu.in', '6362744093', 5, 2023, 'CSE'),
('1RV23CD121', 'Pranav', 'Kumar', 'pranavkumar.cd23@rvce.edu.in', '6362744094', 5, 2023, 'CD');


-- Insert Courses (5 Courses)
INSERT INTO COURSE (course_id, title, credits, dept_id, description) VALUES 
('CD252IA', 'Database Management Systems', 4, 'CD',
 'Database concepts, ER modeling, relational model, SQL, normalization, transactions, concurrency control, NoSQL and big data basics'),

('CS354TA', 'Theory of Computation', 4, 'CSE',
 'Regular languages, finite automata, context-free grammars, pushdown automata, Turing machines, decidability, computability and complexity'),

('IS353IA', 'Artificial Intelligence and Machine Learning', 4, 'CSE',
 'Intelligent agents, search algorithms, supervised and unsupervised learning, decision trees, Naive Bayes, logistic regression, ensemble methods'),

('XX355TBX', 'Cloud Computing', 3, 'CSE',
 'Cloud concepts, service models (IaaS, PaaS, SaaS), cloud architecture, virtualization, storage, capacity planning and cloud application development'),

('HS251TA', 'Principles of Economics and Management', 3, 'IEM',
 'Management principles, planning and organization, leadership and motivation, microeconomics, macroeconomics, markets, GDP and economic indicators');

-- Insert Teaches
INSERT INTO TEACHES (faculty_id, course_id, section_name) VALUES
('F102', 'CS354TA', 'CSE-D'),
('F103', 'XX355TBX', 'CSE-D'),
('F106', 'HS251TA', 'CSE-D'),
('F105', 'CD252IA', 'CSE-D'),
('F104', 'IS353IA', 'CSE-D');

-- Insert Syllabus (5 Units per Course = 25 Total)
INSERT INTO SYLLABUS_UNIT (unit_no, topic, course_id) VALUES
-- Course 1: Database Management Systems (CD252IA)
(1, 'Introduction to Database Systems & ER Modeling', 'CD252IA'),
(2, 'Relational Model and Relational Algebra', 'CD252IA'),
(3, 'SQL and Relational Database Design', 'CD252IA'),
(4, 'Transaction Processing and Concurrency Control', 'CD252IA'),
(5, 'NoSQL Databases and Big Data Concepts', 'CD252IA'),

-- Course 2: Theory of Computation (CS354TA)
(1, 'Regular Languages and Finite Automata', 'CS354TA'),
(2, 'Context Free Grammars and Normal Forms', 'CS354TA'),
(3, 'Pushdown Automata and CFL Properties', 'CS354TA'),
(4, 'Turing Machines and Decidability', 'CS354TA'),
(5, 'Recursive Languages, Chomsky Hierarchy and Complexity', 'CS354TA'),

-- Course 3: Artificial Intelligence and Machine Learning (IS353IA)
(1, 'Intelligent Agents and Uninformed Search', 'IS353IA'),
(2, 'Heuristic Search and Adversarial Search', 'IS353IA'),
(3, 'Supervised Learning and Decision Trees', 'IS353IA'),
(4, 'Nearest Neighbor, Naive Bayes and Logistic Regression', 'IS353IA'),
(5, 'Unsupervised Learning, Clustering and Evaluation', 'IS353IA'),

-- Course 4: Cloud Computing (XX355TBX)
(1, 'Cloud Computing Concepts and Service Models', 'XX355TBX'),
(2, 'Cloud Architecture and Service Oriented Architecture', 'XX355TBX'),
(3, 'Cloud Infrastructure, Storage and Standards', 'XX355TBX'),
(4, 'Virtualization, Hypervisors and Capacity Planning', 'XX355TBX'),
(5, 'Cloud Application Development and Management', 'XX355TBX'),

-- Course 5: Principles of Management & Economics (HS251TA)
(1, 'Introduction to Management and Management Theories', 'HS251TA'),
(2, 'Planning, Organization Structure and Strategy', 'HS251TA'),
(3, 'Motivation and Leadership Theories', 'HS251TA'),
(4, 'Microeconomics and Market Analysis', 'HS251TA'),
(5, 'Macroeconomics, GDP and Economic Indicators', 'HS251TA');


-- Insert Enrollment
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES

-- Student 1
('1RV23CS221', 'CD252IA', NULL, 'Enrolled'),  -- DBMS
('1RV23CS221', 'CS354TA', NULL, 'Enrolled'),  -- TOC
('1RV23CS221', 'IS353IA', NULL, 'Enrolled'),  -- AIML
('1RV23CS221', 'XX355TBX', NULL, 'Enrolled'), -- Cloud Computing
('1RV23CS221', 'HS251TA', NULL, 'Enrolled'),  -- Management & Economics

-- Student 2
('1RV23CS225', 'CD252IA', NULL, 'Enrolled'),
('1RV23CS225', 'CS354TA', NULL, 'Enrolled'),
('1RV23CS225', 'IS353IA', NULL, 'Enrolled'),
('1RV23CS225', 'XX355TBX', NULL, 'Enrolled'),
('1RV23CS225', 'HS251TA', NULL, 'Enrolled'),

-- Student 3
('1RV23CS234', 'CD252IA', NULL, 'Enrolled'),
('1RV23CS234', 'CS354TA', NULL, 'Enrolled'),
('1RV23CS234', 'IS353IA', NULL, 'Enrolled'),
('1RV23CS234', 'XX355TBX', NULL, 'Enrolled'),
('1RV23CS234', 'HS251TA', NULL, 'Enrolled');



-- Insert Schedule 
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES

-- =====================
-- MONDAY
-- =====================
('IS353IA', 'CSE-D', 'Monday', '09:00:00', '10:00:00', 'Chemical 107'), -- AIML
('HS251TA', 'CSE-D', 'Monday', '10:00:00', '11:00:00', 'Chemical 107'), -- POME
('CS354TA', 'CSE-D', 'Monday', '11:30:00', '12:30:00', 'Chemical 107'), -- TOC (T)
('CD252IA', 'CSE-D', 'Monday', '12:30:00', '13:30:00', 'Chemical 107'), -- DBMS

-- =====================
-- TUESDAY
-- =====================
('CD252IA', 'CSE-D', 'Tuesday', '09:00:00', '10:00:00', 'Chemical 107'), -- DBMS
('XX355TBX', 'CSE-D', 'Tuesday', '10:00:00', '11:00:00', 'Chemical 107'), -- Elective-I (Cloud)
('IS353IA', 'CSE-D', 'Tuesday', '11:30:00', '12:30:00', 'Chemical 107'), -- AIML
('CD252IA', 'CSE-D', 'Tuesday', '14:30:00', '16:30:00', 'Lab 3'),       -- DBMS Lab (D2)
('IS353IA', 'CSE-D', 'Tuesday', '14:30:00', '16:30:00', 'Lab 7'),       -- AIML Lab (D1)

-- =====================
-- WEDNESDAY
-- =====================
('CS354TA', 'CSE-D', 'Wednesday', '09:00:00', '10:00:00', 'Chemical 107'), -- TOC
('CD252IA', 'CSE-D', 'Wednesday', '10:00:00', '11:00:00', 'Chemical 107'), -- DBMS
('HS251TA', 'CSE-D', 'Wednesday', '11:30:00', '12:30:00', 'Chemical 107'), -- POME
('HS251TA', 'CSE-D', 'Wednesday', '12:30:00', '13:30:00', 'Chemical 107'), -- POME

-- =====================
-- THURSDAY
-- =====================
('CD252IA', 'CSE-D', 'Thursday', '09:00:00', '11:00:00', 'Lab 3'), -- DBMS Lab (D1)
('IS353IA', 'CSE-D', 'Thursday', '09:00:00', '11:00:00', 'Lab 7'), -- AIML Lab (D2)
('XX355TBX', 'CSE-D', 'Thursday', '11:30:00', '12:30:00', 'Chemical 107'), -- Elective-I
('CS354TA', 'CSE-D', 'Thursday', '12:30:00', '13:30:00', 'Chemical 107'), -- TOC

-- =====================
-- FRIDAY
-- =====================
('XX355TBX', 'CSE-D', 'Friday', '09:00:00', '10:00:00', 'Chemical 107'), -- Elective-I
('CS354TA', 'CSE-D', 'Friday', '10:00:00', '11:00:00', 'Chemical 107'), -- TOC
('IS353IA', 'CSE-D', 'Friday', '11:30:00', '12:30:00', 'Chemical 107'), -- AIML
('CS354TA', 'CSE-D', 'Friday', '12:30:00', '13:30:00', 'Chemical 107'); -- TOC (T)


-- ============================================================
-- 3. PERSONALIZATION TEST DATA (Dynamic Scenarios)
-- ============================================================

-- Alice (The "Senior" - 4th Year, High Grades, Advanced Queries)
-- Joined 2022 (Current Year 2026 - 2022 = 4th Year)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('TEST_SENIOR', 'Alice', 'Senior', 'alice@rvce.edu.in', '9999999991', 7, 2022, 'CSE');

-- Enrollments for Alice (Advanced Courses & Prerequisites Passed)
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('TEST_SENIOR', 'XX355TBX', NULL, 'Enrolled'), -- Start learning Cloud
('TEST_SENIOR', 'CS354TA', 'A', 'Completed'), -- Aced Logic/TOC
('TEST_SENIOR', 'CD252IA', 'A', 'Completed'); -- Aced DBMS

-- Schedule for Alice (Class starting SOON relative to "Monday 22:50" reference)
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('XX355TBX', 'CSE-D', 'Monday', '23:00:00', '23:59:00', 'Late Night Lab');


-- Bob (The "Freshman" - 1st Year, Low Grades/Struggling)
-- Joined 2025 (Current Year 2026 - 2025 = 1st Year)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('TEST_FRESHMAN', 'Bob', 'Freshman', 'bob@rvce.edu.in', '9999999992', 1, 2025, 'CSE');

-- Enrollments for Bob (Intro Courses)
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('TEST_FRESHMAN', 'CD252IA', NULL, 'Enrolled'), -- Enrolled in DBMS
('TEST_FRESHMAN', 'CS354TA', 'D', 'Completed'); -- Failed/Struggled with Logic previously

-- Schedule for Bob (No Class Tonight)
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CD252IA', 'CSE-A', 'Monday', '09:00:00', '10:00:00', 'Morning Hall');


-- Charlie (The "Other" - Data Science Major)
-- To test RAG Filtering (Should ONLY see Data Science stuff)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('TEST_DS', 'Charlie', 'DataSci', 'charlie@rvce.edu.in', '9999999993', 3, 2024, 'CD');

INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('TEST_DS', 'CD252IA', NULL, 'Enrolled'); -- Only DBMS


-- Dave (The "Stressed" - 3rd Year, Poor Grades, Heavy Schedule)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('TEST_STRESSED', 'Dave', 'Stressed', 'dave@rvce.edu.in', '9999999994', 5, 2023, 'CSE');

INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('TEST_STRESSED', 'CS354TA', 'D', 'Completed'), -- Barely passed TOC
('TEST_STRESSED', 'IS353IA', 'C', 'Enrolled'); -- Struggling with AI

-- Schedule: Back to back classes on Monday evening
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('IS353IA', 'CSE-D', 'Monday', '20:00:00', '21:00:00', 'Night Class 1'),
('CS354TA', 'CSE-D', 'Monday', '21:00:00', '22:00:00', 'Night Class 2');


-- Eve (The "Topper" - 2nd Year, O Grade, Ahead of curve)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('TEST_TOPPER', 'Eve', 'Topper', 'eve@rvce.edu.in', '9999999995', 3, 2024, 'CSE');

INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('TEST_TOPPER', 'CD252IA', 'O', 'Completed'), -- Perfect score in DBMS
('TEST_TOPPER', 'CS354TA', 'A+', 'Completed'); -- Near perfect in TOC

-- Schedule: Class tomorrow, free today
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CD252IA', 'CSE-A', 'Tuesday', '09:00:00', '10:00:00', 'Morning Hall');
