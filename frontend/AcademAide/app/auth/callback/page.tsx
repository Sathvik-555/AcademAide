"use client"

import { useEffect, Suspense } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Cookies from "js-cookie"
import { Loader2 } from "lucide-react"

function CallbackContent() {
    const router = useRouter()
    const searchParams = useSearchParams()

    useEffect(() => {
        const token = searchParams.get("token")
        const mode = searchParams.get("mode")
        const studentId = searchParams.get("student_id")
        const email = searchParams.get("email")

        if (token) {
            if (mode === "login") {
                Cookies.set("token", token, { expires: 1 })
                if (studentId) Cookies.set("student_id", studentId, { expires: 1 })
                router.push("/dashboard")
                router.refresh()
            } else if (mode === "signup") {
                // Redirect for onboarding with token in query for easy access in Onboarding page
                // We'll pass it as a query param 'token'
                const onboardingUrl = `/onboarding?token=${token}&email=${email || ""}`
                router.push(onboardingUrl)
            }
        } else {
            router.push("/login?error=auth_failed")
        }
    }, [searchParams, router])

    return (
        <div className="flex h-screen w-full items-center justify-center bg-background">
            <div className="flex flex-col items-center gap-2">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-muted-foreground">Authenticating...</p>
            </div>
        </div>
    )
}

export default function CallbackPage() {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <CallbackContent />
        </Suspense>
    )
}
