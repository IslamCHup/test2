package service

import (
	"testing"

	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
)

// MockDeviceRepository implements DeviceRepository interface for testing
type MockDeviceRepository struct {
	devices map[uint]*model.Device
	nextID  uint
}

func NewMockDeviceRepository() *MockDeviceRepository {
	return &MockDeviceRepository{
		devices: make(map[uint]*model.Device),
		nextID:  1,
	}
}

func (m *MockDeviceRepository) Create(device *model.Device) error {
	device.ID = m.nextID
	m.nextID++
	m.devices[device.ID] = device
	return nil
}

func (m *MockDeviceRepository) GetByID(id uint) (*model.Device, error) {
	if device, ok := m.devices[id]; ok {
		return device, nil
	}
	return nil, nil // Simulate not found
}

func (m *MockDeviceRepository) List(isActive *bool, hostname *string) ([]model.Device, error) {
	var result []model.Device
	for _, device := range m.devices {
		if isActive != nil && device.IsActive != *isActive {
			continue
		}
		if hostname != nil && *hostname != "" {
			// Simple substring match (case-insensitive)
			found := false
			for i := 0; i <= len(device.Hostname)-len(*hostname); i++ {
				if device.Hostname[i:i+len(*hostname)] == *hostname {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, *device)
	}
	return result, nil
}

func (m *MockDeviceRepository) Update(device *model.Device) error {
	if _, ok := m.devices[device.ID]; ok {
		m.devices[device.ID] = device
		return nil
	}
	return nil
}

func (m *MockDeviceRepository) SoftDelete(id uint) error {
	if device, ok := m.devices[id]; ok {
		device.IsActive = false
		return nil
	}
	return nil
}

// TestCreateDevice tests device creation
func TestCreateDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "test-router",
		IP:       "192.168.1.1",
		Location: "datacenter-1",
		IsActive: true,
	}

	err := service.Create(device)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if device.ID == 0 {
		t.Error("Create() did not assign ID to device")
	}

	if device.Hostname != "test-router" {
		t.Errorf("Create() hostname = %v, want test-router", device.Hostname)
	}
}

// TestGetByIDDevice tests getting device by ID
func TestGetByIDDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	// Create a device first
	device := &model.Device{
		Hostname: "test-switch",
		IP:       "192.168.1.2",
		Location: "office",
		IsActive: true,
	}
	repo.Create(device)

	// Get the device
	found, err := service.GetByID(device.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if found.ID != device.ID {
		t.Errorf("GetByID() ID = %v, want %v", found.ID, device.ID)
	}

	if found.Hostname != "test-switch" {
		t.Errorf("GetByID() hostname = %v, want test-switch", found.Hostname)
	}
}

// TestListDevicesWithFilter tests listing devices with filters
func TestListDevicesWithFilter(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	// Create test devices
	devices := []*model.Device{
		{Hostname: "router-main", IP: "192.168.1.1", Location: "dc1", IsActive: true},
		{Hostname: "switch-floor1", IP: "192.168.1.2", Location: "office", IsActive: true},
		{Hostname: "router-backup", IP: "192.168.1.3", Location: "dc2", IsActive: false},
	}

	for _, d := range devices {
		repo.Create(d)
	}

	// Test filtering by is_active
	activeTrue := true
	activeDevices, err := service.List(&activeTrue, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(activeDevices) != 2 {
		t.Errorf("List() returned %d active devices, want 2", len(activeDevices))
	}

	// Test filtering by hostname substring
	searchTerm := "router"
	routerDevices, err := service.List(nil, &searchTerm)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(routerDevices) != 2 {
		t.Errorf("List() returned %d router devices, want 2", len(routerDevices))
	}
}

// TestUpdateDevice tests device update
func TestUpdateDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	// Create a device
	device := &model.Device{
		Hostname: "old-hostname",
		IP:       "192.168.1.1",
		Location: "old-location",
		IsActive: true,
	}
	repo.Create(device)

	// Update the device
	device.Hostname = "new-hostname"
	device.IP = "10.0.0.1"

	updated, err := service.Update(device.ID, device)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Hostname != "new-hostname" {
		t.Errorf("Update() hostname = %v, want new-hostname", updated.Hostname)
	}

	if updated.IP != "10.0.0.1" {
		t.Errorf("Update() ip = %v, want 10.0.0.1", updated.IP)
	}
}

// TestSoftDeleteDevice tests soft delete functionality
func TestSoftDeleteDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	// Create a device
	device := &model.Device{
		Hostname: "to-delete",
		IP:       "192.168.1.1",
		Location: "somewhere",
		IsActive: true,
	}
	repo.Create(device)

	// Soft delete
	err := service.SoftDelete(device.ID)
	if err != nil {
		t.Fatalf("SoftDelete() error = %v", err)
	}

	// Verify device is marked as inactive
	found, _ := repo.GetByID(device.ID)
	if found.IsActive {
		t.Error("SoftDelete() did not set IsActive to false")
	}
}
