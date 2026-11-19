import { api } from './api';

export interface Address {
  id: string;
  label?: string;
  recipient_name: string;
  phone: string;
  address: string;
  province: string;
  city: string;
  district?: string;
  postal_code: string;
  is_default: boolean;
}

export interface CreateAddressPayload {
  label?: string;
  recipient_name: string;
  phone: string;
  address: string;
  province: string;
  city: string;
  district?: string;
  postal_code: string;
  is_default?: boolean;
}

export interface ShippingRateOption {
  courier: string;
  service: string;
  description: string;
  cost: number | string;
  etd: string;
}

export interface CheckoutPayload {
  address_id: string;
  shipping_courier: string;
  shipping_service: string;
  voucher_code?: string;
}

export interface OrderItem {
  id: string;
  product_variant_id: string;
  product_name: string;
  sku: string;
  color: string;
  size: string;
  price: number | string;
  quantity: number;
  subtotal: number | string;
}

export interface Order {
  id: string;
  order_number: string;
  status: string;
  subtotal: number | string;
  shipping_cost: number | string;
  discount: number | string;
  total: number | string;
  shipping_courier: string;
  shipping_service: string;
  shipping_address: string;
  shipping_recipient_name: string;
  shipping_phone: string;
  shipping_city: string;
  shipping_province: string;
  shipping_postal_code: string;
  items: OrderItem[];
  created_at: string;
  payment?: {
    status: string;
    payment_method?: string;
    snap_token?: string;
    redirect_url?: string;
  };
}

export async function getAddresses(): Promise<Address[]> {
  try {
    const res = await api.get<Address[]>('/addresses');
    return res.data || [];
  } catch (e) {
    console.warn('Failed to fetch addresses:', e);
    return [];
  }
}

export async function createAddress(payload: CreateAddressPayload): Promise<Address> {
  const res = await api.post<Address>('/addresses', payload);
  if (res.data) {
    return res.data;
  }
  throw new Error(res.message || 'Failed to save address');
}

export async function getShippingRates(addressId: string, courier?: string): Promise<ShippingRateOption[]> {
  try {
    const query = new URLSearchParams({ address_id: addressId });
    if (courier) query.append('courier', courier);

    const res = await api.get<ShippingRateOption[]>(`/shipping/rates?${query.toString()}`);
    return res.data || [];
  } catch (e) {
    console.warn('Failed to fetch shipping rates from backend:', e);
    // Fallback options
    return [
      { courier: 'JNE', service: 'REG', description: 'Reguler (2-3 Days)', cost: 18000, etd: '2-3 days' },
      { courier: 'SiCepat', service: 'SIUNT', description: 'Best Service (1-2 Days)', cost: 20000, etd: '1-2 days' },
      { courier: 'J&T', service: 'EZ', description: 'Express (1-3 Days)', cost: 19000, etd: '1-3 days' },
    ];
  }
}

export async function checkout(payload: CheckoutPayload): Promise<Order> {
  const res = await api.post<Order>('/orders', payload);
  if (res.data) {
    return res.data;
  }
  throw new Error(res.message || 'Checkout failed');
}

export async function getOrders(): Promise<Order[]> {
  try {
    const res = await api.get<Order[]>('/orders');
    return res.data || [];
  } catch (e) {
    console.warn('Failed to fetch orders:', e);
    return [];
  }
}

export async function getOrderById(id: string): Promise<Order | null> {
  try {
    const res = await api.get<Order>(`/orders/${id}`);
    return res.data || null;
  } catch (e) {
    console.warn('Failed to fetch order details:', e);
    return null;
  }
}
