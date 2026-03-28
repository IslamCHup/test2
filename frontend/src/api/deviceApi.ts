import { Device, CreateDeviceRequest, UpdateDeviceRequest } from '../types/device';

const API_BASE = import.meta.env.VITE_API_URL || '';
const DEVICE_URL = `${API_BASE}/api/devices`;

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 8000);

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        ...(options?.headers || {}),
      },
    });

    if (!response.ok) {
      let message = response.statusText;

      try {
        const data = await response.json();
        message = data.error || message;
      } catch {}

      throw new Error(message);
    }

    // обработка 204 No Content
    if (response.status === 204) {
      return null as T;
    }

    return response.json() as Promise<T>;
  } catch (err) {
    if ((err as Error).name === 'AbortError') {
      throw new Error('Request timeout');
    }
    throw err;
  } finally {
    clearTimeout(timeout);
  }
}

// ================= API =================

export const deviceApi = {
  getAll(params?: { is_active?: boolean; hostname?: string }): Promise<Device[]> {
    const searchParams = new URLSearchParams();

    if (params?.is_active !== undefined) {
      searchParams.set('is_active', String(params.is_active));
    }

    if (params?.hostname) {
      searchParams.set('hostname', params.hostname);
    }

    const query = searchParams.toString();
    const url = query ? `${DEVICE_URL}?${query}` : DEVICE_URL;

    return request<Device[]>(url);
  },

  getById(id: number): Promise<Device> {
    return request<Device>(`${DEVICE_URL}/${id}`);
  },

  create(device: CreateDeviceRequest): Promise<Device> {
    return request<Device>(DEVICE_URL, {
      method: 'POST',
      body: JSON.stringify(device),
    });
  },

  update(id: number, device: UpdateDeviceRequest): Promise<Device> {
    return request<Device>(`${DEVICE_URL}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(device),
    });
  },

  delete(id: number): Promise<void> {
    return request<void>(`${DEVICE_URL}/${id}`, {
      method: 'DELETE',
    });
  },
};