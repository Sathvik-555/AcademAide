import { Button } from "@/components/ui/Button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"

export default function ProfilePage() {
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
                            <Input defaultValue="Sharanya" disabled />
                        </div>
                        <div className="grid gap-2">
                            <Label>Last Name</Label>
                            <Input defaultValue="Student" disabled />
                        </div>
                    </div>
                    <div className="grid gap-2">
                        <Label>Email</Label>
                        <Input defaultValue="student@university.edu" disabled />
                    </div>
                    <div className="grid gap-2">
                        <Label>Student ID</Label>
                        <Input defaultValue="2024CS001" disabled />
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
                            <div className="text-sm font-medium">Computer Science</div>
                        </div>
                        <div className="grid gap-2">
                            <Label>Semester</Label>
                            <div className="text-sm font-medium">5th Semester</div>
                        </div>
                    </div>
                    <div className="grid gap-2">
                        <Label>Current CGPA</Label>
                        <div className="text-sm font-medium">8.9</div>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
