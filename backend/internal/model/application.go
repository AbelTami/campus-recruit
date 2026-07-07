package model

import (
	"time"
)

type Application struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID     uint64     `gorm:"index;not null" json:"studentId"`
	PositionID    uint64     `gorm:"index;not null" json:"positionId"`
	EnterpriseID  uint64     `gorm:"index;not null" json:"enterpriseId"`
	Status        string     `gorm:"size:16;default:pending" json:"status"`
	ResumeURL     string     `gorm:"size:512" json:"resumeUrl"`
	CoverLetter   string     `gorm:"type:text" json:"coverLetter"`
	InterviewAt   *time.Time `json:"interviewAt"`
	InterviewNote string     `gorm:"type:text" json:"interviewNote"`
	OfferSalary   *int       `json:"offerSalary"`
	OfferAt       *time.Time `json:"offerAt"`
	RejectReason  string     `gorm:"size:256" json:"rejectReason"`
	Remark        string     `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Student       *Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Position      *Position  `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Enterprise    *Enterprise `gorm:"foreignKey:EnterpriseID" json:"enterprise,omitempty"`
}

func (Application) TableName() string { return "applications" }

type ApplicationLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ApplicationID uint64    `gorm:"index;not null" json:"applicationId"`
	FromStatus    string    `gorm:"size:16" json:"fromStatus"`
	ToStatus      string    `gorm:"size:16;not null" json:"toStatus"`
	OperatorID    *uint64   `json:"operatorId"`
	OperatorName  string    `gorm:"size:64" json:"operatorName"`
	Note          string    `gorm:"type:text" json:"note"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (ApplicationLog) TableName() string { return "application_logs" }

type EmploymentRecord struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID       uint64    `gorm:"uniqueIndex;not null" json:"studentId"`
	CompanyName     string    `gorm:"size:128;not null" json:"companyName"`
	PositionName    string    `gorm:"size:128" json:"positionName"`
	IndustryID      *uint64   `gorm:"index" json:"industryId"`
	City            string    `gorm:"size:64" json:"city"`
	MonthlySalary   *int      `json:"monthlySalary"`
	AnnualSalary    *int      `json:"annualSalary"`
	ContractType    string    `gorm:"size:32" json:"contractType"`
	EmploymentDate  *time.Time `json:"employmentDate"`
	ProbationMonths *int      `json:"probationMonths"`
	SocialSecurity  *bool     `json:"socialSecurity"`
	DataSource      string    `gorm:"size:16;default:manual" json:"dataSource"`
	Verified        bool      `gorm:"default:false" json:"verified"`
	VerifiedBy      *uint64   `json:"verifiedBy"`
	VerifiedAt      *time.Time `json:"verifiedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (EmploymentRecord) TableName() string { return "employment_records" }
