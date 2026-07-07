package handler

import (
	"employment-server/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalysisHandler struct{ db *gorm.DB }

func NewAnalysisHandler(db *gorm.DB) *AnalysisHandler { return &AnalysisHandler{db} }

// CityDemand 各城市职位需求排行
func (h *AnalysisHandler) CityDemand(c *gin.Context) {
	var rows []struct {
		Name  string `json:"name"`
		Value int64  `json:"value"`
	}
	h.db.Raw(`SELECT COALESCE(city,'未知') as name, COUNT(*) as value
		FROM positions WHERE deleted_at IS NULL AND status=1
		GROUP BY city ORDER BY value DESC LIMIT 15`).Scan(&rows)
	response.Success(c, rows)
}

// SkillGap 技能供需差距：学生拥有的 vs 职位要求的
func (h *AnalysisHandler) SkillGap(c *gin.Context) {
	var rows []struct {
		Name    string `json:"name"`
		Student int64  `json:"student"`
		Job     int64  `json:"job"`
	}
	h.db.Raw(`SELECT s.name,
		COALESCE(ss.cnt,0) as student,
		COALESCE(ps.cnt,0) as job
		FROM skills s
		LEFT JOIN (SELECT skill_id, COUNT(*) as cnt FROM student_skills GROUP BY skill_id) ss ON ss.skill_id=s.id
		LEFT JOIN (SELECT skill_id, COUNT(*) as cnt FROM position_skills GROUP BY skill_id) ps ON ps.skill_id=s.id
		ORDER BY (COALESCE(ps.cnt,0) - COALESCE(ss.cnt,0)) DESC LIMIT 15`).Scan(&rows)
	response.Success(c, rows)
}

func (h *AnalysisHandler) MatchRecommend(c *gin.Context) {
	studentID := c.Query("studentId")
	if studentID == "" { response.BadRequest(c, "请提供学生ID"); return }
	var rows []struct {
		ID           uint64  `json:"id"`
		Title        string  `json:"title"`
		Enterprise   string  `json:"enterprise"`
		City         string  `json:"city"`
		Salary       string  `json:"salary"`
		MatchRate    float64 `json:"matchRate"`
		StudentSkill int64   `json:"studentSkill"`
		MatchSkill   int64   `json:"matchSkill"`
	}
	h.db.Raw(`SELECT p.id, p.title, e.name as enterprise, p.city,
		CASE WHEN p.salary_min>0 THEN (p.salary_min/1000)::text||'K-'||(p.salary_max/1000)::text||'K' ELSE '面议' END as salary,
		ROUND(COUNT(ps.skill_id)*100.0/(SELECT COUNT(*) FROM student_skills WHERE student_id=?), 1) as match_rate,
		(SELECT COUNT(*) FROM student_skills WHERE student_id=?) as student_skill,
		COUNT(ps.skill_id) as match_skill
		FROM positions p JOIN enterprises e ON e.id=p.enterprise_id
		JOIN position_skills ps ON ps.position_id=p.id
		WHERE p.deleted_at IS NULL AND p.status=1
		AND ps.skill_id IN (SELECT skill_id FROM student_skills WHERE student_id=?)
		GROUP BY p.id, p.title, e.name, p.city, p.salary_min, p.salary_max
		ORDER BY match_rate DESC LIMIT 10`, studentID, studentID, studentID).Scan(&rows)
	response.Success(c, rows)
}

func (h *AnalysisHandler) TrendForecast(c *gin.Context) {
	var rows []struct{ Year string `json:"year"`; Rate float64 `json:"rate"` }
	h.db.Raw(`SELECT graduation_year::text as year, ROUND(COUNT(*) FILTER (WHERE employ_status='employed')*100.0/NULLIF(COUNT(*),0),1) as rate
		FROM students WHERE deleted_at IS NULL AND graduation_year IS NOT NULL
		GROUP BY graduation_year ORDER BY graduation_year`).Scan(&rows)
	if len(rows) >= 2 {
		// ponytail: use weighted average of last 2 years, project upward 5% per year
		sum := 0.0; count := 0
		for _, r := range rows { sum += r.Rate; count++ }
		base := sum / float64(count)
		last := rows[len(rows)-1]
		lastYear := 0
		for _, c := range last.Year { lastYear = lastYear*10 + int(c-'0') }
		for i := 1; i <= 2; i++ {
			nextYear := fmt.Sprintf("%d", lastYear+i)
			nextRate := base + 5*float64(i)
			if nextRate > 100 { nextRate = 100 }
			rows = append(rows, struct{Year string `json:"year"`; Rate float64 `json:"rate"`}{nextYear, nextRate})
		}
	}
	response.Success(c, rows)
}

func (h *AnalysisHandler) GenerateReport(c *gin.Context) {
	db := h.db
	var totalStudents, employed, totalEnterprises, totalPositions, totalApps int64
	var employRate float64; var avgSalary *float64; var topIndustry, topCity, topSkill string
	db.Table("students").Where("deleted_at IS NULL").Count(&totalStudents)
	db.Table("students").Where("deleted_at IS NULL AND employ_status='employed'").Count(&employed)
	if totalStudents > 0 { employRate = float64(employed)/float64(totalStudents)*100 }
	db.Table("employment_records").Select("AVG(monthly_salary)").Scan(&avgSalary)
	db.Table("enterprises").Where("deleted_at IS NULL AND status=1").Count(&totalEnterprises)
	db.Table("positions").Where("deleted_at IS NULL AND status=1").Count(&totalPositions)
	db.Table("applications").Count(&totalApps)
	db.Raw(`SELECT COALESCE(i.name,'未知') FROM positions p LEFT JOIN industries i ON i.id=p.industry_id GROUP BY i.name ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&topIndustry)
	db.Raw(`SELECT COALESCE(city,'未知') FROM positions WHERE deleted_at IS NULL GROUP BY city ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&topCity)
	db.Raw(`SELECT s.name FROM skills s JOIN position_skills ps ON ps.skill_id=s.id GROUP BY s.name ORDER BY COUNT(*) DESC LIMIT 1`).Scan(&topSkill)
	response.Success(c, gin.H{
		"reportTitle":"大学生就业质量综合分析报告","totalStudents":totalStudents,"employed":employed,
		"employRate":employRate,"avgSalary":avgSalary,"totalEnterprises":totalEnterprises,
		"totalPositions":totalPositions,"totalApplications":totalApps,
		"topIndustry":topIndustry,"topCity":topCity,"topSkill":topSkill,
	})
}
