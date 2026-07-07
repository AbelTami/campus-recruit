package repository

import (
	"employment-server/internal/model"

	"gorm.io/gorm"
)

type MenuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{db: db}
}

func (r *MenuRepo) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Order("sort_order ASC").Find(&menus).Error
	return menus, err
}

func (r *MenuRepo) Create(m *model.Menu) error { return r.db.Create(m).Error }

func (r *MenuRepo) Update(id uint64, updates map[string]any) error { return r.db.Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error }

func (r *MenuRepo) Delete(id uint64) error { return r.db.Delete(&model.Menu{}, id).Error }

// BuildMenuTree converts a flat menu list into a nested tree.
func BuildMenuTree(menus []model.Menu) []model.Menu {
	menuMap := make(map[uint64]*model.Menu)
	var roots []*model.Menu

	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	for i := range menus {
		if menus[i].ParentID != nil {
			if parent, ok := menuMap[*menus[i].ParentID]; ok {
				parent.Children = append(parent.Children, menus[i])
			}
		} else {
			roots = append(roots, menuMap[menus[i].ID])
		}
	}

	result := make([]model.Menu, len(roots))
	for i, r := range roots {
		result[i] = *r
	}
	return result
}
