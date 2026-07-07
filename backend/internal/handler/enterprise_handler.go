package handler

import (
	"employment-server/internal/model"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EnterpriseHandler struct{ db *gorm.DB }

func NewEnterpriseHandler(db *gorm.DB) *EnterpriseHandler { return &EnterpriseHandler{db} }

func (h *EnterpriseHandler) List(c *gin.Context) {
	var page struct {
		PageIndex  int    `form:"pageIndex"`
		PageSize   int    `form:"pageSize"`
		Keyword    string `form:"keyword"`
		IndustryID string `form:"industryId"`
	}
	if err := c.ShouldBindQuery(&page); err != nil {
		page.PageIndex, page.PageSize = 1, 10
	}
	if page.PageIndex <= 0 { page.PageIndex = 1 }
	if page.PageSize <= 0 || page.PageSize > 100 { page.PageSize = 10 }

	var total int64
	var enterprises []model.Enterprise
	q := h.db.Model(&model.Enterprise{}).Preload("Industry")
	if page.Keyword != "" {
		q = q.Where("name LIKE ? OR short_name LIKE ?", "%"+page.Keyword+"%", "%"+page.Keyword+"%")
	}
	if page.IndustryID != "" {
		q = q.Where("industry_id = ?", page.IndustryID)
	}
	q.Count(&total)
	q.Offset((page.PageIndex-1)*page.PageSize).Limit(page.PageSize).Order("id DESC").Find(&enterprises)

	response.Success(c, gin.H{"list": enterprises, "total": total})
}

func (h *EnterpriseHandler) Get(c *gin.Context) {
	var e model.Enterprise
	if err := h.db.Preload("Industry").First(&e, c.Param("id")).Error; err != nil {
		response.NotFound(c, "企业不存在")
		return
	}
	response.Success(c, e)
}

func (h *EnterpriseHandler) Create(c *gin.Context) {
	var e model.Enterprise
	if err := c.ShouldBindJSON(&e); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.db.Create(&e).Error; err != nil {
		response.InternalError(c, "创建失败")
		return
	}
	response.Success(c, e)
}

func (h *EnterpriseHandler) Update(c *gin.Context) {
	var e model.Enterprise
	if err := h.db.First(&e, c.Param("id")).Error; err != nil {
		response.NotFound(c, "企业不存在")
		return
	}
	var req model.Enterprise
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.db.Model(&e).Updates(map[string]interface{}{
		"name": req.Name, "short_name": req.ShortName, "industry_id": req.IndustryID,
		"scale": req.Scale, "nature": req.Nature, "city": req.City,
		"address": req.Address, "website": req.Website, "description": req.Description,
		"contact_name": req.ContactName, "contact_phone": req.ContactPhone,
		"contact_email": req.ContactEmail, "status": req.Status,
	})
	response.Success(c, e)
}

func (h *EnterpriseHandler) Delete(c *gin.Context) {
	if err := h.db.Delete(&model.Enterprise{}, c.Param("id")).Error; err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *EnterpriseHandler) BatchDelete(c *gin.Context) {
	var req struct{ IDs []uint64 `json:"ids"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.db.Delete(&model.Enterprise{}, req.IDs)
	response.Success(c, gin.H{"deleted": len(req.IDs)})
}
