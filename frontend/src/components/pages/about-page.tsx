'use client'

import Image from 'next/image'
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
  CheckCircle,
  Mail,
  Phone,
  MapPin,
} from 'lucide-react'

export default function AboutPage() {
  const stats = [
    { label: 'Happy Customers', value: '50,000+', icon: Users },
    { label: 'Products Sold', value: '1M+', icon: Award },
    { label: 'Countries Served', value: '25+', icon: Globe },
    { label: 'Years of Excellence', value: '10+', icon: Star },
  ]

  const values = [
    {
      icon: Heart,
      title: 'Customer First',
      description: 'Every decision we make starts with our customers. Your satisfaction is our top priority.',
    },
    {
      icon: Shield,
      title: 'Quality Assurance',
      description: 'We carefully curate every product to ensure it meets our high standards of quality.',
    },
    {
      icon: Zap,
      title: 'Innovation',
      description: 'We continuously evolve our platform to provide the best shopping experience.',
    },
    {
      icon: Truck,
      title: 'Fast Delivery',
      description: 'Quick and reliable shipping to get your orders to you as fast as possible.',
    },
  ]

  const team = [
    {
      name: 'Sarah Johnson',
      role: 'CEO & Founder',
      image: '/team/sarah.jpg',
      description: 'Visionary leader with 15+ years in e-commerce.',
    },
    {
      name: 'Michael Chen',
      role: 'CTO',
      image: '/team/michael.jpg',
      description: 'Tech innovator passionate about user experience.',
    },
    {
      name: 'Emily Rodriguez',
      role: 'Head of Operations',
      image: '/team/emily.jpg',
      description: 'Operations expert ensuring smooth fulfillment.',
    },
    {
      name: 'David Kim',
      role: 'Head of Marketing',
      image: '/team/david.jpg',
      description: 'Creative strategist building brand connections.',
    },
  ]

  const milestones = [
    {
      year: '2014',
      title: 'Company Founded',
      description: 'Started with a vision to revolutionize online shopping.',
    },
    {
      year: '2016',
      title: 'First Million',
      description: 'Reached our first million in revenue and 10,000 customers.',
    },
    {
      year: '2018',
      title: 'Global Expansion',
      description: 'Expanded to serve customers in 15 countries worldwide.',
    },
    {
      year: '2020',
      title: 'Platform Redesign',
      description: 'Launched our modern, mobile-first shopping platform.',
    },
    {
      year: '2022',
      title: 'Sustainability Initiative',
      description: 'Committed to carbon-neutral shipping and eco-friendly packaging.',
    },
    {
      year: '2024',
      title: 'AI Integration',
      description: 'Introduced AI-powered recommendations and customer service.',
    },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background">
      {/* Hero Section */}
      <section className="relative py-24 overflow-hidden">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large">
                <Heart className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">ABOUT US</span>
            </div>
            
            <h1 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
              Building the Future of <span className="text-gradient">E-commerce</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              We're passionate about creating exceptional shopping experiences that connect 
              people with products they love. Our mission is to make online shopping 
              simple, secure, and delightful for everyone.
            </p>
          </div>

          {/* Stats */}
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8 mb-16">
            {stats.map((stat, index) => (
              <Card key={index} variant="elevated" className="border-0 shadow-large text-center">
                <CardContent className="p-8">
                  <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center mx-auto mb-4 shadow-large">
                    <stat.icon className="h-8 w-8 text-white" />
                  </div>
                  <div className="text-3xl font-bold text-foreground mb-2">{stat.value}</div>
                  <div className="text-muted-foreground">{stat.label}</div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Our Story */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
            <div>
              <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
                Our Story
              </h2>
              <div className="space-y-6 text-lg text-muted-foreground leading-relaxed">
                <p>
                  Founded in 2014, EcomStore began as a small startup with a big dream: 
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
              
              <div className="flex items-center gap-4 mt-8">
                <Button size="lg" variant="gradient">
                  Learn More
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                <Button size="lg" variant="outline">
                  Contact Us
                </Button>
              </div>
            </div>
            
            <div className="relative">
              <Card variant="elevated" className="border-0 shadow-xl overflow-hidden">
                <div className="aspect-[4/3] bg-gradient-to-br from-primary-100 to-violet-100 flex items-center justify-center">
                  <div className="text-center">
                    <div className="w-24 h-24 rounded-full bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center mx-auto mb-4 shadow-large">
                      <Target className="h-12 w-12 text-white" />
                    </div>
                    <h3 className="text-2xl font-bold text-foreground mb-2">Our Mission</h3>
                    <p className="text-muted-foreground max-w-sm">
                      To democratize commerce and empower everyone to build thriving businesses online.
                    </p>
                  </div>
                </div>
              </Card>
            </div>
          </div>
        </div>
      </section>

      {/* Our Values */}
      <section className="py-24 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
              Our Values
            </h2>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              These core principles guide everything we do and shape the way we serve our customers.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {values.map((value, index) => (
              <Card key={index} variant="elevated" className="border-0 shadow-large text-center group hover:shadow-xl transition-all duration-300">
                <CardContent className="p-8">
                  <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center mx-auto mb-6 shadow-large group-hover:scale-110 transition-transform duration-300">
                    <value.icon className="h-8 w-8 text-white" />
                  </div>
                  <h3 className="text-xl font-bold text-foreground mb-4">{value.title}</h3>
                  <p className="text-muted-foreground leading-relaxed">{value.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Timeline */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
              Our Journey
            </h2>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              From humble beginnings to global success, here are the key milestones in our story.
            </p>
          </div>

          <div className="relative">
            {/* Timeline line */}
            <div className="absolute left-1/2 transform -translate-x-1/2 w-1 h-full bg-gradient-to-b from-primary-500 to-violet-600 rounded-full"></div>
            
            <div className="space-y-16">
              {milestones.map((milestone, index) => (
                <div key={index} className={`flex items-center ${index % 2 === 0 ? 'flex-row' : 'flex-row-reverse'}`}>
                  <div className={`w-1/2 ${index % 2 === 0 ? 'pr-8 text-right' : 'pl-8 text-left'}`}>
                    <Card variant="elevated" className="border-0 shadow-large">
                      <CardContent className="p-8">
                        <Badge variant="outline" className="mb-4 text-primary border-primary">
                          {milestone.year}
                        </Badge>
                        <h3 className="text-xl font-bold text-foreground mb-3">{milestone.title}</h3>
                        <p className="text-muted-foreground">{milestone.description}</p>
                      </CardContent>
                    </Card>
                  </div>
                  
                  {/* Timeline dot */}
                  <div className="relative z-10">
                    <div className="w-6 h-6 rounded-full bg-gradient-to-br from-primary-500 to-violet-600 border-4 border-background shadow-large"></div>
                  </div>
                  
                  <div className="w-1/2"></div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Team */}
      <section className="py-24 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
              Meet Our Team
            </h2>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              The passionate people behind our success, working every day to make your experience better.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {team.map((member, index) => (
              <Card key={index} variant="elevated" className="border-0 shadow-large text-center group hover:shadow-xl transition-all duration-300">
                <CardContent className="p-8">
                  <div className="w-24 h-24 rounded-full overflow-hidden mx-auto mb-6 bg-gradient-to-br from-primary-100 to-violet-100 flex items-center justify-center">
                    <Users className="h-12 w-12 text-primary" />
                  </div>
                  <h3 className="text-xl font-bold text-foreground mb-2">{member.name}</h3>
                  <p className="text-primary font-semibold mb-3">{member.role}</p>
                  <p className="text-muted-foreground text-sm">{member.description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Contact CTA */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <Card variant="elevated" className="border-0 shadow-xl bg-gradient-to-br from-primary-50 to-violet-50">
            <CardContent className="p-16 text-center">
              <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
                Ready to Start Your Journey?
              </h2>
              <p className="text-xl text-muted-foreground mb-12 max-w-2xl mx-auto">
                Join millions of satisfied customers and discover why EcomStore is the 
                preferred choice for online shopping.
              </p>
              
              <div className="flex flex-col sm:flex-row gap-6 justify-center items-center">
                <Button size="xl" variant="gradient" className="shadow-large">
                  Start Shopping
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
                
                <div className="flex items-center gap-6 text-sm text-muted-foreground">
                  <div className="flex items-center gap-2">
                    <CheckCircle className="h-5 w-5 text-emerald-600" />
                    <span>Free Shipping</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <CheckCircle className="h-5 w-5 text-emerald-600" />
                    <span>30-Day Returns</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <CheckCircle className="h-5 w-5 text-emerald-600" />
                    <span>24/7 Support</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </section>
    </div>
  )
}
