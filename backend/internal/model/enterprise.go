package model

import (
	"time"

	"gorm.io/gorm"
)

type Industry struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Code      string    `gorm:"uniqueIndex;size:32" json:"code"`
	ParentID  *uint64   `gorm:"index" json:"parentId"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	Status    int16     `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Industry) TableName() string { return "industries" }

type Enterprise struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          *uint64        `gorm:"index" json:"userId"`
	Name            string         `gorm:"size:128;not null" json:"name"`
	ShortName       string         `gorm:"size:64" json:"shortName"`
	IndustryID      *uint64        `gorm:"index" json:"industryId"`
	Scale           string         `gorm:"size:16" json:"scale"`
	Nature          string         `gorm:"size:16" json:"nature"`
	City            string         `gorm:"size:64" json:"city"`
	Address         string         `gorm:"size:256" json:"address"`
	Website         string         `gorm:"size:256" json:"website"`
	LogoURL         string         `gorm:"size:512" json:"logoUrl"`
	Description     string         `gorm:"type:text" json:"description"`
	BusinessScope   string         `gorm:"type:text" json:"businessScope"`
	ContactName     string         `gorm:"size:64" json:"contactName"`
	ContactPhone    string         `gorm:"size:20" json:"contactPhone"`
	ContactEmail    string         `gorm:"size:128" json:"contactEmail"`
	ContactPosition string         `gorm:"size:64" json:"contactPosition"`
	Status          int16          `gorm:"default:1" json:"status"`
	Verified        bool           `gorm:"default:false" json:"verified"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Industry        *Industry      `gorm:"foreignKey:IndustryID" json:"industry,omitempty"`
}

func (Enterprise) TableName() string { return "enterprises" }
