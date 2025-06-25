import { loadStripe, Stripe } from '@stripe/stripe-js';

// Get Stripe publishable key from environment
const stripePublishableKey = process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || 
  'pk_test_51Rds2fPkBoVzgojKoocpIByWyKh7xIXomdGV321ZaJwW73xxSOAu0PldZqIDAqSSqiy2CsHTSHLdCuNbrIl7Wsbv00wum3HLfV';

// Initialize Stripe
let stripePromise: Promise<Stripe | null>;

export const getStripe = () => {
  if (!stripePromise) {
    stripePromise = loadStripe(stripePublishableKey);
  }
  return stripePromise;
};

// Stripe configuration
export const stripeConfig = {
  publishableKey: stripePublishableKey,
  appearance: {
    theme: 'stripe' as const,
    variables: {
      colorPrimary: '#0570de',
      colorBackground: '#ffffff',
      colorText: '#30313d',
      colorDanger: '#df1b41',
      fontFamily: 'Inter, system-ui, sans-serif',
      spacingUnit: '4px',
      borderRadius: '8px',
    },
  },
  loader: 'auto' as const,
};

// Helper function to redirect to Stripe Checkout
export const redirectToCheckout = async (sessionId: string) => {
  const stripe = await getStripe();
  if (!stripe) {
    throw new Error('Stripe failed to initialize');
  }

  const { error } = await stripe.redirectToCheckout({
    sessionId,
  });

  if (error) {
    throw error;
  }
};

// Helper function to create payment intent
export const createPaymentIntent = async (amount: number, currency: string = 'usd') => {
  // This would typically call your backend to create a payment intent
  // For now, we'll use the checkout session approach
  throw new Error('Direct payment intents not implemented - use checkout sessions');
};

export default getStripe;
