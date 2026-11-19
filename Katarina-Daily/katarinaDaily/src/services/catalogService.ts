import { api } from './api';
import {
  fallbackProducts,
  categories as fallbackCategoryNames,
  type Product,
  type Category,
  type ProductVariant,
  type ProductImage,
} from '@/data/products';

export type { Product, Category, ProductVariant, ProductImage };

export function formatIDR(amount: number | string | undefined): string {
  if (amount === undefined || amount === null || amount === '') return 'Rp0';
  if (typeof amount === 'string' && amount.startsWith('Rp')) return amount;
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  if (isNaN(num)) return 'Rp0';
  return 'Rp' + Math.round(num).toLocaleString('id-ID');
}

// Backend product response structure
interface BackendCategory {
  id: string;
  name: string;
  slug: string;
  description?: string;
  image_url?: string;
}

interface BackendProductImage {
  id: string;
  image_url: string;
  is_primary: boolean;
  sort_order: number;
}

interface BackendProductVariant {
  id: string;
  product_id: string;
  sku: string;
  color: string;
  size: string;
  price: string;
  stock: number;
  weight: number;
  status: string;
}

export interface BackendProduct {
  id: string;
  category_id: string;
  category?: BackendCategory;
  name: string;
  slug: string;
  description: string;
  status: string;
  min_price: string;
  max_price: string;
  total_stock: number;
  variants?: BackendProductVariant[];
  images?: BackendProductImage[];
}

export function mapBackendProduct(p: BackendProduct): Product {
  const minNum = parseFloat(p.min_price || '0');
  const maxNum = parseFloat(p.max_price || '0');

  let priceLabel = formatIDR(minNum);
  if (maxNum > minNum) {
    priceLabel = `${formatIDR(minNum)} - ${formatIDR(maxNum)}`;
  }

  const primaryImage =
    p.images?.find((img) => img.is_primary)?.image_url ||
    p.images?.[0]?.image_url ||
    'https://images.pexels.com/photos/22441297/pexels-photo-22441297.jpeg?auto=compress&cs=tinysrgb&h=900';

  const mappedVariants: ProductVariant[] = (p.variants || []).map((v) => ({
    id: v.id,
    product_id: v.product_id,
    sku: v.sku,
    color: v.color,
    size: v.size,
    price: parseFloat(v.price || '0'),
    stock: v.stock,
    weight: v.weight,
    status: v.status,
  }));

  return {
    id: p.id,
    name: p.name,
    slug: p.slug,
    description: p.description,
    price: priceLabel,
    rawPrice: minNum,
    min_price: minNum,
    max_price: maxNum,
    total_stock: p.total_stock,
    image: primaryImage,
    category: p.category?.name || 'Tops',
    category_id: p.category_id,
    shopeeLink: 'https://shopee.co.id/',
    variants: mappedVariants,
    images: p.images || [{ image_url: primaryImage, is_primary: true, sort_order: 0 }],
    status: p.status,
  };
}

export interface ProductQueryParams {
  category_id?: string;
  category?: string;
  search?: string;
  color?: string;
  size?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  per_page?: number;
}

export async function getProducts(params?: ProductQueryParams): Promise<Product[]> {
  try {
    const query = new URLSearchParams();
    if (params?.category_id) query.append('category_id', params.category_id);
    if (params?.search) query.append('search', params.search);
    if (params?.color) query.append('color', params.color);
    if (params?.size) query.append('size', params.size);
    if (params?.min_price) query.append('min_price', params.min_price.toString());
    if (params?.max_price) query.append('max_price', params.max_price.toString());
    if (params?.page) query.append('page', params.page.toString());
    if (params?.per_page) query.append('per_page', params.per_page.toString());

    const queryString = query.toString();
    const endpoint = `/products${queryString ? `?${queryString}` : ''}`;
    const res = await api.get<BackendProduct[]>(endpoint);

    if (res.data && Array.isArray(res.data) && res.data.length > 0) {
      return res.data.map(mapBackendProduct);
    }
  } catch (error) {
    console.warn('Backend unavailable, using fallback products:', error);
  }

  // Fallback filtering if backend is not running
  let results = [...fallbackProducts];
  if (params?.search) {
    const s = params.search.toLowerCase();
    results = results.filter(
      (p) => p.name.toLowerCase().includes(s) || p.description.toLowerCase().includes(s)
    );
  }
  if (params?.category && params.category !== 'All') {
    results = results.filter((p) => p.category.toLowerCase() === params.category!.toLowerCase());
  }
  return results;
}

export async function getProductById(idOrSlug: number | string): Promise<Product | undefined> {
  const queryStr = String(idOrSlug);

  try {
    const res = await api.get<BackendProduct>(`/products/${queryStr}`);
    if (res.data) {
      return mapBackendProduct(res.data);
    }
  } catch (error) {
    console.warn(`Could not load product ${queryStr} from backend, checking fallback:`, error);
  }

  return fallbackProducts.find(
    (p) => String(p.id) === queryStr || (p.slug && p.slug === queryStr)
  );
}

export async function getCategories(): Promise<string[]> {
  try {
    const res = await api.get<BackendCategory[]>('/categories');
    if (res.data && Array.isArray(res.data) && res.data.length > 0) {
      const names = res.data.map((c) => c.name);
      return ['All', ...names];
    }
  } catch (error) {
    console.warn('Backend categories unavailable, using default category list:', error);
  }
  return fallbackCategoryNames;
}

export async function getCategoryObjects(): Promise<Category[]> {
  try {
    const res = await api.get<BackendCategory[]>('/categories');
    if (res.data && Array.isArray(res.data)) {
      return res.data;
    }
  } catch (error) {
    console.warn('Backend category objects unavailable:', error);
  }
  return [
    { id: 'cat-1', name: 'Tops', slug: 'tops' },
    { id: 'cat-2', name: 'Outerwear', slug: 'outerwear' },
    { id: 'cat-3', name: 'Dresses', slug: 'dresses' },
    { id: 'cat-4', name: 'Bottoms', slug: 'bottoms' },
  ];
}
