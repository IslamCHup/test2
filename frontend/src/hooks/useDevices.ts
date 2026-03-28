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

  // Keep options ref up to date
  useEffect(() => {
    optionsRef.current = options;
  }, [options]);

  // Cleanup on unmount
  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
      // Cancel any pending request on unmount
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
        abortControllerRef.current = null;
      }
    };
  }, []);

  const fetchDevices = useCallback(async (signal?: AbortSignal) => {
    // Cancel previous request if still pending
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
    
    // Create new AbortController for this request
    abortControllerRef.current = new AbortController();
    const currentSignal = signal || abortControllerRef.current.signal;

    // Use current options from ref to avoid stale closures
    const currentOptions = optionsRef.current;
    
    setLoading(true);
    setError(null);
    
    try {
      const data = await deviceApi.getAll({
        is_active: currentOptions.isActive,
        hostname: currentOptions.hostname,
      }, currentSignal);
      
      // Only update state if component is still mounted and request wasn't aborted
      if (isMountedRef.current && !currentSignal.aborted) {
        setDevices(data);
      }
    } catch (err) {
      // Ignore abort errors
      if (err instanceof Error && err.name === 'AbortError') {
        return;
      }
      
      if (isMountedRef.current) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to fetch devices';
        setError(errorMessage);
        console.error('Error fetching devices:', err);
      }
    } finally {
      if (isMountedRef.current && !currentSignal.aborted) {
        setLoading(false);
      }
    }
  }, []);

  // Fetch devices when options change
  useEffect(() => {
    fetchDevices();
  }, [fetchDevices, options.isActive, options.hostname]);

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
