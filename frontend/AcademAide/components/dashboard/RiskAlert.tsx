
import { AlertCircle, AlertTriangle, Info } from "lucide-react"

interface RiskAlertProps {
    type: "Attendance" | "Grades" | "General";
    severity: "High" | "Medium" | "Low";
    message: string;
    suggestion?: string;
    reason?: string;
}

export function RiskAlert({ severity, message, suggestion, reason }: RiskAlertProps) {
    let bgColor = "bg-blue-50 dark:bg-blue-900/20";
    let textColor = "text-blue-800 dark:text-blue-200";
    let iconColor = "text-blue-500";
    let Icon = Info;

    if (severity === "High") {
        Icon = AlertCircle;
        bgColor = "bg-red-50 dark:bg-red-900/20";
        textColor = "text-red-800 dark:text-red-200";
        iconColor = "text-red-500";
    } else if (severity === "Medium") {
        Icon = AlertTriangle;
        bgColor = "bg-amber-50 dark:bg-amber-900/20";
        textColor = "text-amber-800 dark:text-amber-200";
        iconColor = "text-amber-500";
    }

    return (
        <div className={`p-4 rounded-lg border ${bgColor} ${textColor} flex gap-3 items-start`}>
            <Icon className={`h-5 w-5 mt-0.5 ${iconColor}`} />
            <div className="flex-1">
                <h5 className="font-semibold text-sm leading-none tracking-tight mb-1">{message}</h5>
                {suggestion && (
                    <div className="mt-2 text-sm bg-white/50 dark:bg-black/20 p-2 rounded">
                        <p className="font-medium">💡 Recommendation:</p>
                        <p>{suggestion}</p>
                        {reason && (
                            <p className="text-xs opacity-80 mt-1 italic">Why? {reason}</p>
                        )}
                    </div>
                )}
            </div>
        </div>
    )
}
