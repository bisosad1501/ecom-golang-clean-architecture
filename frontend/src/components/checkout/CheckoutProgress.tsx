'use client'

import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Step {
  id: number
  title: string
  description: string
}

interface CheckoutProgressProps {
  steps: Step[]
  currentStep: number
  className?: string
}

export function CheckoutProgress({ steps, currentStep, className }: CheckoutProgressProps) {
  return (
    <div className={cn("w-full", className)}>
      <div className="flex items-center justify-between">
        {steps.map((step, index) => (
          <div key={step.id} className="flex items-center">
            {/* Step Circle */}
            <div className="flex flex-col items-center">
              <div
                className={cn(
                  "w-10 h-10 rounded-full border-2 flex items-center justify-center transition-all duration-300",
                  step.id < currentStep
                    ? "bg-[#ff9000] border-[#ff9000] text-white"
                    : step.id === currentStep
                    ? "border-[#ff9000] text-[#ff9000] bg-[#ff9000]/10"
                    : "border-gray-600 text-gray-400 bg-gray-800/50"
                )}
              >
                {step.id < currentStep ? (
                  <Check className="h-5 w-5" />
                ) : (
                  <span className="text-sm font-semibold">{step.id}</span>
                )}
              </div>
              
              {/* Step Info */}
              <div className="mt-2 text-center">
                <div
                  className={cn(
                    "text-sm font-medium transition-colors duration-300",
                    step.id <= currentStep ? "text-white" : "text-gray-400"
                  )}
                >
                  {step.title}
                </div>
                <div
                  className={cn(
                    "text-xs transition-colors duration-300",
                    step.id <= currentStep ? "text-gray-300" : "text-gray-500"
                  )}
                >
                  {step.description}
                </div>
              </div>
            </div>
            
            {/* Connector Line */}
            {index < steps.length - 1 && (
              <div
                className={cn(
                  "flex-1 h-px mx-4 transition-colors duration-300",
                  step.id < currentStep ? "bg-[#ff9000]" : "bg-gray-600"
                )}
              />
            )}
          </div>
        ))}
      </div>
    </div>
  )
}

// Mobile version with simplified design
export function CheckoutProgressMobile({ steps, currentStep, className }: CheckoutProgressProps) {
  const currentStepData = steps.find(step => step.id === currentStep)
  
  return (
    <div className={cn("w-full", className)}>
      <div className="flex items-center justify-between mb-4">
        <div>
          <div className="text-sm text-gray-400">
            Step {currentStep} of {steps.length}
          </div>
          <div className="text-lg font-semibold text-white">
            {currentStepData?.title}
          </div>
          <div className="text-sm text-gray-300">
            {currentStepData?.description}
          </div>
        </div>
        
        <div className="text-right">
          <div className="text-2xl font-bold text-[#ff9000]">
            {currentStep}/{steps.length}
          </div>
        </div>
      </div>
      
      {/* Progress Bar */}
      <div className="w-full bg-gray-700 rounded-full h-2">
        <div
          className="bg-gradient-to-r from-[#ff9000] to-orange-600 h-2 rounded-full transition-all duration-500 ease-out"
          style={{ width: `${(currentStep / steps.length) * 100}%` }}
        />
      </div>
    </div>
  )
}
