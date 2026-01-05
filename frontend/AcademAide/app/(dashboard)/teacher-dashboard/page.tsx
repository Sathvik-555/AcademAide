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
    const [health, setHealth] = useState<ClassHealth | null>(null)
    const [atRisk, setAtRisk] = useState<AtRiskStats | null>(null)
    const [alerts, setAlerts] = useState<string[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchData = async () => {
            const token = Cookies.get("token")
            if (!token) {
                router.push("/login")
                return
            }

            // Using mock course_id "CS101" for MVP demo as per seeded data
            const courseId = "CS101"
            const headers = { "Authorization": `Bearer ${token}` }

            try {
                // Fetch Class Health
                const healthRes = await fetch(`http://localhost:8080/teacher/class-health?course_id=${courseId}`, { headers })
                if (healthRes.ok) setHealth(await healthRes.json())

                // Fetch At-Risk Stats
                const riskRes = await fetch(`http://localhost:8080/teacher/at-risk?course_id=${courseId}`, { headers })
                if (riskRes.ok) setAtRisk(await riskRes.json())

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

    if (loading) {
        return <div className="flex h-screen items-center justify-center">Loading dashboard...</div>
    }

    return (
        <div className="flex min-h-screen flex-col p-8 gap-8 bg-gray-50 dark:bg-slate-950">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Teacher Dashboard</h1>
                    <p className="text-muted-foreground">Overview for {health?.title || "your classes"}</p>
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
                            Students with attendance &lt;75% or falling grades.
                        </p>
                    </CardContent>
                </Card>

                {/* Class Health Stats */}
                <Card className="glass border-none shadow-md">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium">Attendance Distribution</CardTitle>
                        <Users className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-2">
                            {health && Object.entries(health.attendance_distribution).map(([range, count]) => (
                                <div key={range} className="flex items-center justify-between text-sm">
                                    <span className="text-muted-foreground">{range}</span>
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

            {/* Performance Heatmap (Simplified as List for MVP) */}
            <Card className="glass border-none shadow-md">
                <CardHeader>
                    <CardTitle>Performance Overview</CardTitle>
                    <CardDescription>Grade distribution for {health?.title}</CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                        {health && Object.entries(health.performance_heatmap).map(([grade, count]) => (
                            <div key={grade} className="p-4 rounded-lg bg-secondary/20 flex flex-col items-center justify-center text-center">
                                <span className="text-2xl font-bold">{count}</span>
                                <span className="text-sm text-muted-foreground mt-1">{grade}</span>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
