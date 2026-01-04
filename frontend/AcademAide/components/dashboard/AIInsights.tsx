
"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { BrainCircuit, Sparkles } from "lucide-react"
import { fetchAIInsights, Risk, Suggestion } from "@/lib/ai-insights"
import { RiskAlert } from "./RiskAlert"
import { WhatIfSimulator } from "./WhatIfSimulator"
import Cookies from "js-cookie"

export function AIInsights() {
    const [risks, setRisks] = useState<Risk[]>([])
    const [suggestions, setSuggestions] = useState<Suggestion[]>([])

    useEffect(() => {
        const loadInsights = async () => {
            const studentId = Cookies.get("student_id") || "S1001" // Default for testing if cookie missing
            const data = await fetchAIInsights(studentId)
            setRisks(data.risks)
            setSuggestions(data.suggestions)
        }
        loadInsights()
    }, [])

    return (
        <Card className="glass border-none shadow-lg relative overflow-hidden">
            <div className="absolute top-0 right-0 p-4 opacity-10">
                <BrainCircuit className="h-32 w-32" />
            </div>

            <CardHeader>
                <div className="flex items-center gap-2">
                    <Sparkles className="h-5 w-5 text-violet-600 dark:text-violet-400" />
                    <CardTitle className="text-xl gradient-text">Academic Intelligence</CardTitle>
                </div>
            </CardHeader>

            <CardContent className="grid gap-6 lg:grid-cols-2">
                <div className="space-y-4">
                    <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider">Analysis & Alerts</h3>
                    {risks.length === 0 ? (
                        <div className="p-4 rounded-lg bg-green-50 dark:bg-green-900/20 text-green-700 dark:text-green-300">
                            No immediate risks detected. Keep up the good work! 🌟
                        </div>
                    ) : (
                        risks.map((risk, idx) => {
                            // Find matching suggestion
                            const suggestion = suggestions.find(s =>
                                s.reason.includes(risk.subject || "") ||
                                s.suggestion.includes(risk.subject || "")
                            )
                            return (
                                <RiskAlert
                                    key={idx}
                                    {...risk}
                                    suggestion={suggestion?.suggestion}
                                    reason={suggestion?.reason}
                                />
                            )
                        })
                    )}
                </div>

                <div className="border-l pl-0 lg:pl-6 border-dashed border-gray-200 dark:border-gray-800">
                    <WhatIfSimulator />
                </div>
            </CardContent>
        </Card>
    )
}
