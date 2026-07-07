package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"employment-server/internal/config"
	"employment-server/internal/model"
	"employment-server/internal/router"
	hashPkg "employment-server/pkg/hash"
	"employment-server/pkg/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化日志
	zapLogger, err := logger.New(cfg.Log.Level, cfg.Log.Output)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer zapLogger.Sync()

	// 3. 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn), // ponytail: only warn+error, no query spam
	})
	if err != nil {
		zapLogger.Fatal("数据库连接失败: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		zapLogger.Fatal("获取数据库连接失败: " + err.Error())
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// 4. 自动迁移（开发环境，生产环境用 golang-migrate）
	if cfg.Server.Mode == "debug" {
		if err := db.AutoMigrate(
			&model.User{}, &model.Role{}, &model.Menu{},
			&model.UserRole{}, &model.RoleMenu{},
			&model.College{}, &model.Major{},
			&model.Student{}, &model.Skill{}, &model.StudentSkill{},
			&model.Industry{}, &model.Enterprise{},
			&model.Position{}, &model.PositionSkill{},
			&model.Application{}, &model.ApplicationLog{},
			&model.EmploymentRecord{},
			&model.OperationLog{}, &model.DictData{},
			&model.AnalysisReport{}, &model.MarketDataSnapshot{},
			&model.EmploymentSurvey{},
		); err != nil {
			zapLogger.Fatal("自动迁移失败: " + err.Error())
		}
		zapLogger.Info("数据库表自动迁移完成")

		// 种子数据
		seed(db, cfg)
	}

	// 5. 连接 Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		zapLogger.Warn("Redis 连接失败（将无法使用 Token 刷新和缓存）: " + err.Error())
	} else {
		zapLogger.Info("Redis 连接成功")
	}

	// 6. 设置 Gin 模式
	// gin.SetMode(cfg.Server.Mode)

	// 7. 注册路由
	r := router.Setup(db, redisClient, cfg, zapLogger)

	// 8. 启动服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		zapLogger.Info(fmt.Sprintf("🚀 服务启动: http://0.0.0.0:%s", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("服务启动失败: " + err.Error())
		}
	}()

	// 9. 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("服务关闭异常: " + err.Error())
	}
	zapLogger.Info("服务已关闭")
}

// seed 初始化种子数据
func seed(db *gorm.DB, cfg *config.Config) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return // 已有数据，跳过
	}

	// 创建默认角色
	adminRole := model.Role{Name: "超级管理员", Code: "super_admin", Status: 1}
	teacherRole := model.Role{Name: "就业指导老师", Code: "career_teacher", Status: 1}
	leaderRole := model.Role{Name: "校领导", Code: "school_leader", Status: 1}
	studentRole := model.Role{Name: "学生", Code: "student", Status: 1}
	db.Create(&adminRole)
	db.Create(&teacherRole)
	db.Create(&leaderRole)
	db.Create(&studentRole)

	// 创建默认管理员（密码: admin123，运行时用 bcrypt+pepper 生成哈希）
	passwordHash, err := hashPkg.HashPassword("admin123", cfg.Encrypt.Pepper)
	if err != nil {
		log.Fatalf("生成密码哈希失败: %v", err)
	}
	adminUser := model.User{
		Username:     "admin",
		PasswordHash: passwordHash,
		Nickname:     "系统管理员",
		Status:       1,
	}
	db.Create(&adminUser)
	db.Model(&adminUser).Association("Roles").Append(&adminRole)

	// 种子菜单（前端动态路由需要）
	seedMenus(db)

	log.Println("✅ 默认管理员: admin / admin123")
	log.Println("⚠️  生产环境请立即修改默认密码！")
}

func seedMenus(db *gorm.DB) {
	type menuSeed struct {
		ID         uint64
		ParentID   *uint64
		Name       string
		Path       string
		Component  string
		Icon       string
		SortOrder  int
		Permission string
		MenuType   int16
	}
	ptr := func(v uint64) *uint64 { return &v }

	menus := []menuSeed{
		{ID: 1, Name: "就业管理", Path: "/employment", Component: "#", Icon: "ri:bar-chart-box-line", SortOrder: 1, MenuType: 1},
		{ID: 2, ParentID: ptr(1), Name: "数据大盘", Path: "dashboard", Component: "views/Dashboard/Analysis", Icon: "ri:dashboard-line", SortOrder: 1, MenuType: 2},
		{ID: 3, ParentID: ptr(1), Name: "学生管理", Path: "students", Component: "##", Icon: "ri:user-line", SortOrder: 2, MenuType: 1},
		{ID: 4, ParentID: ptr(3), Name: "学生列表", Path: "list", Component: "views/Error/404", SortOrder: 1, MenuType: 2, Permission: "student:list"},
		{ID: 5, ParentID: ptr(1), Name: "企业管理", Path: "enterprises", Component: "views/Error/404", Icon: "ri:building-line", SortOrder: 3, MenuType: 2, Permission: "enterprise:list"},
		{ID: 6, ParentID: ptr(1), Name: "就业分析", Path: "analysis", Component: "##", Icon: "ri:line-chart-line", SortOrder: 6, MenuType: 1},
		{ID: 7, ParentID: ptr(6), Name: "就业率分析", Path: "employment-rate", Component: "views/Error/404", SortOrder: 3, MenuType: 2, Permission: "analysis:employ-rate"},
		{ID: 8, Name: "系统管理", Path: "/system", Component: "#", Icon: "ri:settings-line", SortOrder: 99, MenuType: 1},
		{ID: 9, ParentID: ptr(8), Name: "用户管理", Path: "users", Component: "views/Authorization/User/User", SortOrder: 1, MenuType: 2, Permission: "system:user"},
	}

	for _, m := range menus {
		db.Create(&model.Menu{
			ID:         m.ID,
			ParentID:   m.ParentID,
			Name:       m.Name,
			Path:       m.Path,
			Component:  m.Component,
			Icon:       m.Icon,
			SortOrder:  m.SortOrder,
			Permission: m.Permission,
			MenuType:   m.MenuType,
		})
	}

	// 为 admin 角色分配所有菜单
	var menusAll []model.Menu
	db.Find(&menusAll)
	var adminRole model.Role
	db.Where("code = ?", "super_admin").First(&adminRole)
	for _, m := range menusAll {
		db.Create(&model.RoleMenu{RoleID: adminRole.ID, MenuID: m.ID})
	}
}
