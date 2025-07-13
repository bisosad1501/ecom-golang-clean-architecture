import { useState, useCallback, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'

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

interface UseSmartAutocompleteOptions {
  types?: string[]
  limit?: number
  includeTrending?: boolean
  includePersonalized?: boolean
  includeHistory?: boolean
  includePopular?: boolean
  language?: string
  region?: string
  debounceMs?: number
  enabled?: boolean
}

interface UseSmartAutocompleteResult {
  suggestions: SmartAutocompleteResponse | null
  isLoading: boolean
  error: Error | null
  trackInteraction: (suggestion: SmartAutocompleteSuggestion, type: 'click' | 'impression', position?: number) => Promise<void>
  refetch: () => void
}

export function useSmartAutocomplete(
  query: string,
  options: UseSmartAutocompleteOptions = {}
): UseSmartAutocompleteResult {
  const {
    types = ['product', 'category', 'brand', 'query'],
    limit = 10,
    includeTrending = true,
    includePersonalized = true,
    includeHistory = true,
    includePopular = true,
    language = 'en',
    region,
    enabled = true,
  } = options

  const sessionIdRef = useRef(Math.random().toString(36).substring(7))

  // Build query parameters
  const buildParams = useCallback(() => {
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

    return params.toString()
  }, [query, types, limit, includeTrending, includePersonalized, includeHistory, includePopular, language, region])

  // Fetch suggestions using React Query
  const {
    data: suggestions,
    isLoading,
    error,
    refetch,
  } = useQuery({
    queryKey: ['smart-autocomplete', query, buildParams()],
    queryFn: async (): Promise<SmartAutocompleteResponse> => {
      if (!query.trim()) {
        return {
          suggestions: [],
          total: 0,
          has_more: false,
          query_time_ms: 0,
        }
      }

      const response = await apiClient.get(`/search/autocomplete/smart?${buildParams()}`)
      return response.data
    },
    enabled: enabled && query.length > 0,
    staleTime: 1 * 60 * 1000, // 1 minute
    gcTime: 5 * 60 * 1000, // 5 minutes
  })

  // Track user interactions with suggestions
  const trackInteraction = useCallback(async (
    suggestion: SmartAutocompleteSuggestion,
    type: 'click' | 'impression',
    position?: number
  ) => {
    try {
      await apiClient.post('/search/autocomplete/track', {
        entry_id: suggestion.id,
        interaction_type: type,
        session_id: sessionIdRef.current,
        query,
        position,
      })
    } catch (error) {
      console.error('Failed to track autocomplete interaction:', error)
    }
  }, [query])

  return {
    suggestions: suggestions || null,
    isLoading,
    error: error as Error | null,
    trackInteraction,
    refetch,
  }
}

// Hook for managing autocomplete state
export function useAutocompleteState(initialValue = '') {
  const [value, setValue] = useState(initialValue)
  const [isOpen, setIsOpen] = useState(false)
  const [selectedIndex, setSelectedIndex] = useState(-1)

  const handleSelect = useCallback((suggestion: SmartAutocompleteSuggestion) => {
    setValue(suggestion.value)
    setIsOpen(false)
    setSelectedIndex(-1)
  }, [])

  const handleKeyDown = useCallback((
    e: React.KeyboardEvent,
    suggestions: SmartAutocompleteSuggestion[] = []
  ) => {
    if (!isOpen || !suggestions.length) return

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault()
        setSelectedIndex(prev => 
          prev < suggestions.length - 1 ? prev + 1 : prev
        )
        break
      case 'ArrowUp':
        e.preventDefault()
        setSelectedIndex(prev => prev > 0 ? prev - 1 : -1)
        break
      case 'Enter':
        e.preventDefault()
        if (selectedIndex >= 0 && suggestions[selectedIndex]) {
          handleSelect(suggestions[selectedIndex])
        }
        break
      case 'Escape':
        setIsOpen(false)
        setSelectedIndex(-1)
        break
    }
  }, [isOpen, selectedIndex, handleSelect])

  const reset = useCallback(() => {
    setValue('')
    setIsOpen(false)
    setSelectedIndex(-1)
  }, [])

  return {
    value,
    setValue,
    isOpen,
    setIsOpen,
    selectedIndex,
    setSelectedIndex,
    handleSelect,
    handleKeyDown,
    reset,
  }
}

// Hook for autocomplete analytics
export function useAutocompleteAnalytics() {
  const trackSearch = useCallback(async (query: string, resultCount: number) => {
    try {
      await apiClient.post('/search/record', {
        query,
        result_count: resultCount,
        search_type: 'autocomplete_search',
      })
    } catch (error) {
      console.error('Failed to track search:', error)
    }
  }, [])

  const trackSuggestionClick = useCallback(async (
    suggestion: SmartAutocompleteSuggestion,
    query: string,
    position: number
  ) => {
    try {
      await apiClient.post('/search/autocomplete/track', {
        entry_id: suggestion.id,
        interaction_type: 'click',
        query,
        position,
      })
    } catch (error) {
      console.error('Failed to track suggestion click:', error)
    }
  }, [])

  const trackSuggestionImpression = useCallback(async (
    suggestions: SmartAutocompleteSuggestion[],
    query: string
  ) => {
    try {
      // Track impressions for all visible suggestions
      const impressions = suggestions.map((suggestion, index) => ({
        entry_id: suggestion.id,
        interaction_type: 'impression',
        query,
        position: index,
      }))

      // Send batch impression tracking
      await Promise.all(
        impressions.map(impression =>
          apiClient.post('/search/autocomplete/track', impression)
        )
      )
    } catch (error) {
      console.error('Failed to track suggestion impressions:', error)
    }
  }, [])

  return {
    trackSearch,
    trackSuggestionClick,
    trackSuggestionImpression,
  }
}
