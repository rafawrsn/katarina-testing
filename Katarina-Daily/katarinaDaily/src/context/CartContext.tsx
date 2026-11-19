import React, { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import {
  fetchCart,
  addCartItem,
  updateCartItemQuantity,
  removeCartItem,
  clearCart as apiClearCart,
  type Cart,
  type CartItem,
} from '@/services/cartService';
import type { ProductVariant } from '@/services/catalogService';
import { useAuth } from './AuthContext';

interface CartContextType {
  cart: Cart;
  loading: boolean;
  isCartOpen: boolean;
  openCart: () => void;
  closeCart: () => void;
  toggleCart: () => void;
  addToCart: (
    variantId: string,
    quantity?: number,
    metadata?: {
      productName: string;
      productImage?: string;
      variant?: ProductVariant;
      unitPrice: number;
    }
  ) => Promise<void>;
  updateQuantity: (itemId: string, quantity: number) => Promise<void>;
  removeItem: (itemId: string) => Promise<void>;
  clear: () => Promise<void>;
  refreshCart: () => Promise<void>;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [cart, setCart] = useState<Cart>({
    items: [],
    total_items: 0,
    total_weight_grams: 0,
    subtotal: 0,
  });
  const [loading, setLoading] = useState(false);
  const [isCartOpen, setIsCartOpen] = useState(false);

  const refreshCart = async () => {
    setLoading(true);
    try {
      const data = await fetchCart();
      setCart(data);
    } catch (e) {
      console.warn('Failed to load cart:', e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refreshCart();
  }, [user]);

  const openCart = () => setIsCartOpen(true);
  const closeCart = () => setIsCartOpen(false);
  const toggleCart = () => setIsCartOpen((prev) => !prev);

  const addToCart = async (
    variantId: string,
    quantity = 1,
    metadata?: {
      productName: string;
      productImage?: string;
      variant?: ProductVariant;
      unitPrice: number;
    }
  ) => {
    setLoading(true);
    try {
      const updated = await addCartItem(variantId, quantity, metadata);
      setCart(updated);
      setIsCartOpen(true);
    } catch (e) {
      console.error('Failed to add item to cart:', e);
    } finally {
      setLoading(false);
    }
  };

  const updateQuantity = async (itemId: string, quantity: number) => {
    setLoading(true);
    try {
      const updated = await updateCartItemQuantity(itemId, quantity);
      setCart(updated);
    } catch (e) {
      console.error('Failed to update cart quantity:', e);
    } finally {
      setLoading(false);
    }
  };

  const removeItem = async (itemId: string) => {
    setLoading(true);
    try {
      const updated = await removeCartItem(itemId);
      setCart(updated);
    } catch (e) {
      console.error('Failed to remove cart item:', e);
    } finally {
      setLoading(false);
    }
  };

  const clear = async () => {
    setLoading(true);
    try {
      const updated = await apiClearCart();
      setCart(updated);
    } catch (e) {
      console.error('Failed to clear cart:', e);
    } finally {
      setLoading(false);
    }
  };

  return (
    <CartContext.Provider
      value={{
        cart,
        loading,
        isCartOpen,
        openCart,
        closeCart,
        toggleCart,
        addToCart,
        updateQuantity,
        removeItem,
        clear,
        refreshCart,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
}
