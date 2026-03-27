import { useState } from 'react';
import { Device } from './types/device';
import { useDevices } from './hooks/useDevices';
import { DeviceFilters } from './components/DeviceFilters';
import { DeviceForm } from './components/DeviceForm';
import { DeviceTable } from './components/DeviceTable';

export default function App() {
  const [hostnameFilter, setHostnameFilter] = useState('');
  const [isActiveFilter, setIsActiveFilter] = useState('');
  const [editingDevice, setEditingDevice] = useState<Device | null>(null);
  const [showForm, setShowForm] = useState(false);

  const isActiveParam = isActiveFilter === '' ? undefined : isActiveFilter === 'true';

  const {
    devices,
    loading,
    error,
    createDevice,
    updateDevice,
    deleteDevice,
  } = useDevices({
    hostname: hostnameFilter || undefined,
    isActive: isActiveParam,
  });

  const handleCreate = async (data: { hostname: string; ip: string; location?: string; is_active?: boolean }) => {
    await createDevice(data);
    setShowForm(false);
  };

  const handleUpdate = async (data: { hostname: string; ip: string; location?: string; is_active?: boolean }) => {
    if (editingDevice) {
      await updateDevice(editingDevice.id, data);
      setEditingDevice(null);
      setShowForm(false);
    }
  };

  const handleCancel = () => {
    setEditingDevice(null);
    setShowForm(false);
  };

  const handleEdit = (device: Device) => {
    setEditingDevice(device);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this device?')) {
      await deleteDevice(id);
    }
  };

  return (
    <div className="app">
      <header>
        <h1>Device Management</h1>
        <button onClick={() => setShowForm(true)}>Add Device</button>
      </header>

      <DeviceFilters
        hostname={hostnameFilter}
        isActiveFilter={isActiveFilter}
        onHostnameChange={setHostnameFilter}
        onIsActiveFilterChange={setIsActiveFilter}
      />

      {loading && <p>Loading...</p>}
      {error && <p className="error">{error}</p>}

      {showForm ? (
        <DeviceForm
          device={editingDevice}
          onSubmit={editingDevice ? handleUpdate : handleCreate}
          onCancel={handleCancel}
        />
      ) : (
        <DeviceTable devices={devices} onEdit={handleEdit} onDelete={handleDelete} />
      )}
    </div>
  );
}
