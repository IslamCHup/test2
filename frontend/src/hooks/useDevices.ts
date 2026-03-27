import { useState, useEffect, useCallback } from 'react';
import { Device } from '../types/device';
import { deviceApi } from '../api/deviceApi';

interface UseDevicesOptions {
  isActive?: boolean;
  hostname?: string;
}

export function useDevices(options: UseDevicesOptions = {}) {
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchDevices = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await deviceApi.getAll({
        is_active: options.isActive,
        hostname: options.hostname,
      });
      setDevices(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch devices');
    } finally {
      setLoading(false);
    }
  }, [options.isActive, options.hostname]);

  useEffect(() => {
    fetchDevices();
  }, [fetchDevices]);

  const createDevice = async (device: { hostname: string; ip: string; location?: string }) => {
    const newDevice = await deviceApi.create(device);
    setDevices((prev) => [...prev, newDevice]);
    return newDevice;
  };

  const updateDevice = async (id: number, updates: { hostname?: string; ip?: string; location?: string; is_active?: boolean }) => {
    const updatedDevice = await deviceApi.update(id, updates);
    setDevices((prev) => prev.map((d) => (d.id === id ? updatedDevice : d)));
    return updatedDevice;
  };

  const deleteDevice = async (id: number) => {
    await deviceApi.delete(id);
    setDevices((prev) => prev.filter((d) => d.id !== id));
  };

  return {
    devices,
    loading,
    error,
    refetch: fetchDevices,
    createDevice,
    updateDevice,
    deleteDevice,
  };
}
