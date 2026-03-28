import { memo } from 'react';

interface DeviceFiltersProps {
  hostname: string;
  isActiveFilter: string;
  onHostnameChange: (value: string) => void;
  onIsActiveFilterChange: (value: string) => void;
}

export const DeviceFilters = memo(function DeviceFilters({
  hostname,
  isActiveFilter,
  onHostnameChange,
  onIsActiveFilterChange,
}: DeviceFiltersProps) {
  return (
    <div className="filters" role="search">
      <div className="filter-group">
        <label htmlFor="hostname-filter">Search by Hostname:</label>
        <input
          id="hostname-filter"
          type="text"
          value={hostname}
          onChange={(e) => onHostnameChange(e.target.value)}
          placeholder="Enter hostname..."
          aria-label="Filter by hostname"
        />
      </div>
      
      <div className="filter-group">
        <label htmlFor="status-filter">Status:</label>
        <select
          id="status-filter"
          value={isActiveFilter}
          onChange={(e) => onIsActiveFilterChange(e.target.value)}
          aria-label="Filter by status"
        >
          <option value="">All</option>
          <option value="true">Active</option>
          <option value="false">Inactive</option>
        </select>
      </div>
    </div>
  );
});
