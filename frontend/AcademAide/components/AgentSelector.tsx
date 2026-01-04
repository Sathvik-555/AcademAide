"use client"

import { Bot, Brain, Code, BookOpen, ListTodo, Zap, LucideIcon } from "lucide-react"
import { cn } from "@/lib/utils"

export type AgentType = "general" | "socratic" | "code_reviewer" | "research" | "exam" | "motivational"

interface AgentSelectorProps {
    selectedAgent: AgentType
    onSelect: (agent: AgentType) => void
}

interface AgentOption {
    id: AgentType
    label: string
    icon: LucideIcon
    description: string
    color: string
}

const agents: AgentOption[] = [
    {
        id: "general",
        label: "General",
        icon: Bot,
        description: "Your standard academic assistant",
        color: "text-blue-500 bg-blue-50 dark:bg-blue-950/30"
    },
    {
        id: "socratic",
        label: "Socratic",
        icon: Brain,
        description: "Guides you with questions",
        color: "text-purple-500 bg-purple-50 dark:bg-purple-950/30"
    },
    {
        id: "code_reviewer",
        label: "Code",
        icon: Code,
        description: "Reviews your code",
        color: "text-emerald-500 bg-emerald-50 dark:bg-emerald-950/30"
    },
    {
        id: "research",
        label: "Research",
        icon: BookOpen,
        description: "Finds resources & papers",
        color: "text-amber-500 bg-amber-50 dark:bg-amber-950/30"
    },
    {
        id: "exam",
        label: "Exam",
        icon: ListTodo,
        description: "Strategizes for tests",
        color: "text-red-500 bg-red-50 dark:bg-red-950/30"
    },
    {
        id: "motivational",
        label: "Coach",
        icon: Zap,
        description: "Motivates & encourages",
        color: "text-yellow-500 bg-yellow-50 dark:bg-yellow-950/30"
    },
]

export function AgentSelector({ selectedAgent, onSelect }: AgentSelectorProps) {
    return (
        <div className="flex gap-2 overflow-x-auto pb-2 mb-2 no-scrollbar">
            {agents.map((agent) => (
                <button
                    key={agent.id}
                    onClick={() => onSelect(agent.id)}
                    className={cn(
                        "flex items-center gap-2 px-3 py-1.5 rounded-full border transition-all whitespace-nowrap",
                        selectedAgent === agent.id
                            ? cn("border-transparent shadow-sm ring-1 ring-inset", agent.color)
                            : "bg-white dark:bg-slate-900 border-gray-200 dark:border-slate-800 hover:border-gray-300 dark:hover:border-slate-700 opacity-70 hover:opacity-100"
                    )}
                    title={agent.description}
                >
                    <agent.icon className="w-4 h-4" />
                    <span className="text-sm font-medium">{agent.label}</span>
                </button>
            ))}
        </div>
    )
}
