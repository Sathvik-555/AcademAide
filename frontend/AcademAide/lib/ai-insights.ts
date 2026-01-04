
export interface StudentData {
    attendance: {
        subject: string;
        attended: number;
        total: number;
    }[];
    grades: {
        subject: string;
        current: number; // Percentage or GPA equivalent
        previous: number;
    }[];
    cgpa: number;
}

export interface Risk {
    type: "Attendance" | "Grades" | "General";
    severity: "High" | "Medium" | "Low";
    message: string;
    subject?: string;
}

export interface Suggestion {
    suggestion: string;
    reason: string;
}

export interface SimulationResult {
    projectedAttendance: number;
    attendanceDrop: number;
    riskLevel: string;
}


export interface AIInsightsResponse {
    risks: Risk[];
    suggestions: Suggestion[];
}

export interface WhatIfScenario {
    projected_attendance: number;
    percentage_drop: number;
    risk_level: string;
}

// API Service Functions
const API_BASE_URL = "http://localhost:8080"; // Using local backend

export async function fetchAIInsights(studentId: string): Promise<AIInsightsResponse> {
    try {
        const res = await fetch(`${API_BASE_URL}/ai/insights?student_id=${studentId}`);
        if (!res.ok) throw new Error("Failed to fetch insights");
        return await res.json();
    } catch (error) {
        console.error("AI API Error:", error);
        return { risks: [], suggestions: [] };
    }
}

export async function simulateWhatIf(studentId: string, missedClasses: number): Promise<WhatIfScenario> {
    try {
        const res = await fetch(`${API_BASE_URL}/ai/what-if`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ student_id: studentId, missed_classes: missedClasses })
        });
        if (!res.ok) throw new Error("Simulation failed");
        return await res.json();
    } catch (error) {
        console.error("Simulation Error:", error);
        // Fallback mock to allow UI to work if backend is down
        return {
            projected_attendance: 0,
            percentage_drop: 0,
            risk_level: "Low"
        };
    }
}
