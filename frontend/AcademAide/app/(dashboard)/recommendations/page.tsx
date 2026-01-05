"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import Cookies from "js-cookie"
import { BookOpen, Video, FileText, Download, ExternalLink } from "lucide-react"

interface Resource {
    resource_id: number
    title: string
    description: string
    type: string
    course_id: string
    link: string
}

export default function RecommendationsPage() {
    const [resources, setResources] = useState<Resource[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchResources = async () => {
            const token = Cookies.get("token")
            try {
                const res = await fetch("http://localhost:8080/student/resources", {
                    headers: { "Authorization": `Bearer ${token}` }
                })
                if (res.ok) {
                    const data = await res.json()
                    setResources(data || [])
                }
            } catch (err) {
                console.error(err)
            } finally {
                setLoading(false)
            }
        }
        fetchResources()
    }, [])

    const getIcon = (type: string) => {
        switch (type.toLowerCase()) {
            case 'video': return <Video className="h-4 w-4" />
            case 'pdf': return <FileText className="h-4 w-4" />
            default: return <BookOpen className="h-4 w-4" />
        }
    }

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <h1 className="text-3xl font-bold tracking-tight">Recommendations</h1>
                <p className="text-muted-foreground">Curated materials for your enrolled courses</p>
            </div>

            {loading ? (
                <div className="flex items-center justify-center p-12">
                    <div className="animate-pulse text-muted-foreground">Loading personalized resources...</div>
                </div>
            ) : resources.length === 0 ? (
                <div className="text-center p-10 text-muted-foreground border-2 border-dashed rounded-lg">
                    No active recommendations found for your current enrollment.
                </div>
            ) : (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    {resources.map((res) => (
                        <Card key={res.resource_id} className="flex flex-col glass hover:shadow-lg transition-all hover:-translate-y-1">
                            <CardHeader>
                                <div className="flex justify-between items-start">
                                    <div className="space-y-1">
                                        <CardTitle className="text-lg leading-tight">{res.title}</CardTitle>
                                        <div className="text-xs font-semibold text-violet-600 bg-violet-100 dark:bg-violet-900/30 dark:text-violet-300 px-2 py-0.5 rounded w-fit">
                                            {res.course_id}
                                        </div>
                                    </div>
                                    <div className="text-muted-foreground p-2 bg-secondary/50 rounded-full">
                                        {getIcon(res.type)}
                                    </div>
                                </div>
                                <CardDescription className="capitalize">{res.type} Resource</CardDescription>
                            </CardHeader>
                            <CardContent className="flex-1">
                                <p className="text-sm text-muted-foreground line-clamp-3">{res.description}</p>
                            </CardContent>
                            <CardFooter>
                                <Button className="w-full gap-2" variant="outline" onClick={() => window.open(res.link, '_blank')}>
                                    <ExternalLink className="h-4 w-4" /> Access Material
                                </Button>
                            </CardFooter>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    )
}
