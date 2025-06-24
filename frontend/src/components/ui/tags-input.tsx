'use client'

import React from 'react'
import { FormField } from '@/components/ui/form-field'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { X, Plus } from 'lucide-react'

interface TagsInputProps {
  tags: string[]
  onTagsChange: (tags: string[]) => void
  placeholder?: string
  maxTags?: number
  error?: string
  className?: string
}

export function TagsInput({
  tags,
  onTagsChange,
  placeholder = "Add tags...",
  maxTags = 10,
  error,
  className,
}: TagsInputProps) {
  const [inputValue, setInputValue] = React.useState('')

  const addTag = () => {
    const trimmedValue = inputValue.trim()
    if (trimmedValue && !tags.includes(trimmedValue) && tags.length < maxTags) {
      onTagsChange([...tags, trimmedValue])
      setInputValue('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    onTagsChange(tags.filter(tag => tag !== tagToRemove))
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      addTag()
    } else if (e.key === 'Backspace' && !inputValue && tags.length > 0) {
      removeTag(tags[tags.length - 1])
    }
  }

  return (
    <FormField
      label="Tags"
      error={error}
      hint={`Add up to ${maxTags} tags to help categorize your product`}
      className={className}
    >
      <div className="space-y-3">
        {/* Tags display */}
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {tags.map((tag) => (
              <Badge key={tag} variant="secondary" className="flex items-center gap-1">
                {tag}
                <button
                  type="button"
                  onClick={() => removeTag(tag)}
                  className="ml-1 hover:bg-gray-200 rounded-full p-0.5"
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
            ))}
          </div>
        )}

        {/* Input for adding new tags */}
        {tags.length < maxTags && (
          <div className="flex gap-2">
            <Input
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={placeholder}
              className="flex-1"
            />
            <Button
              type="button"
              onClick={addTag}
              disabled={!inputValue.trim() || tags.includes(inputValue.trim())}
              size="sm"
              variant="outline"
            >
              <Plus className="h-4 w-4" />
            </Button>
          </div>
        )}

        {tags.length >= maxTags && (
          <p className="text-xs text-amber-600">
            Maximum {maxTags} tags reached
          </p>
        )}
      </div>
    </FormField>
  )
}
