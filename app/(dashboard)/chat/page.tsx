"use client"

import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Send, Loader2 } from "lucide-react"
import { useState, useRef, useEffect } from "react"
import { cn } from "@/lib/utils"
import Cookies from "js-cookie"

type Message = {
    id: string
    role: "user" | "assistant"
    content: string
    timestamp: Date
}

export default function ChatPage() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            role: "assistant",
            content: "Hello! I'm AcademAide, your personal academic assistant. How can I help you today?",
            timestamp: new Date(),
        },
    ])
    const [input, setInput] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages])

    const handleSendMessage = async (e: React.FormEvent) => {
        e.preventDefault()
        if (!input.trim() || isLoading) return

        const studentId = Cookies.get("student_id")
        if (!studentId) {
            alert("Please login first")
            return
        }

        const newMessage: Message = {
            id: Date.now().toString(),
            role: "user",
            content: input,
            timestamp: new Date(),
        }

        setMessages((prev) => [...prev, newMessage])
        setInput("")
        setIsLoading(true)

        try {
            const res = await fetch("http://localhost:8000/chat/message", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ student_id: studentId, message: newMessage.content }),
            })

            if (!res.ok) throw new Error("Processing failed")

            const data = await res.json()

            const responseMessage: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: data.response || "Sorry, I understand.",
                timestamp: new Date(),
            }
            setMessages((prev) => [...prev, responseMessage])
        } catch (error) {
            const errorMessage: Message = {
                id: (Date.now() + 1).toString(),
                role: "assistant",
                content: "Sorry, something went wrong. Please try again.",
                timestamp: new Date(),
            }
            setMessages((prev) => [...prev, errorMessage])
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="flex flex-col h-[calc(100vh-8rem)]">
            <div className="flex-1 overflow-y-auto mb-4 space-y-4 p-4 rounded-lg border bg-card/50">
                {messages.map((message) => (
                    <div
                        key={message.id}
                        className={cn(
                            "flex w-max max-w-[80%] flex-col gap-2 rounded-2xl px-4 py-3 text-sm shadow-sm",
                            message.role === "user"
                                ? "ml-auto bg-gradient-to-br from-indigo-500 to-violet-600 text-white rounded-br-sm"
                                : "bg-white text-gray-800 dark:bg-slate-800 dark:text-gray-100 rounded-bl-sm border border-gray-100 dark:border-white/10"
                        )}
                    >
                        {message.content}
                    </div>
                ))}
                {isLoading && (
                    <div className="bg-white dark:bg-slate-800 border border-gray-100 dark:border-white/10 w-max rounded-2xl rounded-bl-sm px-4 py-3 text-sm flex items-center gap-2 text-muted-foreground">
                        <Loader2 className="h-3 w-3 animate-spin text-indigo-500" />
                        <span>Thinking...</span>
                    </div>
                )}
                <div ref={messagesEndRef} />
            </div>
            <form onSubmit={handleSendMessage} className="flex gap-2">
                <Input
                    placeholder="Type your message..."
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    disabled={isLoading}
                    className="flex-1 rounded-full border-gray-200 bg-white shadow-sm focus-visible:ring-indigo-500 dark:bg-slate-900/50 dark:border-white/10 pl-6"
                />
                <Button
                    type="submit"
                    size="icon"
                    disabled={isLoading}
                    className="rounded-full h-10 w-10 bg-indigo-600 hover:bg-indigo-700 shadow-md shadow-indigo-500/20"
                >
                    <Send className="h-4 w-4" />
                    <span className="sr-only">Send</span>
                </Button>
            </form>
        </div>
    )
}
