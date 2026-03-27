import { Device } from '../types/device';

interface DeviceTableProps {
  devices: Device[];
  onEdit: (device: Device) => void;
  onDelete: (id: number) => void;
}

export function DeviceTable({ devices, onEdit, onDelete }: DeviceTableProps) {
  if (devices.length === 0) {
    return <p className="no-data">No devices found.</p>;
  }

  return (
    <table className="device-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Hostname</th>
          <th>IP Address</th>
          <th>Location</th>
          <th>Status</th>
          <th>Created At</th>
          <th>Actions</th>
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
              <span className={`status ${device.is_active ? 'active' : 'inactive'}`}>
                {device.is_active ? 'Active' : 'Inactive'}
              </span>
            </td>
            <td>{new Date(device.created_at).toLocaleString()}</td>
            <td className="actions">
              <button onClick={() => onEdit(device)}>Edit</button>
              <button onClick={() => onDelete(device.id)} className="delete">
                Delete
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
