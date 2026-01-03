"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import { LayoutDashboard, MessageSquare, Calendar, BookOpen, Users, Settings, LogOut } from "lucide-react"
import { Button } from "@/components/ui/Button"

const sidebarItems = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: LayoutDashboard,
    },
    {
        title: "Chat Assistant",
        href: "/chat",
        icon: MessageSquare,
    },
    {
        title: "Timetable",
        href: "/timetable",
        icon: Calendar,
    },
    {
        title: "Recommendations",
        href: "/recommendations",
        icon: BookOpen,
    },
    {
        title: "Profile",
        href: "/profile",
        icon: Users,
    },
]

import { useRouter } from "next/navigation"
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

export function Sidebar() {
    const pathname = usePathname()
    const router = useRouter()

    const handleLogout = () => {
        Cookies.remove("token")
        router.push("/login")
        router.refresh()
    }

    return (
        <div className="flex h-full w-64 flex-col border-r bg-white/50 backdrop-blur-xl dark:bg-slate-950 dark:border-white/10 text-card-foreground">
            <div className="flex h-14 items-center border-b border-gray-200/50 dark:border-white/10 px-4">
                <Link href="/dashboard" className="flex items-center gap-2 font-semibold tracking-tight">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                        <BookOpen className="h-5 w-5" />
                    </div>
                    <span>AcademAide</span>
                </Link>
            </div>
            <div className="flex-1 overflow-auto py-4">
                <nav className="grid items-start px-2 text-sm font-medium lg:px-4 space-y-1">
                    {sidebarItems.map((item) => (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={cn(
                                "flex items-center gap-3 rounded-full px-4 py-3 transition-all duration-200",
                                pathname === item.href
                                    ? "bg-gradient-to-r from-indigo-500 to-violet-500 text-white shadow-md shadow-indigo-500/20 font-semibold"
                                    : "text-slate-600 dark:text-white hover:bg-white hover:text-indigo-600 dark:hover:bg-white/20 dark:hover:text-white"
                            )}
                        >
                            <item.icon className={cn("h-4 w-4", pathname === item.href ? "text-white" : "")} />
                            {item.title}
                        </Link>
                    ))}
                </nav>
            </div>
            <div className="border-t border-gray-200/50 dark:border-white/10 p-4">
                <Dialog>
                    <DialogTrigger asChild>
                        <Button
                            variant="ghost"
                            className="w-full justify-start gap-2 text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/10"
                        >
                            <LogOut className="h-4 w-4" />
                            Logout
                        </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-md">
                        <DialogHeader>
                            <DialogTitle>Sign out</DialogTitle>
                            <DialogDescription>
                                Are you sure you want to sign out of your account?
                            </DialogDescription>
                        </DialogHeader>
                        <DialogFooter className="flex flex-col sm:flex-row gap-2">
                            <DialogClose asChild>
                                <Button variant="outline" className="sm:flex-1">Cancel</Button>
                            </DialogClose>
                            <Button
                                variant="destructive"
                                className="sm:flex-1 bg-red-600 hover:bg-red-700"
                                onClick={handleLogout}
                            >
                                Sign out
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>
        </div>
    )
}
