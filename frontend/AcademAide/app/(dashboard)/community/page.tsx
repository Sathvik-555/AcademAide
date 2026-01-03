"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Users2, Search, Plus, UserPlus, BookOpen } from "lucide-react"
import Cookies from "js-cookie"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogClose
} from "@/components/ui/Dialog"

interface Student {
    student_id: string
    first_name: string
    last_name: string
    email: string
}

interface StudyGroup {
    group_id: number
    course_id: string
    group_name: string
    description: string
    created_by: string
    created_at: string
    member_count: number
}

export default function CommunityPage() {
    const [activeTab, setActiveTab] = useState<"groups" | "peers">("groups")
    const [courseId, setCourseId] = useState("CS101")
    const [loading, setLoading] = useState(false)

    // Data
    const [peers, setPeers] = useState<Student[]>([])
    const [groups, setGroups] = useState<StudyGroup[]>([])

    // Create Group Form
    const [newGroupName, setNewGroupName] = useState("")
    const [newGroupDesc, setNewGroupDesc] = useState("")
    const [isCreateOpen, setIsCreateOpen] = useState(false)

    const fetchPeers = async () => {
        setLoading(true)
        try {
            const token = Cookies.get("token")
            const res = await fetch(`http://localhost:8080/groups/peers?course_id=${courseId}`, {
                headers: { "Authorization": `Bearer ${token}` }
            })
            if (res.ok) {
                const data = await res.json()
                setPeers(data || [])
            }
        } catch (e) { console.error(e) }
        finally { setLoading(false) }
    }

    const fetchGroups = async () => {
        setLoading(true)
        try {
            const token = Cookies.get("token")
            const res = await fetch(`http://localhost:8080/groups/list?course_id=${courseId}`, {
                headers: { "Authorization": `Bearer ${token}` }
            })
            if (res.ok) {
                const data = await res.json()
                setGroups(data || [])
            }
        } catch (e) { console.error(e) }
        finally { setLoading(false) }
    }

    const createGroup = async () => {
        try {
            const token = Cookies.get("token")
            const res = await fetch(`http://localhost:8080/groups/create`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify({
                    course_id: courseId,
                    group_name: newGroupName,
                    description: newGroupDesc
                })
            })
            if (res.ok) {
                setIsCreateOpen(false)
                fetchGroups() // Refresh
                setNewGroupName("")
                setNewGroupDesc("")
            } else {
                alert("Failed to create group")
            }
        } catch (e) {
            console.error(e)
            alert("Error creating group")
        }
    }

    // Initial Fetch when tab or course changes
    useEffect(() => {
        if (activeTab === "groups") fetchGroups()
        else fetchPeers()
    }, [activeTab, courseId])

    return (
        <div className="flex flex-col gap-6 max-w-6xl mx-auto pb-10">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight gradient-text flex items-center gap-2">
                    <Users2 className="h-8 w-8 text-primary" />
                    Community Hub
                </h1>
                <p className="text-muted-foreground">Connect with peers and join study groups for your courses.</p>
            </div>

            {/* Controls */}
            <div className="flex flex-col md:flex-row gap-4 items-center justify-between">
                <div className="flex bg-muted p-1 rounded-lg">
                    <button
                        onClick={() => setActiveTab("groups")}
                        className={`px-4 py-2 rounded-md text-sm font-medium transition-all ${activeTab === "groups" ? "bg-white dark:bg-slate-800 shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"}`}
                    >
                        Study Groups
                    </button>
                    <button
                        onClick={() => setActiveTab("peers")}
                        className={`px-4 py-2 rounded-md text-sm font-medium transition-all ${activeTab === "peers" ? "bg-white dark:bg-slate-800 shadow-sm text-foreground" : "text-muted-foreground hover:text-foreground"}`}
                    >
                        Find Peers
                    </button>
                </div>

                <div className="flex gap-2 w-full md:w-auto">
                    <div className="relative flex-1 md:w-64">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                        <input
                            value={courseId}
                            onChange={(e) => setCourseId(e.target.value)}
                            className="pl-9 flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                            placeholder="Filter by Course ID..."
                        />
                    </div>
                    {activeTab === "groups" && (
                        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
                            <DialogTrigger asChild>
                                <Button className="gap-2 shrink-0">
                                    <Plus className="h-4 w-4" /> Create Group
                                </Button>
                            </DialogTrigger>
                            <DialogContent>
                                <DialogHeader>
                                    <DialogTitle>Create Study Group</DialogTitle>
                                    <DialogDescription>Start a new learning circle for {courseId}.</DialogDescription>
                                </DialogHeader>
                                <div className="grid gap-4 py-4">
                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium">Group Name</label>
                                        <input
                                            value={newGroupName}
                                            onChange={(e) => setNewGroupName(e.target.value)}
                                            className="flex h-10 w-full rounded-md border border-input bg-background px-3 text-sm"
                                            placeholder="e.g., Exam Prep Squad"
                                        />
                                    </div>
                                    <div className="grid gap-2">
                                        <label className="text-sm font-medium">Description</label>
                                        <textarea
                                            value={newGroupDesc}
                                            onChange={(e) => setNewGroupDesc(e.target.value)}
                                            className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                                            placeholder="What will you study?"
                                        />
                                    </div>
                                </div>
                                <DialogFooter>
                                    <Button onClick={createGroup}>Create Group</Button>
                                </DialogFooter>
                            </DialogContent>
                        </Dialog>
                    )}
                </div>
            </div>

            {/* Content Area */}
            {loading ? (
                <div className="h-64 flex items-center justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                </div>
            ) : (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    {activeTab === "groups" ? (
                        groups.length > 0 ? groups.map(group => (
                            <Card key={group.group_id} className="glass border-none hover:shadow-lg transition-all">
                                <CardHeader>
                                    <div className="flex justify-between items-start">
                                        <div className="px-2.5 py-0.5 rounded-full bg-primary/10 text-primary text-xs font-bold inline-block mb-2">
                                            {group.course_id}
                                        </div>
                                        <span className="text-xs text-muted-foreground">{new Date(group.created_at).toLocaleDateString()}</span>
                                    </div>
                                    <CardTitle>{group.group_name}</CardTitle>
                                    <CardDescription className="line-clamp-2">{group.description}</CardDescription>
                                </CardHeader>
                                <CardContent>
                                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                        <Users2 className="h-4 w-4" />
                                        {group.member_count} Members
                                    </div>
                                </CardContent>
                                <CardFooter>
                                    <Button variant="outline" className="w-full">Join Group</Button>
                                </CardFooter>
                            </Card>
                        )) : (
                            <div className="col-span-full h-64 flex flex-col items-center justify-center text-muted-foreground border-2 border-dashed rounded-lg">
                                <Users2 className="h-10 w-10 mb-4 opacity-20" />
                                <p>No study groups found for {courseId}.</p>
                                <Button variant="link" onClick={() => setIsCreateOpen(true)}>Create the first one!</Button>
                            </div>
                        )
                    ) : (
                        peers.length > 0 ? peers.map(peer => (
                            <Card key={peer.student_id} className="glass border-none flex flex-row items-center p-4 gap-4">
                                <div className="h-12 w-12 rounded-full bg-gradient-to-br from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold text-lg shrink-0">
                                    {peer.first_name[0]}{peer.last_name[0]}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <h3 className="font-semibold truncate">{peer.first_name} {peer.last_name}</h3>
                                    <p className="text-sm text-muted-foreground truncate">{peer.email}</p>
                                    <div className="flex items-center gap-1 mt-1 text-xs text-muted-foreground">
                                        <BookOpen className="h-3 w-3" />
                                        <span>Enrolled in {courseId}</span>
                                    </div>
                                </div>
                                <Button size="icon" variant="ghost">
                                    <UserPlus className="h-4 w-4" />
                                </Button>
                            </Card>
                        )) : (
                            <div className="col-span-full h-64 flex flex-col items-center justify-center text-muted-foreground border-2 border-dashed rounded-lg">
                                <Search className="h-10 w-10 mb-4 opacity-20" />
                                <p>No peers found for {courseId}.</p>
                            </div>
                        )
                    )}
                </div>
            )}
        </div>
    )
}
