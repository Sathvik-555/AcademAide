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

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        setError("")

        const formData = new FormData(e.currentTarget as HTMLFormElement)
        const studentId = formData.get("student_id") as string
        // Password is not used in the API request body per contract, but usually required for login. 
        // The contract says Request Body: { "student_id": "string" }. I will follow the contract strictly.

        try {
            const res = await fetch("http://localhost:8000/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ student_id: studentId }),
            })

            const data = await res.json()

            if (!res.ok) {
                throw new Error(data.error || "Login failed")
            }

            // Store token and student_id
            Cookies.set("token", data.token, { expires: 1 }) // Expires in 1 day
            Cookies.set("student_id", data.student_id, { expires: 1 })

            router.push("/dashboard")
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
                            <Label htmlFor="student_id">Student ID</Label>
                            <Input
                                id="student_id"
                                name="student_id"
                                type="text"
                                placeholder="2024CS123"
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
