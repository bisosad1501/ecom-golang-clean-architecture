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
    <footer className="bg-black text-white border-t border-gray-800 relative overflow-hidden">
      {/* Enhanced Background Effects */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-gradient-to-t from-[#FF9000]/5 to-transparent"></div>
        <div className="absolute bottom-0 left-0 w-full h-px bg-gradient-to-r from-transparent via-[#FF9000]/30 to-transparent"></div>
      </div>

      {/* Main footer content */}
      <div className="container mx-auto px-4 py-8 relative z-10">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Quick Links */}
          <div className="space-y-4">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <div className="w-2 h-2 bg-[#FF9000] rounded-full"></div>
              Quick Links
            </h3>
            <div className="grid grid-cols-2 gap-2">
              <Link href="/products" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                Products
              </Link>
              <Link href="/categories" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                Categories
              </Link>
              <Link href="/deals" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                Deals
              </Link>
              <Link href="/new-arrivals" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform">
                New Arrivals
              </Link>
            </div>
          </div>

          {/* Support */}
          <div className="space-y-4">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <div className="w-2 h-2 bg-[#FF9000] rounded-full"></div>
              Support
            </h3>
            <div className="space-y-3">
              <Link href="/help" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform flex items-center gap-2">
                <Mail className="h-4 w-4" />
                Help Center
              </Link>
              <Link href="/contact" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform flex items-center gap-2">
                <MessageCircle className="h-4 w-4" />
                Contact Us
              </Link>
              <Link href="/shipping" className="text-gray-300 hover:text-[#FF9000] transition-colors text-sm font-medium hover:translate-x-1 inline-block transition-transform flex items-center gap-2">
                <Truck className="h-4 w-4" />
                Shipping
              </Link>
            </div>
          </div>

          {/* Newsletter & Social */}
          <div className="space-y-4">
            <h3 className="text-lg font-bold text-white flex items-center gap-2">
              <div className="w-2 h-2 bg-[#FF9000] rounded-full"></div>
              Stay Connected
            </h3>
            
            {/* Newsletter */}
            <div className="space-y-3">
              <p className="text-gray-300 text-sm">
                Get exclusive deals & updates
              </p>
              <form className="space-y-2">
                <Input
                  type="email"
                  placeholder="Enter your email"
                  className="bg-gray-900 border-gray-700 text-white placeholder-gray-400 h-10 rounded-lg focus:ring-2 focus:ring-[#FF9000]/50 focus:border-[#FF9000] text-sm"
                />
                <Button className="w-full h-10 rounded-lg text-sm bg-[#FF9000] hover:bg-[#FF9000]/90 text-white">
                  Subscribe
                </Button>
              </form>
            </div>

            {/* Social Links */}
            <div className="space-y-2">
              <p className="text-gray-300 text-sm">Follow us</p>
              <div className="flex space-x-3">
                <Link href={SOCIAL_LINKS.facebook} className="w-9 h-9 rounded-lg bg-gray-800 hover:bg-[#FF9000] flex items-center justify-center transition-all duration-300 hover:scale-110 group">
                  <Globe className="h-4 w-4 text-gray-400 group-hover:text-white transition-colors" />
                </Link>
                <Link href={SOCIAL_LINKS.twitter} className="w-9 h-9 rounded-lg bg-gray-800 hover:bg-[#FF9000] flex items-center justify-center transition-all duration-300 hover:scale-110 group">
                  <MessageCircle className="h-4 w-4 text-gray-400 group-hover:text-white transition-colors" />
                </Link>
                <Link href={SOCIAL_LINKS.instagram} className="w-9 h-9 rounded-lg bg-gray-800 hover:bg-[#FF9000] flex items-center justify-center transition-all duration-300 hover:scale-110 group">
                  <Send className="h-4 w-4 text-gray-400 group-hover:text-white transition-colors" />
                </Link>
                <Link href={SOCIAL_LINKS.youtube} className="w-9 h-9 rounded-lg bg-gray-800 hover:bg-[#FF9000] flex items-center justify-center transition-all duration-300 hover:scale-110 group">
                  <Video className="h-4 w-4 text-gray-400 group-hover:text-white transition-colors" />
                </Link>
              </div>
            </div>
          </div>        </div>

        {/* Trust Badges & Contact Info */}
        <div className="mt-8 pt-6 border-t border-gray-800 flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
          {/* Trust badges */}
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 bg-gray-900 rounded-full px-3 py-1.5 border border-gray-700">
              <Shield className="h-4 w-4 text-green-400" />
              <span className="text-xs font-medium text-gray-300">Secure</span>
            </div>
            <div className="flex items-center gap-2 bg-gray-900 rounded-full px-3 py-1.5 border border-gray-700">
              <Award className="h-4 w-4 text-[#FF9000]" />
              <span className="text-xs font-medium text-gray-300">Verified</span>
            </div>
            <div className="flex items-center gap-2 bg-gray-900 rounded-full px-3 py-1.5 border border-gray-700">
              <Truck className="h-4 w-4 text-blue-400" />
              <span className="text-xs font-medium text-gray-300">Fast Ship</span>
            </div>
          </div>

          {/* Contact info */}
          <div className="flex items-center space-x-6 text-sm text-gray-400">
            <div className="flex items-center space-x-2">
              <Mail className="h-4 w-4 text-[#FF9000]" />
              <span>{CONTACT_INFO.email}</span>
            </div>
            <div className="flex items-center space-x-2">
              <Phone className="h-4 w-4 text-[#FF9000]" />
              <span>{CONTACT_INFO.phone}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom bar */}
      <div className="border-t border-gray-800 bg-gray-950">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-3 md:space-y-0">
            <div className="text-sm text-gray-400">
              © 2024 All rights reserved.
            </div>

            <div className="flex items-center space-x-6">
              <Link href="/privacy" className="text-sm text-gray-400 hover:text-[#FF9000] transition-colors">
                Privacy
              </Link>
              <Link href="/terms" className="text-sm text-gray-400 hover:text-[#FF9000] transition-colors">
                Terms
              </Link>
              <Link href="/cookies" className="text-sm text-gray-400 hover:text-[#FF9000] transition-colors">
                Cookies
              </Link>
            </div>

            <div className="flex items-center space-x-3">
              <span className="text-sm text-gray-400">Payment:</span>
              <div className="flex space-x-2">
                <div className="w-8 h-6 bg-gray-800 rounded border border-gray-700 flex items-center justify-center">
                  <CreditCard className="h-3 w-3 text-gray-400" />
                </div>
                <div className="w-8 h-6 bg-gray-800 rounded border border-gray-700 flex items-center justify-center">
                  <span className="text-xs text-gray-400">💳</span>
                </div>
                <div className="w-8 h-6 bg-gray-800 rounded border border-gray-700 flex items-center justify-center">
                  <span className="text-xs text-gray-400">🏦</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
