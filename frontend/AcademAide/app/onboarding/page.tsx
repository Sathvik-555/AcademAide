"use client"

import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Button } from "@/components/ui/Button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"
import { BookOpen, Loader2 } from "lucide-react"
import { Suspense } from "react"
import Cookies from "js-cookie" // Make sure to install if not present, but it was used in login/page.tsx

function OnboardingContent() {
    const router = useRouter()
    const searchParams = useSearchParams()
    const token = searchParams.get("token")
    const email = searchParams.get("email")

    // We can also extract name from token claims if we wanted, but let's ask user.

    const [loading, setLoading] = useState(false)
    const [error, setError] = useState("")

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault()
        setLoading(true)
        setError("")

        const formData = new FormData(e.currentTarget as HTMLFormElement)

        const payload = {
            first_name: formData.get("first_name"),
            last_name: formData.get("last_name"),
            dept_id: formData.get("dept_id"),
            semester: parseInt(formData.get("semester") as string),
            year_of_joining: parseInt(formData.get("year_of_joining") as string),
            phone_no: formData.get("phone_no"),
        }

        try {
            const res = await fetch("http://localhost:8080/auth/complete-registration", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify(payload),
            })

            const data = await res.json()

            if (!res.ok) {
                throw new Error(data.error || "Registration failed")
            }

            // Store final token and student_id
            Cookies.set("token", data.token, { expires: 1 })
            Cookies.set("student_id", data.student_id, { expires: 1 })

            router.push("/dashboard")
            router.refresh()
        } catch (err: any) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    if (!token) {
        return <div className="p-10 text-center text-red-500">Invalid Session. Please login again.</div>
    }

    return (
        <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-background">
            {/* Background Gradients */}
            <div className="absolute -top-[20%] -left-[10%] h-[500px] w-[500px] rounded-full bg-violet-500/30 blur-[100px]" />
            <div className="absolute top-[40%] -right-[10%] h-[400px] w-[400px] rounded-full bg-fuchsia-500/30 blur-[100px]" />

            <Card className="glass z-10 w-full max-w-lg border-white/20">
                <CardHeader className="space-y-1 text-center">
                    <CardTitle className="text-3xl font-bold tracking-tight">Complete Profile</CardTitle>
                    <CardDescription className="text-base">
                        Please provide a few more details to set up your account for {email}.
                    </CardDescription>
                </CardHeader>
                <form onSubmit={handleRegister}>
                    <CardContent className="grid gap-4">
                        {error && (
                            <div className="p-3 rounded-md bg-red-50 text-red-500 text-sm font-medium dark:bg-red-900/10">
                                {error}
                            </div>
                        )}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="first_name">First Name</Label>
                                <Input id="first_name" name="first_name" required placeholder="John" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="last_name">Last Name</Label>
                                <Input id="last_name" name="last_name" required placeholder="Doe" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="dept_id">Department ID</Label>
                                <Input id="dept_id" name="dept_id" required placeholder="CS" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="semester">Semester</Label>
                                <Input id="semester" name="semester" type="number" required placeholder="1" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div className="grid gap-2">
                                <Label htmlFor="year_of_joining">Year of Joining</Label>
                                <Input id="year_of_joining" name="year_of_joining" type="number" required placeholder="2024" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                            <div className="grid gap-2">
                                <Label htmlFor="phone_no">Phone Number</Label>
                                <Input id="phone_no" name="phone_no" required placeholder="555-0123" className="bg-white/50 border-gray-200 focus:border-primary/50 focus:ring-primary/20 dark:bg-slate-950/50 dark:border-slate-800" />
                            </div>
                        </div>

                    </CardContent>
                    <CardFooter>
                        <Button className="w-full bg-primary hover:bg-primary/90 shadow-lg shadow-primary/25 h-10 transition-all hover:scale-[1.02]" disabled={loading}>
                            {loading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Setting up...
                                </>
                            ) : (
                                "Complete Registration"
                            )}
                        </Button>
                    </CardFooter>
                </form>
            </Card>
        </div>
    )
}

export default function OnboardingPage() {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <OnboardingContent />
        </Suspense>
    )
}
