import { api, setToken, getToken } from './api';

export interface User {
  id: string;
  name: string;
  email: string;
  phone?: string;
  role: 'customer' | 'admin';
  created_at?: string;
}

export interface AuthResponseData {
  token: string;
  user: User;
}

export interface RegisterPayload {
  name: string;
  email: string;
  password: string;
  phone?: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export async function register(payload: RegisterPayload): Promise<AuthResponseData> {
  const response = await api.post<AuthResponseData>('/auth/register', payload);
  if (response.data) {
    setToken(response.data.token);
    return response.data;
  }
  throw new Error(response.message || 'Registration failed');
}

export async function login(payload: LoginPayload): Promise<AuthResponseData> {
  const response = await api.post<AuthResponseData>('/auth/login', payload);
  if (response.data) {
    setToken(response.data.token);
    return response.data;
  }
  throw new Error(response.message || 'Login failed');
}

export async function getCurrentUser(): Promise<User | null> {
  if (!getToken()) {
    return null;
  }
  try {
    const response = await api.get<User>('/auth/me');
    return response.data || null;
  } catch {
    setToken(null);
    return null;
  }
}

export function logout(): void {
  setToken(null);
}
