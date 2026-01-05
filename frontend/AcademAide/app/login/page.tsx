"use client"

import { useRouter } from "next/navigation"
import Link from "next/link"
import { Button } from "@/components/ui/Button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"
import { BookOpen, Loader2 } from "lucide-react"
import Cookies from "js-cookie"
import { useState } from "react"
import { cn } from "@/lib/utils"

export default function LoginPage() {
    const router = useRouter()
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState("")
    const [role, setRole] = useState<"student" | "teacher">("student")

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        setError("")

        const formData = new FormData(e.currentTarget as HTMLFormElement)
        const id = formData.get("id") as string
        const password = formData.get("password") as string
        const password = formData.get("password") as string

        try {
            const res = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ id, password, role }),
            })
        })

        const data = await res.json()

        if (!res.ok) {
            throw new Error(data.error || "Login failed")
        }

        // Store token and role-specific ID
        Cookies.set("token", data.token, { expires: 1 })
        Cookies.set("user_id", data.user_id, { expires: 1 })
        Cookies.set("role", role, { expires: 1 }) // Store role for middleware checks

        if (role === "student") {
            Cookies.set("student_id", data.user_id, { expires: 1 })
            Cookies.set("wallet_address", data.wallet_address || "", { expires: 1 })
            router.push("/dashboard")
        } else {
            Cookies.set("faculty_id", data.user_id, { expires: 1 })
            router.push("/teacher-dashboard")
        }

        router.refresh()
    } catch (err: any) {
        setError(err.message)
    } finally {
        setLoading(false)
    }
}

return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background">
        {/* Background Gradients */}
        <div className="absolute -top-[20%] -left-[10%] h-[500px] w-[500px] rounded-full bg-violet-500/30 blur-[100px]" />
        <div className="absolute top-[40%] -right-[10%] h-[400px] w-[400px] rounded-full bg-fuchsia-500/30 blur-[100px]" />
        <div className="absolute -bottom-[10%] left-[20%] h-[600px] w-[600px] rounded-full bg-orange-500/20 blur-[100px]" />

        <div className="absolute top-8 left-8 z-20">
            <Link className="flex items-center gap-2 font-bold text-xl tracking-tight text-foreground" href="/">
                <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30">
                    <BookOpen className="h-5 w-5" />
                </div>
                <span>AcademAide</span>
            </Link>
        </div>

        <Card className="glass z-10 w-full max-w-md border-white/20">
            <CardHeader className="space-y-1 text-center">
                <CardTitle className="text-3xl font-bold tracking-tight">Welcome back</CardTitle>
                <CardDescription className="text-base">
                    Login to access your personalized academic assistant.
                </CardDescription>
            </CardHeader>
            <form onSubmit={handleLogin}>
                <CardContent className="grid gap-4">
                    {error && (
                        <div className="p-3 rounded-md bg-red-50 text-red-500 text-sm font-medium dark:bg-red-900/10">
                            {error}
                        </div>
                    )}
                    <div className="grid gap-2">
                        <Button variant="outline" type="button" onClick={() => window.location.href = "http://localhost:8080/auth/google/login"} className="w-full bg-white dark:bg-slate-950 text-foreground border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-900 shadow-sm h-10">
                            <svg className="mr-2 h-4 w-4" aria-hidden="true" focusable="false" data-prefix="fab" data-icon="google" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 488 512">
                                <path fill="currentColor" d="M488 261.8C488 403.3 391.1 504 248 504 110.8 504 0 393.2 0 256S110.8 8 248 8c66.8 0 123 24.5 166.3 64.9l-67.5 64.9C258.5 52.6 94.3 116.6 94.3 256c0 86.5 69.1 156.6 153.7 156.6 98.2 0 135-70.4 140.8-106.9H248v-85.3h236.1c2.3 12.7 3.9 24.9 3.9 41.4z"></path>
                            </svg>
                            Sign in with Google
                        </Button>
                        <div className="relative">
                            <div className="absolute inset-0 flex items-center">
                                <span className="w-full border-t border-muted" />
                            </div>
                            <div className="relative flex justify-center text-xs uppercase">
                                <span className="bg-background px-2 text-muted-foreground">Or continue with</span>
                            </div>
                        </div>
                    </div>
                    <div className="flex p-1 bg-muted rounded-lg mb-4">
                        <button
                            type="button"
                            onClick={() => setRole("student")}
                            className={cn(
                                "flex-1 py-1.5 text-sm font-medium rounded-md transition-all",
                                role === "student" ? "bg-white dark:bg-slate-800 shadow-sm" : "text-muted-foreground hover:text-foreground"
                            )}
                        >
                            Student
                        </button>
                        <button
                            type="button"
                            onClick={() => setRole("teacher")}
                            className={cn(
                                "flex-1 py-1.5 text-sm font-medium rounded-md transition-all",
                                role === "teacher" ? "bg-white dark:bg-slate-800 shadow-sm" : "text-muted-foreground hover:text-foreground"
                            )}
                        >
                            Teacher
                        </button>
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="id">{role === "student" ? "Student ID" : "Faculty ID"}</Label>
                        <Input
                            id="id"
                            name="id"
                            type="text"
                            placeholder={role === "student" ? "S1001" : "F001"}
                            required
                            className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800"
                        />
                    </div>
                    <div className="grid gap-2">
                        <Label htmlFor="password">Password</Label>
                        <Input
                            id="password"
                            name="password"
                            type="password"
                            required
                            className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800"
                        />
                    </div>
                </CardContent>
                <CardFooter className="flex flex-col gap-4">
                    <Button className="w-full bg-primary hover:bg-primary/90 shadow-lg shadow-primary/25 h-10 transition-all hover:scale-[1.02]" disabled={loading}>
                        {loading ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                Signing in...
                            </>
                        ) : (
                            "Sign in"
                        )}
                    </Button>
                    <div className="text-center text-sm text-muted-foreground">
                        Don't have an account?{" "}
                        <Link href="/signup" className="underline underline-offset-4 hover:text-primary font-medium transition-colors">
                            Sign up
                        </Link>
                    </div>
                </CardFooter>
            </form>
        </Card>
    </div>
)
}
