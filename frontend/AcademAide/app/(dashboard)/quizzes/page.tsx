"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Loader2, BrainCircuit, CheckCircle, XCircle } from "lucide-react"
import Cookies from "js-cookie"

interface Question {
    id: number
    text: string
    options: string[]
    correct_option: number
}

interface Quiz {
    id: string
    course_id: string
    topic: string
    questions: Question[]
}

export default function QuizzesPage() {
    const [courseId, setCourseId] = useState("CS101")
    const [loading, setLoading] = useState(false)
    const [quiz, setQuiz] = useState<Quiz | null>(null)
    const [answers, setAnswers] = useState<Record<number, number>>({})
    const [submitted, setSubmitted] = useState(false)
    const [score, setScore] = useState(0)

    const generateQuiz = async () => {
        setLoading(true)
        setQuiz(null)
        setSubmitted(false)
        setAnswers({})

        try {
            const token = Cookies.get("token")
            const res = await fetch("http://localhost:8080/quiz/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${token}`
                },
                body: JSON.stringify({ course_id: courseId })
            })

            if (res.ok) {
                const data = await res.json()
                setQuiz(data)
            } else {
                alert("Failed to generate quiz")
            }
        } catch (error) {
            console.error(error)
            alert("Error generating quiz")
        } finally {
            setLoading(false)
        }
    }

    const handleOptionSelect = (qId: number, optionIdx: number) => {
        if (submitted) return
        setAnswers(prev => ({ ...prev, [qId]: optionIdx }))
    }

    const handleSubmit = () => {
        if (!quiz) return
        let correctCount = 0
        quiz.questions.forEach(q => {
            if (answers[q.id] === q.correct_option) {
                correctCount++
            }
        })
        setScore(correctCount)
        setSubmitted(true)
    }

    return (
        <div className="flex flex-col gap-6 max-w-4xl mx-auto pb-10">
            <div className="flex flex-col gap-2">
                <h1 className="text-3xl font-bold tracking-tight gradient-text flex items-center gap-2">
                    <BrainCircuit className="h-8 w-8 text-primary" />
                    AI Quiz Generator
                </h1>
                <p className="text-muted-foreground">Test your knowledge with personalized AI-generated quizzes.</p>
            </div>

            <Card className="glass border-none">
                <CardHeader>
                    <CardTitle>Generate New Quiz</CardTitle>
                    <CardDescription>Select a course code to instantly generate a practice quiz.</CardDescription>
                </CardHeader>
                <CardContent className="flex gap-4">
                    <input
                        type="text"
                        value={courseId}
                        onChange={(e) => setCourseId(e.target.value)}
                        placeholder="Course ID (e.g., CS101)"
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm md:w-[200px]"
                    />
                    <Button onClick={generateQuiz} disabled={loading} className="gap-2">
                        {loading && <Loader2 className="h-4 w-4 animate-spin" />}
                        {loading ? "Generating..." : "Generate Quiz"}
                    </Button>
                </CardContent>
            </Card>

            {quiz && (
                <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                    {submitted && (
                        <Card className="border-green-200 bg-green-50 dark:bg-green-900/20 dark:border-green-900">
                            <CardContent className="pt-6 text-center">
                                <h2 className="text-2xl font-bold text-green-700 dark:text-green-300">
                                    Score: {score} / {quiz.questions.length}
                                </h2>
                                <p className="text-green-600 dark:text-green-400">
                                    {score === quiz.questions.length ? "Perfect Score! 🎉" : "Good effort! Keep practicing."}
                                </p>
                            </CardContent>
                        </Card>
                    )}

                    {quiz.questions.map((q, idx) => (
                        <Card key={q.id} className="glass border-none">
                            <CardHeader>
                                <CardTitle className="text-lg">
                                    {idx + 1}. {q.text}
                                </CardTitle>
                            </CardHeader>
                            <CardContent className="grid gap-3">
                                {q.options.map((opt, optIdx) => {
                                    let variant = "outline"
                                    let className = "justify-start text-left h-auto py-3 px-4"

                                    if (submitted) {
                                        if (optIdx === q.correct_option) {
                                            className += " border-green-500 bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-300 ring-2 ring-green-500 ring-offset-2"
                                        } else if (answers[q.id] === optIdx) {
                                            className += " border-red-500 bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-300"
                                        } else {
                                            className += " opacity-50"
                                        }
                                    } else {
                                        if (answers[q.id] === optIdx) {
                                            variant = "default"
                                            className += " bg-primary text-primary-foreground"
                                        }
                                    }

                                    return (
                                        <Button
                                            key={optIdx}
                                            variant={submitted ? "outline" : (answers[q.id] === optIdx ? "default" : "outline")}
                                            className={className}
                                            onClick={() => handleOptionSelect(q.id, optIdx)}
                                        >
                                            <span className="mr-2 font-bold">{String.fromCharCode(65 + optIdx)}.</span>
                                            {opt}
                                            {submitted && optIdx === q.correct_option && <CheckCircle className="ml-auto h-4 w-4 text-green-600" />}
                                            {submitted && answers[q.id] === optIdx && optIdx !== q.correct_option && <XCircle className="ml-auto h-4 w-4 text-red-600" />}
                                        </Button>
                                    )
                                })}
                            </CardContent>
                        </Card>
                    ))}

                    {!submitted && (
                        <div className="flex justify-end">
                            <Button size="lg" onClick={handleSubmit} disabled={Object.keys(answers).length < quiz.questions.length}>
                                Submit Quiz
                            </Button>
                        </div>
                    )}
                </div>
            )}
        </div>
    )
}
