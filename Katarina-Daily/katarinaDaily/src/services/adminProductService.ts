import { api } from './api';
import type { BackendProduct } from './catalogService';

export interface CreateProductPayload {
  category_id: string;
  name: string;
  slug?: string;
  description?: string;
  status?: 'active' | 'inactive' | 'draft';
}

export interface UpdateProductPayload {
  category_id?: string;
  name?: string;
  slug?: string;
  description?: string;
  status?: 'active' | 'inactive' | 'draft';
}

export interface CreateVariantPayload {
  sku: string;
  color: string;
  size: string;
  price: number | string;
  stock: number;
  weight: number;
  status?: 'active' | 'inactive';
}

export interface UpdateVariantPayload {
  sku?: string;
  color?: string;
  size?: string;
  price?: number | string;
  stock?: number;
  weight?: number;
  status?: 'active' | 'inactive';
}

export interface AddImagePayload {
  image_url: string;
  is_primary?: boolean;
  sort_order?: number;
}

export interface CreateCategoryPayload {
  name: string;
  slug?: string;
  description?: string;
  image_url?: string;
}

/**
 * Creates a product in the backend catalog.
 */
export async function createProduct(data: CreateProductPayload): Promise<BackendProduct> {
  const res = await api.post<BackendProduct>('/admin/products', data);
  if (!res.data) throw new Error(res.message || 'Failed to create product');
  return res.data;
}

/**
 * Updates basic product metadata (name, category, description, status, slug).
 */
export async function updateProduct(id: string, data: UpdateProductPayload): Promise<BackendProduct> {
  const res = await api.put<BackendProduct>(`/admin/products/${id}`, data);
  if (!res.data) throw new Error(res.message || 'Failed to update product');
  return res.data;
}

/**
 * Deletes a product permanently.
 */
export async function deleteProduct(id: string): Promise<void> {
  await api.delete(`/admin/products/${id}`);
}

/**
 * Adds an image to a product.
 */
export async function addProductImage(productId: string, data: AddImagePayload) {
  const res = await api.post(`/products/${productId}/images`, data);
  return res.data;
}

/**
 * Adds a variant to a product.
 */
export async function addProductVariant(productId: string, data: CreateVariantPayload) {
  const payload = {
    ...data,
    price: String(data.price),
  };
  const res = await api.post(`/products/${productId}/variants`, payload);
  return res.data;
}

/**
 * Lists all variants of a product.
 */
export async function getProductVariants(productId: string) {
  const res = await api.get<any[]>(`/products/${productId}/variants`);
  return res.data || [];
}

/**
 * Updates a specific variant (stock, price, SKU, etc.).
 */
export async function updateProductVariant(variantId: string, data: UpdateVariantPayload) {
  const payload: any = { ...data };
  if (data.price !== undefined) {
    payload.price = String(data.price);
  }
  const res = await api.put(`/admin/variants/${variantId}`, payload);
  return res.data;
}

/**
 * Deletes a variant.
 */
export async function deleteProductVariant(variantId: string): Promise<void> {
  await api.delete(`/admin/variants/${variantId}`);
}

/**
 * Creates a new category.
 */
export async function createCategory(data: CreateCategoryPayload) {
  const res = await api.post('/admin/categories', data);
  return res.data;
}
