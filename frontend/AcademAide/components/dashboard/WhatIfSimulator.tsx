
import { useEffect, useState } from "react"
import { simulateWhatIf, WhatIfScenario } from "@/lib/ai-insights"
import { Slider } from "@/components/ui/Slider"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { ArrowRight, Calculator, Loader2 } from "lucide-react"
import Cookies from "js-cookie"

export function WhatIfSimulator() {
    const [missedClasses, setMissedClasses] = useState(0)
    const [simulation, setSimulation] = useState<WhatIfScenario | null>(null)
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        const fetchSimulation = async () => {
            setLoading(true)
            const studentId = Cookies.get("student_id") || "S1001"
            const result = await simulateWhatIf(studentId, missedClasses)
            setSimulation(result)
            setLoading(false)
        }

        // Simple debounce
        const timer = setTimeout(() => {
            fetchSimulation()
        }, 300)

        return () => clearTimeout(timer)
    }, [missedClasses])

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2 mb-2">
                <Calculator className="h-5 w-5 text-primary" />
                <h3 className="font-semibold">What-If Simulator</h3>
            </div>

            <div className="p-4 rounded-lg border bg-card text-card-foreground shadow-sm">
                <p className="text-sm font-medium mb-4">Scenario: Missing next {missedClasses} classes</p>
                <Slider
                    defaultValue={[0]}
                    max={10}
                    step={1}
                    value={[missedClasses]}
                    onValueChange={(vals) => setMissedClasses(vals[0])}
                    className="mb-6"
                />

                <div className="flex items-center justify-between p-3 bg-muted/50 rounded-md">
                    <div className="text-center">
                        <p className="text-xs text-muted-foreground">Current</p>
                        <p className="font-bold text-lg text-emerald-600">85.0%</p>
                    </div>
                    <ArrowRight className="h-4 w-4 text-muted-foreground" />
                    <div className="text-center">
                        <p className="text-xs text-muted-foreground">Projected</p>
                        {loading || !simulation ? (
                            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground mx-auto" />
                        ) : (
                            <p className={`font-bold text-lg ${simulation.risk_level === 'High' ? 'text-red-500' : simulation.risk_level === 'Medium' ? 'text-amber-500' : 'text-emerald-500'}`}>
                                {simulation.projected_attendance}%
                            </p>
                        )}
                    </div>
                </div>

                {simulation && simulation.percentage_drop > 0 && (
                    <p className="text-xs text-center mt-2 text-destructive">
                        Attendance will drop by {simulation.percentage_drop}%
                    </p>
                )}
            </div>
        </div>
    )
}
