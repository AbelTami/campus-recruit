package handler

import (
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct{ db *gorm.DB }

func NewDashboardHandler(db *gorm.DB) *DashboardHandler { return &DashboardHandler{db} }

func (h *DashboardHandler) Overview(c *gin.Context) {
	var totalStudents, employedCount int64
	var employmentRate float64
	var avgSalary *float64

	h.db.Table("students").Where("deleted_at IS NULL").Count(&totalStudents)
	h.db.Table("students").Where("deleted_at IS NULL AND employ_status = ?", "employed").Count(&employedCount)
	if totalStudents > 0 { employmentRate = float64(employedCount) / float64(totalStudents) * 100 }
	h.db.Table("employment_records").Select("AVG(monthly_salary)").Scan(&avgSalary)

	var totalEnterprises, totalPositions, totalApplications int64
	h.db.Table("enterprises").Where("deleted_at IS NULL AND status = 1").Count(&totalEnterprises)
	h.db.Table("positions").Where("deleted_at IS NULL AND status = 1").Count(&totalPositions)
	h.db.Table("applications").Count(&totalApplications)

	response.Success(c, gin.H{
		"totalStudents": totalStudents, "employedCount": employedCount,
		"employmentRate": employmentRate, "avgSalary": avgSalary,
		"totalEnterprises": totalEnterprises, "totalPositions": totalPositions,
		"totalApplications": totalApplications,
	})
}

func (h *DashboardHandler) IndustryDist(c *gin.Context) {
	var rows []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	// ponytail: left join to show all industries even with 0 positions
	h.db.Raw(`SELECT COALESCE(i.name, '未知') as name, COUNT(p.id) as value
		FROM industries i LEFT JOIN positions p ON i.id = p.industry_id AND p.deleted_at IS NULL
		GROUP BY i.name ORDER BY value DESC`).Scan(&rows)
	response.Success(c, rows)
}

func (h *DashboardHandler) EmploymentTrend(c *gin.Context) {
	var rows []struct {
		Year  string  `json:"year"`
		Rate  float64 `json:"rate"`
	}
	h.db.Raw(`SELECT graduation_year::text as year,
		ROUND(COUNT(*) FILTER (WHERE employ_status = 'employed') * 100.0 / NULLIF(COUNT(*), 0), 1) as rate
		FROM students WHERE deleted_at IS NULL AND graduation_year IS NOT NULL
		GROUP BY graduation_year ORDER BY year`).Scan(&rows)
	response.Success(c, rows)
}

func (h *DashboardHandler) SalaryAnalysis(c *gin.Context) {
	var rows []struct {
		Name string  `json:"name"`
		Avg  float64 `json:"avg"`
	}
	h.db.Raw(`SELECT COALESCE(i.name,'未知') as name, ROUND(AVG(p.salary_min+p.salary_max)/2000,1) as avg
		FROM positions p LEFT JOIN industries i ON i.id=p.industry_id
		WHERE p.deleted_at IS NULL AND p.salary_min>0 GROUP BY i.name ORDER BY avg DESC`).Scan(&rows)
	response.Success(c, rows)
}

func (h *DashboardHandler) CollegeEmployment(c *gin.Context) {
	var rows []struct {
		Name  string  `json:"name"`
		Rate  float64 `json:"rate"`
		Total int64   `json:"total"`
	}
	h.db.Raw(`SELECT c.name, COUNT(s.id) as total,
		ROUND(COUNT(*) FILTER (WHERE s.employ_status='employed')*100.0/NULLIF(COUNT(*),0),1) as rate
		FROM colleges c LEFT JOIN students s ON s.college_id=c.id AND s.deleted_at IS NULL
		GROUP BY c.name HAVING COUNT(s.id)>0 ORDER BY rate DESC`).Scan(&rows)
	response.Success(c, rows)
}
