package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/islamchupanov/tz1/internal/dto"
	apperrors "github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceHandler struct {
	service service.DeviceService
	logger  *logger.Logger
}

func NewDeviceHandler(service service.DeviceService, logger *logger.Logger) *DeviceHandler {
	return &DeviceHandler{service: service, logger: logger}
}

// CreateDevice godoc
// @Summary Create a new device
// @Description Create a new network device
// @Tags devices
// @Accept json
// @Produce json
// @Param device body dto.CreateDeviceRequest true "Device data"
// @Success 201 {object} dto.DeviceResponse
// @Failure 400 {object} map[string]string
// @Router /devices [post]
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	device := &model.Device{
		Hostname: req.Hostname,
		IP:       req.IP,
		Location: req.Location,
		IsActive: true,
	}

	h.logger.Info("handler: creating device", "hostname", device.Hostname, "ip", device.IP)
	if err := h.service.Create(device); err != nil {
		h.logger.Error("handler: failed to create device", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device: " + err.Error()})
		return
	}
	h.logger.Info("handler: device created successfully", "id", device.ID)

	response := dto.DeviceResponse{
		ID:        device.ID,
		Hostname:  device.Hostname,
		IP:        device.IP,
		Location:  device.Location,
		IsActive:  device.IsActive,
		CreatedAt: device.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// ListDevices godoc
// @Summary List all devices
// @Description Get list of devices with optional filtering
// @Tags devices
// @Accept json
// @Produce json
// @Param is_active query string false "Filter by is_active (true/false)"
// @Param hostname query string false "Search by hostname (substring)"
// @Success 200 {array} dto.DeviceResponse
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	var isActive *bool
	var hostname *string

	if isActiveParam := c.Query("is_active"); isActiveParam != "" {
		val := isActiveParam == "true"
		isActive = &val
	}

	if hostnameParam := c.Query("hostname"); hostnameParam != "" {
		hostname = &hostnameParam
	}

	h.logger.Info("handler: fetching devices", "is_active", isActive, "hostname", hostname)
	devices, err := h.service.List(isActive, hostname)
	if err != nil {
		h.logger.Error("handler: failed to fetch devices", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices: " + err.Error()})
		return
	}
	h.logger.Info("handler: devices fetched successfully", "count", len(devices))

	response := make([]dto.DeviceResponse, len(devices))
	for i, d := range devices {
		response[i] = dto.DeviceResponse{
			ID:        d.ID,
			Hostname:  d.Hostname,
			IP:        d.IP,
			Location:  d.Location,
			IsActive:  d.IsActive,
			CreatedAt: d.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetDevice godoc
// @Summary Get device by ID
// @Description Get a specific device by its ID
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} dto.DeviceResponse
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("invalid device id", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	h.logger.Info("handler: fetching device", "id", id)
	device, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("handler: device not found", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		h.logger.Error("handler: failed to fetch device", "id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch device: " + err.Error()})
		return
	}
	h.logger.Info("handler: device fetched successfully", "id", device.ID)

	response := dto.DeviceResponse{
		ID:        device.ID,
		Hostname:  device.Hostname,
		IP:        device.IP,
		Location:  device.Location,
		IsActive:  device.IsActive,
		CreatedAt: device.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateDevice godoc
// @Summary Update device
// @Description Update an existing device
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Param device body dto.UpdateDeviceRequest true "Device data"
// @Success 200 {object} dto.DeviceResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [put]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("invalid device id", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	h.logger.Info("handler: fetching device for update", "id", id)
	existingDevice, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("handler: device not found for update", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		h.logger.Error("handler: failed to fetch device for update", "id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch device: " + err.Error()})
		return
	}

	// Update fields if provided
	if req.Hostname != nil {
		existingDevice.Hostname = *req.Hostname
	}
	if req.IP != nil {
		existingDevice.IP = *req.IP
	}
	if req.Location != nil {
		existingDevice.Location = *req.Location
	}
	if req.IsActive != nil {
		existingDevice.IsActive = *req.IsActive
	}

	h.logger.Info("handler: updating device", "id", id)
	updatedDevice, err := h.service.Update(uint(id), existingDevice)
	if err != nil {
		h.logger.Error("handler: failed to update device", "id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device: " + err.Error()})
		return
	}
	h.logger.Info("handler: device updated successfully", "id", updatedDevice.ID)

	response := dto.DeviceResponse{
		ID:        updatedDevice.ID,
		Hostname:  updatedDevice.Hostname,
		IP:        updatedDevice.IP,
		Location:  updatedDevice.Location,
		IsActive:  updatedDevice.IsActive,
		CreatedAt: updatedDevice.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteDevice godoc
// @Summary Delete device (soft delete)
// @Description Soft delete a device by marking it as inactive
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("invalid device id", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	h.logger.Info("handler: deleting device", "id", id)
	if err := h.service.SoftDelete(uint(id)); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Warn("handler: device not found for delete", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		h.logger.Error("handler: failed to delete device", "id", id, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device: " + err.Error()})
		return
	}
	h.logger.Info("handler: device deleted successfully", "id", id)

	c.Status(http.StatusNoContent)
}
