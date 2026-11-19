import { motion, AnimatePresence } from 'framer-motion';
import { X, Plus, Minus, Trash2, ShoppingBag, ArrowRight } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '@/context/CartContext';
import { formatIDR } from '@/services/catalogService';

export function CartDrawer() {
  const { cart, isCartOpen, closeCart, updateQuantity, removeItem } = useCart();
  const navigate = useNavigate();

  const handleCheckout = () => {
    closeCart();
    navigate('/checkout');
  };

  return (
    <AnimatePresence>
      {isCartOpen && (
        <div className="fixed inset-0 z-[90] flex justify-end">
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={closeCart}
            className="fixed inset-0 bg-dark/40 backdrop-blur-sm"
          />

          {/* Drawer panel */}
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
            className="relative flex h-full w-full max-w-md flex-col bg-primary shadow-2xl z-10 border-l border-dark/10"
          >
            {/* Header */}
            <div className="flex items-center justify-between border-b border-dark/10 p-6">
              <div className="flex items-center gap-3">
                <ShoppingBag size={20} className="text-secondary" />
                <h3 className="font-heading text-xl text-dark">Your Shopping Bag</h3>
                <span className="bg-secondary/10 text-secondary px-2 py-0.5 text-[11px] font-medium">
                  {cart.total_items}
                </span>
              </div>
              <button
                onClick={closeCart}
                className="text-dark/50 transition-colors hover:text-dark"
                aria-label="Close cart"
              >
                <X size={20} />
              </button>
            </div>

            {/* Content */}
            {cart.items.length === 0 ? (
              <div className="flex flex-1 flex-col items-center justify-center p-8 text-center">
                <ShoppingBag size={48} className="text-dark/20 mb-4 stroke-1" />
                <p className="font-heading text-lg text-dark">Your bag is empty</p>
                <p className="mt-2 text-xs font-light text-dark/60 max-w-xs">
                  Looks like you haven't added any minimalist pieces to your cart yet.
                </p>
                <Link
                  to="/shop"
                  onClick={closeCart}
                  className="mt-6 inline-flex items-center gap-2 bg-secondary px-6 py-3 text-[11px] uppercase tracking-editorial text-primary transition-colors hover:bg-secondary-600"
                >
                  Explore Collection <ArrowRight size={13} />
                </Link>
              </div>
            ) : (
              <>
                {/* Free shipping banner */}
                <div className="bg-primary-50 px-6 py-3 border-b border-dark/5 text-center">
                  <p className="text-[11px] font-light text-dark/70">
                    {cart.subtotal >= 300000 ? (
                      <span className="text-emerald-700 font-normal">
                        ✨ Congratulations! You've unlocked Free Standard Delivery.
                      </span>
                    ) : (
                      <span>
                        Add{' '}
                        <strong className="text-secondary font-medium">
                          {formatIDR(300000 - cart.subtotal)}
                        </strong>{' '}
                        more to qualify for free delivery!
                      </span>
                    )}
                  </p>
                </div>

                {/* Items List */}
                <div className="flex-1 overflow-y-auto p-6 space-y-6 divide-y divide-dark/5">
                  {cart.items.map((item) => (
                    <div key={item.id} className="pt-6 first:pt-0 flex gap-4">
                      {/* Image */}
                      <div className="h-24 w-20 shrink-0 overflow-hidden bg-primary-100 border border-dark/5">
                        <img
                          src={
                            item.product_image ||
                            'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900'
                          }
                          alt={item.product_name}
                          className="h-full w-full object-cover"
                        />
                      </div>

                      {/* Details */}
                      <div className="flex flex-1 flex-col justify-between">
                        <div>
                          <div className="flex justify-between items-start gap-2">
                            <h4 className="font-heading text-sm text-dark line-clamp-1">
                              {item.product_name}
                            </h4>
                            <button
                              onClick={() => removeItem(item.id)}
                              className="text-dark/40 hover:text-red-500 transition-colors"
                              aria-label="Remove item"
                            >
                              <Trash2 size={14} />
                            </button>
                          </div>

                          {/* Variant specs */}
                          <div className="mt-1 flex flex-wrap gap-2 text-[10px] text-dark/60">
                            {item.color && (
                              <span className="bg-dark/5 px-2 py-0.5 rounded-sm">
                                {item.color}
                              </span>
                            )}
                            {item.size && (
                              <span className="bg-dark/5 px-2 py-0.5 rounded-sm">
                                {item.size}
                              </span>
                            )}
                          </div>
                        </div>

                        {/* Quantity & Price */}
                        <div className="mt-3 flex items-center justify-between">
                          <div className="flex items-center border border-dark/15">
                            <button
                              onClick={() => updateQuantity(item.id, item.quantity - 1)}
                              className="p-1.5 text-dark/60 hover:text-dark hover:bg-dark/5 transition-colors"
                              aria-label="Decrease quantity"
                            >
                              <Minus size={12} />
                            </button>
                            <span className="px-3 text-xs font-light text-dark">
                              {item.quantity}
                            </span>
                            <button
                              onClick={() => updateQuantity(item.id, item.quantity + 1)}
                              className="p-1.5 text-dark/60 hover:text-dark hover:bg-dark/5 transition-colors"
                              aria-label="Increase quantity"
                            >
                              <Plus size={12} />
                            </button>
                          </div>

                          <p className="text-xs font-medium text-dark">
                            {formatIDR(item.subtotal)}
                          </p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Footer / Summary */}
                <div className="border-t border-dark/10 p-6 bg-white/40">
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-xs uppercase tracking-editorial text-dark/60">
                      Subtotal
                    </span>
                    <span className="font-heading text-lg text-dark">
                      {formatIDR(cart.subtotal)}
                    </span>
                  </div>

                  <p className="text-[10px] font-light text-dark/50 mb-5 text-center">
                    Taxes and shipping calculated at checkout.
                  </p>

                  <button
                    onClick={handleCheckout}
                    className="flex w-full items-center justify-center gap-3 bg-secondary py-4 text-[11px] uppercase tracking-editorial text-primary transition-all duration-300 hover:bg-secondary-600"
                  >
                    Proceed to Checkout <ArrowRight size={14} />
                  </button>
                </div>
              </>
            )}
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}
