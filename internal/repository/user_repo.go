package repository

import (
	"employment-server/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").Preload("College").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindAll(pageIndex, pageSize int, keyword, deptID string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{}).Preload("Roles").Preload("College")
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if deptID != "" && deptID != "0" {
		query = query.Where("college_id = ?", deptID)
	}
	query.Count(&total)
	err := query.Offset((pageIndex-1)*pageSize).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepo) Create(user *model.User) error { return r.db.Create(user).Error }

func (r *UserRepo) Update(user *model.User) error { return r.db.Save(user).Error }

func (r *UserRepo) Delete(id uint64) error { return r.db.Delete(&model.User{}, id).Error }

func (r *UserRepo) UpdateLoginInfo(id uint64, ip string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at":  time.Now(),
		"last_login_ip":  ip,
		"login_attempts": 0,
		"locked_until":   nil,
	}).Error
}

func (r *UserRepo) IncrementLoginAttempts(id uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("login_attempts", gorm.Expr("login_attempts + 1")).Error
}

func (r *UserRepo) LockAccount(id uint64, until time.Time) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      2,
		"locked_until": until,
	}).Error
}
