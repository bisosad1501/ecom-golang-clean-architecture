'use client'

import { ReactNode } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface AuthLayoutProps {
  children: ReactNode
  title: string
  description: string
  heroTitle: string
  heroSubtitle: string
  heroDescription: string
  features: Array<{
    icon: ReactNode
    title: string
    description: string
  }>
  bottomContent?: ReactNode
}

export function AuthLayout({
  children,
  title,
  description,
  heroTitle,
  heroSubtitle,
  heroDescription,
  features,
  bottomContent
}: AuthLayoutProps) {
  return (
    <div className="min-h-screen bg-slate-950 relative overflow-hidden">
      {/* Enhanced Dynamic Background */}
      <div className="absolute inset-0">
        {/* Primary Background Gradient */}
        <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-slate-950 to-black"></div>
        
        {/* Orange Accent Layers */}
        <div className="absolute inset-0 bg-gradient-to-tr from-transparent via-transparent to-[#FF9000]/8"></div>
        <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-bl from-[#FF9000]/5 via-transparent to-transparent"></div>
        
        {/* Large Animated Geometric Shapes */}
        <div className="absolute top-20 left-20 w-96 h-96 bg-[#FF9000]/12 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute bottom-32 right-32 w-[28rem] h-[28rem] bg-[#FF9000]/10 rounded-full blur-3xl animate-pulse" style={{animationDelay: '2s'}}></div>
        <div className="absolute top-1/2 left-1/3 w-80 h-80 bg-[#FF9000]/8 rounded-full blur-2xl animate-bounce" style={{animationDuration: '6s'}}></div>
        
        {/* Medium Floating Elements */}
        <div className="absolute top-1/4 right-1/4 w-6 h-6 bg-[#FF9000]/50 rotate-45 animate-spin" style={{animationDuration: '12s'}}></div>
        <div className="absolute bottom-1/3 left-1/5 w-8 h-8 bg-[#FF9000]/40 rotate-45 animate-pulse"></div>
        <div className="absolute top-2/3 right-1/6 w-5 h-5 bg-[#FF9000]/60 rounded-full animate-bounce" style={{animationDelay: '1.5s'}}></div>
        <div className="absolute top-1/6 left-3/4 w-7 h-7 bg-[#FF9000]/45 rotate-45 animate-pulse" style={{animationDelay: '3s'}}></div>
        
        {/* Dynamic Light Streaks */}
        <div className="absolute top-0 left-1/4 w-px h-full bg-gradient-to-b from-transparent via-[#FF9000]/25 to-transparent opacity-60 animate-pulse"></div>
        <div className="absolute top-0 right-1/3 w-px h-full bg-gradient-to-b from-transparent via-[#FF9000]/20 to-transparent opacity-40 animate-pulse" style={{animationDelay: '2s'}}></div>
        
        {/* Diagonal Light Beams */}
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-1/4 left-0 w-full h-px bg-gradient-to-r from-transparent via-[#FF9000]/15 to-transparent rotate-12 opacity-30"></div>
          <div className="absolute bottom-1/4 right-0 w-full h-px bg-gradient-to-l from-transparent via-[#FF9000]/10 to-transparent -rotate-12 opacity-20"></div>
        </div>
        
        {/* Static Particle Effect */}
        <div className="absolute inset-0">
          <div className="absolute w-1 h-1 bg-[#FF9000]/30 rounded-full animate-pulse" style={{left: '15%', top: '20%', animationDelay: '0s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/25 rounded-full animate-pulse" style={{left: '85%', top: '30%', animationDelay: '1s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/35 rounded-full animate-pulse" style={{left: '25%', top: '70%', animationDelay: '2s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/20 rounded-full animate-pulse" style={{left: '75%', top: '80%', animationDelay: '3s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/40 rounded-full animate-pulse" style={{left: '45%', top: '15%', animationDelay: '1.5s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/30 rounded-full animate-pulse" style={{left: '65%', top: '60%', animationDelay: '2.5s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/25 rounded-full animate-pulse" style={{left: '35%', top: '45%', animationDelay: '0.5s'}} />
          <div className="absolute w-1 h-1 bg-[#FF9000]/35 rounded-full animate-pulse" style={{left: '55%', top: '25%', animationDelay: '3.5s'}} />
        </div>
      </div>

      <div className="relative z-10 h-screen flex">
        {/* Left Side - Enhanced Branding */}
        <div className="hidden lg:flex lg:w-2/5 flex-col justify-center px-6 xl:px-8 relative">
          {/* Enhanced Gradient Orbs */}
          <div className="absolute top-1/4 left-1/4 w-40 h-40 bg-[#FF9000]/20 rounded-full blur-3xl animate-pulse"></div>
          <div className="absolute bottom-1/3 right-1/4 w-32 h-32 bg-[#FF9000]/15 rounded-full blur-2xl animate-pulse" style={{animationDelay: '1s'}}></div>
          <div className="absolute top-1/2 left-1/2 w-24 h-24 bg-[#FF9000]/10 rounded-full blur-xl animate-bounce" style={{animationDuration: '4s'}}></div>

          <div className="relative z-10 max-w-md mx-auto">
            {/* Enhanced Hero Section - No Logo */}
            <div className="text-center mb-6">
              <h1 className="text-3xl font-bold text-white mb-3 leading-tight">
                {heroTitle}
                <span className="text-transparent bg-gradient-to-r from-[#FF9000] to-[#e67e00] bg-clip-text block mt-1 text-2xl">
                  {heroSubtitle}
                </span>
              </h1>

              <p className="text-base text-gray-300 leading-relaxed mb-6 max-w-sm mx-auto">
                {heroDescription}
              </p>
            </div>

            {/* Enhanced Feature Grid with Colors */}
            <div className="grid grid-cols-1 gap-3">
              {features.map((feature, index) => {
                const colors = [
                  'from-purple-500/20 to-purple-600/10 border-purple-400/30',
                  'from-blue-500/20 to-blue-600/10 border-blue-400/30',
                  'from-green-500/20 to-green-600/10 border-green-400/30',
                  'from-pink-500/20 to-pink-600/10 border-pink-400/30'
                ]
                const iconColors = [
                  'from-purple-500 to-purple-600',
                  'from-blue-500 to-blue-600',
                  'from-green-500 to-green-600',
                  'from-pink-500 to-pink-600'
                ]
                return (
                  <div
                    key={index}
                    className={`flex items-center gap-3 p-3 rounded-xl bg-gradient-to-r ${colors[index % colors.length]} border backdrop-blur-sm hover:scale-[1.02] transition-all duration-300 hover:shadow-lg`}
                  >
                    <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${iconColors[index % iconColors.length]} flex items-center justify-center shadow-lg`}>
                      {feature.icon}
                    </div>
                    <div className="flex-1">
                      <h3 className="text-white font-semibold text-sm mb-1">{feature.title}</h3>
                      <p className="text-gray-400 text-xs leading-relaxed">{feature.description}</p>
                    </div>
                  </div>
                )
              })}
            </div>

            {/* Compact Trust Indicators */}
            <div className="mt-6 pt-4 border-t border-gray-700/50">
              <div className="flex items-center justify-center gap-4 text-gray-400 text-xs">
                <div className="flex items-center gap-1">
                  <div className="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></div>
                  <span>SSL Secured</span>
                </div>
                <div className="flex items-center gap-1">
                  <div className="w-1.5 h-1.5 bg-blue-500 rounded-full animate-pulse"></div>
                  <span>GDPR Compliant</span>
                </div>
                <div className="flex items-center gap-1">
                  <div className="w-1.5 h-1.5 bg-[#FF9000] rounded-full animate-pulse"></div>
                  <span>24/7 Support</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Right Side - Form */}
        <div className="flex-1 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-4">
          <div className="max-w-sm w-full">
            <Card className="border-0 shadow-2xl bg-gray-900/90 backdrop-blur-md border border-gray-700/50 rounded-2xl overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-br from-gray-800/20 to-gray-900/40 pointer-events-none"></div>

              <CardHeader className="text-center pb-4 pt-6 relative">
                <CardTitle className="text-xl font-bold text-white mb-1">{title}</CardTitle>
                <CardDescription className="text-sm text-gray-300 leading-relaxed">
                  {description}
                </CardDescription>
              </CardHeader>

              <CardContent className="px-6 pb-6 relative">
                {children}
                {bottomContent}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
