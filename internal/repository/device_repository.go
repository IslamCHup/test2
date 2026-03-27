package repository

import (
"time"

"github.com/islamchupanov/tz1/internal/logger"
"github.com/islamchupanov/tz1/internal/model"
"gorm.io/gorm"
)

type DeviceRepository interface {
Create(device *model.Device) error
GetByID(id uint) (*model.Device, error)
List(isActive *bool, hostname *string) ([]model.Device, error)
Update(device *model.Device) error
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
r.logger.Info("creating device: hostname=%s ip=%s location=%s", device.Hostname, device.IP, device.Location)

if err := r.db.Create(device).Error; err != nil {
r.logger.Error("failed to create device: %v", err)
return err
}

r.logger.Info("device created: id=%d", device.ID)
return nil
}

func (r *deviceRepo) GetByID(id uint) (*model.Device, error) {
var device model.Device

r.logger.Info("fetching device by id: id=%d", id)

err := r.db.
Where("id = ? AND deleted_at IS NULL", id).
First(&device).Error

if err != nil {
r.logger.Error("failed to fetch device: id=%d error=%v", id, err)
return nil, err
}

r.logger.Info("device fetched: id=%d hostname=%s", device.ID, device.Hostname)

return &device, nil
}

func (r *deviceRepo) List(isActive *bool, hostname *string) ([]model.Device, error) {
var devices []model.Device

r.logger.Info("fetching devices list: is_active=%v hostname=%v", isActive, hostname)

query := r.db.Where("deleted_at IS NULL")

if isActive != nil {
query = query.Where("is_active = ?", *isActive)
}

if hostname != nil && *hostname != "" {
query = query.Where("LOWER(hostname) LIKE ?", "%"+*hostname+"%")
}

if err := query.
Order("id asc").
Find(&devices).Error; err != nil {

r.logger.Error("failed to fetch devices: error=%v", err)
return nil, err
}

r.logger.Info("devices fetched: count=%d", len(devices))
return devices, nil
}

func (r *deviceRepo) Update(device *model.Device) error {
r.logger.Info("updating device: id=%d", device.ID)

res := r.db.Model(&model.Device{}).
Where("id = ? AND deleted_at IS NULL", device.ID).
Updates(map[string]interface{}{
"hostname":  device.Hostname,
"ip":        device.IP,
"location":  device.Location,
"is_active": device.IsActive,
})

if res.Error != nil {
r.logger.Error("failed to update device: id=%d error=%v", device.ID, res.Error)
return res.Error
}

if res.RowsAffected == 0 {
r.logger.Warn("device not found for update: id=%d", device.ID)
return gorm.ErrRecordNotFound
}

r.logger.Info("device updated: id=%d", device.ID)
return nil
}

func (r *deviceRepo) SoftDelete(id uint) error {
r.logger.Info("soft deleting device: id=%d", id)

res := r.db.Model(&model.Device{}).
Where("id = ? AND deleted_at IS NULL", id).
Updates(map[string]interface{}{
"deleted_at": time.Now(),
"is_active":  false,
})

if res.Error != nil {
r.logger.Error("failed to soft delete device: id=%d error=%v", id, res.Error)
return res.Error
}

if res.RowsAffected == 0 {
r.logger.Warn("device not found for delete: id=%d", id)
return gorm.ErrRecordNotFound
}

r.logger.Info("device soft deleted: id=%d", id)
return nil
}
