'use client'

import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
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

export default function ContactPage() {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    subject: '',
    category: '',
    message: '',
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

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

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background">
      {/* Hero Section */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large">
                <MessageCircle className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">CONTACT US</span>
            </div>
            
            <h1 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
              Get in <span className="text-gradient">Touch</span>
            </h1>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              Have a question, need support, or want to share feedback? We're here to help! 
              Choose the best way to reach us and we'll get back to you as soon as possible.
            </p>
          </div>

          {/* Contact Methods */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-16">
            {contactMethods.map((method, index) => (
              <Card key={index} variant="elevated" className="border-0 shadow-large text-center group hover:shadow-xl transition-all duration-300">
                <CardContent className="p-8">
                  <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center mx-auto mb-6 shadow-large group-hover:scale-110 transition-transform duration-300">
                    <method.icon className="h-8 w-8 text-white" />
                  </div>
                  <h3 className="text-xl font-bold text-foreground mb-2">{method.title}</h3>
                  <p className="text-muted-foreground mb-4">{method.description}</p>
                  <p className="text-primary font-semibold mb-4">{method.value}</p>
                  <Badge variant="outline" className="mb-4">
                    <Clock className="h-3 w-3 mr-1" />
                    {method.available}
                  </Badge>
                  <Button variant="outline" size="sm" className="w-full">
                    {method.action}
                  </Button>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Contact Form & Info */}
      <section className="py-24 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-16">
            {/* Contact Form */}
            <div>
              <Card variant="elevated" className="border-0 shadow-xl">
                <CardContent className="p-8">
                  <h2 className="text-2xl font-bold text-foreground mb-6">Send us a Message</h2>
                  
                  <form onSubmit={handleSubmit} className="space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Full Name *
                        </label>
                        <Input
                          placeholder="Your full name"
                          value={formData.name}
                          onChange={(e) => handleInputChange('name', e.target.value)}
                          required
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Email Address *
                        </label>
                        <Input
                          type="email"
                          placeholder="your.email@example.com"
                          value={formData.email}
                          onChange={(e) => handleInputChange('email', e.target.value)}
                          required
                        />
                      </div>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Category
                      </label>
                      <Select value={formData.category} onValueChange={(value) => handleInputChange('category', value)}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select a category" />
                        </SelectTrigger>
                        <SelectContent>
                          {categories.map((category) => (
                            <SelectItem key={category.value} value={category.value}>
                              <div className="flex items-center gap-2">
                                <category.icon className="h-4 w-4" />
                                {category.label}
                              </div>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Subject *
                      </label>
                      <Input
                        placeholder="Brief description of your inquiry"
                        value={formData.subject}
                        onChange={(e) => handleInputChange('subject', e.target.value)}
                        required
                      />
                    </div>
                    
                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Message *
                      </label>
                      <Textarea
                        placeholder="Please provide details about your inquiry..."
                        rows={6}
                        value={formData.message}
                        onChange={(e) => handleInputChange('message', e.target.value)}
                        required
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
                </CardContent>
              </Card>
            </div>

            {/* Office Locations */}
            <div>
              <h2 className="text-2xl font-bold text-foreground mb-8">Our Offices</h2>
              
              <div className="space-y-6">
                {officeLocations.map((office, index) => (
                  <Card key={index} variant="elevated" className="border-0 shadow-large">
                    <CardContent className="p-6">
                      <div className="flex items-start gap-4">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-medium">
                          <MapPin className="h-6 w-6 text-white" />
                        </div>
                        
                        <div className="flex-1">
                          <h3 className="text-lg font-bold text-foreground mb-2">{office.city}</h3>
                          <div className="space-y-2 text-sm text-muted-foreground">
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
                    </CardContent>
                  </Card>
                ))}
              </div>

              {/* Quick Tips */}
              <Card variant="elevated" className="border-0 shadow-large mt-8">
                <CardContent className="p-6">
                  <h3 className="text-lg font-bold text-foreground mb-4">Quick Tips</h3>
                  <div className="space-y-3">
                    <div className="flex items-start gap-3">
                      <CheckCircle className="h-5 w-5 text-emerald-600 mt-0.5" />
                      <p className="text-sm text-muted-foreground">
                        For order-related inquiries, please include your order number
                      </p>
                    </div>
                    <div className="flex items-start gap-3">
                      <CheckCircle className="h-5 w-5 text-emerald-600 mt-0.5" />
                      <p className="text-sm text-muted-foreground">
                        Check our Help Center for instant answers to common questions
                      </p>
                    </div>
                    <div className="flex items-start gap-3">
                      <CheckCircle className="h-5 w-5 text-emerald-600 mt-0.5" />
                      <p className="text-sm text-muted-foreground">
                        We typically respond to emails within 24 hours
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ Section */}
      <section className="py-24">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-6">
              Frequently Asked Questions
            </h2>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
              Find quick answers to the most common questions our customers ask.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
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
              <Card key={index} variant="elevated" className="border-0 shadow-large">
                <CardContent className="p-6">
                  <h3 className="text-lg font-bold text-foreground mb-3">{faq.question}</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed">{faq.answer}</p>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="text-center mt-12">
            <Button size="lg" variant="outline">
              View All FAQs
            </Button>
          </div>
        </div>
      </section>
    </div>
  )
}
