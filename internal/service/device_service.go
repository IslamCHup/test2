package service

import (
	"github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/repository"
	"gorm.io/gorm"
)

type DeviceService interface {
	Create(device *model.Device) error
	GetByID(id uint) (*model.Device, error)
	List(isActive *bool, hostname *string) ([]model.Device, error)
	Update(id uint, device *model.Device) (*model.Device, error)
	SoftDelete(id uint) error
}

type deviceService struct {
	repo   repository.DeviceRepository
	logger *logger.Logger
}

func NewDeviceService(repo repository.DeviceRepository, logger *logger.Logger) DeviceService {
	return &deviceService{repo: repo, logger: logger}
}

func (s *deviceService) Create(device *model.Device) error {
	s.logger.Info("service: creating device", "hostname", device.Hostname, "ip", device.IP)
	err := s.repo.Create(device)
	if err != nil {
		s.logger.Error("service: failed to create device", "error", err)
	} else {
		s.logger.Info("service: device created successfully", "id", device.ID)
	}
	return err
}

func (s *deviceService) GetByID(id uint) (*model.Device, error) {
	s.logger.Info("service: fetching device by id", "id", id)
	d, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn("service: device not found", "id", id)
			return nil, errors.ErrNotFound
		}
		s.logger.Error("service: failed to fetch device", "id", id, "error", err)
		return nil, err
	}
	s.logger.Info("service: device fetched successfully", "id", id, "hostname", d.Hostname)
	return d, nil
}

func (s *deviceService) List(isActive *bool, hostname *string) ([]model.Device, error) {
	s.logger.Info("service: fetching devices list", "is_active", isActive, "hostname", hostname)
	devices, err := s.repo.List(isActive, hostname)
	if err != nil {
		s.logger.Error("service: failed to fetch devices", "error", err)
	} else {
		s.logger.Info("service: devices fetched successfully", "count", len(devices))
	}
	return devices, err
}

func (s *deviceService) Update(id uint, device *model.Device) (*model.Device, error) {
	s.logger.Info("service: updating device", "id", id)
	dbDevice, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Warn("service: device not found for update", "id", id)
		return nil, errors.ErrNotFound
	}

	dbDevice.Hostname = device.Hostname
	dbDevice.IP = device.IP
	dbDevice.Location = device.Location
	dbDevice.IsActive = device.IsActive

	if err := s.repo.Update(dbDevice); err != nil {
		s.logger.Error("service: failed to update device", "id", id, "error", err)
		return nil, err
	}

	s.logger.Info("service: device updated successfully", "id", id)
	return dbDevice, nil
}

func (s *deviceService) SoftDelete(id uint) error {
	s.logger.Info("service: soft deleting device", "id", id)
	if _, err := s.repo.GetByID(id); err != nil {
		s.logger.Warn("service: device not found for delete", "id", id)
		return errors.ErrNotFound
	}
	err := s.repo.SoftDelete(id)
	if err != nil {
		s.logger.Error("service: failed to soft delete device", "id", id, "error", err)
	} else {
		s.logger.Info("service: device soft deleted successfully", "id", id)
	}
	return err
}
