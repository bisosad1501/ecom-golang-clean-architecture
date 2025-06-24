'use client'

import { useEffect, useState } from 'react'

interface AnimatedBackgroundProps {
  className?: string
  variant?: 'hero' | 'section'
}

interface Particle {
  id: number
  x: number
  y: number
  size: number
  speedX: number
  speedY: number
  opacity: number
  color: string
}

export function AnimatedBackground({ className, variant = 'section' }: AnimatedBackgroundProps) {
  const [particles, setParticles] = useState<Particle[]>([])
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 })

  useEffect(() => {
    // Generate particles
    const newParticles: Particle[] = []
    const particleCount = variant === 'hero' ? 50 : 20

    for (let i = 0; i < particleCount; i++) {
      newParticles.push({
        id: i,
        x: Math.random() * 100,
        y: Math.random() * 100,
        size: Math.random() * 4 + 1,
        speedX: (Math.random() - 0.5) * 0.5,
        speedY: (Math.random() - 0.5) * 0.5,
        opacity: Math.random() * 0.5 + 0.1,
        color: Math.random() > 0.5 ? '#6366f1' : '#8b5cf6'
      })
    }
    setParticles(newParticles)

    // Mouse move handler
    const handleMouseMove = (e: MouseEvent) => {
      setMousePosition({
        x: (e.clientX / window.innerWidth) * 100,
        y: (e.clientY / window.innerHeight) * 100
      })
    }

    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [variant])

  return (
    <div className={`absolute inset-0 overflow-hidden ${className}`}>
      {/* Dynamic gradient orbs */}
      <div
        className="absolute w-96 h-96 bg-gradient-to-br from-primary-400/30 to-violet-500/30 rounded-full blur-3xl transition-all duration-1000 ease-out"
        style={{
          top: `${20 + mousePosition.y * 0.1}%`,
          right: `${10 + mousePosition.x * 0.1}%`,
          transform: `scale(${1 + mousePosition.x * 0.001})`,
        }}
      />

      <div
        className="absolute w-80 h-80 bg-gradient-to-br from-violet-400/25 to-primary-500/25 rounded-full blur-2xl transition-all duration-1500 ease-out"
        style={{
          bottom: `${15 + mousePosition.y * 0.05}%`,
          left: `${15 + mousePosition.x * 0.05}%`,
          transform: `scale(${1.2 - mousePosition.x * 0.001})`,
        }}
      />

      <div
        className="absolute w-64 h-64 bg-gradient-to-br from-primary-300/20 to-violet-400/20 rounded-full blur-3xl transition-all duration-2000 ease-out"
        style={{
          top: `${40 + mousePosition.y * 0.08}%`,
          left: `${40 + mousePosition.x * 0.08}%`,
          transform: `rotate(${mousePosition.x * 0.1}deg)`,
        }}
      />

      {/* Floating Particles */}
      <div className="absolute inset-0">
        {particles.map((particle) => (
          <div
            key={particle.id}
            className="absolute rounded-full animate-float"
            style={{
              left: `${particle.x}%`,
              top: `${particle.y}%`,
              width: `${particle.size}px`,
              height: `${particle.size}px`,
              backgroundColor: particle.color,
              opacity: particle.opacity,
              animationDelay: `${particle.id * 0.1}s`,
              animationDuration: `${3 + particle.id * 0.1}s`,
            }}
          />
        ))}
      </div>

      {/* Advanced SVG Pattern */}
      <svg
        className="absolute inset-0 w-full h-full opacity-20 text-white"
        viewBox="0 0 1200 800"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <linearGradient id="heroGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#6366f1" stopOpacity="0.3" />
            <stop offset="50%" stopColor="#8b5cf6" stopOpacity="0.1" />
            <stop offset="100%" stopColor="#6366f1" stopOpacity="0.3" />
          </linearGradient>

          <radialGradient id="radialGlow" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#ffffff" stopOpacity="0.1" />
            <stop offset="100%" stopColor="#6366f1" stopOpacity="0.05" />
          </radialGradient>

          <filter id="glow">
            <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
            <feMerge>
              <feMergeNode in="coloredBlur"/>
              <feMergeNode in="SourceGraphic"/>
            </feMerge>
          </filter>

          <pattern id="hexPattern" x="0" y="0" width="60" height="52" patternUnits="userSpaceOnUse">
            <polygon points="30,2 50,17 50,35 30,50 10,35 10,17"
                     fill="none"
                     stroke="currentColor"
                     strokeWidth="0.5"
                     opacity="0.1"/>
          </pattern>
        </defs>

        {/* Hexagonal grid */}
        <rect width="100%" height="100%" fill="url(#hexPattern)" />

        {/* Flowing curves with glow */}
        <path
          d="M-100 200 Q 200 100 400 200 T 800 200 Q 1000 150 1300 200"
          stroke="url(#heroGradient)"
          strokeWidth="3"
          fill="none"
          filter="url(#glow)"
          className="animate-pulse"
          style={{ animationDuration: '4s' }}
        />

        <path
          d="M-100 400 Q 300 300 500 400 T 900 400 Q 1100 350 1300 400"
          stroke="url(#heroGradient)"
          strokeWidth="2"
          fill="none"
          filter="url(#glow)"
          className="animate-pulse"
          style={{ animationDelay: '1s', animationDuration: '5s' }}
        />

        <path
          d="M-100 600 Q 250 500 450 600 T 850 600 Q 1050 550 1300 600"
          stroke="url(#heroGradient)"
          strokeWidth="2"
          fill="none"
          filter="url(#glow)"
          className="animate-pulse"
          style={{ animationDelay: '2s', animationDuration: '6s' }}
        />

        {/* Geometric shapes with animations */}
        <g className="animate-spin" style={{ transformOrigin: '200px 150px', animationDuration: '20s' }}>
          <polygon points="200,130 220,150 200,170 180,150"
                   fill="url(#radialGlow)"
                   stroke="#6366f1"
                   strokeWidth="1"
                   opacity="0.3"/>
        </g>

        <g className="animate-spin" style={{ transformOrigin: '1000px 300px', animationDuration: '25s', animationDirection: 'reverse' }}>
          <circle cx="1000" cy="300" r="15"
                  fill="none"
                  stroke="url(#heroGradient)"
                  strokeWidth="2"
                  opacity="0.4"/>
          <circle cx="1000" cy="300" r="8"
                  fill="url(#radialGlow)"
                  opacity="0.2"/>
        </g>

        <g className="animate-bounce" style={{ animationDuration: '3s' }}>
          <rect x="750" y="500" width="20" height="20"
                fill="url(#heroGradient)"
                opacity="0.3"
                rx="4"/>
        </g>

        {/* Constellation effect */}
        <g opacity="0.4">
          <circle cx="150" cy="100" r="2" fill="#ffffff"/>
          <circle cx="300" cy="180" r="1.5" fill="#6366f1"/>
          <circle cx="450" cy="120" r="2" fill="#8b5cf6"/>
          <circle cx="600" cy="200" r="1" fill="#ffffff"/>
          <circle cx="750" cy="140" r="1.5" fill="#6366f1"/>
          <circle cx="900" cy="160" r="2" fill="#8b5cf6"/>

          <line x1="150" y1="100" x2="300" y2="180" stroke="#6366f1" strokeWidth="0.5" opacity="0.3"/>
          <line x1="300" y1="180" x2="450" y2="120" stroke="#8b5cf6" strokeWidth="0.5" opacity="0.3"/>
          <line x1="450" y1="120" x2="600" y2="200" stroke="#6366f1" strokeWidth="0.5" opacity="0.3"/>
          <line x1="600" y1="200" x2="750" y2="140" stroke="#8b5cf6" strokeWidth="0.5" opacity="0.3"/>
          <line x1="750" y1="140" x2="900" y2="160" stroke="#6366f1" strokeWidth="0.5" opacity="0.3"/>
        </g>
      </svg>

      {/* Dynamic gradient overlays */}
      <div
        className="absolute top-0 left-0 w-full h-1/2 bg-gradient-to-br from-primary-500/10 via-transparent to-violet-500/5 transition-all duration-1000"
        style={{
          opacity: 0.3 + mousePosition.x * 0.002,
        }}
      />
      <div
        className="absolute bottom-0 right-0 w-2/3 h-2/3 bg-gradient-to-tl from-violet-600/8 via-primary-400/5 to-transparent transition-all duration-1500"
        style={{
          opacity: 0.4 + mousePosition.y * 0.002,
        }}
      />
      <div
        className="absolute inset-0 bg-gradient-to-r from-transparent via-white/5 to-transparent transition-all duration-2000"
        style={{
          transform: `translateX(${mousePosition.x * 0.1}px)`,
        }}
      />

      {/* Ambient light effect */}
      {variant === 'hero' && (
        <div
          className="absolute inset-0 bg-radial-gradient from-white/10 via-transparent to-transparent transition-all duration-1000"
          style={{
            background: `radial-gradient(circle at ${mousePosition.x}% ${mousePosition.y}%, rgba(255,255,255,0.1) 0%, transparent 50%)`,
          }}
        />
      )}
    </div>
  )
}
