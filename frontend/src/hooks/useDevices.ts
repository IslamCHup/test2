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
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch devices';
      setError(errorMessage);
      console.error('Error fetching devices:', err);
    } finally {
      setLoading(false);
    }
  }, [options.isActive, options.hostname]);

  useEffect(() => {
    fetchDevices();
  }, [fetchDevices]);

  const createDevice = async (device: { hostname: string; ip: string; location?: string }) => {
    try {
      const newDevice = await deviceApi.create(device);
      setDevices((prev) => [...prev, newDevice]);
      return newDevice;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create device';
      setError(errorMessage);
      console.error('Error creating device:', err);
      throw err;
    }
  };

  const updateDevice = async (id: number, updates: { hostname?: string; ip?: string; location?: string; is_active?: boolean }) => {
    try {
      const updatedDevice = await deviceApi.update(id, updates);
      setDevices((prev) => prev.map((d) => (d.id === id ? updatedDevice : d)));
      return updatedDevice;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update device';
      setError(errorMessage);
      console.error('Error updating device:', err);
      throw err;
    }
  };

  const deleteDevice = async (id: number) => {
    try {
      await deviceApi.delete(id);
      setDevices((prev) => prev.filter((d) => d.id !== id));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to delete device';
      setError(errorMessage);
      console.error('Error deleting device:', err);
      throw err;
    }
  };

  return {
    devices,
    loading,
    error,
    refetch: fetchDevices,
    createDevice,
    updateDevice,
    deleteDevice,
    clearError: () => setError(null),
  };
}
