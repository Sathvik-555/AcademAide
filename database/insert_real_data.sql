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
('F101', 'CD252IA', 'CSE-A'), -- Dr. Pratiba teaching DBMS to CSE-A (Sathviks Section)
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
-- Insert Enrollment (Empty - populated in Test Data section below)
-- INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES





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
-- 2.5 INSERT RESOURCES
-- ============================================================
TRUNCATE TABLE RESOURCE RESTART IDENTITY CASCADE;

INSERT INTO RESOURCE (title, description, type, course_id, link) VALUES
-- DBMS Resources
('DBMS SQL Tutorial', 'Comprehensive guide to SQL commands', 'Article', 'CD252IA', 'https://www.geeksforgeeks.org/sql-tutorial/'),
('Normalization Guide', '1NF, 2NF, 3NF and BCNF explained', 'Article', 'CD252IA', 'https://www.geeksforgeeks.org/database-normalization-introduction-normal-forms/'),
('ACID Properties', 'Transaction Management Deep Dive', 'Video', 'CD252IA', 'https://www.youtube.com/watch?v=5W_E44y-Eio'),

-- TOC Resources
('Finite Automata Intro', 'DFA and NFA concepts', 'Article', 'CS354TA', 'https://www.geeksforgeeks.org/introduction-of-finite-automata/'),
('Turing Machines', 'Introduction to Turing Machines', 'Video', 'CS354TA', 'https://www.youtube.com/watch?v=Qa6csfkK7_I'),
('Undecidability', 'Halting Problem and Decidability', 'Article', 'CS354TA', 'https://www.geeksforgeeks.org/decidability-and-undecidability-in-toc/'),

-- AIML Resources
('AI Search Algorithms', 'BFS, DFS, A* Search Guide', 'Article', 'IS353IA', 'https://www.geeksforgeeks.org/search-algorithms-in-ai/'),
('Machine Learning Full Course', 'Complete ML Course for Beginners', 'Video', 'IS353IA', 'https://www.youtube.com/watch?v=GwIo3gDZCVQ'),
('Neural Networks Visualized', 'Deep Learning Concepts', 'Video', 'IS353IA', 'https://www.youtube.com/watch?v=aircAruvnKk'),

-- Cloud Resources
('AWS Cloud Practitioner', 'Official Certification Page', 'Web', 'XX355TBX', 'https://aws.amazon.com/certification/certified-cloud-practitioner/'),
('Virtualization Explained', 'Hypervisors and VMs', 'Video', 'XX355TBX', 'https://www.youtube.com/watch?v=3hTk4cO8hQo'),

-- Management Resources
('Microeconomics 101', 'Supply and Demand Fundamentals', 'Article', 'HS251TA', 'https://www.investopedia.com/terms/m/microeconomics.asp');


-- ============================================================
-- 3. PERSONALIZATION TEST DATA (Dynamic Scenarios)
-- ============================================================

-- Sharanya (Replaces Alice as the "Senior/Topper" persona for testing)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('1RV23CS225', 'Sharanya', 'Narendran', 'sharanyan.cs23@rvce.edu.in', '9980706884', 5, 2023, 'CSE');

-- Enrollments for Sharanya (Matching Sathvik's Schedule + High Grades)
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('1RV23CS225', 'CD252IA', 'A', 'Enrolled'),    -- DBMS (A Grade)
('1RV23CS225', 'CS354TA', 'A', 'Enrolled'),    -- TOC (A Grade)
('1RV23CS225', 'IS353IA', NULL, 'Enrolled'),   -- AIML
('1RV23CS225', 'XX355TBX', NULL, 'Enrolled'),  -- Cloud
('1RV23CS225', 'HS251TA', NULL, 'Enrolled');   -- Management

-- Schedule: Sharanya follows the standard class schedule defined above (no custom inserts needed).


-- Sathvik (Replaces Bob as the "Freshman/Struggling" persona for testing)
-- Joined 2025 (Current Year 2026 - 2025 = 1st Year)
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('1RV23CS221', 'Sathvik', 'Vasudeva', 'sathvikvasudeva.cs23@rvce.edu.in', '7019865562', 1, 2025, 'CSE');

-- Enrollments for Sathvik
INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('1RV23CS221', 'CD252IA', NULL, 'Enrolled'), -- Enrolled in DBMS
('1RV23CS221', 'CS354TA', 'D', 'Completed'); -- Failed/Struggled with Logic previously

-- Schedule for Sathvik (No Class Tonight)
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CD252IA', 'CSE-A', 'Monday', '09:00:00', '10:00:00', 'Morning Hall');


-- Shrinivas (Replaces Charlie)
-- Testing: Enrollments EXCEPT TOC, new timetable
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('1RV23CS234', 'Shrinivas', 'Deshpande', 'shrinivasdeshpande.cs23@rvce.edu.in', '6362744093', 5, 2023, 'CSE');

INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('1RV23CS234', 'CD252IA', NULL, 'Enrolled'),   -- DBMS
('1RV23CS234', 'IS353IA', NULL, 'Enrolled'),   -- AIML
('1RV23CS234', 'XX355TBX', NULL, 'Enrolled'),  -- Cloud
('1RV23CS234', 'HS251TA', NULL, 'Enrolled');   -- Management

-- Schedule for Shrinivas (Custom Timetable)
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
-- Monday
('CD252IA', 'CSE-A', 'Monday', '14:00:00', '15:00:00', 'Lab 1'),
('IS353IA', 'CSE-A', 'Monday', '15:00:00', '16:00:00', 'Lab 1'),
-- Tuesday
('XX355TBX', 'CSE-A', 'Tuesday', '10:00:00', '11:00:00', 'Classroom 202'),
-- Wednesday
('CD252IA', 'CSE-A', 'Wednesday', '09:00:00', '11:00:00', 'Lab 3'), -- DBMS Lab
('HS251TA', 'CSE-A', 'Wednesday', '11:30:00', '12:30:00', 'Classroom 202'),
-- Thursday
('XX355TBX', 'CSE-A', 'Thursday', '09:00:00', '10:00:00', 'Classroom 202'),
('HS251TA', 'CSE-A', 'Thursday', '10:00:00', '11:00:00', 'Classroom 202'),
-- Friday
('IS353IA', 'CSE-A', 'Friday', '11:00:00', '12:00:00', 'Classroom 202'),
('XX355TBX', 'CSE-A', 'Friday', '14:00:00', '15:00:00', 'Classroom 202');


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


-- Sai Abhiram (Replaces Eve) - Assigned D Grade in DBMS
INSERT INTO STUDENT (student_id, s_first_name, s_last_name, s_email, s_phone_no, semester, year_of_joining, dept_id) 
VALUES ('1RV23CS211', 'Sai', 'Abhiram', 'saiabhiram.cs23@rvce.edu.in', '9999999995', 3, 2024, 'CSE');

INSERT INTO ENROLLS_IN (student_id, course_id, grade, status) VALUES
('1RV23CS211', 'CD252IA', 'D', 'Enrolled'); -- Struggling with DBMS (D Grade)

-- Schedule: Class tomorrow, free today
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CD252IA', 'CSE-A', 'Tuesday', '09:00:00', '11:00:00', 'Room 101');
INSERT INTO SCHEDULE (course_id, section_name, day_of_week, start_time, end_time, room_number) VALUES
('CD252IA', 'CSE-A', 'Tuesday', '09:00:00', '10:00:00', 'Morning Hall');
