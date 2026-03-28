import { useState, useCallback, useMemo } from 'react';
import { Device } from './types/device';
import { useDevices } from './hooks/useDevices';
import { DeviceFilters } from './components/DeviceFilters';
import { DeviceForm } from './components/DeviceForm';
import { DeviceTable } from './components/DeviceTable';

// Simple debounce hook for search inputs
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useState(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  });

  return debouncedValue;
}

export default function App() {
  const [hostnameFilter, setHostnameFilter] = useState('');
  const [isActiveFilter, setIsActiveFilter] = useState('');
  const [editingDevice, setEditingDevice] = useState<Device | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  // Debounce hostname filter to avoid excessive API calls (300ms delay)
  const debouncedHostnameFilter = useDebounce(hostnameFilter, 300);

  const isActiveParam = isActiveFilter === '' ? undefined : isActiveFilter === 'true';

  const {
    devices,
    loading,
    error,
    createDevice,
    updateDevice,
    deleteDevice,
    refetch,
    clearError,
  } = useDevices({
    hostname: debouncedHostnameFilter || undefined,
    isActive: isActiveParam,
  });

  const handleCreate = async (data: { hostname: string; ip: string; location?: string; is_active?: boolean }) => {
    try {
      await createDevice(data);
      setShowForm(false);
      setFormError(null);
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'Failed to create device');
      throw err;
    }
  };

  const handleUpdate = async (data: { hostname: string; ip: string; location?: string; is_active?: boolean }) => {
    if (editingDevice) {
      try {
        await updateDevice(editingDevice.id, data);
        setEditingDevice(null);
        setShowForm(false);
        setFormError(null);
      } catch (err) {
        setFormError(err instanceof Error ? err.message : 'Failed to update device');
        throw err;
      }
    }
  };

  const handleCancel = () => {
    setEditingDevice(null);
    setShowForm(false);
    setFormError(null);
  };

  const handleEdit = (device: Device) => {
    setEditingDevice(device);
    setShowForm(true);
    setFormError(null);
  };

  const handleDelete = useCallback(async (id: number) => {
    if (window.confirm('Are you sure you want to delete this device?')) {
      try {
        await deleteDevice(id);
      } catch (err) {
        console.error('Failed to delete device:', err);
      }
    }
  }, [deleteDevice]);

  const handleRefetch = useCallback(() => {
    refetch();
  }, [refetch]);

  // Memoize filtered devices count for display
  const devicesCount = useMemo(() => devices.length, [devices]);

  return (
    <div className="app">
      <header>
        <h1>Device Management</h1>
        <button onClick={() => setShowForm(true)} disabled={loading}>
          Add Device
        </button>
      </header>

      <DeviceFilters
        hostname={hostnameFilter}
        isActiveFilter={isActiveFilter}
        onHostnameChange={(value) => {
          setHostnameFilter(value);
          clearError();
        }}
        onIsActiveFilterChange={(value) => {
          setIsActiveFilter(value);
          clearError();
        }}
      />

      {loading && <p className="loading" role="status">Loading devices...</p>}
      
      {error && (
        <div className="error" role="alert">
          <p>{error}</p>
          <button onClick={handleRefetch}>Retry</button>
        </div>
      )}

      {!loading && !error && (
        <p className="devices-count">
          Showing {devicesCount} {devicesCount === 1 ? 'device' : 'devices'}
        </p>
      )}

      {showForm ? (
        <DeviceForm
          device={editingDevice}
          onSubmit={editingDevice ? handleUpdate : handleCreate}
          onCancel={handleCancel}
          error={formError}
        />
      ) : (
        <DeviceTable 
          devices={devices} 
          onEdit={handleEdit} 
          onDelete={handleDelete}
          isLoading={loading}
        />
      )}
    </div>
  );
}
