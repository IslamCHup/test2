import { Device, CreateDeviceRequest, UpdateDeviceRequest } from '../types/device';

const API_BASE = '/api/devices';

export const deviceApi = {
  async getAll(params?: { is_active?: boolean; hostname?: string }): Promise<Device[]> {
    const searchParams = new URLSearchParams();
    if (params?.is_active !== undefined) {
      searchParams.set('is_active', String(params.is_active));
    }
    if (params?.hostname) {
      searchParams.set('hostname', params.hostname);
    }

    const query = searchParams.toString();
    const url = query ? `${API_BASE}?${query}` : API_BASE;

    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch devices: ${response.statusText}`);
    }
    return response.json();
  },

  async getById(id: number): Promise<Device> {
    const response = await fetch(`${API_BASE}/${id}`);
    if (!response.ok) {
      throw new Error(`Failed to fetch device: ${response.statusText}`);
    }
    return response.json();
  },

  async create(device: CreateDeviceRequest): Promise<Device> {
    const response = await fetch(API_BASE, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(device),
    });
    if (!response.ok) {
      throw new Error(`Failed to create device: ${response.statusText}`);
    }
    return response.json();
  },

  async update(id: number, device: UpdateDeviceRequest): Promise<Device> {
    const response = await fetch(`${API_BASE}/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(device),
    });
    if (!response.ok) {
      throw new Error(`Failed to update device: ${response.statusText}`);
    }
    return response.json();
  },

  async delete(id: number): Promise<void> {
    const response = await fetch(`${API_BASE}/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error(`Failed to delete device: ${response.statusText}`);
    }
  },
};
