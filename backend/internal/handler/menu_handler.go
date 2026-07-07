package handler

import (
	"employment-server/internal/model"
	"employment-server/internal/repository"
	dto "employment-server/internal/dto/response"
	"employment-server/pkg/response"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repository.MenuRepo
}

func NewMenuHandler(menuRepo *repository.MenuRepo) *MenuHandler {
	return &MenuHandler{menuRepo: menuRepo}
}

func (h *MenuHandler) GetMenus(c *gin.Context) {
	menus, err := h.menuRepo.FindAll()
	if err != nil {
		response.InternalError(c, "获取菜单失败")
		return
	}

	tree := repository.BuildMenuTree(menus)
	result := make([]dto.MenuTreeResponse, 0, len(tree))
	for _, m := range tree {
		result = append(result, toRouteDTO(m))
	}

	response.Success(c, result)
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req struct {
		ParentID  *uint64 `json:"parentId"`
		Name      string  `json:"meta.title"`
		Title     string  `json:"title"`
		Path      string  `json:"path"`
		Component string  `json:"component"`
		Icon      string  `json:"meta.icon"`
		Type      int16   `json:"type"`
		Status    int16   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	m := model.Menu{
		ParentID: req.ParentID, Name: req.Name, Path: req.Path,
		Component: req.Component, Icon: req.Icon, MenuType: req.Type + 1, Status: req.Status,
	}
	h.menuRepo.Create(&m)
	response.Success(c, m)
}

func (h *MenuHandler) Update(c *gin.Context) {
	var req struct {
		ParentID  *uint64 `json:"parentId"`
		Name      *string `json:"meta.title"`
		Title     *string `json:"title"`
		Path      *string `json:"path"`
		Component *string `json:"component"`
		Icon      *string `json:"meta.icon"`
		Type      *int16  `json:"type"`
		Status    *int16  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
	updates := map[string]any{}
	if req.Name != nil { updates["name"] = *req.Name }
	if req.ParentID != nil { updates["parent_id"] = *req.ParentID }
	if req.Path != nil { updates["path"] = *req.Path }
	if req.Component != nil { updates["component"] = *req.Component }
	if req.Icon != nil { updates["icon"] = *req.Icon }
	if req.Type != nil { updates["menu_type"] = *req.Type + 1 }
	if req.Status != nil { updates["status"] = *req.Status }
	var uid uint64
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &uid); err != nil { response.BadRequest(c, "参数错误"); return }
	h.menuRepo.Update(uid, updates)
	response.Success(c, gin.H{"message": "updated"})
}

func (h *MenuHandler) Delete(c *gin.Context) {
	var uid uint64
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &uid); err != nil { response.BadRequest(c, "参数错误"); return }
	h.menuRepo.Delete(uid)
	response.Success(c, nil)
}

func toRouteDTO(m model.Menu) dto.MenuTreeResponse {
	var perm []string
	if m.Permission != "" {
		perm = strings.Split(m.Permission, ",")
	}

	// DB menuType: 1=目录, 2=菜单 → 前端 type: 0=目录, 1=菜单
	menuType := m.MenuType - 1

	d := dto.MenuTreeResponse{
		ID:        m.ID,
		ParentID:  m.ParentID,
		Type:      menuType,
		Path:      m.Path,
		Name:      m.Name,
		Title:     m.Name,
		Component: m.Component,
		Status:    m.Status,
		Meta: dto.RouteMeta{
			Title:      m.Name,
			Icon:       m.Icon,
			Hidden:     !m.Visible,
			AlwaysShow: len(m.Children) > 0,
			Permission: perm,
		},
		Children: make([]dto.MenuTreeResponse, 0, len(m.Children)),
	}

	for _, child := range m.Children {
		d.Children = append(d.Children, toRouteDTO(child))
	}
	return d
}
