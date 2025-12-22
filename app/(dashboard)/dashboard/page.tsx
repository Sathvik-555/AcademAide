import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { BookOpen, Trophy, Clock, AlertCircle } from "lucide-react"

export default function DashboardPage() {
    return (
        <div className="flex flex-col gap-8">
            <div className="flex flex-col gap-2">
                <h1 className="text-4xl font-extrabold tracking-tight gradient-text">Hello, Sharanya! 👋</h1>
                <p className="text-lg text-muted-foreground">Ready to conquer your academic goals today?</p>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
                {/* Colorful Stats Cards */}
                <Card className="glass border-none shadow-lg hover:shadow-xl transition-all hover:scale-[1.02] bg-gradient-to-br from-violet-500/10 to-purple-500/10 dark:from-violet-500/20 dark:to-purple-500/20">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-bold text-violet-700 dark:text-violet-300">Courses Enrolled</CardTitle>
                        <div className="h-8 w-8 rounded-full bg-violet-100 dark:bg-violet-900/50 flex items-center justify-center text-violet-600 dark:text-violet-300">
                            <BookOpen className="h-4 w-4" />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold text-violet-900 dark:text-violet-100">5</div>
                        <p className="text-xs font-medium text-violet-600/80 dark:text-violet-300/80">Active this semester</p>
                    </CardContent>
                </Card>

                <Card className="glass border-none shadow-lg hover:shadow-xl transition-all hover:scale-[1.02] bg-gradient-to-br from-emerald-500/10 to-teal-500/10 dark:from-emerald-500/20 dark:to-teal-500/20">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-bold text-emerald-700 dark:text-emerald-300">CGPA</CardTitle>
                        <div className="h-8 w-8 rounded-full bg-emerald-100 dark:bg-emerald-900/50 flex items-center justify-center text-emerald-600 dark:text-emerald-300">
                            <Trophy className="h-4 w-4" />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold text-emerald-900 dark:text-emerald-100">8.9</div>
                        <p className="text-xs font-medium text-emerald-600/80 dark:text-emerald-300/80">+0.2 from last semester</p>
                    </CardContent>
                </Card>

                <Card className="glass border-none shadow-lg hover:shadow-xl transition-all hover:scale-[1.02] bg-gradient-to-br from-amber-500/10 to-orange-500/10 dark:from-amber-500/20 dark:to-orange-500/20">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-bold text-amber-700 dark:text-amber-300">Up Next</CardTitle>
                        <div className="h-8 w-8 rounded-full bg-amber-100 dark:bg-amber-900/50 flex items-center justify-center text-amber-600 dark:text-amber-300">
                            <Clock className="h-4 w-4" />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className="text-xl font-bold text-amber-900 dark:text-amber-100 truncate">DBMS Lab</div>
                        <p className="text-xs font-medium text-amber-600/80 dark:text-amber-300/80">In 30 minutes</p>
                    </CardContent>
                </Card>

                <Card className="glass border-none shadow-lg hover:shadow-xl transition-all hover:scale-[1.02] bg-gradient-to-br from-rose-500/10 to-pink-500/10 dark:from-rose-500/20 dark:to-pink-500/20">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-bold text-rose-700 dark:text-rose-300">Assignments</CardTitle>
                        <div className="h-8 w-8 rounded-full bg-rose-100 dark:bg-rose-900/50 flex items-center justify-center text-rose-600 dark:text-rose-300">
                            <AlertCircle className="h-4 w-4" />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold text-rose-900 dark:text-rose-100">3</div>
                        <p className="text-xs font-medium text-rose-600/80 dark:text-rose-300/80">Due this week</p>
                    </CardContent>
                </Card>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-7">
                <Card className="col-span-4 glass border-none">
                    <CardHeader>
                        <CardTitle>Overview</CardTitle>
                    </CardHeader>
                    <CardContent className="pl-2">
                        <div className="h-[200px] flex items-center justify-center text-muted-foreground">
                            Attendance Chart Placeholder
                        </div>
                    </CardContent>
                </Card>
                <Card className="col-span-3">
                    <CardHeader>
                        <CardTitle>Recent Activity</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            <div className="flex items-center gap-4">
                                <div className="h-9 w-9 rounded-full bg-primary/10 flex items-center justify-center text-primary text-xs font-bold">AI</div>
                                <div className="grid gap-1">
                                    <p className="text-sm font-medium">Chatted with Assistant</p>
                                    <p className="text-xs text-muted-foreground">2 hours ago</p>
                                </div>
                            </div>
                            <div className="flex items-center gap-4">
                                <div className="h-9 w-9 rounded-full bg-muted flex items-center justify-center text-muted-foreground text-xs font-bold">DL</div>
                                <div className="grid gap-1">
                                    <p className="text-sm font-medium">Submission: Lab Report</p>
                                    <p className="text-xs text-muted-foreground">Yesterday</p>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    )
}
