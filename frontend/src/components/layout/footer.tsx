import Link from 'next/link'
import { 
  Facebook, 
  Twitter, 
  Instagram, 
  Youtube, 
  Mail, 
  Phone, 
  MapPin,
  CreditCard,
  Shield,
  Truck
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { APP_NAME, SOCIAL_LINKS, CONTACT_INFO } from '@/constants'

export function Footer() {
  return (
    <footer className="bg-gray-900 text-white">
      {/* Main footer content */}
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Company info */}
          <div className="space-y-4">
            <div className="flex items-center space-x-2">
              <div className="h-8 w-8 rounded-lg bg-blue-600 flex items-center justify-center">
                <span className="text-white font-bold text-lg">E</span>
              </div>
              <span className="text-xl font-bold">{APP_NAME}</span>
            </div>
            <p className="text-gray-300 text-sm leading-relaxed">
              Your trusted partner for quality products and exceptional service. 
              We bring you the best shopping experience with fast delivery and 
              excellent customer support.
            </p>
            <div className="flex space-x-4">
              <Link href={SOCIAL_LINKS.facebook} className="text-gray-400 hover:text-white transition-colors">
                <Facebook className="h-5 w-5" />
              </Link>
              <Link href={SOCIAL_LINKS.twitter} className="text-gray-400 hover:text-white transition-colors">
                <Twitter className="h-5 w-5" />
              </Link>
              <Link href={SOCIAL_LINKS.instagram} className="text-gray-400 hover:text-white transition-colors">
                <Instagram className="h-5 w-5" />
              </Link>
              <Link href={SOCIAL_LINKS.youtube} className="text-gray-400 hover:text-white transition-colors">
                <Youtube className="h-5 w-5" />
              </Link>
            </div>
          </div>

          {/* Quick links */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Quick Links</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/products" className="text-gray-300 hover:text-white transition-colors text-sm">
                  All Products
                </Link>
              </li>
              <li>
                <Link href="/categories" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Categories
                </Link>
              </li>
              <li>
                <Link href="/deals" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Special Deals
                </Link>
              </li>
              <li>
                <Link href="/new-arrivals" className="text-gray-300 hover:text-white transition-colors text-sm">
                  New Arrivals
                </Link>
              </li>
              <li>
                <Link href="/bestsellers" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Best Sellers
                </Link>
              </li>
              <li>
                <Link href="/brands" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Brands
                </Link>
              </li>
            </ul>
          </div>

          {/* Customer service */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Customer Service</h3>
            <ul className="space-y-2">
              <li>
                <Link href="/help" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Help Center
                </Link>
              </li>
              <li>
                <Link href="/contact" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Contact Us
                </Link>
              </li>
              <li>
                <Link href="/shipping" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Shipping Info
                </Link>
              </li>
              <li>
                <Link href="/returns" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Returns & Exchanges
                </Link>
              </li>
              <li>
                <Link href="/size-guide" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Size Guide
                </Link>
              </li>
              <li>
                <Link href="/track-order" className="text-gray-300 hover:text-white transition-colors text-sm">
                  Track Your Order
                </Link>
              </li>
            </ul>
          </div>

          {/* Newsletter */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Stay Updated</h3>
            <p className="text-gray-300 text-sm">
              Subscribe to our newsletter for exclusive deals and updates.
            </p>
            <form className="space-y-2">
              <Input
                type="email"
                placeholder="Enter your email"
                className="bg-gray-800 border-gray-700 text-white placeholder-gray-400"
              />
              <Button className="w-full" variant="default">
                Subscribe
              </Button>
            </form>
            
            {/* Contact info */}
            <div className="space-y-2 pt-4">
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <Mail className="h-4 w-4" />
                <span>{CONTACT_INFO.email}</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <Phone className="h-4 w-4" />
                <span>{CONTACT_INFO.phone}</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <MapPin className="h-4 w-4" />
                <span>{CONTACT_INFO.address}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Features bar */}
      <div className="border-t border-gray-800">
        <div className="container mx-auto px-4 py-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <Truck className="h-8 w-8 text-primary-400" />
              </div>
              <div>
                <h4 className="font-semibold text-sm">Free Shipping</h4>
                <p className="text-gray-400 text-xs">On orders over $50</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <Shield className="h-8 w-8 text-primary-400" />
              </div>
              <div>
                <h4 className="font-semibold text-sm">Secure Payment</h4>
                <p className="text-gray-400 text-xs">100% secure transactions</p>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              <div className="flex-shrink-0">
                <CreditCard className="h-8 w-8 text-primary-400" />
              </div>
              <div>
                <h4 className="font-semibold text-sm">Easy Returns</h4>
                <p className="text-gray-400 text-xs">30-day return policy</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom bar */}
      <div className="border-t border-gray-800">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row justify-between items-center space-y-2 md:space-y-0">
            <div className="text-sm text-gray-400">
              © 2024 {APP_NAME}. All rights reserved.
            </div>
            
            <div className="flex items-center space-x-6">
              <Link href="/privacy" className="text-sm text-gray-400 hover:text-white transition-colors">
                Privacy Policy
              </Link>
              <Link href="/terms" className="text-sm text-gray-400 hover:text-white transition-colors">
                Terms of Service
              </Link>
              <Link href="/cookies" className="text-sm text-gray-400 hover:text-white transition-colors">
                Cookie Policy
              </Link>
            </div>
            
            <div className="flex items-center space-x-2">
              <span className="text-sm text-gray-400">We accept:</span>
              <div className="flex space-x-1">
                <div className="w-8 h-5 bg-gray-700 rounded flex items-center justify-center">
                  <span className="text-xs text-white">💳</span>
                </div>
                <div className="w-8 h-5 bg-gray-700 rounded flex items-center justify-center">
                  <span className="text-xs text-white">💳</span>
                </div>
                <div className="w-8 h-5 bg-gray-700 rounded flex items-center justify-center">
                  <span className="text-xs text-white">💳</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  )
}
