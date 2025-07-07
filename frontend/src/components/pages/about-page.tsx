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
import {
  PageWrapper,
  PageHeader,
  PageSection,
  PageContainer,
  PageGrid,
  PageCard
} from '@/components/layout'
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'

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
    <PageWrapper>
      {/* Hero Section */}
      <PageSection size="xl" className="text-center">
        <div className="flex items-center justify-center gap-2 mb-4">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-orange-500 to-orange-400 flex items-center justify-center shadow-lg">
            <Heart className="h-4 w-4 text-white" />
          </div>
          <span className="text-orange-400 font-semibold text-sm">ABOUT US</span>
        </div>

        <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold text-white mb-4">
          Building the Future of <span className="text-orange-400">E-commerce</span>
        </h1>
        <p className="text-base lg:text-lg text-gray-200 max-w-2xl mx-auto leading-relaxed mb-8">
          We're passionate about creating exceptional shopping experiences that connect
          people with products they love. Our mission is to make online shopping
          simple, secure, and delightful for everyone.
        </p>

        {/* Stats */}
        <PageGrid type="features" className="mt-8">
          {stats.map((stat, index) => (
            <PageCard key={index} padding="base" className="text-center bg-gray-900 border-gray-600">
              <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-orange-500 to-orange-400 flex items-center justify-center mx-auto mb-3 shadow-lg">
                <stat.icon className="h-6 w-6 text-white" />
              </div>
              <div className="text-2xl font-bold text-white mb-1">{stat.value}</div>
              <div className="text-gray-200 text-sm">{stat.label}</div>
            </PageCard>
          ))}
        </PageGrid>
      </PageSection>

      {/* Our Story */}
      <PageSection size="lg">
        <PageGrid type="twoColumn" className="items-center">
          <div>
            <h2 className="text-2xl lg:text-3xl font-bold text-white mb-4">
              Our Story
            </h2>
            <div className="space-y-4 text-base text-gray-200 leading-relaxed">
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

              <div className="flex items-center gap-3 mt-6">
                <Button size="default" className="bg-orange-500 hover:bg-orange-600 text-white">
                  Learn More
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Button>
                <Button size="default" variant="outline" className="border-gray-600 text-white hover:bg-gray-800">
                  Contact Us
                </Button>
              </div>
            </div>

            <div className="relative">
              <PageCard padding="lg" className="overflow-hidden">
                <div className="aspect-[4/3] bg-gradient-to-br from-gray-800 to-gray-700 flex items-center justify-center">
                  <div className="text-center">
                    <div className="w-16 h-16 rounded-full bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center mx-auto mb-3 shadow-lg">
                      <Target className="h-8 w-8 text-white" />
                    </div>
                    <h3 className="text-xl font-bold text-white mb-2">Our Mission</h3>
                    <p className="text-gray-300 text-sm max-w-sm">
                      To democratize commerce and empower everyone to build thriving businesses online.
                    </p>
                  </div>
                </div>
              </PageCard>
            </div>
        </PageGrid>
      </PageSection>

      {/* Our Values */}
      <PageSection
        size="lg"
        title="Our Values"
        subtitle="These core principles guide everything we do and shape the way we serve our customers."
        className="bg-gray-900/50"
      >
        <PageGrid type="features">
          {values.map((value, index) => (
            <PageCard key={index} padding="base" className="text-center group hover:border-gray-600 transition-all duration-300">
              <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center mx-auto mb-4 shadow-lg group-hover:scale-110 transition-transform duration-300">
                <value.icon className="h-6 w-6 text-white" />
              </div>
              <h3 className="text-lg font-bold text-white mb-3">{value.title}</h3>
              <p className="text-gray-300 text-sm leading-relaxed">{value.description}</p>
            </PageCard>
          ))}
        </PageGrid>
      </PageSection>

      {/* Timeline */}
      <PageSection
        size="lg"
        title="Our Journey"
        subtitle="From humble beginnings to global success, here are the key milestones in our story."
      >

        <div className="relative">
          {/* Timeline line */}
          <div className="absolute left-1/2 transform -translate-x-1/2 w-1 h-full bg-gradient-to-b from-orange-500 to-orange-600 rounded-full"></div>

          <div className="space-y-12">
            {milestones.map((milestone, index) => (
              <div key={index} className={`flex items-center ${index % 2 === 0 ? 'flex-row' : 'flex-row-reverse'}`}>
                <div className={`w-1/2 ${index % 2 === 0 ? 'pr-6 text-right' : 'pl-6 text-left'}`}>
                  <PageCard padding="base">
                    <Badge className="mb-3 text-orange-400 border-orange-400 bg-transparent text-xs">
                      {milestone.year}
                    </Badge>
                    <h3 className="text-lg font-bold text-white mb-2">{milestone.title}</h3>
                    <p className="text-gray-300 text-sm">{milestone.description}</p>
                  </PageCard>
                </div>

                {/* Timeline dot */}
                <div className="relative z-10">
                  <div className="w-4 h-4 rounded-full bg-gradient-to-br from-orange-500 to-orange-600 border-2 border-black shadow-lg"></div>
                </div>

                <div className="w-1/2"></div>
              </div>
            ))}
          </div>
        </div>
      </PageSection>

      {/* Team */}
      <PageSection
        size="lg"
        title="Meet Our Team"
        subtitle="The passionate people behind our success, working every day to make your experience better."
        className="bg-gray-900/50"
      >
        <PageGrid type="features">
          {team.map((member, index) => (
            <PageCard key={index} padding="base" className="text-center group hover:border-gray-600 transition-all duration-300">
              <div className="w-16 h-16 rounded-full overflow-hidden mx-auto mb-4 bg-gradient-to-br from-gray-700 to-gray-600 flex items-center justify-center">
                <Users className="h-8 w-8 text-orange-400" />
              </div>
              <h3 className="text-lg font-bold text-white mb-1">{member.name}</h3>
              <p className="text-orange-400 font-semibold mb-2 text-sm">{member.role}</p>
              <p className="text-gray-300 text-xs">{member.description}</p>
            </PageCard>
          ))}
        </PageGrid>
      </PageSection>

      {/* Contact CTA */}
      <PageSection size="lg">
        <PageCard padding="lg" className="bg-gradient-to-br from-gray-800 to-gray-700 text-center">
          <h2 className="text-2xl lg:text-3xl font-bold text-white mb-4">
            Ready to Start Your Journey?
          </h2>
          <p className="text-base text-gray-300 mb-8 max-w-xl mx-auto">
            Join millions of satisfied customers and discover why BiHub is the
            preferred choice for online shopping.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <Button size="default" className="bg-orange-500 hover:bg-orange-600 text-white">
              Start Shopping
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>

            <div className="flex items-center gap-4 text-xs text-gray-300">
              <div className="flex items-center gap-1">
                <CheckCircle className="h-4 w-4 text-emerald-400" />
                <span>Free Shipping</span>
              </div>
              <div className="flex items-center gap-1">
                <CheckCircle className="h-4 w-4 text-emerald-400" />
                <span>30-Day Returns</span>
              </div>
              <div className="flex items-center gap-1">
                <CheckCircle className="h-4 w-4 text-emerald-400" />
                <span>24/7 Support</span>
              </div>
            </div>
          </div>
        </PageCard>
      </PageSection>
    </PageWrapper>
  )
}
