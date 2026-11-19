import { api, getToken } from './api';
import { getProducts, mapBackendProduct, type Product } from './catalogService';

const LOCAL_WISHLIST_KEY = 'katarina_local_wishlist';

export function getLocalWishlist(): string[] {
  try {
    const raw = localStorage.getItem(LOCAL_WISHLIST_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed)) {
        return parsed.map((id) => String(id));
      }
    }
  } catch {
    // ignore
  }
  return [];
}

export function saveLocalWishlist(ids: string[]): void {
  try {
    localStorage.setItem(LOCAL_WISHLIST_KEY, JSON.stringify(ids));
  } catch {
    // ignore
  }
}

export function clearLocalWishlist(): void {
  try {
    localStorage.removeItem(LOCAL_WISHLIST_KEY);
  } catch {
    // ignore
  }
}

export async function fetchWishlist(): Promise<string[]> {
  if (getToken()) {
    try {
      const res = await api.get<any[]>('/wishlist');
      if (res.data && Array.isArray(res.data)) {
        return res.data.map((item) => String(item.product_id || item.product?.id || item.id));
      }
      return [];
    } catch (e) {
      console.warn('Backend wishlist fetch failed:', e);
      return [];
    }
  }
  return getLocalWishlist();
}

export async function addToWishlist(productId: string | number): Promise<string[]> {
  const pId = String(productId);
  if (getToken()) {
    try {
      await api.post('/wishlist', { product_id: pId });
      return await fetchWishlist();
    } catch (e) {
      console.warn('Backend wishlist add failed:', e);
      return await fetchWishlist();
    }
  }

  const list = getLocalWishlist();
  if (!list.includes(pId)) {
    list.unshift(pId);
    saveLocalWishlist(list);
  }
  return list;
}

export async function removeFromWishlist(productId: string | number): Promise<string[]> {
  const pId = String(productId);
  if (getToken()) {
    try {
      await api.delete(`/wishlist/${pId}`);
      return await fetchWishlist();
    } catch (e) {
      console.warn('Backend wishlist delete failed:', e);
      return await fetchWishlist();
    }
  }

  const list = getLocalWishlist().filter((id) => id !== pId);
  saveLocalWishlist(list);
  return list;
}

export async function getWishlistProducts(): Promise<Product[]> {
  if (getToken()) {
    try {
      const res = await api.get<any[]>('/wishlist');
      if (res.data && Array.isArray(res.data)) {
        const productsFromBackend: Product[] = [];
        for (const item of res.data) {
          if (item.product) {
            productsFromBackend.push(mapBackendProduct(item.product));
          }
        }
        if (productsFromBackend.length > 0) {
          return productsFromBackend;
        }
      }
      const ids = await fetchWishlist();
      if (ids.length === 0) return [];
      const allProds = await getProducts();
      return allProds.filter((p) => ids.includes(String(p.id)));
    } catch (e) {
      console.warn('Backend wishlist detailed fetch failed:', e);
      return [];
    }
  }

  const ids = getLocalWishlist().map((id) => String(id).trim());
  if (ids.length === 0) {
    return [];
  }

  const allProducts = await getProducts();
  const idSet = new Set(ids);
  return allProducts.filter(
    (p) => idSet.has(String(p.id)) || (p.slug && idSet.has(p.slug))
  );
}


