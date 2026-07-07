package handler

import (
	"employment-server/internal/model"
	"employment-server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CollegeHandler struct {
	db *gorm.DB
}

func NewCollegeHandler(db *gorm.DB) *CollegeHandler {
	return &CollegeHandler{db: db}
}

// 部门树节点（匹配前端 DepartmentItem）
type deptNode struct {
	ID       uint64     `json:"id"`
	Label    string     `json:"label"`
	Children []deptNode `json:"children,omitempty"`
}

func (h *CollegeHandler) List(c *gin.Context) {
	var colleges []model.College
	h.db.Order("sort_order ASC").Find(&colleges)

	children := make([]deptNode, 0, len(colleges))
	for _, c := range colleges {
		children = append(children, deptNode{ID: c.ID, Label: c.Name})
	}

	root := deptNode{
		ID:       0,
		Label:    "全校",
		Children: children,
	}

	response.Success(c, gin.H{"list": []deptNode{root}})
}

func (h *CollegeHandler) Create(c *gin.Context) {
	response.Success(c, gin.H{"message": "created"})
}

func (h *CollegeHandler) Update(c *gin.Context) {
	response.Success(c, gin.H{"message": "updated"})
}

func (h *CollegeHandler) Delete(c *gin.Context) {
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *CollegeHandler) BatchDelete(c *gin.Context) {
	response.Success(c, gin.H{"message": "deleted"})
}
