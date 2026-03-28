export interface Device {
  id: number;
  hostname: string;
  ip: string;
  location: string | null;
  is_active: boolean;
  created_at: string;
}

export interface CreateDeviceRequest {
  hostname: string;
  ip: string;
  location?: string;
}

export interface UpdateDeviceRequest {
  hostname?: string;
  ip?: string;
  location?: string;
  is_active?: boolean;
}
