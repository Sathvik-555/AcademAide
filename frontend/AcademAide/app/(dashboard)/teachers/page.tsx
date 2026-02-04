"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { Mail, BookOpen, User } from "lucide-react"
import { useEffect, useState } from "react"
import Cookies from "js-cookie"

interface Teacher {
    name: string
    email: string
    course: string
    course_id: string
}

export default function TeachersPage() {
    const [teachers, setTeachers] = useState<Teacher[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchTeachers = async () => {
            const token = Cookies.get("token")
            if (!token) return

            try {
                const res = await fetch("http://localhost:8080/student/teachers", {
                    headers: {
                        "Authorization": `Bearer ${token}`
                    }
                })
                if (res.ok) {
                    const data = await res.json()
                    setTeachers(data)
                }
            } catch (err) {
                console.error("Failed to fetch teachers", err)
            } finally {
                setLoading(false)
            }
        }

        fetchTeachers()
    }, [])

    return (
        <div className="flex flex-col gap-6">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight">My Teachers</h1>
                <p className="text-muted-foreground">
                    Connect with the faculty members for your enrolled courses.
                </p>
            </div>

            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {[1, 2, 3].map((i) => (
                        <Card key={i} className="animate-pulse">
                            <CardHeader className="h-24 bg-muted/50" />
                            <CardContent className="h-32 bg-muted/30" />
                        </Card>
                    ))}
                </div>
            ) : teachers.length === 0 ? (
                <Card className="p-8 text-center text-muted-foreground">
                    No teachers found for your enrolled courses.
                </Card>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {teachers.map((teacher, index) => (
                        <Card key={index} className="overflow-hidden hover:shadow-lg transition-all duration-200 border-l-4 border-l-primary/50">
                            <CardHeader className="bg-muted/30 pb-4">
                                <div className="flex items-center gap-4">
                                    <div className="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                                        <User className="h-6 w-6" />
                                    </div>
                                    <div>
                                        <CardTitle className="text-lg">{teacher.name}</CardTitle>
                                        <a
                                            href={`mailto:${teacher.email}`}
                                            className="text-sm text-muted-foreground flex items-center gap-1 hover:text-primary transition-colors"
                                        >
                                            <Mail className="h-3 w-3" />
                                            {teacher.email}
                                        </a>
                                    </div>
                                </div>
                            </CardHeader>
                            <CardContent className="pt-6">
                                <div className="space-y-4">
                                    <div className="flex items-start gap-3">
                                        <div className="mt-1 h-8 w-8 rounded-lg bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center text-indigo-600 dark:text-indigo-400 shrink-0">
                                            <BookOpen className="h-4 w-4" />
                                        </div>
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">Course</p>
                                            <p className="font-semibold text-foreground">{teacher.course}</p>
                                            <p className="text-xs text-muted-foreground font-mono mt-1 px-1.5 py-0.5 bg-muted rounded inline-block">
                                                {teacher.course_id}
                                            </p>
                                        </div>
                                    </div>

                                    <div className="pt-2 flex justify-end">
                                        <a
                                            href={`https://mail.google.com/mail/?view=cm&fs=1&to=${teacher.email}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-secondary text-secondary-foreground hover:bg-secondary/80 h-9 px-4 py-2"
                                        >
                                            Contact Faculty
                                        </a>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
