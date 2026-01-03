"use client"

import { Button } from "@/components/ui/Button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"
import { useEffect, useState } from "react"
import Cookies from "js-cookie"
import { Loader2 } from "lucide-react"

type ProfileData = {
    student_id: string
    first_name: string
    last_name: string
    email: string
    phone_no: string
    semester: number
    year_of_joining: number
    dept_id: string
}

export default function ProfilePage() {
    const [profile, setProfile] = useState<ProfileData | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState("")

    useEffect(() => {
        const fetchProfile = async () => {
            const studentId = Cookies.get("student_id")
            if (!studentId) {
                setError("No student ID found. Please login again.")
                setLoading(false)
                return
            }

            const token = Cookies.get("token")

            try {
                const res = await fetch(`http://localhost:8080/student/profile?student_id=${studentId}`, {
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
            <h1 className="text-3xl font-bold">Profile</h1>

            <Card>
                <CardHeader>
                    <CardTitle>Personal Information</CardTitle>
                    <CardDescription>View and update your personal details.</CardDescription>
                </CardHeader>
                <CardContent className="grid gap-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div className="grid gap-2">
                            <Label>First Name</Label>
                            <Input value={profile?.first_name || ""} disabled />
                        </div>
                        <div className="grid gap-2">
                            <Label>Last Name</Label>
                            <Input value={profile?.last_name || ""} disabled />
                        </div>
                    </div>
                    <div className="grid gap-2">
                        <Label>Email</Label>
                        <Input value={profile?.email || ""} disabled />
                    </div>
                    <div className="grid gap-2">
                        <Label>Student ID</Label>
                        <Input value={profile?.student_id || ""} disabled />
                    </div>
                    <div className="grid gap-2">
                        <Label>Phone</Label>
                        <Input value={profile?.phone_no || ""} disabled />
                    </div>
                </CardContent>
                <CardFooter className="justify-end border-t px-6 py-4">
                    <Button variant="outline">Edit (Disabled)</Button>
                </CardFooter>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Academic Information</CardTitle>
                    <CardDescription>Your current academic standing.</CardDescription>
                </CardHeader>
                <CardContent className="grid gap-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div className="grid gap-2">
                            <Label>Department</Label>
                            <div className="text-sm font-medium">{profile?.dept_id}</div>
                        </div>
                        <div className="grid gap-2">
                            <Label>Semester</Label>
                            <div className="text-sm font-medium">{profile?.semester}th Semester</div>
                        </div>
                    </div>
                    <div className="grid gap-2">
                        <Label>Year of Joining</Label>
                        <div className="text-sm font-medium">{profile?.year_of_joining}</div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
