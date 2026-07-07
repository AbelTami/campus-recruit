package model

import "time"

type OperationLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      *uint64   `json:"userId"`
	Username    string    `gorm:"size:64;not null" json:"username"`
	Role        string    `gorm:"size:64" json:"role"`
	Module      string    `gorm:"size:64;not null" json:"module"`
	Action      string    `gorm:"size:64;not null" json:"action"`
	TargetType  string    `gorm:"size:64" json:"targetType"`
	TargetID    *uint64   `json:"targetId"`
	Description string    `gorm:"size:512" json:"description"`
	Method      string    `gorm:"size:10" json:"method"`
	URL         string    `gorm:"size:256" json:"url"`
	IP          string    `gorm:"size:45;not null" json:"ip"`
	UserAgent   string    `gorm:"type:text" json:"userAgent"`
	RequestBody string    `gorm:"type:text" json:"requestBody"`
	StatusCode  int       `json:"statusCode"`
	CostMs      int       `json:"costMs"`
	IsDeleted   bool      `gorm:"default:false" json:"isDeleted"`
	CreatedAt   time.Time `gorm:"index" json:"createdAt"`
}

func (OperationLog) TableName() string { return "operation_logs" }

type DictData struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	DictType  string    `gorm:"size:64;not null;index" json:"dictType"`
	DictLabel string    `gorm:"size:128;not null" json:"dictLabel"`
	DictValue string    `gorm:"size:128;not null" json:"dictValue"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	Status    int16     `gorm:"default:1" json:"status"`
	Remark    string    `gorm:"size:256" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
}

func (DictData) TableName() string { return "dict_data" }

type AnalysisReport struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:256;not null" json:"title"`
	ReportType  string    `gorm:"size:32;not null" json:"reportType"`
	ReportScope string    `gorm:"size:16;default:university" json:"reportScope"`
	ScopeID     *uint64   `json:"scopeId"`
	ReportYear  *int      `json:"reportYear"`
	ReportPeriod string   `gorm:"size:16" json:"reportPeriod"`
	Content     string    `gorm:"type:jsonb" json:"content"`
	Summary     string    `gorm:"type:text" json:"summary"`
	FileURL     string    `gorm:"size:512" json:"fileUrl"`
	FileFormat  string    `gorm:"size:16" json:"fileFormat"`
	Status      string    `gorm:"size:16;default:draft" json:"status"`
	CreatedBy   *uint64   `json:"createdBy"`
	PublishedAt *time.Time `json:"publishedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (AnalysisReport) TableName() string { return "analysis_reports" }

type MarketDataSnapshot struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Source       string    `gorm:"size:128" json:"source"`
	IndustryID   *uint64   `gorm:"index" json:"industryId"`
	City         string    `gorm:"size:64" json:"city"`
	PositionName string    `gorm:"size:128" json:"positionName"`
	AvgSalary    *int      `json:"avgSalary"`
	SalaryP25    *int      `json:"salaryP25"`
	SalaryP50    *int      `json:"salaryP50"`
	SalaryP75    *int      `json:"salaryP75"`
	DemandCount  *int      `json:"demandCount"`
	SupplyCount  *int      `json:"supplyCount"`
	GrowthRate   *float64  `json:"growthRate"`
	SnapshotDate time.Time `gorm:"not null;index" json:"snapshotDate"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (MarketDataSnapshot) TableName() string { return "market_data_snapshots" }

type EmploymentSurvey struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID      *uint64   `json:"studentId"`
	SurveyYear     int       `gorm:"not null" json:"surveyYear"`
	SurveyType     string    `gorm:"size:16;default:graduate" json:"surveyType"`
	EmployStatus   string    `gorm:"size:16" json:"employStatus"`
	CompanyName    string    `gorm:"size:128" json:"companyName"`
	IndustryID     *uint64   `json:"industryId"`
	City           string    `gorm:"size:64" json:"city"`
	MonthlySalary  *int      `json:"monthlySalary"`
	PositionName   string    `gorm:"size:128" json:"positionName"`
	Satisfaction   *int16    `json:"satisfaction"`
	MajorRelated   *bool     `json:"majorRelated"`
	SkillMatch     *int16    `json:"skillMatch"`
	JobSource      string    `gorm:"size:64" json:"jobSource"`
	Feedback       string    `gorm:"type:text" json:"feedback"`
	SubmittedAt    *time.Time `json:"submittedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (EmploymentSurvey) TableName() string { return "employment_surveys" }
