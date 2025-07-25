'use client'

import { useState } from 'react'
import Link from 'next/link'
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
  ArrowRight,
  Sparkles,
  Headphones,
  Globe,
  Zap
} from 'lucide-react'
import { toast } from 'sonner'
import { AnimatedBackground } from '@/components/ui/animated-background'

export function ContactPage() {
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
      description: 'Get help via email within 24 hours',
      value: 'support@bihub.com',
      action: 'Send Email',
      available: '24/7',
      gradient: 'from-blue-500/20 to-cyan-500/20',
      color: 'from-blue-500 to-cyan-500'
    },
    {
      icon: Phone,
      title: 'Phone Support',
      description: 'Speak directly with our expert team',
      value: '+1 (555) 123-4567',
      action: 'Call Now',
      available: 'Mon-Fri 9AM-6PM EST',
      gradient: 'from-green-500/20 to-emerald-500/20',
      color: 'from-green-500 to-emerald-500'
    },
    {
      icon: MessageCircle,
      title: 'Live Chat',
      description: 'Instant chat with support agents',
      value: 'Available now',
      action: 'Start Chat',
      available: 'Mon-Fri 9AM-6PM EST',
      gradient: 'from-purple-500/20 to-pink-500/20',
      color: 'from-purple-500 to-pink-500'
    },
    {
      icon: HelpCircle,
      title: 'Help Center',
      description: 'Browse FAQs and detailed guides',
      value: 'Self-service support',
      action: 'Visit Help Center',
      available: '24/7',
      gradient: 'from-orange-500/20 to-red-500/20',
      color: 'from-orange-500 to-red-500'
    },
  ]

  const supportCategories = [
    {
      icon: ShoppingBag,
      title: 'Orders & Shopping',
      description: 'Order status, returns, exchanges',
      gradient: 'from-blue-500/20 to-cyan-500/20'
    },
    {
      icon: Truck,
      title: 'Shipping & Delivery',
      description: 'Tracking, delivery issues, shipping info',
      gradient: 'from-green-500/20 to-emerald-500/20'
    },
    {
      icon: CreditCard,
      title: 'Payments & Billing',
      description: 'Payment methods, billing questions',
      gradient: 'from-purple-500/20 to-pink-500/20'
    },
    {
      icon: Users,
      title: 'Account & Profile',
      description: 'Account settings, profile management',
      gradient: 'from-orange-500/20 to-red-500/20'
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
    <div className="min-h-screen bg-black text-white">
      {/* Hero Section - Matching home page style */}
      <section className="relative bg-gradient-to-br from-black via-gray-900 to-black text-white overflow-hidden min-h-[60vh] flex items-center py-16">
        <AnimatedBackground variant="hero" />

        <div className="container mx-auto px-4 relative z-10">
          <div className="max-w-6xl mx-auto">
            <div className="text-center mb-16">
              {/* Badge */}
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6 animate-fade-in">
                <MessageCircle className="h-4 w-4 text-[#ff9000]" />
                <span className="text-[#ff9000] font-semibold text-sm">CONTACT US</span>
              </div>

              {/* Main heading */}
              <h1 className="text-4xl lg:text-6xl font-bold mb-6 animate-scale-in">
                Get in
                <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Touch</span>
              </h1>

              <p className="text-lg text-gray-300 max-w-3xl mx-auto leading-relaxed mb-8 animate-fade-in" style={{ animationDelay: '0.2s' }}>
                Have a question, need support, or want to share feedback? We're here to help!
                Choose the best way to reach us and we'll get back to you as soon as possible.
              </p>

              {/* CTA Buttons */}
              <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12 animate-fade-in" style={{ animationDelay: '0.4s' }}>
                <Button size="lg" className="shadow-lg hover:shadow-xl transition-all duration-300 group text-white" style={{backgroundColor: '#FF9000'}}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
                >
                  <Headphones className="mr-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
                  Start Live Chat
                  <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </Button>
                <Button size="lg" variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}} asChild>
                  <Link href="/about">
                    <Globe className="mr-2 h-5 w-5" />
                    Learn About Us
                  </Link>
                </Button>
              </div>
            </div>

            {/* Contact Methods Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 animate-fade-in" style={{ animationDelay: '0.6s' }}>
              {contactMethods.map((method, index) => (
                <div key={index} className="group">
                  <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300 hover:scale-105">
                    <CardContent className="p-6 text-center">
                      <div className={`w-16 h-16 rounded-2xl bg-gradient-to-br ${method.gradient} border border-white/20 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-300`}>
                        <method.icon className="h-8 w-8 text-white" />
                      </div>
                      <h3 className="text-xl font-bold text-white mb-2 group-hover:text-[#ff9000] transition-colors">
                        {method.title}
                      </h3>
                      <p className="text-gray-300 text-sm mb-3 leading-relaxed">
                        {method.description}
                      </p>
                      <p className="text-[#ff9000] font-semibold mb-3 text-sm">{method.value}</p>
                      <Badge className="mb-4 bg-[#ff9000]/20 text-[#ff9000] border border-[#ff9000]/30">
                        <Clock className="h-3 w-3 mr-1" />
                        {method.available}
                      </Badge>
                      <Button variant="outline" size="sm" className="w-full border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-[#ff9000]/50 hover:text-[#ff9000] transition-all">
                        {method.action}
                      </Button>
                    </CardContent>
                  </Card>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Support Categories Section */}
      <section className="py-20 bg-gradient-to-b from-black to-gray-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <HelpCircle className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">SUPPORT CATEGORIES</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold text-white mb-4">
              How Can We Help You?
            </h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">
              Choose the category that best matches your inquiry for faster assistance
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            {supportCategories.map((category, index) => (
              <div key={index} className="group animate-scale-in" style={{ animationDelay: `${index * 0.1}s` }}>
                <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300 hover:scale-105 cursor-pointer">
                  <CardContent className="p-6 text-center">
                    <div className={`w-16 h-16 rounded-2xl bg-gradient-to-br ${category.gradient} border border-white/20 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-300`}>
                      <category.icon className="h-8 w-8 text-white" />
                    </div>
                    <h3 className="text-xl font-bold text-white mb-3 group-hover:text-[#ff9000] transition-colors">
                      {category.title}
                    </h3>
                    <p className="text-gray-300 text-sm leading-relaxed">
                      {category.description}
                    </p>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Contact Form Section */}
      <section className="py-20 bg-gradient-to-b from-gray-900 to-black">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-start">
              {/* Left Content - Contact Info */}
              <div className="animate-fade-in">
                <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
                  <Send className="h-4 w-4 text-[#ff9000]" />
                  <span className="text-[#ff9000] font-semibold text-sm">SEND MESSAGE</span>
                </div>
                <h2 className="text-3xl lg:text-4xl font-bold text-white mb-6">
                  Send Us a Message
                </h2>
                <p className="text-gray-300 leading-relaxed mb-8">
                  Fill out the form and we'll get back to you within 24 hours.
                  For urgent matters, please use our live chat or phone support.
                </p>

                {/* Contact Info Cards */}
                <div className="space-y-4">
                  <div className="flex items-center gap-4 p-4 rounded-xl bg-white/[0.08] backdrop-blur-sm border border-white/15">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 border border-blue-500/30 flex items-center justify-center">
                      <Mail className="h-6 w-6 text-blue-400" />
                    </div>
                    <div>
                      <h4 className="font-semibold text-white">Email</h4>
                      <p className="text-gray-300 text-sm">support@bihub.com</p>
                    </div>
                  </div>

                  <div className="flex items-center gap-4 p-4 rounded-xl bg-white/[0.08] backdrop-blur-sm border border-white/15">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 border border-green-500/30 flex items-center justify-center">
                      <Phone className="h-6 w-6 text-green-400" />
                    </div>
                    <div>
                      <h4 className="font-semibold text-white">Phone</h4>
                      <p className="text-gray-300 text-sm">+1 (555) 123-4567</p>
                    </div>
                  </div>

                  <div className="flex items-center gap-4 p-4 rounded-xl bg-white/[0.08] backdrop-blur-sm border border-white/15">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 border border-purple-500/30 flex items-center justify-center">
                      <MapPin className="h-6 w-6 text-purple-400" />
                    </div>
                    <div>
                      <h4 className="font-semibold text-white">Address</h4>
                      <p className="text-gray-300 text-sm">123 Business St, City, State 12345</p>
                    </div>
                  </div>
                </div>
              </div>

              {/* Right Content - Contact Form */}
              <div className="animate-scale-in" style={{ animationDelay: '0.2s' }}>
                <Card className="overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15">
                  <CardContent className="p-8">
                    <form onSubmit={handleSubmit} className="space-y-6">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <label className="block text-sm font-medium text-white mb-2">Name *</label>
                          <Input
                            type="text"
                            value={formData.name}
                            onChange={(e) => handleInputChange('name', e.target.value)}
                            className="bg-gray-900 border-gray-700 text-white placeholder:text-gray-400 focus:border-[#ff9000] focus:ring-[#ff9000]"
                            placeholder="Your full name"
                            required
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-white mb-2">Email *</label>
                          <Input
                            type="email"
                            value={formData.email}
                            onChange={(e) => handleInputChange('email', e.target.value)}
                            className="bg-gray-900 border-gray-700 text-white placeholder:text-gray-400 focus:border-[#ff9000] focus:ring-[#ff9000]"
                            placeholder="your@email.com"
                            required
                          />
                        </div>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-white mb-2">Subject *</label>
                        <Input
                          type="text"
                          value={formData.subject}
                          onChange={(e) => handleInputChange('subject', e.target.value)}
                          className="bg-gray-900 border-gray-700 text-white placeholder:text-gray-400 focus:border-[#ff9000] focus:ring-[#ff9000]"
                          placeholder="What's this about?"
                          required
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-white mb-2">Category</label>
                        <select
                          value={formData.category}
                          onChange={(e) => handleInputChange('category', e.target.value)}
                          className="w-full bg-gray-900 border border-gray-700 text-white rounded-md px-3 py-2 focus:border-[#ff9000] focus:ring-[#ff9000] focus:outline-none"
                        >
                          <option value="">Select a category</option>
                          <option value="general">General Inquiry</option>
                          <option value="orders">Order Support</option>
                          <option value="shipping">Shipping & Delivery</option>
                          <option value="payments">Payment Issues</option>
                          <option value="account">Account Support</option>
                          <option value="technical">Technical Support</option>
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-white mb-2">Message *</label>
                        <Textarea
                          value={formData.message}
                          onChange={(e) => handleInputChange('message', e.target.value)}
                          className="bg-gray-900 border-gray-700 text-white placeholder:text-gray-400 focus:border-[#ff9000] focus:ring-[#ff9000] min-h-[120px]"
                          placeholder="Tell us more about your inquiry..."
                          required
                        />
                      </div>

                      <Button
                        type="submit"
                        disabled={!isFormValid || isSubmitting}
                        className="w-full bg-[#ff9000] hover:bg-[#e67e00] text-white disabled:opacity-50 disabled:cursor-not-allowed group"
                      >
                        {isSubmitting ? (
                          <>
                            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                            Sending...
                          </>
                        ) : (
                          <>
                            <Send className="mr-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                            Send Message
                          </>
                        )}
                      </Button>
                    </form>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ Section */}
      <section className="py-20 bg-gradient-to-b from-black to-gray-900">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-[#ff9000]/30 bg-[#ff9000]/10 backdrop-blur-sm mb-6">
              <HelpCircle className="h-4 w-4 text-[#ff9000]" />
              <span className="text-[#ff9000] font-semibold text-sm">QUICK ANSWERS</span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold text-white mb-4">
              Frequently Asked Questions
            </h2>
            <p className="text-lg text-gray-300 max-w-2xl mx-auto">
              Find quick answers to common questions before reaching out
            </p>
          </div>

          <div className="max-w-4xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-8">
            {[
              {
                question: "How long does shipping take?",
                answer: "Standard shipping takes 3-5 business days. Express shipping is available for 1-2 day delivery.",
                icon: Truck
              },
              {
                question: "What's your return policy?",
                answer: "We offer 30-day returns on most items. Items must be in original condition with tags attached.",
                icon: CheckCircle
              },
              {
                question: "How can I track my order?",
                answer: "You'll receive a tracking number via email once your order ships. You can also track in your account.",
                icon: MapPin
              },
              {
                question: "Do you offer customer support?",
                answer: "Yes! We offer 24/7 email support and live chat Monday-Friday 9AM-6PM EST.",
                icon: Headphones
              }
            ].map((faq, index) => (
              <div key={index} className="group animate-scale-in" style={{ animationDelay: `${index * 0.1}s` }}>
                <Card className="h-full overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 hover:bg-white/[0.12] hover:border-[#ff9000]/50 transition-all duration-300">
                  <CardContent className="p-6">
                    <div className="flex items-start gap-4">
                      <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 border border-[#ff9000]/30 flex items-center justify-center flex-shrink-0 group-hover:scale-110 transition-transform duration-300">
                        <faq.icon className="h-6 w-6 text-[#ff9000]" />
                      </div>
                      <div>
                        <h3 className="text-lg font-bold text-white mb-2 group-hover:text-[#ff9000] transition-colors">
                          {faq.question}
                        </h3>
                        <p className="text-gray-300 text-sm leading-relaxed">
                          {faq.answer}
                        </p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>

          <div className="text-center mt-12">
            <Button variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}}>
              <HelpCircle className="mr-2 h-5 w-5" />
              View All FAQs
            </Button>
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
              <span className="text-[#ff9000] font-semibold text-sm">GET STARTED</span>
            </div>

            <h2 className="text-3xl lg:text-5xl font-bold text-white mb-6">
              Ready to Experience
              <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Amazing Support?</span>
            </h2>

            <p className="text-lg text-gray-300 max-w-2xl mx-auto mb-8 leading-relaxed">
              Join thousands of satisfied customers who love our exceptional support and service.
              We're here to help you every step of the way.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="lg" className="shadow-lg hover:shadow-xl transition-all duration-300 group text-white" style={{backgroundColor: '#FF9000'}}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
              >
                <Zap className="mr-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
                Start Shopping
                <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Button>
              <Button size="lg" variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}} asChild>
                <Link href="/about">
                  <Users className="mr-2 h-5 w-5" />
                  About Our Team
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}


