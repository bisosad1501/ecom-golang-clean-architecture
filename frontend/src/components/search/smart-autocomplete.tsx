'use client'

import React, { useState, useEffect, useRef, useCallback } from 'react'
import { Search, Clock, TrendingUp, Star, Tag, Package, Building2, Hash } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import { useDebounce } from '@/hooks/use-debounce'

interface SmartAutocompleteSuggestion {
  id: string
  type: 'product' | 'category' | 'brand' | 'query'
  value: string
  display_text: string
  entity_id?: string
  priority: number
  score: number
  is_trending: boolean
  is_personalized: boolean
  metadata?: Record<string, any>
  synonyms?: string[]
  tags?: string[]
  reason?: string
}

interface SmartAutocompleteResponse {
  suggestions: SmartAutocompleteSuggestion[]
  categories?: SmartAutocompleteSuggestion[]
  brands?: SmartAutocompleteSuggestion[]
  products?: SmartAutocompleteSuggestion[]
  queries?: SmartAutocompleteSuggestion[]
  trending?: SmartAutocompleteSuggestion[]
  popular?: SmartAutocompleteSuggestion[]
  history?: SmartAutocompleteSuggestion[]
  total: number
  has_more: boolean
  query_time_ms: number
}

interface SmartAutocompleteProps {
  value: string
  onChange: (value: string) => void
  onSelect: (suggestion: SmartAutocompleteSuggestion) => void
  placeholder?: string
  className?: string
  includeTrending?: boolean
  includePersonalized?: boolean
  includeHistory?: boolean
  includePopular?: boolean
  types?: string[]
  limit?: number
  language?: string
  region?: string
}

