package service

import (
	"strings"
	"testing"

	appErrors "github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/dto"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
)

// ================= MOCK =================

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
	device, ok := m.devices[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	return device, nil
}

func (m *MockDeviceRepository) List(isActive *bool, hostname *string, limit, offset int) ([]model.Device, error) {
	var result []model.Device

	for _, device := range m.devices {
		if isActive != nil && device.IsActive != *isActive {
			continue
		}

		if hostname != nil && *hostname != "" {
			if !strings.Contains(strings.ToLower(device.Hostname), strings.ToLower(*hostname)) {
				continue
			}
		}

		result = append(result, *device)
	}

	// pagination
	start := offset
	end := offset + limit

	if start > len(result) {
		return []model.Device{}, nil
	}
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

func (m *MockDeviceRepository) Update(id uint, req dto.UpdateDeviceRequest) error {
	device, ok := m.devices[id]
	if !ok {
		return appErrors.ErrNotFound
	}

	if req.Hostname != nil {
		device.Hostname = *req.Hostname
	}
	if req.IP != nil {
		device.IP = *req.IP
	}
	if req.Location != nil {
		device.Location = *req.Location
	}
	if req.IsActive != nil {
		device.IsActive = *req.IsActive
	}

	return nil
}

func (m *MockDeviceRepository) SoftDelete(id uint) error {
	device, ok := m.devices[id]
	if !ok {
		return appErrors.ErrNotFound
	}

	device.IsActive = false
	return nil
}

// ================= TESTS =================

func TestCreateDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "test-router",
		IP:       "192.168.1.1",
		Location: "dc",
		IsActive: true,
	}

	err := service.Create(device)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if device.ID == 0 {
		t.Fatal("expected ID to be set")
	}
}

func TestCreateDevice_InvalidIP(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "test",
		IP:       "invalid-ip",
	}

	err := service.Create(device)
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetByID(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "test",
		IP:       "192.168.1.1",
		Location: "loc",
		IsActive: true,
	}
	repo.Create(device)

	found, err := service.GetByID(device.ID)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}

	if found.ID != device.ID {
		t.Errorf("wrong id")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	_, err := service.GetByID(999)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListDevices(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	devices := []*model.Device{
		{Hostname: "router-main", IP: "192.168.1.1", Location: "dc1", IsActive: true},
		{Hostname: "switch", IP: "192.168.1.2", Location: "office", IsActive: true},
		{Hostname: "router-backup", IP: "192.168.1.3", Location: "dc2", IsActive: false},
	}

	for _, d := range devices {
		repo.Create(d)
	}

	active := true
	list, err := service.List(&active, nil, 10, 0)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 active devices")
	}

	search := "router"
	list, err = service.List(nil, &search, 10, 0)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 router devices")
	}
}

func TestListDevices_Pagination(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	for i := 0; i < 5; i++ {
		repo.Create(&model.Device{
			Hostname: "device",
			IP:       "192.168.1.1",
			IsActive: true,
		})
	}

	list, err := service.List(nil, nil, 2, 0)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 devices, got %d", len(list))
	}
}

func TestUpdateDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "old",
		IP:       "192.168.1.1",
		Location: "loc",
		IsActive: true,
	}
	repo.Create(device)

	newHostname := "new"
	newIP := "192.168.1.2"

	req := dto.UpdateDeviceRequest{
		Hostname: &newHostname,
		IP:       &newIP,
	}

	updated, err := service.Update(device.ID, req)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if updated.Hostname != "new" {
		t.Errorf("hostname not updated")
	}

	if updated.IP != "192.168.1.2" {
		t.Errorf("ip not updated")
	}
}

func TestSoftDeleteDevice(t *testing.T) {
	repo := NewMockDeviceRepository()
	log := logger.InitLog("debug")
	service := NewDeviceService(repo, log)

	device := &model.Device{
		Hostname: "to-delete",
		IP:       "192.168.1.1",
		Location: "loc",
		IsActive: true,
	}
	repo.Create(device)

	err := service.SoftDelete(device.ID)
	if err != nil {
		t.Fatalf("SoftDelete error: %v", err)
	}

	found, _ := repo.GetByID(device.ID)

	if found.IsActive {
		t.Errorf("device not deactivated")
	}
}