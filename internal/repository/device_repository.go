package repository

import (
	"errors"
	"strings"

	appErrors "github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/dto"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device *model.Device) error
	GetByID(id uint) (*model.Device, error)
	List(isActive *bool, hostname *string, limit, offset int) ([]model.Device, error)
	Update(id uint, req dto.UpdateDeviceRequest) error
	SoftDelete(id uint) error
}

type deviceRepo struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewDeviceRepository(db *gorm.DB, logger *logger.Logger) DeviceRepository {
	return &deviceRepo{
		db:     db,
		logger: logger,
	}
}

func (r *deviceRepo) Create(device *model.Device) error {
	if r.logger != nil {
		r.logger.Info("creating device",
			"hostname", device.Hostname,
			"ip", device.IP,
			"location", device.Location,
		)
	}

	if err := r.db.Create(device).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("failed to create device", "error", err)
		}
		return err
	}

	return nil
}

func (r *deviceRepo) GetByID(id uint) (*model.Device, error) {
	var device model.Device

	err := r.db.First(&device, id).Error
	if err != nil {
		if r.logger != nil {
			r.logger.Error("failed to fetch device", "id", id, "error", err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.ErrNotFound
		}

		return nil, err
	}

	return &device, nil
}

func (r *deviceRepo) List(isActive *bool, hostname *string, limit, offset int) ([]model.Device, error) {
	var devices []model.Device

	query := r.db.Model(&model.Device{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if hostname != nil && *hostname != "" {
		query = query.Where("LOWER(hostname) LIKE ?", "%"+strings.ToLower(*hostname)+"%")
	}

	// pagination
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.
		Order("id ASC").
		Find(&devices).Error; err != nil {

		if r.logger != nil {
			r.logger.Error("failed to fetch devices", "error", err)
		}
		return nil, err
	}

	return devices, nil
}

func (r *deviceRepo) Update(id uint, req dto.UpdateDeviceRequest) error {
	updates := map[string]interface{}{}

	if req.Hostname != nil {
		updates["hostname"] = *req.Hostname
	}
	if req.IP != nil {
		updates["ip"] = *req.IP
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return nil
	}

	res := r.db.Model(&model.Device{}).
		Where("id = ?", id).
		Updates(updates)

	if res.Error != nil {
		if r.logger != nil {
			r.logger.Error("failed to update device", "id", id, "error", res.Error)
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return appErrors.ErrNotFound
	}

	return nil
}

func (r *deviceRepo) SoftDelete(id uint) error {
	res := r.db.Delete(&model.Device{}, id)

	if res.Error != nil {
		if r.logger != nil {
			r.logger.Error("failed to delete device", "id", id, "error", res.Error)
		}
		return res.Error
	}

	if res.RowsAffected == 0 {
		return appErrors.ErrNotFound
	}

	return nil
}