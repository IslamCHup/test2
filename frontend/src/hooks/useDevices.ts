import { useState, useEffect, useCallback, useRef } from 'react';
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
  
  // Ref to track if component is mounted to prevent setState after unmount
  const isMountedRef = useRef(true);
  // Ref to store the latest options to avoid stale closures
  const optionsRef = useRef(options);
  // AbortController ref for cancelling pending requests
  const abortControllerRef = useRef<AbortController | null>(null);

  const requestIdRef = useRef(0);
  const isMountedRef = useRef(true);

  useEffect(() => {
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const fetchDevices = useCallback(async () => {
    const requestId = ++requestIdRef.current;

    setLoading(true);
    setError(null);

    try {
      const data = await deviceApi.getAll({
        is_active: options.isActive,
        hostname: options.hostname,
      });

      // защита от race condition
      if (requestId !== requestIdRef.current || !isMountedRef.current) return;

      setDevices(data);
    } catch (err) {
      if (!isMountedRef.current) return;

      const errorMessage =
        err instanceof Error ? err.message : 'Failed to fetch devices';

      setError(errorMessage);
      console.error('Error fetching devices:', err);
    } finally {
      if (isMountedRef.current) {
        setLoading(false);
      }
    }
  }, []);

  // Fetch devices when options change
  useEffect(() => {
    fetchDevices();
  }, [fetchDevices, options.isActive, options.hostname]);

  // ================= CREATE =================

  const createDevice = async (device: {
    hostname: string;
    ip: string;
    location?: string;
  }) => {
    try {
      const newDevice = await deviceApi.create(device);

      if (isMountedRef.current) {
        setDevices((prev) => [...prev, newDevice]);
      }

      return newDevice;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Failed to create device';

      setError(errorMessage);
      console.error('Error creating device:', err);
      throw err;
    }
  };

  // ================= UPDATE =================

  const updateDevice = async (
    id: number,
    updates: {
      hostname?: string;
      ip?: string;
      location?: string;
      is_active?: boolean;
    }
  ) => {
    try {
      const updatedDevice = await deviceApi.update(id, updates);

      if (isMountedRef.current) {
        setDevices((prev) =>
          prev.map((d) => (d.id === id ? updatedDevice : d))
        );
      }

      return updatedDevice;
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Failed to update device';

      setError(errorMessage);
      console.error('Error updating device:', err);
      throw err;
    }
  };

  // ================= DELETE =================

  const deleteDevice = async (id: number) => {
    try {
      await deviceApi.delete(id);

      if (isMountedRef.current) {
        setDevices((prev) => prev.filter((d) => d.id !== id));
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Failed to delete device';

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