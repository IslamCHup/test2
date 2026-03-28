package dto

import "time"

type CreateDeviceRequest struct {
	Hostname string `json:"hostname" binding:"required,min=1,max=255"`
	IP       string `json:"ip" binding:"required,ip"`
	Location string `json:"location" binding:"required"`
}

type UpdateDeviceRequest struct {
	Hostname *string `json:"hostname" binding:"omitempty,min=1,max=255"`
	IP       *string `json:"ip" binding:"omitempty,ip"`
	Location *string `json:"location"`
	IsActive *bool   `json:"is_active"`
}

type DeviceResponse struct {
	ID        uint      `json:"id"`
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