export function SmartAutocomplete({
  value,
  onChange,
  onSelect,
  placeholder = "Search products, categories, brands...",
  className,
  includeTrending = true,
  includePersonalized = true,
  includeHistory = true,
  includePopular = true,
  types = ['product', 'category', 'brand', 'query'],
  limit = 10,
  language = 'en',
  region,
}: SmartAutocompleteProps) {
  const [suggestions, setSuggestions] = useState<SmartAutocompleteResponse | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [isOpen, setIsOpen] = useState(false)
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const [sessionId] = useState(() => Math.random().toString(36).substring(7))
  
  const inputRef = useRef<HTMLInputElement>(null)
  const suggestionsRef = useRef<HTMLDivElement>(null)
  
  const debouncedValue = useDebounce(value, 300)

  // Fetch suggestions
  const fetchSuggestions = useCallback(async (query: string) => {
    if (!query.trim()) {
      setSuggestions(null)
      setIsOpen(false)
      return
    }

    setIsLoading(true)
    try {
      const params = new URLSearchParams({
        q: query,
        types: types.join(','),
        limit: limit.toString(),
        include_trending: includeTrending.toString(),
        include_personalized: includePersonalized.toString(),
        include_history: includeHistory.toString(),
        include_popular: includePopular.toString(),
        language,
      })

      if (region) {
        params.append('region', region)
      }

      const response = await fetch(`/api/v1/search/autocomplete/smart?${params}`, {
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (response.ok) {
        const data = await response.json()
        setSuggestions(data.data)
        setIsOpen(true)
        setSelectedIndex(-1)
      }
    } catch (error) {
      console.error('Failed to fetch suggestions:', error)
    } finally {
      setIsLoading(false)
    }
  }, [types, limit, includeTrending, includePersonalized, includeHistory, includePopular, language, region])

  // Track suggestion impression
  const trackImpression = useCallback(async (suggestion: SmartAutocompleteSuggestion) => {
    try {
      await fetch('/api/v1/search/autocomplete/track', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          entry_id: suggestion.id,
          interaction_type: 'impression',
          session_id: sessionId,
          query: value,
        }),
      })
    } catch (error) {
      console.error('Failed to track impression:', error)
    }
  }, [sessionId, value])

  // Track suggestion click
  const trackClick = useCallback(async (suggestion: SmartAutocompleteSuggestion, position: number) => {
    try {
      await fetch('/api/v1/search/autocomplete/track', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          entry_id: suggestion.id,
          interaction_type: 'click',
          session_id: sessionId,
          query: value,
          position,
        }),
      })
    } catch (error) {
      console.error('Failed to track click:', error)
    }
  }, [sessionId, value])

  // Effect to fetch suggestions when debounced value changes
  useEffect(() => {
    fetchSuggestions(debouncedValue)
  }, [debouncedValue, fetchSuggestions])

  // Track impressions when suggestions are shown
  useEffect(() => {
    if (suggestions?.suggestions) {
      suggestions.suggestions.forEach(trackImpression)
    }
  }, [suggestions, trackImpression])

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen || !suggestions?.suggestions.length) return

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setSelectedIndex(prev => 
          prev < suggestions.suggestions.length - 1 ? prev + 1 : prev
        )
        break
      case 'ArrowUp':
        e.preventDefault()
        setSelectedIndex(prev => prev > 0 ? prev - 1 : -1)
        break
      case 'Enter':
        e.preventDefault()
        if (selectedIndex >= 0 && suggestions.suggestions[selectedIndex]) {
          handleSelect(suggestions.suggestions[selectedIndex], selectedIndex)
        }
        break
      case 'Escape':
        setIsOpen(false)
        setSelectedIndex(-1)
        break
    }
  }

  // Handle suggestion selection
  const handleSelect = (suggestion: SmartAutocompleteSuggestion, position: number) => {
    trackClick(suggestion, position)
    onSelect(suggestion)
    onChange(suggestion.value)
    setIsOpen(false)
    setSelectedIndex(-1)
  }

  // Get icon for suggestion type
  const getTypeIcon = (type: string, suggestion: SmartAutocompleteSuggestion) => {
    const iconClass = "h-4 w-4"
    
    if (suggestion.is_trending) {
      return <TrendingUp className={cn(iconClass, "text-orange-500")} />
    }
    
    if (suggestion.is_personalized) {
      return <Star className={cn(iconClass, "text-blue-500")} />
    }

    switch (type) {
      case 'product':
        return <Package className={cn(iconClass, "text-green-500")} />
      case 'category':
        return <Tag className={cn(iconClass, "text-purple-500")} />
      case 'brand':
        return <Building2 className={cn(iconClass, "text-red-500")} />
      case 'query':
        return <Hash className={cn(iconClass, "text-gray-500")} />
      default:
        return <Search className={cn(iconClass, "text-gray-400")} />
    }
  }

  // Get reason badge
  const getReasonBadge = (suggestion: SmartAutocompleteSuggestion) => {
    if (suggestion.is_trending) {
      return <Badge variant="secondary" className="text-xs bg-orange-100 text-orange-700">Trending</Badge>
    }
    if (suggestion.is_personalized) {
      return <Badge variant="secondary" className="text-xs bg-blue-100 text-blue-700">For You</Badge>
    }
    if (suggestion.reason === 'popular') {
      return <Badge variant="secondary" className="text-xs bg-green-100 text-green-700">Popular</Badge>
    }
    if (suggestion.reason === 'history') {
      return <Badge variant="secondary" className="text-xs bg-gray-100 text-gray-700">Recent</Badge>
    }
    return null
  }

  return (
    <div className={cn("relative", className)}>
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          ref={inputRef}
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => value && setIsOpen(true)}
          onBlur={() => setTimeout(() => setIsOpen(false), 200)}
          placeholder={placeholder}
          className="pl-10 pr-4"
        />
        {isLoading && (
          <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
            <div className="animate-spin h-4 w-4 border-2 border-gray-300 border-t-blue-600 rounded-full"></div>
          </div>
        )}
      </div>

      {isOpen && suggestions && suggestions.suggestions.length > 0 && (
        <div
          ref={suggestionsRef}
          className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-96 overflow-y-auto"
        >
          <div className="py-2">
            {suggestions.suggestions.map((suggestion, index) => (
              <div
                key={suggestion.id}
                className={cn(
                  "flex items-center justify-between px-4 py-3 cursor-pointer transition-colors",
                  "hover:bg-gray-50",
                  selectedIndex === index && "bg-blue-50 border-l-2 border-blue-500"
                )}
                onClick={() => handleSelect(suggestion, index)}
              >
                <div className="flex items-center space-x-3 flex-1 min-w-0">
                  {getTypeIcon(suggestion.type, suggestion)}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-2">
                      <span className="text-sm font-medium text-gray-900 truncate">
                        {suggestion.display_text}
                      </span>
                      {getReasonBadge(suggestion)}
                    </div>
                    {suggestion.metadata?.price && (
                      <div className="text-xs text-gray-500">
                        ${suggestion.metadata.price}
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  {suggestion.score > 0 && (
                    <span className="text-xs text-gray-400">
                      {Math.round(suggestion.score)}
                    </span>
                  )}
                  <Clock className="h-3 w-3 text-gray-300" />
                </div>
              </div>
            ))}
          </div>
          
          {suggestions.query_time_ms && (
            <div className="px-4 py-2 border-t border-gray-100 bg-gray-50">
              <span className="text-xs text-gray-500">
                {suggestions.total} results in {suggestions.query_time_ms}ms
              </span>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
