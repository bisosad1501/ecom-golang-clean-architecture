'use client'

import { useState, useEffect } from 'react'
import { Clock, Package, Truck, CheckCircle, XCircle, ShieldCheck, Zap, MapPin, RefreshCw, User, MessageCircle, Eye, EyeOff } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { formatDate, cn } from '@/lib/utils'
import { OrderEvent } from '@/types'

interface OrderTimelineProps {
  orderId: string
  events?: OrderEvent[]
  showPrivateEvents?: boolean
  className?: string
}

const getEventIcon = (eventType: string) => {
  switch (eventType) {
    case 'created':
      return <Package className="h-4 w-4 text-blue-400" />
    case 'status_changed':
      return <RefreshCw className="h-4 w-4 text-purple-400" />
    case 'payment_received':
      return <CheckCircle className="h-4 w-4 text-green-400" />
    case 'payment_failed':
      return <XCircle className="h-4 w-4 text-red-400" />
    case 'shipped':
      return <Truck className="h-4 w-4 text-indigo-400" />
    case 'delivered':
      return <CheckCircle className="h-4 w-4 text-emerald-400" />
    case 'cancelled':
      return <XCircle className="h-4 w-4 text-red-400" />
    case 'refunded':
      return <RefreshCw className="h-4 w-4 text-orange-400" />
    case 'returned':
      return <RefreshCw className="h-4 w-4 text-yellow-400" />
    case 'note_added':
      return <MessageCircle className="h-4 w-4 text-cyan-400" />
    case 'tracking_updated':
      return <MapPin className="h-4 w-4 text-pink-400" />
    case 'inventory_reserved':
      return <ShieldCheck className="h-4 w-4 text-blue-400" />
    case 'inventory_released':
      return <RefreshCw className="h-4 w-4 text-gray-400" />
    default:
      return <Clock className="h-4 w-4 text-gray-400" />
  }
}

const getEventColor = (eventType: string) => {
  switch (eventType) {
    case 'created':
      return 'border-blue-500/30 bg-blue-500/10'
    case 'status_changed':
      return 'border-purple-500/30 bg-purple-500/10'
    case 'payment_received':
      return 'border-green-500/30 bg-green-500/10'
    case 'payment_failed':
      return 'border-red-500/30 bg-red-500/10'
    case 'shipped':
      return 'border-indigo-500/30 bg-indigo-500/10'
    case 'delivered':
      return 'border-emerald-500/30 bg-emerald-500/10'
    case 'cancelled':
      return 'border-red-500/30 bg-red-500/10'
    case 'refunded':
      return 'border-orange-500/30 bg-orange-500/10'
    case 'returned':
      return 'border-yellow-500/30 bg-yellow-500/10'
    case 'note_added':
      return 'border-cyan-500/30 bg-cyan-500/10'
    case 'tracking_updated':
      return 'border-pink-500/30 bg-pink-500/10'
    case 'inventory_reserved':
      return 'border-blue-500/30 bg-blue-500/10'
    case 'inventory_released':
      return 'border-gray-500/30 bg-gray-500/10'
    default:
      return 'border-gray-500/30 bg-gray-500/10'
  }
}

export function OrderTimeline({ orderId, events = [], showPrivateEvents = false, className }: OrderTimelineProps) {
  const [filteredEvents, setFilteredEvents] = useState<OrderEvent[]>([])
  const [showPrivate, setShowPrivate] = useState(showPrivateEvents)

  useEffect(() => {
    const filtered = showPrivate 
      ? events 
      : events.filter(event => event.is_public)
    
    // Sort by created_at descending (newest first)
    const sorted = filtered.sort((a, b) => 
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
    
    setFilteredEvents(sorted)
  }, [events, showPrivate])

  if (!events.length) {
    return (
      <Card className={cn("bg-gray-900/50 border-gray-800", className)}>
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <Clock className="h-5 w-5" />
            Order Timeline
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center py-8 text-gray-400">
            <Clock className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No events available for this order</p>
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className={cn("bg-gray-900/50 border-gray-800", className)}>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-white flex items-center gap-2">
            <Clock className="h-5 w-5" />
            Order Timeline
          </CardTitle>
          {events.some(e => !e.is_public) && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowPrivate(!showPrivate)}
              className="border-gray-700 text-gray-300 hover:bg-gray-800"
            >
              {showPrivate ? (
                <>
                  <EyeOff className="h-4 w-4 mr-2" />
                  Hide Private
                </>
              ) : (
                <>
                  <Eye className="h-4 w-4 mr-2" />
                  Show All
                </>
              )}
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {filteredEvents.map((event, index) => (
            <div key={event.id} className="relative">
              {/* Timeline line */}
              {index < filteredEvents.length - 1 && (
                <div className="absolute left-6 top-12 w-px h-8 bg-gray-700" />
              )}
              
              {/* Event item */}
              <div className={cn(
                "flex gap-4 p-4 rounded-lg border transition-all duration-200 hover:shadow-lg",
                getEventColor(event.event_type)
              )}>
                {/* Event icon */}
                <div className="flex-shrink-0 w-12 h-12 rounded-full bg-gray-800/50 border border-gray-700 flex items-center justify-center">
                  {getEventIcon(event.event_type)}
                </div>
                
                {/* Event content */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <h4 className="text-white font-medium">{event.title}</h4>
                      <p className="text-gray-300 text-sm mt-1">{event.description}</p>
                      
                      {/* Event data */}
                      {event.data && (
                        <div className="mt-2 p-2 bg-gray-800/50 rounded text-xs text-gray-400 font-mono">
                          {event.data}
                        </div>
                      )}
                      
                      {/* User info */}
                      {event.user && (
                        <div className="flex items-center gap-2 mt-2 text-xs text-gray-400">
                          <User className="h-3 w-3" />
                          <span>{event.user.first_name} {event.user.last_name}</span>
                          <span>({event.user.email})</span>
                        </div>
                      )}
                    </div>
                    
                    <div className="flex flex-col items-end gap-2">
                      <time className="text-xs text-gray-400">
                        {formatDate(event.created_at)}
                      </time>
                      
                      {/* Privacy badge */}
                      <Badge 
                        variant={event.is_public ? "default" : "secondary"}
                        className={cn(
                          "text-xs",
                          event.is_public 
                            ? "bg-green-500/20 text-green-300 border-green-500/40" 
                            : "bg-gray-500/20 text-gray-300 border-gray-500/40"
                        )}
                      >
                        {event.is_public ? "Public" : "Private"}
                      </Badge>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
        
        {filteredEvents.length === 0 && (
          <div className="text-center py-8 text-gray-400">
            <MessageCircle className="h-12 w-12 mx-auto mb-4 opacity-50" />
            <p>No {showPrivate ? '' : 'public '}events to display</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
