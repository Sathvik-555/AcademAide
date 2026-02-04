"use client"

import { Button } from "@/components/ui/Button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"
import { useEffect, useState } from "react"
import Cookies from "js-cookie"
import { Loader2, Mail, Phone, User, Briefcase } from "lucide-react"

type TeacherProfile = {
    faculty_id: string
    first_name: string
    last_name: string
    email: string
    phone_no: string
    departments: string[]
}

export default function TeacherProfilePage() {
    const [profile, setProfile] = useState<TeacherProfile | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")

    useEffect(() => {
        const fetchProfile = async () => {
            const token = Cookies.get("token")
            if (!token) {
                setError("Unauthorized")
                setLoading(false)
                return
            }

            try {
                const res = await fetch("http://localhost:8080/teacher/profile", {
                    headers: {
                        "Authorization": `Bearer ${token}`
                    }
                })
                if (!res.ok) throw new Error("Failed to fetch profile")
                const data = await res.json()
                setProfile(data)
            } catch (err) {
                setError("Could not load profile data")
            } finally {
                setLoading(false)
            }
        }

        fetchProfile()
    }, [])

    if (loading) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
            </div>
        )
    }

    if (error) {
        return (
            <div className="flex h-full items-center justify-center text-red-500 font-medium">
                {error}
            </div>
        )
    }

    return (
        <div className="max-w-4xl mx-auto grid gap-6">
            <h1 className="text-3xl font-bold">Faculty Profile</h1>

            <Card className="border-l-4 border-l-primary/50">
                <CardHeader>
                    <div className="flex items-center gap-4">
                        <div className="h-16 w-16 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                            <User className="h-8 w-8" />
                        </div>
                        <div>
                            <CardTitle className="text-2xl">{profile?.first_name} {profile?.last_name}</CardTitle>
                            <CardDescription>{profile?.faculty_id}</CardDescription>
                        </div>
                    </div>
                </CardHeader>
                <CardContent className="grid gap-6 mt-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <div className="space-y-2">
                            <Label className="flex items-center gap-2 text-muted-foreground">
                                <Mail className="h-4 w-4" /> Email Address
                            </Label>
                            <Input value={profile?.email || ""} disabled className="bg-muted/50" />
                        </div>
                        <div className="space-y-2">
                            <Label className="flex items-center gap-2 text-muted-foreground">
                                <Phone className="h-4 w-4" /> Phone Number
                            </Label>
                            <Input value={profile?.phone_no || ""} disabled className="bg-muted/50" />
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Briefcase className="h-5 w-5 text-primary" />
                        Academic Affiliations
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        <div>
                            <Label>Departments</Label>
                            <div className="mt-2 flex flex-wrap gap-2">
                                {profile?.departments && profile.departments.length > 0 ? (
                                    profile.departments.map((dept, i) => (
                                        <span key={i} className="px-3 py-1 rounded-full bg-secondary text-secondary-foreground text-sm font-medium">
                                            {dept}
                                        </span>
                                    ))
                                ) : (
                                    <span className="text-muted-foreground text-sm">No specific department affiliations found.</span>
                                )}
                            </div>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
