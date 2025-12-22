import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"

export default function RecommendationsPage() {
    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold tracking-tight">Recommendations</h1>
                <p className="text-muted-foreground">Personalized for your academic growth</p>
            </div>

            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                <Card>
                    <CardHeader>
                        <div className="flex justify-between items-start">
                            <CardTitle>Advanced Database Systems</CardTitle>
                            <div className="bg-primary/10 text-primary px-2 py-0.5 rounded text-xs font-semibold">Course</div>
                        </div>
                        <CardDescription>Based on your interest in DBMS</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p className="text-sm text-muted-foreground">Explore NoSQL databases, distributed systems, and advanced indexing techniques. This elective complements your current DBMS coursework.</p>
                    </CardContent>
                    <CardFooter>
                        <Button variant="outline" className="w-full">View Syllabus</Button>
                    </CardFooter>
                </Card>

                <Card>
                    <CardHeader>
                        <div className="flex justify-between items-start">
                            <CardTitle>Machine Learning for Beginners</CardTitle>
                            <div className="bg-primary/10 text-primary px-2 py-0.5 rounded text-xs font-semibold">Course</div>
                        </div>
                        <CardDescription>Prerequisite for final year project</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p className="text-sm text-muted-foreground">Start early with ML foundations. Taking this now will help you in your upcoming capstone project selection.</p>
                    </CardContent>
                    <CardFooter>
                        <Button variant="outline" className="w-full">Enroll Now</Button>
                    </CardFooter>
                </Card>

                <Card>
                    <CardHeader>
                        <div className="flex justify-between items-start">
                            <CardTitle>Study Plan for Finals</CardTitle>
                            <div className="bg-green-100 text-green-700 px-2 py-0.5 rounded text-xs font-semibold">Advice</div>
                        </div>
                        <CardDescription>Generated based on your schedule</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <p className="text-sm text-muted-foreground">You have gaps on Wednesday. Utilize 11:00 AM - 1:00 PM for Computer Networks revision to improve your grade.</p>
                    </CardContent>
                    <CardFooter>
                        <Button className="w-full">Add to Calendar</Button>
                    </CardFooter>
                </Card>
            </div>
        </div>
    )
}
