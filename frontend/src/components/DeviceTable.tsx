import { memo } from 'react';
import { Device } from '../types/device';

interface DeviceTableProps {
  devices: Device[];
  onEdit: (device: Device) => void;
  onDelete: (id: number) => void;
  isLoading?: boolean;
}

export const DeviceTable = memo(function DeviceTable({ 
  devices, 
  onEdit, 
  onDelete,
  isLoading = false 
}: DeviceTableProps) {
  if (isLoading) {
    return null; // Loading state is handled by parent
  }

  if (devices.length === 0) {
    return <p className="no-data" role="status">No devices found.</p>;
  }

  return (
    <table className="device-table" aria-label="Devices table">
      <thead>
        <tr>
          <th scope="col">ID</th>
          <th scope="col">Hostname</th>
          <th scope="col">IP Address</th>
          <th scope="col">Location</th>
          <th scope="col">Status</th>
          <th scope="col">Created At</th>
          <th scope="col">Actions</th>
        </tr>
      </thead>
      <tbody>
        {devices.map((device) => (
          <tr key={device.id}>
            <td>{device.id}</td>
            <td>{device.hostname}</td>
            <td>{device.ip}</td>
            <td>{device.location || '-'}</td>
            <td>
              <span 
                className={`status ${device.is_active ? 'active' : 'inactive'}`}
                role="status"
                aria-label={device.is_active ? 'Active' : 'Inactive'}
              >
                {device.is_active ? 'Active' : 'Inactive'}
              </span>
            </td>
            <td>
              <time dateTime={device.created_at}>
                {new Date(device.created_at).toLocaleString()}
              </time>
            </td>
            <td className="actions">
              <button 
                onClick={() => onEdit(device)}
                aria-label={`Edit device ${device.hostname}`}
              >
                Edit
              </button>
              <button 
                onClick={() => onDelete(device.id)} 
                className="delete"
                aria-label={`Delete device ${device.hostname}`}
              >
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
});
