package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/islamchupanov/tz1/internal/dto"
	apperrors "github.com/islamchupanov/tz1/internal/errors"
	"github.com/islamchupanov/tz1/internal/logger"
	"github.com/islamchupanov/tz1/internal/model"
	"github.com/islamchupanov/tz1/internal/service"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device := &model.Device{
		Hostname: req.Hostname,
		IP:       req.IP,
		Location: req.Location,
		IsActive: true,
	}

	if err := h.service.Create(device); err != nil {
		if err.Error() == "hostname cannot be empty" || err.Error() == "invalid ip address" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("failed to create device", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, toResponse(device))
}

// ListDevices godoc
// @Summary List all devices
// @Description Get list of devices with optional filtering
// @Tags devices
// @Accept json
// @Produce json
// @Param is_active query string false "Filter by is_active (true/false)"
// @Param hostname query string false "Search by hostname (substring)"
// @Param limit query int false "Limit number of results (max 100)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} dto.DeviceResponse
// @Failure 400 {object} map[string]string
// @Router /devices [get]
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	var isActive *bool
	var hostname *string

	if v := c.Query("is_active"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid is_active"})
			return
		}
		isActive = &parsed
	}

	if v := c.Query("hostname"); v != "" {
		v = strings.TrimSpace(v)
		hostname = &v
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	devices, err := h.service.List(isActive, hostname, limit, offset)
	if err != nil {
		h.logger.Error("failed to fetch devices", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch devices"})
		return
	}

	response := make([]dto.DeviceResponse, 0, len(devices))
	for i := range devices {
		response = append(response, toResponse(&devices[i]))
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [get]
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.logger.Warn("invalid id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	device, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		h.logger.Error("failed to fetch device", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch device"})
		return
	}

	c.JSON(http.StatusOK, toResponse(device))
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
	id, err := parseID(c)
	if err != nil {
		h.logger.Warn("invalid id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}

		if err.Error() == "hostname cannot be empty" || err.Error() == "invalid ip address" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error("failed to update device", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update device"})
		return
	}

	c.JSON(http.StatusOK, toResponse(device))
}

// DeleteDevice godoc
// @Summary Delete device (soft delete)
// @Description Soft delete a device by marking it as inactive
// @Tags devices
// @Accept json
// @Produce json
// @Param id path int true "Device ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.logger.Warn("invalid id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return
	}

	if err := h.service.SoftDelete(id); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		h.logger.Error("failed to delete device", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete device"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ================= HELPERS =================

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func toResponse(d *model.Device) dto.DeviceResponse {
	return dto.DeviceResponse{
		ID:        d.ID,
		Hostname:  d.Hostname,
		IP:        d.IP,
		Location:  d.Location,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
	}
}