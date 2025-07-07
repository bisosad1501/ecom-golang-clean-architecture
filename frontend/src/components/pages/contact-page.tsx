'use client'

import { useState, useEffect } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import {
  Mail,
  Phone,
  MapPin,
  Clock,
  MessageCircle,
  Send,
  CheckCircle,
  HelpCircle,
  ShoppingBag,
  Truck,
  CreditCard,
  Users,
} from 'lucide-react'
import { toast } from 'sonner'
import {
  PageWrapper,
  PageHeader,
  PageSection,
  PageContainer,
  PageGrid,
  PageCard
} from '@/components/layout'
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'

function ContactPage() {
  const [isHydrated, setIsHydrated] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    category: '',
    message: '',
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    setIsHydrated(true)
  }, [])

  const contactMethods = [
    {
      icon: Mail,
      title: 'Email Support',
      description: 'Get help via email',
      value: 'support@ecomstore.com',
      action: 'Send Email',
      available: '24/7',
    },
    {
      icon: Phone,
      title: 'Phone Support',
      description: 'Speak with our team',
      value: '+1 (555) 123-4567',
      action: 'Call Now',
      available: 'Mon-Fri 9AM-6PM',
    },
    {
      icon: MessageCircle,
      title: 'Live Chat',
      description: 'Chat with support',
      value: 'Available now',
      action: 'Start Chat',
      available: 'Mon-Fri 9AM-6PM',
    },
    {
      icon: HelpCircle,
      title: 'Help Center',
      description: 'Browse FAQs & guides',
      value: 'Self-service support',
      action: 'Visit Help Center',
      available: '24/7',
    },
  ]

  const categories = [
    { value: 'general', label: 'General Inquiry', icon: HelpCircle },
    { value: 'orders', label: 'Order Support', icon: ShoppingBag },
    { value: 'shipping', label: 'Shipping & Delivery', icon: Truck },
    { value: 'payments', label: 'Payment Issues', icon: CreditCard },
    { value: 'account', label: 'Account Support', icon: Users },
    { value: 'technical', label: 'Technical Support', icon: MessageCircle },
  ]

  const officeLocations = [
    {
      city: 'San Francisco',
      address: '123 Market Street, Suite 100',
      postal: 'San Francisco, CA 94105',
      phone: '+1 (555) 123-4567',
      hours: 'Mon-Fri 9AM-6PM PST',
    },
    {
      city: 'New York',
      address: '456 Broadway, Floor 15',
      postal: 'New York, NY 10013',
      phone: '+1 (555) 987-6543',
      hours: 'Mon-Fri 9AM-6PM EST',
    },
    {
      city: 'London',
      address: '789 Oxford Street',
      postal: 'London W1C 1JN, UK',
      phone: '+44 20 7123 4567',
      hours: 'Mon-Fri 9AM-5PM GMT',
    },
  ]

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      // TODO: Implement form submission
      await new Promise(resolve => setTimeout(resolve, 2000)) // Simulate API call
      
      toast.success('Message sent successfully! We\'ll get back to you soon.')
      setFormData({
        name: '',
        email: '',
        subject: '',
        category: '',
        message: '',
      })
    } catch (error) {
      toast.error('Failed to send message. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  const isFormValid = formData.name && formData.email && formData.subject && formData.message

  // Prevent hydration mismatch by not rendering dynamic content until hydrated
  if (!isHydrated) {
    return (
      <PageWrapper>
        <PageSection size="xl" className="text-center">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-800 rounded w-48 mx-auto mb-4"></div>
            <div className="h-12 bg-gray-800 rounded w-96 mx-auto mb-4"></div>
            <div className="h-6 bg-gray-800 rounded w-80 mx-auto"></div>
          </div>
        </PageSection>
      </PageWrapper>
    )
  }

  return (
    <PageWrapper>
      {/* Hero Section */}
      <PageSection size="xl" className="text-center">
        <div className="flex items-center justify-center gap-2 mb-4">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-orange-500 to-orange-400 flex items-center justify-center shadow-lg">
            <MessageCircle className="h-4 w-4 text-white" />
          </div>
          <span className="text-orange-400 font-semibold text-sm">CONTACT US</span>
        </div>

        <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold text-white mb-4">
          Get in <span className="text-orange-400">Touch</span>
        </h1>
        <p className="text-base lg:text-lg text-gray-200 max-w-2xl mx-auto leading-relaxed mb-8">
          Have a question, need support, or want to share feedback? We're here to help!
          Choose the best way to reach us and we'll get back to you as soon as possible.
        </p>

        {/* Contact Methods */}
        <PageGrid type="features">
          {contactMethods.map((method, index) => (
            <PageCard key={index} padding="base" className="text-center group bg-gray-900 border-gray-600 hover:border-gray-500 transition-all duration-300">
              <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-orange-500 to-orange-400 flex items-center justify-center mx-auto mb-4 shadow-lg group-hover:scale-110 transition-transform duration-300">
                <method.icon className="h-6 w-6 text-white" />
              </div>
              <h3 className="text-lg font-bold text-white mb-2">{method.title}</h3>
              <p className="text-gray-200 mb-3 text-sm">{method.description}</p>
              <p className="text-orange-400 font-semibold mb-3 text-sm">{method.value}</p>
              <Badge className="mb-3 bg-gray-800 text-gray-200 border border-gray-600 text-xs">
                <Clock className="h-3 w-3 mr-1" />
                {method.available}
              </Badge>
              <Button variant="outline" size="sm" className="w-full border-gray-500 text-white hover:bg-gray-800 hover:border-gray-400">
                {method.action}
              </Button>
            </PageCard>
          ))}
        </PageGrid>
      </PageSection>

      {/* Contact Form & Info */}
      <PageSection size="lg" className="bg-gray-900/50">
        <PageGrid type="twoColumn">
          {/* Contact Form */}
          <div>
            <PageCard padding="lg">
              <h2 className="text-xl font-bold text-white mb-4">Send us a Message</h2>

              <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-semibold text-gray-200 mb-2">
                          Full Name *
                        </label>
                        <Input
                          placeholder="Your full name"
                          value={formData.name}
                          onChange={(e) => handleInputChange('name', e.target.value)}
                          required
                          className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500"
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-semibold text-gray-200 mb-2">
                          Email Address *
                        </label>
                        <Input
                          type="email"
                          placeholder="your.email@example.com"
                          value={formData.email}
                          onChange={(e) => handleInputChange('email', e.target.value)}
                          required
                          className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500"
                        />
                      </div>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-gray-200 mb-2">
                        Category
                      </label>
                      <select 
                        value={formData.category} 
                        onChange={(e) => handleInputChange('category', e.target.value)}
                        className="w-full bg-gray-800 border border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 rounded-md p-3"
                      >
                        <option value="">Select a category</option>
                        {categories.map((category) => (
                          <option key={category.value} value={category.value}>
                            {category.label}
                          </option>
                        ))}
                      </select>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-gray-200 mb-2">
                        Subject *
                      </label>
                      <Input
                        placeholder="Brief description of your inquiry"
                        value={formData.subject}
                        onChange={(e) => handleInputChange('subject', e.target.value)}
                        required
                        className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500"
                      />
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-gray-200 mb-2">
                        Message *
                      </label>
                      <Textarea
                        placeholder="Please provide details about your inquiry..."
                        rows={6}
                        value={formData.message}
                        onChange={(e) => handleInputChange('message', e.target.value)}
                        required
                        className="bg-gray-800 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500"
                      />
                    </div>
                    
                    <Button 
                      type="submit" 
                      size="lg" 
                      variant="gradient"
                      disabled={!isFormValid || isSubmitting}
                      className="w-full"
                    >
                      {isSubmitting ? (
                        <>
                          <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                          Sending...
                        </>
                      ) : (
                        <>
                          <Send className="h-5 w-5 mr-2" />
                          Send Message
                        </>
                      )}
                    </Button>
                  </form>
            </PageCard>
          </div>

            {/* Office Locations */}
            <div>
              <h2 className="text-xl font-bold text-white mb-6">Our Offices</h2>
              
              <div className="space-y-6">
                {officeLocations.map((office, index) => (
                  <PageCard key={index} padding="base" className="bg-gray-900 border-gray-600">
                    <div className="flex items-start gap-4">
                      <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500 to-orange-400 flex items-center justify-center shadow-lg">
                        <MapPin className="h-6 w-6 text-white" />
                      </div>
                      
                      <div className="flex-1">
                        <h3 className="text-lg font-bold text-white mb-2">{office.city}</h3>
                        <div className="space-y-2 text-sm text-gray-300">
                          <p>{office.address}</p>
                          <p>{office.postal}</p>
                          <div className="flex items-center gap-2">
                            <Phone className="h-4 w-4" />
                            <span>{office.phone}</span>
                          </div>
                          <div className="flex items-center gap-2">
                            <Clock className="h-4 w-4" />
                            <span>{office.hours}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </PageCard>
                ))}
              </div>

              {/* Quick Tips */}
              <PageCard padding="base" className="bg-gray-900 border-gray-600 mt-6">
                <h3 className="text-lg font-bold text-white mb-4">Quick Tips</h3>
                <div className="space-y-3">
                  <div className="flex items-start gap-3">
                    <CheckCircle className="h-5 w-5 text-emerald-400 mt-0.5" />
                    <p className="text-sm text-gray-300">
                      For order-related inquiries, please include your order number
                    </p>
                  </div>
                  <div className="flex items-start gap-3">
                    <CheckCircle className="h-5 w-5 text-emerald-400 mt-0.5" />
                    <p className="text-sm text-gray-300">
                      Check our Help Center for instant answers to common questions
                    </p>
                  </div>
                  <div className="flex items-start gap-3">
                    <CheckCircle className="h-5 w-5 text-emerald-400 mt-0.5" />
                    <p className="text-sm text-gray-300">
                      We typically respond to emails within 24 hours
                    </p>
                  </div>
                </div>
              </PageCard>
            </div>
        </PageGrid>
      </PageSection>

      {/* FAQ Section */}
      <PageSection size="lg">
        <div className="text-center mb-12">
          <h2 className="text-3xl lg:text-4xl font-bold text-white mb-6">
            Frequently Asked Questions
          </h2>
          <p className="text-lg text-gray-300 max-w-3xl mx-auto">
            Find quick answers to the most common questions our customers ask.
          </p>
        </div>

        <PageGrid type="features">
          {[
            {
              question: 'How can I track my order?',
              answer: 'You can track your order using the tracking number sent to your email, or log into your account to view order status.',
            },
            {
              question: 'What is your return policy?',
              answer: 'We offer a 30-day return policy for most items. Items must be in original condition with tags attached.',
            },
            {
              question: 'Do you offer international shipping?',
              answer: 'Yes, we ship to over 25 countries worldwide. Shipping costs and delivery times vary by location.',
            },
            {
              question: 'How do I change my order?',
              answer: 'Orders can be modified within 1 hour of placement. After that, please contact our support team for assistance.',
            },
            {
              question: 'What payment methods do you accept?',
              answer: 'We accept all major credit cards, PayPal, Apple Pay, Google Pay, and other secure payment methods.',
            },
            {
              question: 'How do I create an account?',
              answer: 'Click the "Sign Up" button in the top right corner and follow the simple registration process.',
            },
          ].map((faq, index) => (
            <PageCard key={index} padding="base" className="bg-gray-900 border-gray-600">
              <h3 className="text-lg font-bold text-white mb-3">{faq.question}</h3>
              <p className="text-gray-300 text-sm leading-relaxed">{faq.answer}</p>
            </PageCard>
          ))}
        </PageGrid>

        <div className="text-center mt-8">
          <Button size="default" variant="outline" className="border-gray-600 text-white hover:bg-gray-800">
            View All FAQs
          </Button>
        </div>
      </PageSection>
    </PageWrapper>
  )
}

export default ContactPage
