package service

import (
	"github.com/islamchupanov/tz1/internal/errors"
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
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) Create(device *model.Device) error {
	return s.repo.Create(device)
}

func (s *deviceService) GetByID(id uint) (*model.Device, error) {
	d, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return d, nil
}

func (s *deviceService) List(isActive *bool, hostname *string) ([]model.Device, error) {
	return s.repo.List(isActive, hostname)
}

func (s *deviceService) Update(id uint, device *model.Device) (*model.Device, error) {
	dbDevice, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	dbDevice.Hostname = device.Hostname
	dbDevice.IP = device.IP
	dbDevice.Location = device.Location
	dbDevice.IsActive = device.IsActive

	if err := s.repo.Update(dbDevice); err != nil {
		return nil, err
	}

	return dbDevice, nil
}

func (s *deviceService) SoftDelete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return errors.ErrNotFound
	}
	return s.repo.SoftDelete(id)
}
