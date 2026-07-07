package handler

import (
	"employment-server/internal/model"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PositionHandler struct{ db *gorm.DB }

func NewPositionHandler(db *gorm.DB) *PositionHandler { return &PositionHandler{db} }

func (h *PositionHandler) List(c *gin.Context) {
	var page struct {
		PageIndex  int    `form:"pageIndex"`; PageSize int `form:"pageSize"`
		Keyword              string `form:"keyword"`; City string `form:"city"`
		IndustryID           string `form:"industryId"`; Status string `form:"status"`
		EducationRequirement string `form:"educationRequirement"`
	}
	if err := c.ShouldBindQuery(&page); err != nil { page.PageIndex, page.PageSize = 1, 10 }
	if page.PageIndex <= 0 { page.PageIndex = 1 }
	if page.PageSize <= 0 || page.PageSize > 100 { page.PageSize = 10 }

	var total int64
	var positions []model.Position
	q := h.db.Model(&model.Position{}).Preload("Enterprise").Preload("Industry")
	if page.Keyword != "" { q = q.Where("title LIKE ?", "%"+page.Keyword+"%") }
	if page.City != "" { q = q.Where("city = ?", page.City) }
	if page.EducationRequirement != "" { q = q.Where("education_requirement = ?", page.EducationRequirement) }
	if page.IndustryID != "" { q = q.Where("industry_id = ?", page.IndustryID) }
	if page.Status != "" { q = q.Where("status = ?", page.Status) }
	q.Count(&total)
	q.Offset((page.PageIndex-1)*page.PageSize).Limit(page.PageSize).Order("id DESC").Find(&positions)
	response.Success(c, gin.H{"list": positions, "total": total})
}

func (h *PositionHandler) Get(c *gin.Context) {
	var p model.Position
	if err := h.db.Preload("Enterprise").Preload("Industry").First(&p, c.Param("id")).Error; err != nil {
		response.NotFound(c, "职位不存在")
		return
	}
	response.Success(c, p)
}

func (h *PositionHandler) Create(c *gin.Context) {
	var p model.Position
	if err := c.ShouldBindJSON(&p); err != nil { response.BadRequest(c, "参数错误"); return }
	if err := h.db.Create(&p).Error; err != nil { response.InternalError(c, "创建失败"); return }
	response.Success(c, p)
}

func (h *PositionHandler) Update(c *gin.Context) {
	var p model.Position
	if err := h.db.First(&p, c.Param("id")).Error; err != nil { response.NotFound(c, "职位不存在"); return }
	var req model.Position
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	h.db.Model(&p).Updates(map[string]interface{}{
		"title": req.Title, "enterprise_id": req.EnterpriseID, "industry_id": req.IndustryID,
		"city": req.City, "education_requirement": req.EducationRequirement,
		"experience_requirement": req.ExperienceRequirement, "salary_min": req.SalaryMin,
		"salary_max": req.SalaryMax, "headcount": req.Headcount, "description": req.Description,
		"requirement": req.Requirement, "welfare": req.Welfare, "status": req.Status,
	})
	response.Success(c, p)
}

func (h *PositionHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&model.Position{}, c.Param("id")).Error; err != nil { response.InternalError(c, "删除失败"); return }
	response.Success(c, nil)
}

func (h *PositionHandler) BatchDelete(c *gin.Context) {
	var req struct{ IDs []uint64 `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	h.db.Delete(&model.Position{}, req.IDs)
	response.Success(c, gin.H{"deleted": len(req.IDs)})
}
