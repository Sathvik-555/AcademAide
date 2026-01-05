"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { Users, AlertTriangle, TrendingUp, BookOpen, LogOut } from "lucide-react"
import Cookies from "js-cookie"
import { useRouter } from "next/navigation"

interface ClassHealth {
    course_id: string
    title: string
    attendance_distribution: Record<string, number>
    performance_heatmap: Record<string, number>
}

interface AtRiskStats {
    high_risk_count: number
    medium_risk_count: number
    low_risk_count: number
}

export default function TeacherDashboard() {
    const router = useRouter()
    const [course, setCourse] = useState<{ course_id: string; title: string } | null>(null)
    const [students, setStudents] = useState<any[]>([])

    const [health, setHealth] = useState<ClassHealth | null>(null)
    const [atRisk, setAtRisk] = useState<AtRiskStats | null>(null)
    const [alerts, setAlerts] = useState<string[]>([])
    const [loading, setLoading] = useState(true)

    const [selectedStudent, setSelectedStudent] = useState<any | null>(null)
    const [detailLoading, setDetailLoading] = useState(false)
    const [announcement, setAnnouncement] = useState("")

    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get("token")
            if (!token) {
                router.push("/login")
                return
            }
            const headers = { "Authorization": `Bearer ${token}` }

            try {
                // 1. Fetch Teacher's Courses
                const courseRes = await fetch("http://localhost:8080/teacher/courses", { headers })
                if (!courseRes.ok) return
                const courses = await courseRes.json()

                if (courses && courses.length > 0) {
                    const firstCourse = courses[0]
                    setCourse(firstCourse)
                    const courseId = firstCourse.course_id

                    // 2. Fetch Health & Risk
                    const healthRes = await fetch(`http://localhost:8080/teacher/class-health?course_id=${courseId}`, { headers })
                    if (healthRes.ok) setHealth(await healthRes.json())

                    const riskRes = await fetch(`http://localhost:8080/teacher/at-risk?course_id=${courseId}`, { headers })
                    if (riskRes.ok) setAtRisk(await riskRes.json())

                    // 3. Fetch Students List
                    const studRes = await fetch(`http://localhost:8080/teacher/students?course_id=${courseId}`, { headers })
                    if (studRes.ok) setStudents(await studRes.json())
                }

                // Fetch Alerts
                const alertsRes = await fetch(`http://localhost:8080/teacher/alerts`, { headers })
                if (alertsRes.ok) {
                    const data = await alertsRes.json()
                    setAlerts(data.alerts || [])
                }
            } catch (err) {
                console.error("Failed to fetch teacher data", err)
            } finally {
                setLoading(false)
            }
        }
        fetchData()
    }, [router])

    const handleLogout = () => {
        Cookies.remove("token")
        Cookies.remove("user_id")
        Cookies.remove("role")
        Cookies.remove("faculty_id")
        router.push("/login")
    }

    const fetchStudentDetails = async (studentId: string) => {
        if (!course) return
        setDetailLoading(true)
        setSelectedStudent(null) // Reset
        try {
            const token = Cookies.get("token")
            const res = await fetch(`http://localhost:8080/teacher/student-details?student_id=${studentId}&course_id=${course.course_id}`, {
                headers: { "Authorization": `Bearer ${token}` }
            })
            if (res.ok) {
                const data = await res.json()
                setSelectedStudent(data)
            }
        } catch (e) {
            console.error("Error details", e)
        } finally {
            setDetailLoading(false)
        }
    }

    const postAnnouncement = async () => {
        if (!course || !announcement) return
        try {
            const token = Cookies.get("token")
            const res = await fetch(`http://localhost:8080/teacher/announce`, {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ course_id: course.course_id, content: announcement })
            })
            if (res.ok) {
                alert("Announcement Posted!")
                setAnnouncement("")
            }
        } catch (e) {
            console.error(e)
        }
    }

    if (loading) return <div className="flex h-screen items-center justify-center">Loading dashboard...</div>

    return (
        <div className="flex min-h-screen flex-col p-8 gap-8 bg-gray-50 dark:bg-slate-950">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Teacher Dashboard</h1>
                    <p className="text-muted-foreground">
                        {course ? `Overview for ${course.title} (${course.course_id})` : "No classes assigned."}
                    </p>
                </div>
                <button onClick={handleLogout} className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-md transition-colors">
                    <LogOut className="h-4 w-4" />
                    Logout
                </button>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                {/* At Risk Summary Card */}
                <Card className="glass border-none shadow-md bg-gradient-to-br from-red-50 to-orange-50 dark:from-red-950/20 dark:to-orange-950/20">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium text-red-600 dark:text-red-400">At-Risk Students</CardTitle>
                        <AlertTriangle className="h-5 w-5 text-red-600 dark:text-red-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="flex items-baseline space-x-2">
                            <div className="text-3xl font-bold text-red-700 dark:text-red-300">
                                {atRisk?.high_risk_count || 0}
                            </div>
                            <span className="text-sm text-muted-foreground">High Risk</span>
                        </div>
                        <div className="mt-2 text-sm text-orange-600 dark:text-orange-400">
                            + {atRisk?.medium_risk_count || 0} Medium Risk
                        </div>
                        <p className="text-xs text-muted-foreground mt-4 italic">
                            Students with Low Grades (D/F) needs attention.
                        </p>
                    </CardContent>
                </Card>

                {/* Class Health Stats */}
                <Card className="glass border-none shadow-md">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium">Performance Heatmap</CardTitle>
                        <Users className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            {health && Object.entries(health.performance_heatmap).map(([grade, count]) => (
                                <div key={grade} className="flex items-center justify-between text-sm">
                                    <span className="text-muted-foreground">{grade}</span>
                                    <span className="font-medium">{count} Students</span>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>

                {/* Alerts Section */}
                <Card className="glass border-none shadow-md md:col-span-2 lg:col-span-1">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium">Early Warning Alerts</CardTitle>
                        <TrendingUp className="h-5 w-5 text-amber-600 dark:text-amber-400" />
                    </CardHeader>
                    <CardContent>
                        <ul className="space-y-3">
                            {alerts.map((alert, idx) => (
                                <li key={idx} className="flex gap-2 p-2 rounded bg-amber-50 dark:bg-amber-900/10 text-sm text-amber-900 dark:text-amber-200 border border-amber-100 dark:border-amber-800/50">
                                    <AlertTriangle className="h-4 w-4 shrink-0 mt-0.5" />
                                    <span>{alert}</span>
                                </li>
                            ))}
                            {alerts.length === 0 && <p className="text-sm text-muted-foreground">No active alerts.</p>}
                        </ul>
                    </CardContent>
                </Card>
            </div>

            <div className="grid gap-6 md:grid-cols-3">
                {/* Students List Section */}
                {students.length > 0 && (
                    <Card className="glass border-none shadow-md md:col-span-2">
                        <CardHeader>
                            <CardTitle>My Students ({students.length})</CardTitle>
                            <CardDescription>Enrolled in {course?.title}. Click to view details.</CardDescription>
                        </CardHeader>
                        <CardContent>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                {students.map((student) => (
                                    <div
                                        key={student.student_id}
                                        onClick={() => fetchStudentDetails(student.student_id)}
                                        className="p-4 rounded-lg bg-card border flex items-center gap-3 cursor-pointer hover:bg-secondary/50 transition-colors"
                                    >
                                        <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold">
                                            {student.name.charAt(0)}
                                        </div>
                                        <div>
                                            <p className="font-medium">{student.name}</p>
                                            <p className="text-xs text-muted-foreground">{student.student_id}</p>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                )}

                {/* Announcement Box */}
                <Card className="glass border-none shadow-md">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <BookOpen className="h-4 w-4" />
                            Broadcast
                        </CardTitle>
                        <CardDescription>Post to enrolled students</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <textarea
                            className="flex min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                            placeholder="Type announcement..."
                            value={announcement}
                            onChange={(e) => setAnnouncement(e.target.value)}
                        />
                        <button
                            onClick={postAnnouncement}
                            disabled={!announcement}
                            className="inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground shadow hover:bg-primary/90 h-9 px-4 py-2 w-full"
                        >
                            Post Announcement
                        </button>
                    </CardContent>
                </Card>
            </div>

            {/* Student Details Dialog/Modal */}
            {selectedStudent && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
                    <Card className="w-full max-w-md bg-background border-none shadow-xl">
                        <CardHeader className="flex flex-row items-center justify-between">
                            <div>
                                <CardTitle>{selectedStudent.name}</CardTitle>
                                <CardDescription>{selectedStudent.student_id}</CardDescription>
                            </div>
                            <button onClick={() => setSelectedStudent(null)} className="text-muted-foreground hover:text-foreground">X</button>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <p className="text-muted-foreground">Email</p>
                                    <p className="font-medium truncate">{selectedStudent.email}</p>
                                </div>
                                <div>
                                    <p className="text-muted-foreground">Course Grade</p>
                                    <p className={`font-bold ${['A', 'A+'].includes(selectedStudent.current_grade) ? 'text-green-600' : 'text-foreground'}`}>
                                        {selectedStudent.current_grade}
                                    </p>
                                </div>
                                <div className="col-span-2">
                                    <p className="text-muted-foreground">Risk Status</p>
                                    <div className={`mt-1 inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ${selectedStudent.risk_status === 'High Risk'
                                            ? 'border-transparent bg-red-500 text-white shadow hover:bg-red-600'
                                            : 'border-transparent bg-green-500 text-white shadow hover:bg-green-600'
                                        }`}>
                                        {selectedStudent.risk_status}
                                    </div>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            )}
        </div>
    )
}
