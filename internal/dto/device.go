package dto

import "time"

type CreateDeviceRequest struct {
	Hostname string `json:"hostname" binding:"required"`
	IP       string `json:"ip" binding:"required,ip"`
	Location string `json:"location"`
}

type UpdateDeviceRequest struct {
	Hostname *string `json:"hostname"`
	IP       *string `json:"ip"`
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
