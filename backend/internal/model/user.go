package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Username           string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash       string         `gorm:"size:256;not null" json:"-"`
	Nickname           string         `gorm:"size:64" json:"nickname"`
	Email              []byte         `json:"-"`
	EmailPlain         string         `gorm:"-" json:"email,omitempty"`
	Phone              []byte         `json:"-"`
	PhonePlain         string         `gorm:"-" json:"phone,omitempty"`
	Avatar             string         `gorm:"size:512" json:"avatar"`
	Status             int16          `gorm:"not null;default:1" json:"status"` // 1=启用 0=禁用 2=锁定
	LoginAttempts      int16          `gorm:"not null;default:0" json:"-"`
	LockedUntil        *time.Time     `json:"-"`
	LastLoginAt        *time.Time     `json:"lastLoginAt"`
	LastLoginIP        string         `gorm:"size:45" json:"-"`
	PasswordChangedAt  *time.Time     `json:"-"`
	MustChangePassword bool           `gorm:"default:false" json:"-"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	CollegeID          *uint64        `gorm:"index" json:"departmentId"`
	College            *College       `gorm:"foreignKey:CollegeID" json:"department,omitempty"`
	Roles              []Role         `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

func (User) TableName() string { return "users" }

type Role struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Code        string         `gorm:"uniqueIndex;size:64;not null" json:"code"`
	Description string         `gorm:"size:256" json:"description"`
	Status      int16          `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Menus       []Menu         `gorm:"many2many:role_menus" json:"menus,omitempty"`
}

func (Role) TableName() string { return "roles" }

type Menu struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID   *uint64        `gorm:"index" json:"parentId"`
	Name       string         `gorm:"size:64;not null" json:"name"`
	Path       string         `gorm:"size:256" json:"path"`
	Component  string         `gorm:"size:256" json:"component"`
	Icon       string         `gorm:"size:64" json:"icon"`
	SortOrder  int            `gorm:"not null;default:0" json:"sortOrder"`
	Permission string         `gorm:"size:128" json:"permission"`
	MenuType   int16          `gorm:"not null;default:2" json:"menuType"` // 1=目录 2=菜单 3=按钮
	Status     int16          `gorm:"not null;default:1" json:"status"`
	Visible    bool           `gorm:"not null;default:true" json:"visible"`
	KeepAlive  bool           `gorm:"not null;default:false" json:"keepAlive"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Children   []Menu         `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string { return "menus" }

type UserRole struct {
	UserID uint64 `gorm:"primaryKey"`
	RoleID uint64 `gorm:"primaryKey"`
}

func (UserRole) TableName() string { return "user_roles" }

type RoleMenu struct {
	RoleID uint64 `gorm:"primaryKey"`
	MenuID uint64 `gorm:"primaryKey"`
}

func (RoleMenu) TableName() string { return "role_menus" }
