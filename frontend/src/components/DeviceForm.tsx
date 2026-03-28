import { useState, useEffect } from 'react';
import { Device } from '../types/device';

interface DeviceFormProps {
  device?: Device | null;
  onSubmit: (data: { hostname: string; ip: string; location?: string; is_active?: boolean }) => Promise<void>;
  onCancel: () => void;
  error?: string | null;
}

export function DeviceForm({ device, onSubmit, onCancel, error }: DeviceFormProps) {
  const [hostname, setHostname] = useState('');
  const [ip, setIp] = useState('');
  const [location, setLocation] = useState('');
  const [isActive, setIsActive] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (device) {
      setHostname(device.hostname);
      setIp(device.ip);
      setLocation(device.location ?? '');
      setIsActive(device.is_active);
    } else {
      // Reset form when creating new device
      setHostname('');
      setIp('');
      setLocation('');
      setIsActive(true);
    }
  }, [device]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      await onSubmit({ hostname, ip, location: location || undefined, is_active: isActive });
    } catch (err) {
      // Error is handled by parent component via the error prop
      console.error('Form submission error:', err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="device-form" aria-label={device ? 'Edit device form' : 'Create device form'}>
      <h2>{device ? 'Edit Device' : 'Add Device'}</h2>
      
      {error && <div className="error" role="alert">{error}</div>}
      
      <div className="form-group">
        <label htmlFor="hostname">Hostname *</label>
        <input
          id="hostname"
          type="text"
          value={hostname}
          onChange={(e) => setHostname(e.target.value)}
          required
          placeholder="e.g., router-01"
          disabled={submitting}
          autoComplete="off"
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="ip">IP Address *</label>
        <input
          id="ip"
          type="text"
          value={ip}
          onChange={(e) => setIp(e.target.value)}
          required
          placeholder="e.g., 192.168.1.1"
          disabled={submitting}
          autoComplete="off"
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="location">Location</label>
        <input
          id="location"
          type="text"
          value={location}
          onChange={(e) => setLocation(e.target.value)}
          placeholder="e.g., Data Center A"
          disabled={submitting}
          autoComplete="off"
        />
      </div>
      
      <div className="form-group checkbox">
        <label>
          <input
            type="checkbox"
            checked={isActive}
            onChange={(e) => setIsActive(e.target.checked)}
            disabled={submitting}
          />
          Active
        </label>
      </div>
      
      <div className="form-actions">
        <button type="button" onClick={onCancel} disabled={submitting}>
          Cancel
        </button>
        <button type="submit" disabled={submitting}>
          {submitting ? 'Saving...' : device ? 'Update' : 'Create'}
        </button>
      </div>
    </form>
  );
}
