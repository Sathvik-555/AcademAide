"use client"

import { useState, useEffect } from "react"
import Cookies from "js-cookie"
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/Card"
import { BookOpen, Users } from "lucide-react"

interface Course {
    course_id: string
    title: string
    section: string
}

export default function MyCoursesPage() {
    const [courses, setCourses] = useState<Course[]>([])
    const [loading, setLoading] = useState(true)

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
                }
            } catch (error) {
                console.error("Failed to fetch courses", error)
            } finally {
                setLoading(false)
            }
        }
        fetchCourses()
    }, [])

    if (loading) return <div>Loading courses...</div>

    return (
        <div className="max-w-6xl mx-auto p-6 space-y-6">
            <h1 className="text-3xl font-bold flex items-center gap-2">
                <BookOpen className="h-8 w-8 text-primary" />
                My Courses
            </h1>
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {courses.length > 0 ? (
                    courses.map((course) => (
                        <Card key={course.course_id + course.section} className="hover:shadow-lg transition-all">
                            <CardHeader>
                                <div className="flex justify-between items-start">
                                    <div className="bg-primary/10 text-primary px-2 py-1 rounded text-xs font-bold">
                                        {course.course_id}
                                    </div>
                                    <div className="bg-secondary/20 text-secondary-foreground px-2 py-1 rounded text-xs font-semibold">
                                        Section {course.section}
                                    </div>
                                </div>
                                <CardTitle className="mt-2 line-clamp-1">{course.title}</CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-sm text-muted-foreground flex items-center gap-2">
                                    <Users className="h-4 w-4" />
                                    <span>View Enrolled Students</span>
                                </div>
                            </CardContent>
                        </Card>
                    ))
                ) : (
                    <div className="col-span-full text-center text-muted-foreground">
                        No courses found assigned to you.
                    </div>
                )}
            </div>
        </div>
    )
}
