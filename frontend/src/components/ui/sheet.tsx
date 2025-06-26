"use client"

import * as React from "react"
import { X } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

interface SheetProps {
  children: React.ReactNode
  defaultOpen?: boolean
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

interface SheetContextType {
  open: boolean
  onOpenChange: (open: boolean) => void
}

const SheetContext = React.createContext<SheetContextType | undefined>(undefined)

const useSheet = () => {
  const context = React.useContext(SheetContext)
  if (!context) {
    throw new Error("useSheet must be used within a Sheet")
  }
  return context
}

const Sheet = React.forwardRef<HTMLDivElement, SheetProps>(
  ({ children, defaultOpen = false, open, onOpenChange, ...props }, ref) => {
    const [internalOpen, setInternalOpen] = React.useState(defaultOpen)
    
    const isControlled = open !== undefined
    const isOpen = isControlled ? open : internalOpen
    
    const handleOpenChange = React.useCallback((newOpen: boolean) => {
      if (isControlled) {
        onOpenChange?.(newOpen)
      } else {
        setInternalOpen(newOpen)
        onOpenChange?.(newOpen)
      }
    }, [isControlled, onOpenChange])

    return (
      <SheetContext.Provider value={{ open: isOpen, onOpenChange: handleOpenChange }}>
        <div ref={ref} {...props}>
          {children}
        </div>
      </SheetContext.Provider>
    )
  }
)
Sheet.displayName = "Sheet"

const SheetTrigger = React.forwardRef<
  HTMLButtonElement,
  React.ComponentPropsWithoutRef<"button">
>(({ className, children, ...props }, ref) => {
  const { onOpenChange } = useSheet()
  
  return (
    <button
      ref={ref}
      className={cn("", className)}
      onClick={() => onOpenChange(true)}
      {...props}
    >
      {children}
    </button>
  )
})
SheetTrigger.displayName = "SheetTrigger"

interface SheetContentProps extends React.ComponentPropsWithoutRef<"div"> {
  side?: "left" | "right" | "top" | "bottom"
}

const SheetContent = React.forwardRef<HTMLDivElement, SheetContentProps>(
  ({ className, children, side = "right", ...props }, ref) => {
    const { open, onOpenChange } = useSheet()
    
    if (!open) return null
    
    const sideClasses = {
      left: "left-0 top-0 h-full w-80 transform transition-transform duration-300",
      right: "right-0 top-0 h-full w-80 transform transition-transform duration-300", 
      top: "top-0 left-0 w-full h-80 transform transition-transform duration-300",
      bottom: "bottom-0 left-0 w-full h-80 transform transition-transform duration-300"
    }
    
    return (
      <>
        {/* Overlay */}
        <div 
          className="fixed inset-0 bg-black/50 z-50"
          onClick={() => onOpenChange(false)}
        />
        
        {/* Sheet Content */}
        <div
          ref={ref}
          className={cn(
            "fixed z-50 bg-background border shadow-lg",
            sideClasses[side],
            className
          )}
          {...props}
        >
          <div className="flex flex-col h-full">
            <div className="flex items-center justify-end p-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onOpenChange(false)}
                className="h-8 w-8 p-0"
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
            <div className="flex-1 overflow-y-auto p-4">
              {children}
            </div>
          </div>
        </div>
      </>
    )
  }
)
SheetContent.displayName = "SheetContent"

const SheetHeader = React.forwardRef<
  HTMLDivElement,
  React.ComponentPropsWithoutRef<"div">
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("flex flex-col space-y-2 text-center sm:text-left", className)}
    {...props}
  />
))
SheetHeader.displayName = "SheetHeader"

const SheetTitle = React.forwardRef<
  HTMLDivElement,
  React.ComponentPropsWithoutRef<"div">
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("text-lg font-semibold text-foreground", className)}
    {...props}
  />
))
SheetTitle.displayName = "SheetTitle"

export { Sheet, SheetTrigger, SheetContent, SheetHeader, SheetTitle }
