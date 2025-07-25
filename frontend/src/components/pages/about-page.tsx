'use client'

import Link from 'next/link'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Heart,
  Users,
  Globe,
  Award,
  Target,
  Zap,
  Shield,
  Truck,
  ArrowRight,
  Star,
  Mail,
  Sparkles,
  TrendingUp,
  Eye,
  Lightbulb,
  Handshake
} from 'lucide-react'
import { AnimatedBackground } from '@/components/ui/animated-background'

export function AboutPage() {
  const stats = [
    { label: 'Happy Customers', value: '50K+', icon: Users, color: 'from-blue-500 to-cyan-500' },
    { label: 'Products Sold', value: '1M+', icon: Award, color: 'from-purple-500 to-pink-500' },
    { label: 'Countries Served', value: '25+', icon: Globe, color: 'from-green-500 to-emerald-500' },
    { label: 'Years of Excellence', value: '10+', icon: Star, color: 'from-orange-500 to-red-500' },
  ]

  const values = [
    {
      icon: Heart,
      title: 'Customer First',
      description: 'Every decision we make starts with our customers. Your satisfaction drives our innovation and excellence.',
      gradient: 'from-red-500/20 to-pink-500/20'
    },
    {
      icon: Shield,
      title: 'Quality Assurance',
      description: 'We carefully curate every product to ensure it meets our high standards of quality and reliability.',
      gradient: 'from-blue-500/20 to-cyan-500/20'
    },
    {
      icon: Zap,
      title: 'Innovation',
      description: 'We continuously evolve our platform to provide the best shopping experience possible.',
      gradient: 'from-yellow-500/20 to-orange-500/20'
    },
    {
      icon: Truck,
      title: 'Fast Delivery',
      description: 'Quick and reliable shipping to get your orders to you as fast as possible worldwide.',
      gradient: 'from-green-500/20 to-teal-500/20'
    },
  ]

  const team = [
    {
      name: 'Sarah Johnson',
      role: 'CEO & Founder',
      description: 'Visionary leader with 15+ years in e-commerce innovation.',
      gradient: 'from-purple-500/20 to-pink-500/20'
    },
    {
      name: 'Michael Chen',
      role: 'CTO',
      description: 'Tech innovator passionate about creating seamless user experiences.',
      gradient: 'from-blue-500/20 to-cyan-500/20'
    },
    {
      name: 'Emily Rodriguez',
      role: 'Head of Operations',
      description: 'Operations expert ensuring smooth and efficient order fulfillment.',
      gradient: 'from-green-500/20 to-emerald-500/20'
    },
    {
      name: 'David Kim',
      role: 'Head of Marketing',
      description: 'Creative strategist building meaningful brand connections worldwide.',
      gradient: 'from-orange-500/20 to-red-500/20'
    },
  ]

  const milestones = [
    {
      year: '2014',
      title: 'Company Founded',
      description: 'Started with a vision to revolutionize online shopping experience.',
      icon: Lightbulb
    },
    {
      year: '2018',
      title: 'Global Expansion',
      description: 'Expanded to serve customers in 25+ countries worldwide.',
      icon: Globe
    },
    {
      year: '2020',
      title: 'Platform Redesign',
      description: 'Launched our modern, mobile-first shopping platform.',
      icon: Zap
    },
    {
      year: '2024',
      title: 'AI Integration',
      description: 'Introduced AI-powered recommendations and customer service.',
      icon: Target
    },
  ]

  return (
    <div className="min-h-screen bg-black text-white">
      {/* Hero Section - Matching home page style */}
      <section className="relative bg-gradient-to-br from-black via-gray-900 to-black text-white overflow-hidden min-h-[60vh] flex items-center py-16">
        <AnimatedBackground variant="hero" />

        <div className="container mx-auto px-4 relative z-10">
          <div className="max-w-6xl mx-auto">
            <div className="text-center mb-16">
              {/* Badge */}
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6 animate-fade-in">
                <Heart className="h-4 w-4 text-[#ff9000]" />
                <span className="text-[#ff9000] font-semibold text-sm">ABOUT BIHUB</span>
              </div>

              {/* Main heading */}
              <h1 className="text-4xl lg:text-6xl font-bold mb-6 animate-scale-in">
                Building the Future of
                <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> E-commerce</span>
              </h1>

              <p className="text-lg text-gray-300 max-w-3xl mx-auto leading-relaxed mb-8 animate-fade-in" style={{ animationDelay: '0.2s' }}>
                We're passionate about creating exceptional shopping experiences that connect
                people with products they love. Our mission is to make online shopping
                simple, secure, and delightful for everyone.
              </p>

              {/* CTA Buttons */}
              <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12 animate-fade-in" style={{ animationDelay: '0.4s' }}>
                <Button size="lg" className="shadow-lg hover:shadow-xl transition-all duration-300 group text-white" style={{backgroundColor: '#FF9000'}} asChild
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
                >
                  <Link href="/contact">
                    <Mail className="mr-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
                    Get in Touch
                    <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                  </Link>
                </Button>
                <Button size="lg" variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}} asChild>
                  <Link href="/products">
                    <Sparkles className="mr-2 h-5 w-5" />
                    Explore Products
                  </Link>
                </Button>
              </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-6 animate-fade-in" style={{ animationDelay: '0.6s' }}>
              {stats.map((stat, index) => (
                <div key={index} className="text-center group">
                  <div className="glass-effect p-6 rounded-2xl border border-white/10 hover:border-[#ff9000]/30 transition-all duration-300 hover:scale-105">
                    <div className={`w-16 h-16 rounded-2xl bg-gradient-to-br ${stat.color} flex items-center justify-center mx-auto mb-4 shadow-lg group-hover:scale-110 transition-transform duration-300`}>
                      <stat.icon className="h-8 w-8 text-white" />
                    </div>
                    <div className="text-3xl font-bold text-white mb-2">{stat.value}</div>
                    <div className="text-gray-300 text-sm font-medium">{stat.label}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Our Values Section */}
      <section className="py-20 bg-gradient-to-b from-black to-gray-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <Target className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">OUR VALUES</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold text-white mb-4">
              What Drives Us Forward
            </h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">
              Our core values shape every decision we make and every product we deliver
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {values.map((value, index) => (
              <div key={index} className="group animate-scale-in" style={{ animationDelay: `${index * 0.1}s` }}>
                <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300 hover:scale-105">
                  <CardContent className="p-6 text-center">
                    <div className={`w-16 h-16 rounded-2xl bg-gradient-to-br ${value.gradient} border border-white/20 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-300`}>
                      <value.icon className="h-8 w-8 text-white" />
                    </div>
                    <h3 className="text-xl font-bold text-white mb-3 group-hover:text-[#ff9000] transition-colors">
                      {value.title}
                    </h3>
                    <p className="text-gray-300 text-sm leading-relaxed">
                      {value.description}
                    </p>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Our Story Section */}
      <section className="py-20 bg-gradient-to-b from-gray-900 to-black">
        <div className="container mx-auto px-4">
          <div className="max-w-6xl mx-auto">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
              {/* Left Content */}
              <div className="animate-fade-in">
                <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
                  <Eye className="h-4 w-4 text-[#ff9000]" />
                  <span className="text-[#ff9000] font-semibold text-sm">OUR STORY</span>
                </div>
                <h2 className="text-3xl lg:text-4xl font-bold text-white mb-6">
                  From Vision to Reality
                </h2>
                <div className="space-y-4 text-gray-300 leading-relaxed">
                  <p>
                    Founded in 2014, BiHub began as a small startup with a big dream:
                    to create the most customer-centric e-commerce platform in the world.
                    What started in a garage has grown into a global marketplace serving
                    millions of customers.
                  </p>
                  <p>
                    We believe that shopping should be more than just a transaction—it should
                    be an experience that brings joy, discovery, and convenience to your life.
                    That's why we've built our platform with cutting-edge technology,
                    intuitive design, and a deep understanding of what customers truly want.
                  </p>
                  <p>
                    Today, we're proud to be a trusted partner for both shoppers and sellers,
                    creating a vibrant ecosystem where great products meet great people.
                  </p>
                </div>

                <div className="flex flex-col sm:flex-row gap-4 mt-8">
                  <Button size="lg" className="bg-[#ff9000] hover:bg-[#e67e00] text-white group">
                    <Handshake className="mr-2 h-5 w-5 group-hover:scale-110 transition-transform" />
                    Join Our Team
                    <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                  </Button>
                  <Button size="lg" variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800" asChild>
                    <Link href="/contact">
                      Learn More
                    </Link>
                  </Button>
                </div>
              </div>

              {/* Right Content - Mission Card */}
              <div className="animate-scale-in" style={{ animationDelay: '0.2s' }}>
                <Card className="overflow-hidden bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 border border-[#ff9000]/30 backdrop-blur-sm">
                  <CardContent className="p-8 text-center">
                    <div className="w-20 h-20 rounded-full bg-gradient-to-br from-[#ff9000] to-orange-600 flex items-center justify-center mx-auto mb-6 shadow-xl">
                      <Target className="h-10 w-10 text-white" />
                    </div>
                    <h3 className="text-2xl font-bold text-white mb-4">Our Mission</h3>
                    <p className="text-gray-200 leading-relaxed mb-6">
                      To democratize commerce and empower everyone to build thriving businesses online,
                      while creating exceptional shopping experiences that bring joy to millions of customers worldwide.
                    </p>
                    <Badge className="bg-[#ff9000]/20 text-[#ff9000] border border-[#ff9000]/30">
                      Trusted by 50K+ Customers
                    </Badge>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Team Section */}
      <section className="py-20 bg-gradient-to-b from-black to-gray-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <Users className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">OUR TEAM</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold text-white mb-4">
              Meet the Visionaries
            </h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">
              Passionate professionals dedicated to revolutionizing your shopping experience
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {team.map((member, index) => (
              <div key={index} className="group animate-scale-in" style={{ animationDelay: `${index * 0.1}s` }}>
                <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300 hover:scale-105">
                  <CardContent className="p-6 text-center">
                    <div className={`w-20 h-20 rounded-full bg-gradient-to-br ${member.gradient} border-2 border-white/20 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-300`}>
                      <Users className="h-10 w-10 text-white" />
                    </div>
                    <h3 className="text-xl font-bold text-white mb-2 group-hover:text-[#ff9000] transition-colors">
                      {member.name}
                    </h3>
                    <Badge className="mb-3 bg-[#ff9000]/20 text-[#ff9000] border border-[#ff9000]/30">
                      {member.role}
                    </Badge>
                    <p className="text-gray-300 text-sm leading-relaxed">
                      {member.description}
                    </p>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Timeline Section */}
      <section className="py-20 bg-gradient-to-b from-gray-900 to-black">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <TrendingUp className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">OUR JOURNEY</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold text-white mb-4">
              Milestones & Achievements
            </h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">
              Key moments that shaped our journey to becoming a leading e-commerce platform
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {milestones.map((milestone, index) => (
              <div key={index} className="group animate-scale-in" style={{ animationDelay: `${index * 0.1}s` }}>
                <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300 hover:scale-105">
                  <CardContent className="p-6 text-center">
                    <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 border border-[#ff9000]/30 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-300">
                      <milestone.icon className="h-8 w-8 text-[#ff9000]" />
                    </div>
                    <div className="text-2xl font-bold text-[#ff9000] mb-2">{milestone.year}</div>
                    <h3 className="text-lg font-bold text-white mb-3 group-hover:text-[#ff9000] transition-colors">
                      {milestone.title}
                    </h3>
                    <p className="text-gray-300 text-sm leading-relaxed">
                      {milestone.description}
                    </p>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-br from-black via-gray-900 to-black relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />

        <div className="container mx-auto px-4 relative z-10">
          <div className="max-w-4xl mx-auto text-center">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <Sparkles className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">JOIN US</span>
            </div>

            <h2 className="text-3xl lg:text-5xl font-bold text-white mb-6">
              Ready to Start Your
              <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Shopping Journey?</span>
            </h2>

            <p className="text-lg text-gray-300 max-w-2xl mx-auto mb-8 leading-relaxed">
              Join millions of satisfied customers who trust BiHub for their online shopping needs.
              Discover amazing products, unbeatable prices, and exceptional service.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" className="shadow-lg hover:shadow-xl transition-all duration-300 group text-white" style={{backgroundColor: '#FF9000'}} asChild
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
              >
                <Link href="/products">
                  <Sparkles className="mr-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
                  Start Shopping
                  <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}} asChild>
                <Link href="/contact">
                  <Mail className="mr-2 h-5 w-5" />
                  Contact Us
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
