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
    <footer className="bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
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
      <div className="container mx-auto px-4 py-12 relative z-10">{/* Reduced padding from py-20 to py-12 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">{/* Reduced gap from 12 to 8 */}
          {/* Company info */}
          <div className="space-y-6 lg:col-span-2">{/* Reduced spacing from 8 to 6 */}
            <div className="flex items-center space-x-3">{/* Reduced spacing */}
              <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-primary-500 via-primary-600 to-violet-600 flex items-center justify-center shadow-2xl">{/* Smaller logo */}
                <span className="text-white font-bold text-lg">E</span>{/* Smaller text */}
              </div>
              <span className="text-2xl font-bold text-gradient bg-gradient-to-r from-white via-violet-200 to-white bg-clip-text text-transparent">{APP_NAME}</span>{/* Smaller text */}
            </div>
            <p className="text-slate-300 text-base leading-relaxed max-w-md">{/* Smaller text */}
              Elevating your shopping experience with premium products, exceptional service,
              and innovative solutions that exceed expectations.
            </p>

            {/* Trust badges */}
            <div className="flex items-center gap-4">{/* Reduced gap */}
              <div className="flex items-center gap-1.5 bg-slate-800/50 rounded-full px-3 py-1.5 border border-slate-700/50">{/* Smaller padding */}
                <Shield className="h-3.5 w-3.5 text-green-400" />{/* Smaller icon */}
                <span className="text-xs font-medium text-slate-300">SSL Secured</span>{/* Smaller text */}
              </div>
              <div className="flex items-center gap-1.5 bg-slate-800/50 rounded-full px-3 py-1.5 border border-slate-700/50">
                <Award className="h-3.5 w-3.5 text-yellow-400" />
                <span className="text-xs font-medium text-slate-300">Verified Store</span>
              </div>
            </div>
            <div className="flex space-x-3">{/* Reduced spacing */}
              <Link href={SOCIAL_LINKS.facebook} className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-blue-600 hover:to-blue-700 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">{/* Smaller icons */}
                <Globe className="h-4 w-4 group-hover:text-white transition-colors" />{/* Smaller icon */}
              </Link>
              <Link href={SOCIAL_LINKS.twitter} className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-sky-500 hover:to-sky-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <MessageCircle className="h-4 w-4 group-hover:text-white transition-colors" />
              </Link>
              <Link href={SOCIAL_LINKS.instagram} className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-pink-500 hover:to-purple-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <Send className="h-4 w-4 group-hover:text-white transition-colors" />
              </Link>
              <Link href={SOCIAL_LINKS.youtube} className="w-10 h-10 rounded-xl bg-gradient-to-br from-slate-800 to-slate-700 hover:from-red-500 hover:to-red-600 flex items-center justify-center transition-all duration-300 hover:scale-110 shadow-large hover:shadow-xl group">
                <Video className="h-4 w-4 group-hover:text-white transition-colors" />
              </Link>
            </div>
          </div>

          {/* Quick links */}
          <div className="space-y-4">{/* Reduced spacing */}
            <h3 className="text-lg font-bold text-white">Quick Links</h3>{/* Smaller heading */}
            <ul className="space-y-2">{/* Reduced spacing */}
              <li>
                <Link href="/products" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">{/* Smaller text */}
                  All Products
                </Link>
              </li>
              <li>
                <Link href="/categories" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Categories
                </Link>
              </li>
              <li>
                <Link href="/deals" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Special Deals
                </Link>
              </li>
              <li>
                <Link href="/new-arrivals" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  New Arrivals
                </Link>
              </li>
              <li>
                <Link href="/bestsellers" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Best Sellers
                </Link>
              </li>
              <li>
                <Link href="/brands" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Brands
                </Link>
              </li>
            </ul>
          </div>

          {/* Customer service */}
          <div className="space-y-4">{/* Reduced spacing */}
            <h3 className="text-lg font-bold text-white">Customer Service</h3>{/* Smaller heading */}
            <ul className="space-y-2">{/* Reduced spacing */}
              <li>
                <Link href="/help" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">{/* Smaller text */}
                  Help Center
                </Link>
              </li>
              <li>
                <Link href="/contact" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Contact Us
                </Link>
              </li>
              <li>
                <Link href="/shipping" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Shipping Info
                </Link>
              </li>
              <li>
                <Link href="/returns" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Returns & Exchanges
                </Link>
              </li>
              <li>
                <Link href="/size-guide" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Size Guide
                </Link>
              </li>
              <li>
                <Link href="/track-order" className="text-slate-300 hover:text-primary-400 transition-colors text-xs font-medium hover:translate-x-1 inline-block transition-transform">
                  Track Your Order
                </Link>
              </li>
            </ul>
          </div>

          {/* Newsletter */}
          <div className="space-y-4">{/* Reduced spacing */}
            <h3 className="text-lg font-bold text-white">Stay Updated</h3>{/* Smaller heading */}
            <p className="text-slate-300 text-xs leading-relaxed">{/* Smaller text */}
              Subscribe to our newsletter for exclusive deals and updates.
            </p>
            <form className="space-y-2">{/* Reduced spacing */}
              <Input
                type="email"
                placeholder="Enter your email"
                className="bg-slate-800/50 border-slate-600 text-white placeholder-slate-400 h-10 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary text-sm"
              />
              <Button className="w-full h-10 rounded-lg text-sm" variant="gradient">{/* Smaller button */}
                Subscribe
              </Button>
            </form>

            {/* Contact info */}
            <div className="space-y-2 pt-4 border-t border-slate-700">{/* Reduced spacing and padding */}
              <div className="flex items-center space-x-2 text-xs text-slate-300">{/* Smaller spacing and text */}
                <div className="w-6 h-6 rounded-lg bg-slate-800 flex items-center justify-center">{/* Smaller icon container */}
                  <Mail className="h-3 w-3 text-primary-400" />{/* Smaller icon */}
                </div>
                <span>{CONTACT_INFO.email}</span>
              </div>
              <div className="flex items-center space-x-2 text-xs text-slate-300">
                <div className="w-6 h-6 rounded-lg bg-slate-800 flex items-center justify-center">
                  <Phone className="h-3 w-3 text-primary-400" />
                </div>
                <span>{CONTACT_INFO.phone}</span>
              </div>
              <div className="flex items-center space-x-2 text-xs text-slate-300">
                <div className="w-6 h-6 rounded-lg bg-slate-800 flex items-center justify-center">
                  <MapPin className="h-3 w-3 text-primary-400" />
                </div>
                <span>{CONTACT_INFO.address}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      
      {/* Bottom bar */}
      <div className="border-t border-slate-700/50 bg-slate-900/50">
        <div className="container mx-auto px-4 py-4">{/* Reduced padding */}
          <div className="flex flex-col md:flex-row justify-between items-center space-y-3 md:space-y-0">{/* Reduced spacing */}
            <div className="text-xs text-slate-400 font-medium">{/* Smaller text */}
              © 2024 {APP_NAME}. All rights reserved.
            </div>

            <div className="flex items-center space-x-6">{/* Reduced spacing */}
              <Link href="/privacy" className="text-xs text-slate-400 hover:text-primary-400 transition-colors font-medium">{/* Smaller text */}
                Privacy Policy
              </Link>
              <Link href="/terms" className="text-xs text-slate-400 hover:text-primary-400 transition-colors font-medium">
                Terms of Service
              </Link>
              <Link href="/cookies" className="text-xs text-slate-400 hover:text-primary-400 transition-colors font-medium">
                Cookie Policy
              </Link>
            </div>

            <div className="flex items-center space-x-2">{/* Reduced spacing */}
              <span className="text-xs text-slate-400 font-medium">We accept:</span>{/* Smaller text */}
              <div className="flex space-x-1.5">{/* Reduced spacing */}
                <div className="w-8 h-5 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">{/* Smaller payment icons */}
                  <span className="text-xs text-white font-bold">💳</span>
                </div>
                <div className="w-8 h-5 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">
                  <span className="text-xs text-white font-bold">💳</span>
                </div>
                <div className="w-8 h-5 bg-slate-800 rounded-lg flex items-center justify-center shadow-soft hover:shadow-medium transition-all">
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
