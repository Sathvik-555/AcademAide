"use client"

import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Send, Loader2, Trash2 } from "lucide-react"
import { useState, useRef, useEffect, Suspense } from "react"
import { cn } from "@/lib/utils"
import { AgentSelector, AgentType } from "@/components/AgentSelector"
import ReactMarkdown from "react-markdown"
import remarkGfm from "remark-gfm"
import Cookies from "js-cookie"
import { useSearchParams, useRouter } from "next/navigation"

type Message = {
    id: string
    role: "user" | "assistant"
    content: string
    timestamp: Date
}

function ChatContent() {
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

    const [selectedAgent, setSelectedAgent] = useState<AgentType>("general")
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const searchParams = useSearchParams()
    const router = useRouter()
    const query = searchParams.get("q")
    const queryProcessed = useRef(false)

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages])

    const sendMessage = async (text: string) => {
        if (!text.trim() || isLoading) return

        const studentId = Cookies.get("student_id")
        if (!studentId) {
            alert("Please login first")
            return
        }

        const newMessage: Message = {
            id: Date.now().toString(),
            role: "user",
            content: text,
            timestamp: new Date(),
        }

        setMessages((prev) => [...prev, newMessage])
        setIsLoading(true)

        const token = Cookies.get("token")
        if (!token) {
            alert("No auth token found")
            setIsLoading(false)
            return
        }

        try {
            const res = await fetch("http://localhost:8080/chat/message", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify({
                    student_id: studentId,
                    message: newMessage.content,
                    agent_id: selectedAgent
                }),
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

    const clearChat = async () => {
        if (!confirm("Are you sure you want to clear the chat history? This cannot be undone.")) return

        const studentId = Cookies.get("student_id")
        const token = Cookies.get("token")
        if (!studentId || !token) return

        try {
            const res = await fetch(`http://localhost:8080/chat/history?student_id=${studentId}`, {
                method: "DELETE",
                headers: {
                    "Authorization": `Bearer ${token}`
                }
            })

            if (res.ok) {
                // Reset to initial state
                setMessages([{
                    id: "1",
                    role: "assistant",
                    content: "Chat history cleared. How can I help you freshly?",
                    timestamp: new Date(),
                }])
            } else {
                alert("Failed to clear chat")
            }
        } catch (e) {
            alert("Error clearing chat")
        }
    }

    // Handle search query param
    useEffect(() => {
        if (query && !queryProcessed.current) {
            queryProcessed.current = true
            sendMessage(query)
            // Remove query param
            router.replace('/chat')
        }
    }, [query])

    const handleFormSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        sendMessage(input)
        setInput("")
    }

    return (
        <div className="flex flex-col h-[calc(100vh-8rem)]">
            <div className="flex justify-between items-center mb-2 px-1">
                <AgentSelector selectedAgent={selectedAgent} onSelect={setSelectedAgent} />
                <Button
                    variant="ghost"
                    size="icon"
                    onClick={clearChat}
                    title="Clear Chat History"
                    className="text-muted-foreground hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/30"
                >
                    <Trash2 className="h-4 w-4" />
                </Button>
            </div>
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
                        <div className={cn(
                            "prose text-sm max-w-none break-words",
                            message.role === "user"
                                ? "prose-invert"
                                : "dark:prose-invert"
                        )}>
                            <ReactMarkdown
                                remarkPlugins={[remarkGfm]}
                                components={{
                                    pre: ({ node, ...props }) => <div className="overflow-auto w-full my-2 bg-black/10 dark:bg-black/30 p-2 rounded-lg" {...props} />,
                                    code: ({ node, ...props }) => <code className="bg-black/10 dark:bg-black/30 rounded px-1" {...props} />
                                }}
                            >
                                {message.content}
                            </ReactMarkdown>
                        </div>
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
            <form onSubmit={handleFormSubmit} className="flex gap-2">
                <Input
                    placeholder="Type your message..."
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    disabled={isLoading}
                    className="flex-1 rounded-full border-gray-200 bg-white shadow-sm focus-visible:ring-indigo-500 dark:bg-slate-900/50 dark:border-white/10 pl-6 text-slate-900 dark:text-slate-100"
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

export default function ChatPage() {
    return (
        <Suspense fallback={<div>Loading chat...</div>}>
            <ChatContent />
        </Suspense>
    )
}
