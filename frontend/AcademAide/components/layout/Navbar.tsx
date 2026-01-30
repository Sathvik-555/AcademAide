"use client"

import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Search, Bell, PanelLeft } from "lucide-react"
import { ModeToggle } from "@/components/mode-toggle"

interface NavbarProps {
    onToggleSidebar?: () => void;
}

export function Navbar({ onToggleSidebar }: NavbarProps) {
    return (
        <header className="flex h-14 items-center gap-4 border-b border-gray-200/50 bg-white/50 backdrop-blur-xl dark:bg-slate-900/50 dark:border-white/10 px-4 lg:h-[60px] lg:px-6">
            <Button variant="ghost" size="icon" className="hidden md:flex" onClick={onToggleSidebar}>
                <PanelLeft className="h-5 w-5" />
                <span className="sr-only">Toggle Sidebar</span>
            </Button>
            <div className="w-full flex-1">
                <form onSubmit={(e) => {
                    e.preventDefault()
                    const formData = new FormData(e.currentTarget)
                    const query = formData.get("q")
                    if (query) {
                        window.location.href = `/chat?q=${encodeURIComponent(query.toString())}`
                    }
                }}>
                    <div className="relative">
                        <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                        <Input
                            name="q"
                            type="search"
                            placeholder="Search courses, faculty..."
                            className="w-full appearance-none bg-white/50 pl-8 shadow-none md:w-2/3 lg:w-1/3 dark:bg-slate-950/50 dark:text-slate-100"
                        />
                    </div>
                </form>
            </div>
            <Button variant="ghost" size="icon" className="h-8 w-8 rounded-full">
                <Bell className="h-4 w-4" />
                <span className="sr-only">Toggle notifications</span>
            </Button>
            <ModeToggle />
            <Button variant="ghost" size="icon" className="rounded-full">
                {/* Placeholder for user avatar */}
                <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold">
                    S
                </div>
                <span className="sr-only">Toggle user menu</span>
            </Button>
        </header>
    )
}
