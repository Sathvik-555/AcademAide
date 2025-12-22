"use client"

import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { cn } from "@/lib/utils"
import { useEffect, useState } from "react"
import Cookies from "js-cookie"
import { Loader2 } from "lucide-react"

type TimetableItem = {
    course_id: string
    title: string
    section_name: string
    day_of_week: string
    start_time: string
    end_time: string
    room_number: string
}

const baseSchedule = [
    { time: "09:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "10:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "11:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "12:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "13:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "14:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
    { time: "15:00", mon: "", tue: "", wed: "", thu: "", fri: "" },
]

// Helper to determine cell color based on subject or content
function getSubjectStyle(subject: string) {
    if (!subject) return ""
    // Fixed slots
    if (subject === "Break" || subject === "Lunch") return "bg-slate-100 text-slate-500 italic dark:bg-slate-800/50 dark:text-slate-400"

    // Check for course codes or names
    const text = subject.toUpperCase()
    if (text.includes("DBMS") || text.includes("CS")) return "bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300 border-l-4 border-violet-500"
    if (text.includes("CN") || text.includes("NET")) return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 border-l-4 border-blue-500"
    if (text.includes("OS") || text.includes("SYS")) return "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300 border-l-4 border-emerald-500"
    if (text.includes("DAA") || text.includes("ALG")) return "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 border-l-4 border-amber-500"
    if (text.includes("AI") || text.includes("INT")) return "bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300 border-l-4 border-rose-500"
    if (text.includes("PROJ")) return "bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300 border-l-4 border-indigo-500"

    // Default
    return "bg-slate-50 text-slate-700 dark:bg-slate-800/30 dark:text-slate-300 border-l-4 border-gray-300"
}

export default function TimetablePage() {
    const [schedule, setSchedule] = useState(baseSchedule)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")

    useEffect(() => {
        const fetchTimetable = async () => {
            const studentId = Cookies.get("student_id")
            if (!studentId) {
                setError("No student ID found.")
                setLoading(false)
                return
            }

            try {
                const res = await fetch(`http://localhost:8000/student/timetable?student_id=${studentId}`)
                if (!res.ok) throw new Error("Failed to fetch timetable")
                const result = await res.json()
                const data: TimetableItem[] = result.data

                // Clone base schedule
                const newSchedule = JSON.parse(JSON.stringify(baseSchedule))

                // Map API data to table
                const dayMap: { [key: string]: string } = {
                    "Monday": "mon", "Tuesday": "tue", "Wednesday": "wed", "Thursday": "thu", "Friday": "fri"
                }

                data.forEach(item => {
                    const dayKey = dayMap[item.day_of_week]
                    if (!dayKey) return

                    // Find row by start time (assuming "HH:MM" format match)
                    // API might return "10:00:00", simplify to "10:00"
                    const startTimeShort = item.start_time.substring(0, 5)

                    const row = newSchedule.find((r: any) => r.time === startTimeShort)
                    if (row) {
                        row[dayKey] = item.course_id
                    }
                })

                // Insert standard breaks if logic allows, or just leave as is.
                // For now, let's inject Lunch at 13:00 if empty
                const lunchRow = newSchedule.find((r: any) => r.time === "13:00")
                if (lunchRow) {
                    ["mon", "tue", "wed", "thu", "fri"].forEach(d => {
                        if (!lunchRow[d]) lunchRow[d] = "Lunch"
                    })
                }

                setSchedule(newSchedule)
            } catch (err) {
                console.error(err)
                setError("Could not load timetable")
            } finally {
                setLoading(false)
            }
        }

        fetchTimetable()
    }, [])

    if (loading) {
        return (
            <div className="flex h-full items-center justify-center p-12">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        )
    }

    return (
        <div className="flex flex-col gap-8">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-4xl font-extrabold tracking-tight gradient-text">Timetable</h1>
                    <p className="text-muted-foreground mt-1">Stay organized with your weekly schedule.</p>
                </div>
            </div>

            <Card className="glass border-2 border-white/50 dark:border-white/10 shadow-xl overflow-hidden">
                <CardHeader className="bg-gradient-to-r from-indigo-50/50 to-purple-50/50 dark:from-indigo-900/20 dark:to-purple-900/20 border-b border-gray-100 dark:border-white/5">
                    <CardTitle className="text-xl text-indigo-900 dark:text-indigo-100">Weekly Schedule</CardTitle>
                    <CardDescription>View your classes and labs for the semester.</CardDescription>
                </CardHeader>
                <CardContent className="p-0">
                    <div className="relative overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
                        <table className="w-full text-sm text-left rtl:text-right">
                            <thead className="text-xs uppercase bg-indigo-50/80 dark:bg-indigo-950/50 text-indigo-600 dark:text-indigo-300 font-bold tracking-wider">
                                <tr>
                                    <th scope="col" className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 w-24">Time</th>
                                    <th scope="col" className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 min-w-[120px]">Monday</th>
                                    <th scope="col" className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 min-w-[120px]">Tuesday</th>
                                    <th scope="col" className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 min-w-[120px]">Wednesday</th>
                                    <th scope="col" className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 min-w-[120px]">Thursday</th>
                                    <th scope="col" className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 min-w-[120px]">Friday</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-indigo-50 dark:divide-white/5">
                                {schedule.map((slot, index) => (
                                    <tr key={index} className="bg-white/40 dark:bg-slate-900/40 hover:bg-white/80 dark:hover:bg-slate-900/80 transition-colors last:border-b-0 border-b border-gray-200 dark:border-gray-700">
                                        <td className="px-6 py-4 font-medium whitespace-nowrap text-slate-500 dark:text-slate-400 border-r border-gray-200 dark:border-gray-700 last:border-r-0">
                                            {slot.time}
                                        </td>
                                        {/* @ts-ignore */}
                                        {["mon", "tue", "wed", "thu", "fri"].map((day) => (
                                            <td key={day} className="px-2 py-2 border-r border-gray-200 dark:border-gray-700 last:border-r-0 h-16">
                                                {/* @ts-ignore */}
                                                <div className={cn("h-full w-full rounded-md flex items-center justify-center p-2 text-center text-xs font-semibold shadow-sm transition-all hover:scale-[1.02]", getSubjectStyle(slot[day]))}>
                                                    {/* @ts-ignore */}
                                                    {slot[day]}
                                                </div>
                                            </td>
                                        ))}
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
