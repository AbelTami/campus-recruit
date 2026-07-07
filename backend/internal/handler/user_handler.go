package handler

import (
	"employment-server/internal/model"
	"employment-server/internal/repository"
	"employment-server/pkg/response"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo *repository.UserRepo
	db       *gorm.DB
}

func NewUserHandler(userRepo *repository.UserRepo, db *gorm.DB) *UserHandler {
	return &UserHandler{userRepo: userRepo, db: db}
}

func (h *UserHandler) List(c *gin.Context) {
	var page struct {
		PageIndex int    `form:"pageIndex"`
		PageSize  int    `form:"pageSize"`
		ID        string `form:"id"`
		Keyword   string `form:"keyword"`
	}
	if err := c.ShouldBindQuery(&page); err != nil {
		page.PageIndex = 1
		page.PageSize = 10
	}
	if page.PageIndex <= 0 { page.PageIndex = 1 }
	if page.PageSize <= 0 { page.PageSize = 10 }

	users, total, err := h.userRepo.FindAll(page.PageIndex, page.PageSize, page.Keyword, page.ID)
	if err != nil {
		response.InternalError(c, "查询用户失败")
		return
	}

	type userItem struct {
		ID       uint64   `json:"id"`
		Username string   `json:"username"`
		Nickname string   `json:"nickname"`
		Status   int16    `json:"status"`
		Roles    []string `json:"roles"`
		College  *struct{ ID uint64 `json:"id"`; Name string `json:"label"` } `json:"department"`
	}
	list := make([]userItem, 0, len(users))
	for _, u := range users {
		roles := make([]string, len(u.Roles))
		for i, r := range u.Roles { roles[i] = r.Code }
		item := userItem{u.ID, u.Username, u.Nickname, u.Status, roles, nil}
		if u.College != nil { item.College = &struct{ ID uint64 `json:"id"`; Name string `json:"label"` }{u.College.ID, u.College.Name} }
		list = append(list, item)
	}

	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *UserHandler) Get(c *gin.Context) {
	id := c.Param("id")
	response.Success(c, gin.H{"message": "user " + id})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: real user creation with password hashing
	response.Success(c, gin.H{"message": "created"})
}

func (h *UserHandler) Update(c *gin.Context) {
	var req struct {
		Username     string   `json:"username"`
		Nickname     string   `json:"nickname"`
		Status       *int16   `json:"status"`
		Role         []string `json:"role"`
		DepartmentID *uint64  `json:"departmentId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	var uid uint64
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &uid); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	user, err := h.userRepo.FindByID(uid)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	if req.Username != "" { user.Username = req.Username }
	if req.Nickname != "" { user.Nickname = req.Nickname }
	if req.Status != nil { user.Status = *req.Status }
	if req.DepartmentID != nil { user.CollegeID = req.DepartmentID; user.College = nil }
	if req.Role != nil {
		var roles []model.Role
		for _, code := range req.Role {
			var r model.Role
			if err := h.db.Where("code = ? OR id::text = ?", code, code).First(&r).Error; err == nil {
				roles = append(roles, r)
			}
		}
		if len(roles) > 0 { h.db.Model(user).Association("Roles").Replace(roles) }
	}
	h.userRepo.Update(user)
	response.Success(c, gin.H{"message": "updated"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var uid uint64
	if _, err := fmt.Sscanf(id, "%d", &uid); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.userRepo.Delete(uid); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *UserHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	response.Success(c, gin.H{"deleted": len(req.IDs)})
}
