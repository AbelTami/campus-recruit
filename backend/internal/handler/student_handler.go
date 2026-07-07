package handler

import (
	"employment-server/internal/model"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StudentHandler struct{ db *gorm.DB }

func NewStudentHandler(db *gorm.DB) *StudentHandler { return &StudentHandler{db} }

func (h *StudentHandler) List(c *gin.Context) {
	var page struct {
		PageIndex    int    `form:"pageIndex"`
		PageSize     int    `form:"pageSize"`
		Keyword      string `form:"keyword"`
		Name         string `form:"name"`
		EmployStatus string `form:"employStatus"`
	}
	if err := c.ShouldBindQuery(&page); err != nil {
		page.PageIndex, page.PageSize = 1, 10
	}
	if page.PageIndex <= 0 { page.PageIndex = 1 }
	if page.PageSize <= 0 || page.PageSize > 100 { page.PageSize = 10 }

	var total int64
	var students []model.Student
	q := h.db.Model(&model.Student{}).Preload("College").Preload("Major")
	kw := page.Keyword
	if kw == "" { kw = page.Name }
	if kw != "" {
		q = q.Where("name LIKE ? OR student_no LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	if page.EmployStatus != "" {
		q = q.Where("employ_status = ?", page.EmployStatus)
	}
	q.Count(&total)
	q.Offset((page.PageIndex-1)*page.PageSize).Limit(page.PageSize).Order("id DESC").Find(&students)

	response.Success(c, gin.H{"list": students, "total": total})
}

func (h *StudentHandler) Get(c *gin.Context) {
	var s model.Student
	if err := h.db.Preload("College").Preload("Major").Preload("Skills").First(&s, c.Param("id")).Error; err != nil {
		response.NotFound(c, "学生不存在")
		return
	}
	response.Success(c, s)
}

func (h *StudentHandler) Create(c *gin.Context) {
	var s model.Student
	if err := c.ShouldBindJSON(&s); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.db.Create(&s).Error; err != nil {
		response.InternalError(c, "创建失败")
		return
	}
	response.Success(c, s)
}

func (h *StudentHandler) Update(c *gin.Context) {
	var s model.Student
	if err := h.db.First(&s, c.Param("id")).Error; err != nil {
		response.NotFound(c, "学生不存在")
		return
	}
	var req model.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// ponytail: update only non-zero fields
	h.db.Model(&s).Updates(map[string]interface{}{
		"name": req.Name, "gender": req.Gender, "college_id": req.CollegeID,
		"major_id": req.MajorID, "grade": req.Grade, "education_level": req.EducationLevel,
		"graduation_date": req.GraduationDate, "phone": req.Phone, "email": req.Email,
		"expected_city": req.ExpectedCity, "expected_salary_min": req.ExpectedSalaryMin,
		"expected_salary_max": req.ExpectedSalaryMax, "employ_status": req.EmployStatus,
		"remark": req.Remark, "student_no": req.StudentNo,
	})
	response.Success(c, s)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&model.Student{}, c.Param("id")).Error; err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *StudentHandler) BatchDelete(c *gin.Context) {
	var req struct{ IDs []uint64 `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.db.Delete(&model.Student{}, req.IDs)
	response.Success(c, gin.H{"deleted": len(req.IDs)})
}
