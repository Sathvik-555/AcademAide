"use client"

import { useState, useEffect } from "react"
import Cookies from "js-cookie"
import { Users, Search, Mail } from "lucide-react"

interface Course {
    course_id: string
    title: string
    section: string
}

interface Student {
    student_id: string
    name: string
    email: string
}

export default function StudentsPage() {
    const [courses, setCourses] = useState<Course[]>([])
    const [selectedCourse, setSelectedCourse] = useState<string>("")
    const [students, setStudents] = useState<Student[]>([])
    const [loading, setLoading] = useState(true)
    const [loadingStudents, setLoadingStudents] = useState(false)

    useEffect(() => {
        const fetchCourses = async () => {
            try {
                const token = Cookies.get("token")
                const res = await fetch("http://localhost:8080/teacher/courses", {
                    headers: { "Authorization": `Bearer ${token}` }
                })
                if (res.ok) {
                    const data = await res.json()
                    setCourses(data || [])
                    if (data && data.length > 0) {
                        setSelectedCourse(data[0].course_id)
                    }
                }
            } catch (error) {
                console.error("Failed to fetch courses", error)
            } finally {
                setLoading(false)
            }
        }
        fetchCourses()
    }, [])

    useEffect(() => {
        if (!selectedCourse) return

        const fetchStudents = async () => {
            setLoadingStudents(true)
            try {
                const token = Cookies.get("token")
                const res = await fetch(`http://localhost:8080/teacher/students?course_id=${selectedCourse}`, {
                    headers: { "Authorization": `Bearer ${token}` }
                })
                if (res.ok) {
                    const data = await res.json()
                    setStudents(data || [])
                }
            } catch (error) {
                console.error("Failed to fetch students", error)
            } finally {
                setLoadingStudents(false)
            }
        }
        fetchStudents()
    }, [selectedCourse])

    if (loading) return <div>Loading...</div>

    return (
        <div className="max-w-6xl mx-auto p-6 space-y-6">
            <h1 className="text-3xl font-bold flex items-center gap-2">
                <Users className="h-8 w-8 text-primary" />
                Enrolled Students
            </h1>

            {/* Course Filter */}
            <div className="flex items-center gap-4 bg-muted/50 p-4 rounded-lg">
                <label className="text-sm font-medium whitespace-nowrap">Select Course:</label>
                <select
                    value={selectedCourse}
                    onChange={(e) => setSelectedCourse(e.target.value)}
                    className="flex h-10 w-full md:w-[300px] rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                >
                    {courses.map(c => (
                        <option key={c.course_id} value={c.course_id}>
                            {c.course_id} - {c.title} ({c.section})
                        </option>
                    ))}
                </select>
            </div>

            {/* Students List */}
            <div className="border rounded-md">
                <div className="relative w-full overflow-auto">
                    <table className="w-full caption-bottom text-sm">
                        <thead className="[&_tr]:border-b">
                            <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">ID</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Name</th>
                                <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Email</th>
                                <th className="h-12 px-4 text-right align-middle font-medium text-muted-foreground">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="[&_tr:last-child]:border-0">
                            {loadingStudents ? (
                                <tr>
                                    <td colSpan={4} className="p-4 text-center">Loading students...</td>
                                </tr>
                            ) : students.length > 0 ? (
                                students.map((student) => (
                                    <tr key={student.student_id} className="border-b transition-colors hover:bg-muted/50">
                                        <td className="p-4 align-middle">{student.student_id}</td>
                                        <td className="p-4 align-middle font-medium">{student.name}</td>
                                        <td className="p-4 align-middle text-muted-foreground">{student.email}</td>
                                        <td className="p-4 align-middle text-right">
                                            <a
                                                href={`https://mail.google.com/mail/?view=cm&fs=1&to=${student.email}`}
                                                target="_blank"
                                                rel="noreferrer"
                                                className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 hover:bg-accent hover:text-accent-foreground h-9 w-9"
                                                title="Send email via Gmail"
                                            >
                                                <Mail className="h-4 w-4" />
                                            </a>
                                        </td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan={4} className="p-4 text-center text-muted-foreground">No students enrolled in this course.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    )
}
