
"use client"

import * as React from "react"
import { cn } from "@/lib/utils"

interface SliderProps extends React.InputHTMLAttributes<HTMLInputElement> {
    value?: number[]
    defaultValue?: number[]
    max?: number
    step?: number
    onValueChange?: (value: number[]) => void
}

const Slider = React.forwardRef<HTMLInputElement, SliderProps>(
    ({ className, value, defaultValue, max = 100, step = 1, onValueChange, ...props }, ref) => {

        const [localValue, setLocalValue] = React.useState(defaultValue ? defaultValue[0] : 0)

        const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
            const val = parseFloat(e.target.value)
            setLocalValue(val)
            if (onValueChange) {
                onValueChange([val])
            }
        }

        return (
            <input
                type="range"
                className={cn(
                    "w-full h-2 bg-secondary rounded-lg appearance-none cursor-pointer accent-primary",
                    className
                )}
                min={0}
                max={max}
                step={step}
                value={value ? value[0] : localValue}
                onChange={handleChange}
                ref={ref}
                {...props}
            />
        )
    }
)
Slider.displayName = "Slider"

export { Slider }
