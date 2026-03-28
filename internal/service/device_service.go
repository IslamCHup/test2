package service

import (
	"errors"
	"net"
	"strings"

	appErrors "github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/dto"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/repository"
)

type DeviceService interface {
	Create(device *model.Device) error
	GetByID(id uint) (*model.Device, error)
	List(isActive *bool, hostname *string, limit, offset int) ([]model.Device, error)
	Update(id uint, req dto.UpdateDeviceRequest) (*model.Device, error)
	SoftDelete(id uint) error
}

type deviceService struct {
	repo   repository.DeviceRepository
	logger *logger.Logger
}

func NewDeviceService(repo repository.DeviceRepository, logger *logger.Logger) DeviceService {
	return &deviceService{
		repo:   repo,
		logger: logger,
	}
}

// ================= VALIDATION =================

func validateHostname(hostname string) error {
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return errors.New("hostname cannot be empty")
	}
	return nil
}

func validateIP(ip string) error {
	ip = strings.TrimSpace(ip)
	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip address")
	}
	return nil
}

// ================= CREATE =================

func (s *deviceService) Create(device *model.Device) error {
	device.Hostname = strings.TrimSpace(device.Hostname)
	device.IP = strings.TrimSpace(device.IP)
	device.Location = strings.TrimSpace(device.Location)

	if err := validateHostname(device.Hostname); err != nil {
		return err
	}
	if err := validateIP(device.IP); err != nil {
		return err
	}

	if err := s.repo.Create(device); err != nil {
		s.logger.Error("failed to create device", "error", err)
		return err
	}
	return nil
}

// ================= GET =================

func (s *deviceService) GetByID(id uint) (*model.Device, error) {
	device, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil, appErrors.ErrNotFound
		}
		s.logger.Error("failed to get device", "id", id, "error", err)
		return nil, err
	}
	return device, nil
}

// ================= LIST =================

func (s *deviceService) List(isActive *bool, hostname *string, limit, offset int) ([]model.Device, error) {
	devices, err := s.repo.List(isActive, hostname, limit, offset)
	if err != nil {
		s.logger.Error("failed to list devices", "error", err)
		return nil, err
	}
	return devices, nil
}

// ================= UPDATE =================

func (s *deviceService) Update(id uint, req dto.UpdateDeviceRequest) (*model.Device, error) {

	if req.Hostname != nil {
		trimmed := strings.TrimSpace(*req.Hostname)
		if err := validateHostname(trimmed); err != nil {
			return nil, err
		}
		req.Hostname = &trimmed
	}

	if req.IP != nil {
		trimmed := strings.TrimSpace(*req.IP)
		if err := validateIP(trimmed); err != nil {
			return nil, err
		}
		req.IP = &trimmed
	}

	if req.Location != nil {
		trimmed := strings.TrimSpace(*req.Location)
		req.Location = &trimmed
	}

	err := s.repo.Update(id, req)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return nil, appErrors.ErrNotFound
		}
		s.logger.Error("failed to update device", "id", id, "error", err)
		return nil, err
	}

	device, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to fetch updated device", "id", id, "error", err)
		return nil, err
	}

	return device, nil
}

// ================= DELETE =================

func (s *deviceService) SoftDelete(id uint) error {
	err := s.repo.SoftDelete(id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return appErrors.ErrNotFound
		}
		s.logger.Error("failed to delete device", "id", id, "error", err)
		return err
	}

	return nil
}