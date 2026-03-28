package model

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Hostname  string         `gorm:"size:255;not null" json:"hostname"`
	IP        string         `gorm:"type:inet;not null" json:"ip"`
	Location  string         `gorm:"type:text;not null" json:"location"`
	IsActive  bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
