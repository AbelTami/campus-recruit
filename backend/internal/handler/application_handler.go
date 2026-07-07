package handler

import (
	"employment-server/internal/model"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationHandler struct{ db *gorm.DB }

func NewApplicationHandler(db *gorm.DB) *ApplicationHandler { return &ApplicationHandler{db} }

func (h *ApplicationHandler) List(c *gin.Context) {
	var page struct {
		PageIndex int    `form:"pageIndex"`; PageSize int `form:"pageSize"`
		Status    string `form:"status"`; Keyword string `form:"keyword"`
	}
	if err := c.ShouldBindQuery(&page); err != nil { page.PageIndex, page.PageSize = 1, 10 }
	if page.PageIndex <= 0 { page.PageIndex = 1 }
	if page.PageSize <= 0 || page.PageSize > 100 { page.PageSize = 10 }

	var total int64
	var apps []model.Application
	q := h.db.Model(&model.Application{}).Preload("Student").Preload("Position").Preload("Enterprise")
	if page.Status != "" { q = q.Where("status = ?", page.Status) }
	if page.Keyword != "" {
		kw := "%" + page.Keyword + "%"
		q = q.Where(
			"EXISTS (SELECT 1 FROM students s WHERE s.id = applications.student_id AND s.name LIKE ?) OR "+
			"EXISTS (SELECT 1 FROM positions p WHERE p.id = applications.position_id AND p.title LIKE ?) OR "+
			"EXISTS (SELECT 1 FROM enterprises e WHERE e.id = applications.enterprise_id AND e.name LIKE ?)",
			kw, kw, kw,
		)
	}
	q.Count(&total)
	q.Offset((page.PageIndex-1)*page.PageSize).Limit(page.PageSize).Order("created_at DESC").Find(&apps)
	response.Success(c, gin.H{"list": apps, "total": total})
}

func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	var app model.Application
	if err := h.db.First(&app, c.Param("id")).Error; err != nil { response.NotFound(c, "投递不存在"); return }
	oldStatus := app.Status
	h.db.Model(&app).Update("status", req.Status)
	if req.Note != "" {
		h.db.Create(&model.ApplicationLog{
			ApplicationID: app.ID, FromStatus: oldStatus, ToStatus: req.Status,
			OperatorName: c.GetString("username"), Note: req.Note,
		})
	}
	response.Success(c, app)
}

func (h *ApplicationHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&model.Application{}, c.Param("id")).Error; err != nil { response.InternalError(c, "删除失败"); return }
	response.Success(c, nil)
}
