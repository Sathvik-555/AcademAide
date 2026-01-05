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

const timeSlots = [
    "09:00 - 10:00",
    "10:00 - 11:00",
    "11:00 - 11:30", // Break
    "11:30 - 12:30",
    "12:30 - 13:30", // Lunch
    "13:30 - 14:30", // Is this right? Wait, user asked for 11:30 to 12:30.
    // User Request: "9:00 to 10:00 and then 10:00 to 11:00 and then 11:30 to 12:30"
    // My previous baseSchedule had 11:00, 11:30...
    // Let's adjust to be specific.
    // But this array is also used for MAPPING data "startTime".
    // If I change strings here, I break the mapping logic: `if (timeSlots.includes(startTimeShort))`
    // I should separate Labels from Keys.
]

// Updated Approach: Use Object for mapping
const timeSlotConfig = [
    { key: "09:00", label: "09:00 - 10:00" },
    { key: "10:00", label: "10:00 - 11:00" },
    { key: "11:00", label: "11:00 - 11:30" },
    { key: "11:30", label: "11:30 - 12:30" },
    { key: "12:30", label: "12:30 - 13:30" },
    { key: "13:30", label: "13:30 - 14:30" }, // Adjusted based on previous 13:30 start, but data says 14:30 lab
    { key: "14:30", label: "14:30 - 15:30" },
    { key: "15:30", label: "15:30 - 16:30" },
]

const daysOfWeek = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]

// Helper to determine cell color based on subject
function getSubjectStyle(subject: string) {
    if (!subject) return ""
    if (subject === "Lunch" || subject === "Break") return "bg-slate-100 text-slate-500 italic dark:bg-slate-800/50 dark:text-slate-400"

    const text = subject.toUpperCase()
    if (text.includes("DBMS") || text.includes("CD252")) return "bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-300 border-l-4 border-violet-500"
    if (text.includes("TOC") || text.includes("CS354")) return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300 border-l-4 border-blue-500"
    if (text.includes("AIML") || text.includes("IS353")) return "bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300 border-l-4 border-emerald-500"
    if (text.includes("CLOUD") || text.includes("XX355")) return "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300 border-l-4 border-amber-500"
    if (text.includes("ECON") || text.includes("HS251")) return "bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300 border-l-4 border-rose-500"

    return "bg-slate-50 text-slate-700 dark:bg-slate-800/30 dark:text-slate-300 border-l-4 border-gray-300"
}

export default function TimetablePage() {
    // Transposed structure: object where keys are Time Slots
    // We will render Rows as Days.
    // So schedule[Day][Time] = Data
    const [scheduleMap, setScheduleMap] = useState<Record<string, Record<string, any>>>({})
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

            const token = Cookies.get("token")

            try {
                const res = await fetch(`http://localhost:8080/student/timetable?student_id=${studentId}`, {
                    headers: {
                        "Authorization": `Bearer ${token}`
                    }
                })
                if (!res.ok) throw new Error("Failed to fetch timetable")
                const result = await res.json()
                const data: TimetableItem[] = result.data

                // Build Map: Day -> Time Key -> Data
                const newMap: Record<string, Record<string, any>> = {}

                // Initialize empty
                daysOfWeek.forEach(day => {
                    newMap[day] = {}
                    timeSlotConfig.forEach(slot => {
                        newMap[day][slot.key] = null
                    })
                })

                data.forEach(item => {
                    const day = item.day_of_week
                    if (!newMap[day]) return

                    const startTimeShort = item.start_time.substring(0, 5) // "09:00"

                    // Assign to slot
                    // Check if exists in config
                    const slotFound = timeSlotConfig.find(s => s.key === startTimeShort)

                    if (slotFound) {
                        newMap[day][startTimeShort] = {
                            course_id: item.course_id,
                            room_number: item.room_number,
                            title: item.title,
                            end_time: item.end_time.substring(0, 5),
                            type: "class"
                        }
                    }
                })

                setScheduleMap(newMap)
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
                <CardContent className="p-0 overflow-x-auto">
                    <div className="min-w-max border border-gray-200 dark:border-gray-700">
                        <table className="w-full text-sm text-left">
                            <thead className="text-xs uppercase bg-indigo-50/80 dark:bg-indigo-950/50 text-indigo-600 dark:text-indigo-300 font-bold tracking-wider">
                                <tr>
                                    <th className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 w-32 sticky left-0 bg-indigo-50/95 dark:bg-slate-900/95 z-10">Day / Time</th>
                                    {timeSlotConfig.map(slot => (
                                        <th key={slot.key} className="px-6 py-4 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0 min-w-[140px] text-center">
                                            {slot.label}
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-indigo-50 dark:divide-white/5">
                                {daysOfWeek.map(day => {
                                    const skippedIndices = new Set<number>()

                                    return (
                                        <tr key={day} className="bg-white/40 dark:bg-slate-900/40 hover:bg-white/80 dark:hover:bg-slate-900/80 transition-colors border-b border-gray-200 dark:border-gray-700 last:border-b-0">
                                            <td className="px-6 py-4 font-bold text-slate-700 dark:text-slate-300 border-r border-gray-200 dark:border-gray-700 sticky left-0 bg-white/95 dark:bg-slate-900/95 z-10 shadow-[2px_0_5px_-2px_rgba(0,0,0,0.1)]">
                                                {day}
                                            </td>
                                            {timeSlotConfig.map((slot, index) => {
                                                if (skippedIndices.has(index)) return null

                                                const cellData = scheduleMap[day]?.[slot.key]
                                                let colSpan = 1

                                                if (cellData) {
                                                    const endTime = cellData.end_time

                                                    for (let i = index + 1; i < timeSlotConfig.length; i++) {
                                                        const nextSlotKey = timeSlotConfig[i].key
                                                        if (nextSlotKey < endTime) {
                                                            colSpan++
                                                            skippedIndices.add(i)
                                                        } else {
                                                            break
                                                        }
                                                    }
                                                }

                                                return (
                                                    <td key={`${day}-${slot.key}`} colSpan={colSpan} className="px-2 py-2 border-r border-gray-200 dark:border-gray-700 last:border-r-0 h-24">
                                                        {cellData ? (
                                                            <div className={cn("h-full w-full rounded-md flex flex-col items-center justify-center p-2 text-center shadow-sm cursor-help transition-all hover:scale-[1.02] hover:shadow-md", getSubjectStyle(cellData.course_id))} title={cellData.title}>
                                                                <div className="text-xs font-bold break-words">{cellData.course_id}</div>
                                                                <div className="text-[10px] opacity-90 mt-1 flex items-center gap-1 justify-center">
                                                                    <span>📍 {cellData.room_number}</span>
                                                                </div>
                                                            </div>
                                                        ) : (
                                                            <div className="h-full w-full"></div>
                                                        )}
                                                    </td>
                                                )
                                            })}
                                        </tr>
                                    )
                                })}
                            </tbody>
                        </table>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
