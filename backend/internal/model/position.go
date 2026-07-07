package model

import (
	"time"

	"gorm.io/gorm"
)

type Position struct {
	ID                    uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	EnterpriseID          uint64         `gorm:"index;not null" json:"enterpriseId"`
	Title                 string         `gorm:"size:128;not null" json:"title"`
	IndustryID            *uint64        `gorm:"index" json:"industryId"`
	City                  string         `gorm:"size:64" json:"city"`
	District              string         `gorm:"size:64" json:"district"`
	Address               string         `gorm:"size:256" json:"address"`
	EducationRequirement  string         `gorm:"size:16" json:"educationRequirement"`
	ExperienceRequirement *int           `json:"experienceRequirement"`
	SalaryMin             *int           `json:"salaryMin"`
	SalaryMax             *int           `json:"salaryMax"`
	SalaryType            string         `gorm:"size:16;default:monthly" json:"salaryType"`
	Headcount             int            `gorm:"default:1" json:"headcount"`
	JobType               string         `gorm:"size:16;default:fulltime" json:"jobType"`
	Description           string         `gorm:"type:text" json:"description"`
	Requirement           string         `gorm:"type:text" json:"requirement"`
	Welfare               string         `gorm:"type:text" json:"welfare"`
	ContactInfo           string         `gorm:"size:256" json:"contactInfo"`
	Status                int16          `gorm:"default:1" json:"status"` // 1=招聘中 0=下架 2=过期
	PublishAt             *time.Time     `json:"publishAt"`
	ExpireAt              *time.Time     `json:"expireAt"`
	ViewCount             int            `gorm:"default:0" json:"viewCount"`
	ApplyCount            int            `gorm:"default:0" json:"applyCount"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
	Enterprise            *Enterprise    `gorm:"foreignKey:EnterpriseID" json:"enterprise,omitempty"`
	Industry              *Industry      `gorm:"foreignKey:IndustryID" json:"industry,omitempty"`
	Skills                []Skill        `gorm:"many2many:position_skills" json:"skills,omitempty"`
}

func (Position) TableName() string { return "positions" }

type PositionSkill struct {
	PositionID uint64 `gorm:"primaryKey" json:"positionId"`
	SkillID    uint64 `gorm:"primaryKey" json:"skillId"`
	Importance int16  `gorm:"default:1;check:importance between 1 and 5" json:"importance"`
}

func (PositionSkill) TableName() string { return "position_skills" }
