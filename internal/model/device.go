package model

import "time"

type Device struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Hostname  string     `gorm:"type:text;not null" json:"hostname"`
	IP        string     `gorm:"type:text;not null" json:"ip"`
	Location  string     `gorm:"type:text;not null" json:"location"`
	IsActive  bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
