import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"

const schedule = [
    { time: "09:00", mon: "DBMS (L)", tue: "", wed: "CN (L)", thu: "DBMS (T)", fri: "OS (L)" },
    { time: "10:00", mon: "DBMS (L)", tue: "DAA (L)", wed: "CN (L)", thu: "", fri: "OS (L)" },
    { time: "11:00", mon: "Break", tue: "DAA (L)", wed: "Break", thu: "AI (L)", fri: "Break" },
    { time: "12:00", mon: "OS (L)", tue: "AI (L)", wed: "DBMS Lab", thu: "AI (L)", fri: "DAA (T)" },
    { time: "13:00", mon: "Lunch", tue: "Lunch", wed: "Lunch", thu: "Lunch", fri: "Lunch" },
    { time: "14:00", mon: "CN (L)", tue: "Project", wed: "DBMS Lab", thu: "Project", fri: "Seminar" },
    { time: "15:00", mon: "", tue: "Project", wed: "", thu: "", fri: "" },
]

export default function TimetablePage() {
    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold tracking-tight">Timetable</h1>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Weekly Schedule</CardTitle>
                    <CardDescription>View your classes and labs for the semester.</CardDescription>
                </CardHeader>
                <CardContent>
                    <div className="relative overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
                        <table className="w-full text-sm text-left rtl:text-right">
                            <thead className="text-xs uppercase bg-muted">
                                <tr>
                                    <th scope="col" className="px-6 py-3 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0">Time</th>
                                    <th scope="col" className="px-6 py-3 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0">Monday</th>
                                    <th scope="col" className="px-6 py-3 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0">Tuesday</th>
                                    <th scope="col" className="px-6 py-3 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0">Wednesday</th>
                                    <th scope="col" className="px-6 py-3 border-b border-r border-gray-200 dark:border-gray-700 last:border-r-0">Thursday</th>
                                    <th scope="col" className="px-6 py-3 border-b border-gray-200 dark:border-gray-700">Friday</th>
                                </tr>
                            </thead>
                            <tbody>
                                {schedule.map((slot, index) => (
                                    <tr key={index} className="bg-background hover:bg-muted/50 transition-colors last:border-b-0 border-b border-gray-200 dark:border-gray-700">
                                        <td className="px-6 py-4 font-medium whitespace-nowrap border-r border-gray-200 dark:border-gray-700 last:border-r-0">{slot.time}</td>
                                        <td className="px-6 py-4 border-r border-gray-200 dark:border-gray-700 last:border-r-0">{slot.mon}</td>
                                        <td className="px-6 py-4 border-r border-gray-200 dark:border-gray-700 last:border-r-0">{slot.tue}</td>
                                        <td className="px-6 py-4 border-r border-gray-200 dark:border-gray-700 last:border-r-0">{slot.wed}</td>
                                        <td className="px-6 py-4 border-r border-gray-200 dark:border-gray-700 last:border-r-0">{slot.thu}</td>
                                        <td className="px-6 py-4">{slot.fri}</td>
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
