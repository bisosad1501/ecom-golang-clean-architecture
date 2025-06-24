import Link from 'next/link'
import {
  Mail,
  Phone,
  MapPin,
  CreditCard,
  Shield,
  Truck,
  Award,
  ExternalLink,
  Globe,
  MessageCircle,
  Video,
  Send
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { APP_NAME, SOCIAL_LINKS, CONTACT_INFO } from '@/constants'

export function Footer() {
  return (
    <footer className="bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 text-white relative overflow-hidden">
      {/* Background pattern */}
      <div className="absolute inset-0 opacity-5">
        <svg className="w-full h-full" viewBox="0 0 100 100" fill="none">
          <pattern id="footerGrid" width="10" height="10" patternUnits="userSpaceOnUse">
            <path d="M 10 0 L 0 0 0 10" fill="none" stroke="currentColor" strokeWidth="0.5"/>
          </pattern>
          <rect width="100%" height="100%" fill="url(#footerGrid)" />
        </svg>
      </div>

      {/* Main footer content */}
      <div className="container mx-auto px-4 py-20 relative z-10">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12">
          {/* Company info */}
          <div className="space-y-8 lg:col-span-2">
            <div className="flex items-center space-x-4">
              <div className="h-14 w-14 rounded-2xl bg-gradient-to-br from-primary-500 via-primary-600 to-violet-600 flex items-center justify-center shadow-2xl">
                <span className="text-white font-bold text-2xl">E</span>
              </div>
              <span className="text-3xl font-bold text-gradient bg-gradient-to-r from-white via-violet-200 to-white bg-clip-text text-transparent">{APP_NAME}</span>
            </div>
            <p className="text-slate-300 text-lg leading-relaxed max-w-md">
              Elevating your shopping experience with premium products, exceptional service,
              and innovative solutions that exceed expectations.
            </p>

            {/* Trust badges */}
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-2 bg-slate-800/50 rounded-full px-4 py-2 border border-slate-700/50">
                <Shield className="h-4 w-4 text-green-400" />
                <span className="text-sm font-medium text-slate-300">SSL Secured</span>
              </div>
              <div className="flex items-center gap-2 bg-slate-800/50 rounded-full px-4 py-2 border border-slate-700/50">
                <Award className="h-4 w-4 text-yellow-400" />
                <span className="text-sm font-medium text-slate-300">Verified Store</span>
              </div>
            </div>
            <div className="flex space-x-4">
              <Link href={SOCIAL_LINKS.facebook} className="w-12 h-12 rounded-2xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-blue-600 hover:to-blue-700 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <Globe className="h-5 w-5 group-hover:text-white transition-colors" />
              </Link>
              <Link href={SOCIAL_LINKS.twitter} className="w-12 h-12 rounded-2xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-sky-500 hover:to-sky-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <MessageCircle className="h-5 w-5 group-hover:text-white transition-colors" />
              </Link>
              <Link href={SOCIAL_LINKS.instagram} className="w-12 h-12 rounded-2xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-pink-500 hover:to-purple-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <Send className="h-5 w-5 group-hover:text-white transition-colors" />
              </Link>
              <Link href={SOCIAL_LINKS.youtube} className="w-12 h-12 rounded-2xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-red-500 hover:to-red-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <Video className="h-5 w-5 group-hover:text-white transition-colors" />
              </Link>
            </div>
          </div>

          {/* Quick links */}
          <div className="space-y-6">
            <h3 className="text-xl font-bold text-white">Quick Links</h3>
            <ul className="space-y-3">
              <li>
                <Link href="/products" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  All Products
                </Link>
              </li>
              <li>
                <Link href="/categories" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Categories
                </Link>
              </li>
              <li>
                <Link href="/deals" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Special Deals
                </Link>
              </li>
              <li>
                <Link href="/new-arrivals" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  New Arrivals
                </Link>
              </li>
              <li>
                <Link href="/bestsellers" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Best Sellers
                </Link>
              </li>
              <li>
                <Link href="/brands" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Brands
                </Link>
              </li>
            </ul>
          </div>

          {/* Customer service */}
          <div className="space-y-6">
            <h3 className="text-xl font-bold text-white">Customer Service</h3>
            <ul className="space-y-3">
              <li>
                <Link href="/help" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Help Center
                </Link>
              </li>
              <li>
                <Link href="/contact" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Contact Us
                </Link>
              </li>
              <li>
                <Link href="/shipping" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Shipping Info
                </Link>
              </li>
              <li>
                <Link href="/returns" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Returns & Exchanges
                </Link>
              </li>
              <li>
                <Link href="/size-guide" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Size Guide
                </Link>
              </li>
              <li>
                <Link href="/track-order" className="text-slate-300 hover:text-primary-400 transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                  Track Your Order
                </Link>
              </li>
            </ul>
          </div>

          {/* Newsletter */}
          <div className="space-y-6">
            <h3 className="text-xl font-bold text-white">Stay Updated</h3>
            <p className="text-slate-300 text-sm leading-relaxed">
              Subscribe to our newsletter for exclusive deals and updates.
            </p>
            <form className="space-y-3">
              <Input
                type="email"
                placeholder="Enter your email"
                className="bg-slate-800/50 border-slate-600 text-white placeholder-slate-400 h-12 rounded-xl focus:ring-2 focus:ring-primary/50 focus:border-primary"
              />
              <Button className="w-full h-12 rounded-xl" variant="gradient">
                Subscribe
              </Button>
            </form>

            {/* Contact info */}
            <div className="space-y-3 pt-6 border-t border-slate-700">
              <div className="flex items-center space-x-3 text-sm text-slate-300">
                <div className="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center">
                  <Mail className="h-4 w-4 text-primary-400" />
                </div>
                <span>{CONTACT_INFO.email}</span>
              </div>
              <div className="flex items-center space-x-3 text-sm text-slate-300">
                <div className="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center">
                  <Phone className="h-4 w-4 text-primary-400" />
                </div>
                <span>{CONTACT_INFO.phone}</span>
              </div>
              <div className="flex items-center space-x-3 text-sm text-slate-300">
                <div className="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center">
                  <MapPin className="h-4 w-4 text-primary-400" />
                </div>
                <span>{CONTACT_INFO.address}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Company Stats & Achievements */}
      <div className="border-t border-slate-700/50 bg-gradient-to-r from-slate-800/50 to-slate-900/50">
        <div className="container mx-auto px-4 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            <div className="text-center group">
              <div className="text-3xl font-bold text-gradient bg-gradient-to-r from-primary-400 to-violet-400 bg-clip-text text-transparent mb-2 group-hover:scale-110 transition-transform duration-300">
                10K+
              </div>
              <div className="text-slate-400 text-sm font-medium">Happy Customers</div>
            </div>

            <div className="text-center group">
              <div className="text-3xl font-bold text-gradient bg-gradient-to-r from-emerald-400 to-green-400 bg-clip-text text-transparent mb-2 group-hover:scale-110 transition-transform duration-300">
                5K+
              </div>
              <div className="text-slate-400 text-sm font-medium">Products Sold</div>
            </div>

            <div className="text-center group">
              <div className="text-3xl font-bold text-gradient bg-gradient-to-r from-yellow-400 to-orange-400 bg-clip-text text-transparent mb-2 group-hover:scale-110 transition-transform duration-300">
                99%
              </div>
              <div className="text-slate-400 text-sm font-medium">Satisfaction Rate</div>
            </div>

            <div className="text-center group">
              <div className="text-3xl font-bold text-gradient bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent mb-2 group-hover:scale-110 transition-transform duration-300">
                24/7
              </div>
              <div className="text-slate-400 text-sm font-medium">Support Available</div>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom bar */}
      <div className="border-t border-slate-700/50 bg-slate-900/50">
        <div className="container mx-auto px-4 py-6">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
            <div className="text-sm text-slate-400 font-medium">
              © 2024 {APP_NAME}. All rights reserved.
            </div>

            <div className="flex items-center space-x-8">
              <Link href="/privacy" className="text-sm text-slate-400 hover:text-primary-400 transition-colors font-medium">
                Privacy Policy
              </Link>
              <Link href="/terms" className="text-sm text-slate-400 hover:text-primary-400 transition-colors font-medium">
                Terms of Service
              </Link>
              <Link href="/cookies" className="text-sm text-slate-400 hover:text-primary-400 transition-colors font-medium">
                Cookie Policy
              </Link>
            </div>

            <div className="flex items-center space-x-3">
              <span className="text-sm text-slate-400 font-medium">We accept:</span>
              <div className="flex space-x-2">
                <div className="w-10 h-6 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">
                  <span className="text-xs text-white font-bold">💳</span>
                </div>
                <div className="w-10 h-6 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">
                  <span className="text-xs text-white font-bold">💳</span>
                </div>
                <div className="w-10 h-6 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">
                  <span className="text-xs text-white font-bold">💳</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
