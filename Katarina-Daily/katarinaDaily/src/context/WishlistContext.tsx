import React, { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import {
  fetchWishlist,
  addToWishlist as apiAddToWishlist,
  removeFromWishlist as apiRemoveFromWishlist,
} from '@/services/wishlistService';
import { useAuth } from './AuthContext';

interface WishlistContextType {
  wishlist: string[];
  wishlistCount: number;
  isInWishlist: (productId: string | number) => boolean;
  toggleWishlist: (productId: string | number) => Promise<void>;
  refreshWishlist: () => Promise<void>;
}

const WishlistContext = createContext<WishlistContextType | undefined>(undefined);

export function WishlistProvider({ children }: { children: ReactNode }) {
  const { user } = useAuth();
  const [wishlist, setWishlist] = useState<string[]>([]);

  const refreshWishlist = async () => {
    try {
      const ids = await fetchWishlist();
      setWishlist(ids);
    } catch (e) {
      console.warn('Failed to load wishlist:', e);
    }
  };

  useEffect(() => {
    refreshWishlist();
  }, [user]);

  const isInWishlist = (productId: string | number) => {
    return wishlist.includes(String(productId));
  };

  const toggleWishlist = async (productId: string | number) => {
    const idStr = String(productId);
    if (isInWishlist(idStr)) {
      const updated = await apiRemoveFromWishlist(idStr);
      setWishlist(updated);
    } else {
      const updated = await apiAddToWishlist(idStr);
      setWishlist(updated);
    }
  };

  return (
    <WishlistContext.Provider
      value={{
        wishlist,
        wishlistCount: wishlist.length,
        isInWishlist,
        toggleWishlist,
        refreshWishlist,
      }}
    >
      {children}
    </WishlistContext.Provider>
  );
}

export function useWishlist() {
  const context = useContext(WishlistContext);
  if (!context) {
    throw new Error('useWishlist must be used within a WishlistProvider');
  }
  return context;
}
