package model

import (
	"time"

	"gorm.io/gorm"
)

type College struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Code      string         `gorm:"uniqueIndex;size:32" json:"code"`
	DeanName  string         `gorm:"size:64" json:"deanName"`
	Contact   string         `gorm:"size:128" json:"contact"`
	SortOrder int            `gorm:"not null;default:0" json:"sortOrder"`
	Status    int16          `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Majors    []Major        `gorm:"foreignKey:CollegeID" json:"majors,omitempty"`
}

func (College) TableName() string { return "colleges" }

type Major struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CollegeID      uint64         `gorm:"index;not null" json:"collegeId"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	Code           string         `gorm:"uniqueIndex;size:32" json:"code"`
	Category       string         `gorm:"size:32" json:"category"`
	EducationLevel string         `gorm:"size:16" json:"educationLevel"`
	DurationYears  int16          `gorm:"default:4" json:"durationYears"`
	Status         int16          `gorm:"default:1" json:"status"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	College        *College       `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
}

func (Major) TableName() string { return "majors" }

type Student struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            *uint64        `gorm:"index" json:"userId"`
	StudentNo         string         `gorm:"uniqueIndex;size:32;not null" json:"studentNo"`
	Name              string         `gorm:"size:64;not null" json:"name"`
	Gender            *int16         `json:"gender"`
	BirthDate         *time.Time     `json:"birthDate"`
	CollegeID         *uint64        `gorm:"index" json:"collegeId"`
	MajorID           *uint64        `gorm:"index" json:"majorId"`
	Grade             string         `gorm:"size:4" json:"grade"`
	EducationLevel    string         `gorm:"size:16" json:"educationLevel"`
	GraduationDate    string         `gorm:"size:16" json:"graduationDate"`
	PoliticalStatus   string         `gorm:"size:16" json:"politicalStatus"`
	HometownCity      string         `gorm:"size:64" json:"hometownCity"`
	Phone             []byte         `json:"-"`
	PhonePlain        string         `gorm:"-" json:"phone,omitempty"`
	Email             []byte         `json:"-"`
	EmailPlain        string         `gorm:"-" json:"email,omitempty"`
	Wechat            string         `gorm:"size:64" json:"wechat"`
	QQ                string         `gorm:"size:20" json:"qq"`
	IDCardHash        string         `gorm:"size:64" json:"-"`
	ResumeURL         string         `gorm:"size:512" json:"resumeUrl"`
	ExpectedCity      string         `gorm:"size:64" json:"expectedCity"`
	ExpectedSalaryMin *int           `json:"expectedSalaryMin"`
	ExpectedSalaryMax *int           `json:"expectedSalaryMax"`
	ExpectedIndustry  string         `gorm:"size:64" json:"expectedIndustry"`
	ExpectedPositions string         `gorm:"type:text" json:"expectedPositions"` // JSON array text
	EmployStatus      string         `gorm:"size:16;default:unemployed" json:"employStatus"`
	EmployCompany     string         `gorm:"size:128" json:"employCompany"`
	EmployPosition    string         `gorm:"size:128" json:"employPosition"`
	EmploySalary      *int           `json:"employSalary"`
	EmployCity        string         `gorm:"size:64" json:"employCity"`
	EmployDate        *time.Time     `json:"employDate"`
	DataSource        string         `gorm:"size:16;default:manual" json:"dataSource"`
	Verified          bool           `gorm:"default:false" json:"verified"`
	Remark            string         `gorm:"type:text" json:"remark"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	College           *College       `gorm:"foreignKey:CollegeID" json:"college,omitempty"`
	Major             *Major         `gorm:"foreignKey:MajorID" json:"major,omitempty"`
	Skills            []Skill        `gorm:"many2many:student_skills" json:"skills,omitempty"`
}

func (Student) TableName() string { return "students" }

type Skill struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Category  string    `gorm:"size:32" json:"category"`
	SortOrder int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
}

func (Skill) TableName() string { return "skills" }

type StudentSkill struct {
	StudentID   uint64 `gorm:"primaryKey" json:"studentId"`
	SkillID     uint64 `gorm:"primaryKey" json:"skillId"`
	Proficiency int16  `gorm:"default:1;check:proficiency between 1 and 5" json:"proficiency"`
}

func (StudentSkill) TableName() string { return "student_skills" }
