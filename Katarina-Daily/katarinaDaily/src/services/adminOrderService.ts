import { api } from './api';
import type { Order } from '@/services/orderService';

/**
 * Fetch all orders (admin view).
 * Returns an array of Order objects. If the request fails, an empty array is returned.
 */
export async function getAllOrders(): Promise<Order[]> {
  try {
    const res = await api.get<Order[]>('/admin/orders');
    return res.data || [];
  } catch (e) {
    console.warn('Failed to fetch all orders (admin):', e);
    return [];
  }
}

/**
 * Update the status of a specific order.
 * @param orderId - ID of the order to update.
 * @param payload - Object containing the new status.
 */
export async function updateOrderStatus(
  orderId: string,
  payload: { status: string }
): Promise<void> {
  await api.put(`/admin/orders/${orderId}/status`, payload);
}

/**
 * Delete an order permanently.
 * @param orderId - ID of the order to delete.
 */
export async function deleteOrder(orderId: string): Promise<void> {
  await api.delete(`/admin/orders/${orderId}`);
}
