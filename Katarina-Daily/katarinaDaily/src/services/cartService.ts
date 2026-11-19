import { api, getToken } from './api';
import type { ProductVariant } from './catalogService';

export interface CartItem {
  id: string;
  product_variant_id: string;
  product_id?: string;
  product_name: string;
  product_image?: string;
  variant?: ProductVariant;
  color?: string;
  size?: string;
  sku?: string;
  quantity: number;
  unit_price: number;
  subtotal: number;
  available?: boolean;
}

export interface Cart {
  id?: string;
  items: CartItem[];
  total_items: number;
  total_weight_grams: number;
  subtotal: number;
}

const LOCAL_CART_KEY = 'katarina_local_cart';

function getLocalCart(): Cart {
  try {
    const raw = localStorage.getItem(LOCAL_CART_KEY);
    if (raw) {
      return JSON.parse(raw);
    }
  } catch {
    // ignore
  }
  return { items: [], total_items: 0, total_weight_grams: 0, subtotal: 0 };
}

function saveLocalCart(cart: Cart): void {
  try {
    localStorage.setItem(LOCAL_CART_KEY, JSON.stringify(cart));
  } catch {
    // ignore
  }
}

function recalculateCart(items: CartItem[]): Cart {
  let totalItems = 0;
  let totalWeight = 0;
  let subtotal = 0;

  items.forEach((item) => {
    totalItems += item.quantity;
    subtotal += item.unit_price * item.quantity;
    totalWeight += (item.variant?.weight || 250) * item.quantity;
  });

  return {
    items,
    total_items: totalItems,
    total_weight_grams: totalWeight,
    subtotal,
  };
}

export async function fetchCart(): Promise<Cart> {
  if (getToken()) {
    try {
      const res = await api.get<any>('/cart');
      if (res.data) {
        const items: CartItem[] = (res.data.items || []).map((i: any) => ({
          id: i.id,
          product_variant_id: i.variant?.id || i.product_variant_id,
          product_id: i.product_id,
          product_name: i.product_name,
          color: i.variant?.color,
          size: i.variant?.size,
          sku: i.variant?.sku,
          variant: i.variant
            ? {
                ...i.variant,
                price: parseFloat(i.variant.price || '0'),
              }
            : undefined,
          quantity: i.quantity,
          unit_price: parseFloat(i.unit_price || '0'),
          subtotal: parseFloat(i.subtotal || '0'),
          available: i.available,
        }));

        return {
          id: res.data.id,
          items,
          total_items: res.data.total_items || items.reduce((acc, it) => acc + it.quantity, 0),
          total_weight_grams: res.data.total_weight_grams || 0,
          subtotal: parseFloat(res.data.subtotal || '0'),
        };
      }
    } catch (e) {
      console.warn('Failed to fetch remote cart, falling back to local cart:', e);
    }
  }
  return getLocalCart();
}

export async function addCartItem(
  productVariantId: string,
  quantity = 1,
  fallbackMetadata?: {
    productName: string;
    productImage?: string;
    variant?: ProductVariant;
    unitPrice: number;
  }
): Promise<Cart> {
  if (getToken()) {
    try {
      await api.post('/cart/items', {
        product_variant_id: productVariantId,
        quantity,
      });
      return await fetchCart();
    } catch (e) {
      console.warn('Failed to add to remote cart, saving locally:', e);
    }
  }

  // Local fallback
  const cart = getLocalCart();
  const existingIdx = cart.items.findIndex(
    (i) => i.product_variant_id === productVariantId
  );

  if (existingIdx >= 0) {
    cart.items[existingIdx].quantity += quantity;
    cart.items[existingIdx].subtotal =
      cart.items[existingIdx].quantity * cart.items[existingIdx].unit_price;
  } else {
    cart.items.push({
      id: `local-item-${Date.now()}`,
      product_variant_id: productVariantId,
      product_name: fallbackMetadata?.productName || 'Fashion Item',
      product_image: fallbackMetadata?.productImage,
      variant: fallbackMetadata?.variant,
      color: fallbackMetadata?.variant?.color,
      size: fallbackMetadata?.variant?.size,
      sku: fallbackMetadata?.variant?.sku,
      quantity,
      unit_price: fallbackMetadata?.unitPrice || 129000,
      subtotal: (fallbackMetadata?.unitPrice || 129000) * quantity,
      available: true,
    });
  }

  const updated = recalculateCart(cart.items);
  saveLocalCart(updated);
  return updated;
}

export async function updateCartItemQuantity(
  itemId: string,
  quantity: number
): Promise<Cart> {
  if (quantity <= 0) {
    return removeCartItem(itemId);
  }

  if (getToken() && !itemId.startsWith('local-')) {
    try {
      await api.put(`/cart/items/${itemId}`, { quantity });
      return await fetchCart();
    } catch (e) {
      console.warn('Remote cart update failed:', e);
    }
  }

  const cart = getLocalCart();
  const item = cart.items.find((i) => i.id === itemId);
  if (item) {
    item.quantity = quantity;
    item.subtotal = item.quantity * item.unit_price;
  }
  const updated = recalculateCart(cart.items);
  saveLocalCart(updated);
  return updated;
}

export async function removeCartItem(itemId: string): Promise<Cart> {
  if (getToken() && !itemId.startsWith('local-')) {
    try {
      await api.delete(`/cart/items/${itemId}`);
      return await fetchCart();
    } catch (e) {
      console.warn('Remote cart remove failed:', e);
    }
  }

  const cart = getLocalCart();
  cart.items = cart.items.filter((i) => i.id !== itemId);
  const updated = recalculateCart(cart.items);
  saveLocalCart(updated);
  return updated;
}

export async function clearCart(): Promise<Cart> {
  if (getToken()) {
    try {
      await api.delete('/cart');
    } catch (e) {
      console.warn('Remote cart clear failed:', e);
    }
  }
  const empty: Cart = { items: [], total_items: 0, total_weight_grams: 0, subtotal: 0 };
  saveLocalCart(empty);
  return empty;
}
