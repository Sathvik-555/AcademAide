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
    reference?: string
}

interface Quiz {
    id: string
    course_id: string
    topic: string
    questions: Question[]
}

const COURSES = [
    { id: "CD252IA", name: "Database Management Systems" },
    { id: "CS354TA", name: "Theory of Computation" },
    { id: "IS353IA", name: "Artificial Intelligence & ML" },
    { id: "XX355TBX", name: "Cloud Computing" },
    { id: "HS251TA", name: "Economics & Management" },
]

export default function QuizzesPage() {
    const [courseId, setCourseId] = useState(COURSES[0].id)
    const [unitId, setUnitId] = useState<number>(0) // 0 means All Units
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
                body: JSON.stringify({
                    course_id: courseId,
                    unit: unitId // 0 for all, 1-5 for specific
                })
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

    const [analysis, setAnalysis] = useState<{ weak_areas: string[], study_priorities: { topic: string, priority: string, reason: string }[] } | null>(null)
    const [analyzing, setAnalyzing] = useState(false)

    const handleSubmit = async () => {
        if (!quiz) return
        let correctCount = 0
        const wrongQuestions: any[] = []

        quiz.questions.forEach(q => {
            if (answers[q.id] === q.correct_option) {
                correctCount++
            } else {
                wrongQuestions.push({
                    question_text: q.text,
                    correct_answer: q.options[q.correct_option],
                    user_answer: q.options[answers[q.id]] || "Skipped",
                    reference: q.reference || ""
                })
            }
        })
        setScore(correctCount)
        setSubmitted(true)

        // Trigger Analysis
        if (wrongQuestions.length > 0) {
            setAnalyzing(true)
            try {
                const token = Cookies.get("token")
                const res = await fetch("http://localhost:8080/ai/quiz-analysis", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${token}`
                    },
                    body: JSON.stringify({
                        course_id: quiz.course_id,
                        wrong_questions: wrongQuestions,
                        total_questions: quiz.questions.length,
                        score: correctCount
                    })
                })
                if (res.ok) {
                    const data = await res.json()
                    setAnalysis(data)
                }
            } catch (e) {
                console.error("Analysis failed", e)
            } finally {
                setAnalyzing(false)
            }
        }
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
                    <CardDescription>Select a course and optional unit to generate a practice quiz.</CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col md:flex-row gap-4">
                    <select
                        value={courseId}
                        onChange={(e) => setCourseId(e.target.value)}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm md:w-[300px]"
                    >
                        {COURSES.map(c => (
                            <option key={c.id} value={c.id}>{c.id} - {c.name}</option>
                        ))}
                    </select>

                    <select
                        value={unitId}
                        onChange={(e) => setUnitId(Number(e.target.value))}
                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm md:w-[150px]"
                    >
                        <option value={0}>All Units</option>
                        <option value={1}>Unit 1</option>
                        <option value={2}>Unit 2</option>
                        <option value={3}>Unit 3</option>
                        <option value={4}>Unit 4</option>
                        <option value={5}>Unit 5</option>
                    </select>

                    <Button onClick={generateQuiz} disabled={loading} className="gap-2 flex-1 md:flex-none">
                        {loading && <Loader2 className="h-4 w-4 animate-spin" />}
                        {loading ? "Generating..." : "Generate Quiz"}
                    </Button>
                </CardContent>
            </Card>

            {quiz && (
                <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-500">
                    {submitted && (
                        <div className="space-y-6">
                            <Card className="border-green-200 bg-green-50 dark:bg-green-900/20 dark:border-green-900">
                                <CardContent className="pt-6 text-center">
                                    <h2 className="text-2xl font-bold text-green-700 dark:text-green-300">
                                        Score: {score} / {quiz.questions.length}
                                    </h2>
                                    <p className="text-green-600 dark:text-green-400">
                                        {score === quiz.questions.length ? "Perfect Score! 🎉" : "Good effort! Check the analysis below."}
                                    </p>
                                </CardContent>
                            </Card>

                            {/* Analysis Section */}
                            {(analyzing || analysis) && (
                                <Card className="border-indigo-200 bg-indigo-50 dark:bg-indigo-900/20 dark:border-indigo-900">
                                    <CardHeader>
                                        <CardTitle className="flex items-center gap-2 text-indigo-800 dark:text-indigo-200">
                                            <BrainCircuit className="h-5 w-5" />
                                            AI Performance Analysis
                                        </CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        {analyzing ? (
                                            <div className="flex items-center gap-2 text-indigo-600">
                                                <Loader2 className="h-4 w-4 animate-spin" />
                                                Analyzing your mistakes...
                                            </div>
                                        ) : analysis ? (
                                            <div className="space-y-4">
                                                <div>
                                                    <h4 className="font-semibold text-indigo-900 dark:text-indigo-100 mb-2">Weak Areas Detected:</h4>
                                                    <div className="flex flex-wrap gap-2">
                                                        {analysis.weak_areas.map((area, i) => (
                                                            <span key={i} className="px-2 py-1 bg-red-100 text-red-700 dark:bg-red-900/50 dark:text-red-200 rounded-md text-sm font-medium">
                                                                {area}
                                                            </span>
                                                        ))}
                                                    </div>
                                                </div>
                                                <div>
                                                    <h4 className="font-semibold text-indigo-900 dark:text-indigo-100 mb-2">Prioritized Study Plan:</h4>
                                                    <ul className="space-y-3">
                                                        {analysis.study_priorities.map((item, i) => (
                                                            <li key={i} className="flex gap-3 items-start bg-white/50 dark:bg-black/20 p-3 rounded-lg">
                                                                <div className={`mt-1 h-2 w-2 rounded-full shrink-0 ${item.priority === "High" ? "bg-red-500" :
                                                                    item.priority === "Medium" ? "bg-amber-500" : "bg-green-500"
                                                                    }`} />
                                                                <div>
                                                                    <div className="font-medium text-slate-900 dark:text-slate-100">
                                                                        {item.topic}
                                                                        <span className="ml-2 text-xs opacity-70 uppercase tracking-wider border px-1 rounded">
                                                                            {item.priority} Priority
                                                                        </span>
                                                                    </div>
                                                                    <p className="text-sm text-slate-600 dark:text-slate-400 mt-1">{item.reason}</p>
                                                                </div>
                                                            </li>
                                                        ))}
                                                    </ul>
                                                </div>
                                            </div>
                                        ) : null}
                                    </CardContent>
                                </Card>
                            )}
                        </div>
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
