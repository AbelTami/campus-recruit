package router

import (
	"employment-server/internal/config"
	"employment-server/internal/handler"
	"employment-server/internal/middleware"
	"employment-server/internal/model"
	"employment-server/internal/repository"
	"employment-server/internal/service"
	crypto "employment-server/pkg/crypto"
	hash "employment-server/pkg/hash"
	"employment-server/pkg/jwt"
	"employment-server/pkg/response"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, logger *zap.Logger) *gin.Engine {
	jwtService := jwt.NewService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	userRepo := repository.NewUserRepo(db); menuRepo := repository.NewMenuRepo(db)
	authService := service.NewAuthService(userRepo, jwtService, cfg, redisClient)

	authH := handler.NewAuthHandler(authService); menuH := handler.NewMenuHandler(menuRepo)
	userH := handler.NewUserHandler(userRepo, db); collegeH := handler.NewCollegeHandler(db)
	studentH := handler.NewStudentHandler(db); enterpriseH := handler.NewEnterpriseHandler(db)
	positionH := handler.NewPositionHandler(db); dashboardH := handler.NewDashboardHandler(db)
	applicationH := handler.NewApplicationHandler(db)
		analysisH := handler.NewAnalysisHandler(db)

	
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Recovery(logger), middleware.Logger(logger), middleware.CORS(), middleware.GlobalRateLimit(200, 400))
	r.GET("/api/v1/health", func(c *gin.Context) { response.Success(c, gin.H{"status": "ok"}) })
	r.GET("/api/v1/portal/stats", func(c *gin.Context) {
		var totalStudents, employed, totalPositions, totalEnterprises int64
		var employRate float64
		db.Table("students").Where("deleted_at IS NULL").Count(&totalStudents)
		db.Table("students").Where("deleted_at IS NULL AND employ_status = ?", "employed").Count(&employed)
		if totalStudents > 0 { employRate = float64(employed) / float64(totalStudents) * 100 }
		db.Table("positions").Where("deleted_at IS NULL AND status = 1").Count(&totalPositions)
		db.Table("enterprises").Where("deleted_at IS NULL AND status = 1").Count(&totalEnterprises)
		response.Success(c, gin.H{
			"totalStudents": totalStudents, "totalPositions": totalPositions,
			"totalEnterprises": totalEnterprises, "employmentRate": employRate,
		})
	})

	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authH.Login); public.POST("/auth/refresh", authH.RefreshToken)
		public.GET("/portal/positions", positionH.List)
		public.GET("/portal/positions/:id", positionH.Get)
		public.GET("/portal/enterprises", enterpriseH.List)
		public.GET("/portal/enterprises/:id", enterpriseH.Get)
		public.POST("/auth/register", func(c *gin.Context) {
			var req struct {
				Username string `json:"username"`; Password string `json:"password"`
				Name     string `json:"name"`; StudentNo string `json:"studentNo"`
			}
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			// pony: minimal student registration
			passwordHash, _ := hash.HashPassword(req.Password, cfg.Encrypt.Pepper)
			user := model.User{Username: req.Username, PasswordHash: passwordHash, Nickname: req.Name, Status: 1}
			if err := db.Create(&user).Error; err != nil { response.Error(c, response.CodeConflict, "用户名已存在"); return }
			if err := db.Create(&model.Student{UserID: &user.ID, Name: req.Name, StudentNo: req.StudentNo}).Error; err != nil {
				response.Error(c, response.CodeInternal, "创建学生档案失败: "+err.Error()); return
			}
			response.Success(c, gin.H{"message": "注册成功"})
		})
	}

	secure := r.Group("/api/v1")
	secure.Use(middleware.JWTAuth(jwtService))
	{
		secure.POST("/auth/logout", authH.Logout); secure.GET("/auth/userinfo", authH.GetUserInfo)
		secure.GET("/admin/menus", menuH.GetMenus); secure.POST("/admin/menus", menuH.Create)
		secure.PUT("/admin/menus/:id", menuH.Update); secure.DELETE("/admin/menus/:id", menuH.Delete)
		secure.GET("/admin/industries", func(c *gin.Context) {
			var industries []model.Industry
			db.Order("sort_order").Find(&industries)
			response.Success(c, industries)
		})
		secure.GET("/portal/student/applications", func(c *gin.Context) {
			userID := c.GetUint64("user_id")
			var student model.Student
			if err := db.Where("user_id = ?", userID).First(&student).Error; err != nil { response.NotFound(c, "未找到学生信息"); return }
			var apps []model.Application
			q := db.Preload("Position").Preload("Enterprise").Where("student_id = ?", student.ID)
			if pid := c.Query("positionId"); pid != "" {
				q = q.Where("position_id = ?", pid)
			}
			q.Order("created_at DESC").Find(&apps)
			response.Success(c, gin.H{"list": apps})
		})
		secure.POST("/portal/student/apply", func(c *gin.Context) {
			userID := c.GetUint64("user_id")
			var student model.Student
			if err := db.Where("user_id = ?", userID).First(&student).Error; err != nil { response.NotFound(c, "未找到学生信息"); return }
			var req struct{ PositionID uint64 `json:"positionId"` }
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			var pos model.Position
			if err := db.First(&pos, req.PositionID).Error; err != nil { response.NotFound(c, "职位不存在"); return }
			app := model.Application{StudentID: student.ID, PositionID: pos.ID, EnterpriseID: pos.EnterpriseID, Status: "pending"}
			if err := db.Create(&app).Error; err != nil { response.Error(c, response.CodeConflict, "已投递过该职位"); return }
			response.Success(c, app)
		})
		secure.GET("/portal/student/profile", func(c *gin.Context) {
			var s model.Student
			if err := db.Preload("College").Preload("Major").Where("user_id = ?", c.GetUint64("user_id")).First(&s).Error; err != nil { response.NotFound(c, "未找到学生信息"); return }
			if p, err := crypto.Decrypt(s.Phone, cfg.Encrypt.AESKey); err == nil { s.PhonePlain = p }
			if e, err := crypto.Decrypt(s.Email, cfg.Encrypt.AESKey); err == nil { s.EmailPlain = e }
			response.Success(c, s)
		})
		secure.PUT("/admin/users/password", func(c *gin.Context) {
			var req struct{ OldPassword string `json:"oldPassword"`; NewPassword string `json:"newPassword"` }
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			user, err := userRepo.FindByID(c.GetUint64("user_id"))
			if err != nil { response.NotFound(c, "用户不存在"); return }
			if !hash.VerifyPassword(req.OldPassword, user.PasswordHash, cfg.Encrypt.Pepper) {
				response.BadRequest(c, "旧密码错误"); return
			}
			pw, _ := hash.HashPassword(req.NewPassword, cfg.Encrypt.Pepper)
			userRepo.Update(&model.User{ID: user.ID, PasswordHash: pw})
			response.Success(c, nil)
		})
		secure.POST("/upload", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil { response.BadRequest(c, "请选择文件"); return }
			if file.Size > int64(cfg.Upload.MaxSizeMB)*1024*1024 { response.BadRequest(c, "文件过大"); return }
			ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" { response.BadRequest(c, "仅支持 JPG/PNG/WebP"); return }
			name := fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)
			if err := c.SaveUploadedFile(file, filepath.Join(cfg.Upload.Path, name)); err != nil { response.InternalError(c, "上传失败"); return }
			response.Success(c, gin.H{"url": "/uploads/" + name})
		})
		secure.PUT("/portal/student/profile", func(c *gin.Context) {
			var s model.Student
			if err := db.Where("user_id = ?", c.GetUint64("user_id")).First(&s).Error; err != nil { response.NotFound(c, "未找到学生信息"); return }
			var req struct {
				Gender            *int16  `json:"gender"`
				Grade             *string `json:"grade"`
				EducationLevel    *string `json:"educationLevel"`
				GraduationDate    *string `json:"graduationDate"`
				Phone             *string `json:"phone"`
				Email             *string `json:"email"`
				Wechat            *string `json:"wechat"`
				QQ                *string `json:"qq"`
				HometownCity      *string `json:"hometownCity"`
				ExpectedCity      *string `json:"expectedCity"`
				ExpectedIndustry  *string `json:"expectedIndustry"`
				ExpectedSalaryMin *int    `json:"expectedSalaryMin"`
				ExpectedSalaryMax *int    `json:"expectedSalaryMax"`
				ResumeURL         *string `json:"resumeUrl"`
			}
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			updates := map[string]any{}
			if req.Gender != nil { updates["gender"] = *req.Gender }
			if req.Grade != nil { updates["grade"] = *req.Grade }
			if req.EducationLevel != nil { updates["education_level"] = *req.EducationLevel }
			if req.GraduationDate != nil { updates["graduation_date"] = *req.GraduationDate }
			if req.Wechat != nil { updates["wechat"] = *req.Wechat }
			if req.QQ != nil { updates["qq"] = *req.QQ }
			if req.HometownCity != nil { updates["hometown_city"] = *req.HometownCity }
			if req.ExpectedCity != nil { updates["expected_city"] = *req.ExpectedCity }
			if req.ExpectedIndustry != nil { updates["expected_industry"] = *req.ExpectedIndustry }
			if req.ExpectedSalaryMin != nil { updates["expected_salary_min"] = *req.ExpectedSalaryMin }
			if req.ExpectedSalaryMax != nil { updates["expected_salary_max"] = *req.ExpectedSalaryMax }
			if req.ResumeURL != nil { updates["resume_url"] = *req.ResumeURL }
			if req.Phone != nil {
				enc, err := crypto.Encrypt(*req.Phone, cfg.Encrypt.AESKey)
				if err == nil { updates["phone"] = enc }
			}
			if req.Email != nil {
				enc, err := crypto.Encrypt(*req.Email, cfg.Encrypt.AESKey)
				if err == nil { updates["email"] = enc }
			}
			if len(updates) > 0 { db.Model(&s).Updates(updates) }
			db.Preload("College").Preload("Major").First(&s, s.ID)
			if p, err := crypto.Decrypt(s.Phone, cfg.Encrypt.AESKey); err == nil { s.PhonePlain = p }
			if e, err := crypto.Decrypt(s.Email, cfg.Encrypt.AESKey); err == nil { s.EmailPlain = e }
			response.Success(c, s)
		})
	}

	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(jwtService))
	{
		admin.GET("/dashboard/overview", dashboardH.Overview)
		admin.GET("/dashboard/industry-dist", dashboardH.IndustryDist)
		admin.GET("/dashboard/employment-trend", dashboardH.EmploymentTrend)
		admin.GET("/dashboard/salary-analysis", dashboardH.SalaryAnalysis)
		admin.GET("/dashboard/college-employment", dashboardH.CollegeEmployment)

		admin.GET("/users", userH.List); admin.GET("/users/:id", userH.Get)
		admin.POST("/users", userH.Create); admin.PUT("/users/:id", userH.Update)
		admin.DELETE("/users/:id", userH.Delete); admin.POST("/users/delete", userH.BatchDelete)

		admin.GET("/students", studentH.List); admin.GET("/students/:id", studentH.Get)
		admin.POST("/students", studentH.Create); admin.PUT("/students/:id", studentH.Update)
		admin.DELETE("/students/:id", studentH.Delete); admin.POST("/students/delete", studentH.BatchDelete)

		admin.GET("/enterprises", enterpriseH.List); admin.GET("/enterprises/:id", enterpriseH.Get)
		admin.POST("/enterprises", enterpriseH.Create); admin.PUT("/enterprises/:id", enterpriseH.Update)
		admin.DELETE("/enterprises/:id", enterpriseH.Delete); admin.POST("/enterprises/delete", enterpriseH.BatchDelete)

		admin.GET("/positions", positionH.List); admin.GET("/positions/:id", positionH.Get)
		admin.POST("/positions", positionH.Create); admin.PUT("/positions/:id", positionH.Update)
		admin.DELETE("/positions/:id", positionH.Delete); admin.POST("/positions/delete", positionH.BatchDelete)

		admin.GET("/applications", applicationH.List)
		admin.PUT("/applications/:id/status", applicationH.UpdateStatus)
		admin.DELETE("/applications/:id", applicationH.Delete)

		admin.GET("/analysis/city-demand", analysisH.CityDemand)
		admin.GET("/analysis/skill-gap", analysisH.SkillGap)
		admin.GET("/analysis/match-recommend", analysisH.MatchRecommend)
		admin.GET("/analysis/trend-forecast", analysisH.TrendForecast)
		admin.GET("/analysis/report", analysisH.GenerateReport)

		admin.GET("/colleges", collegeH.List); admin.POST("/colleges", collegeH.Create)
		admin.PUT("/colleges/:id", collegeH.Update); admin.DELETE("/colleges/:id", collegeH.Delete)
		admin.POST("/colleges/delete", collegeH.BatchDelete)

		admin.GET("/menus/flat", func(c *gin.Context) {
			var menus []model.Menu; db.Where("status = 1").Order("sort_order ASC").Find(&menus)
			type fm struct{ ID uint64 `json:"id"`; ParentID *uint64 `json:"parentId"`; Label,Title,Path,Component string; Status int16 `json:"status"`; Children []int `json:"children"`; Meta struct{ Title,Icon string; Permission []string `json:"permission"` } `json:"meta"` }
			list := make([]fm, 0, len(menus))
			for _, m := range menus {
				var p []string; if m.Permission != "" { p = strings.Split(m.Permission, ",") }
				list = append(list, fm{m.ID, m.ParentID, m.Name, m.Name, m.Path, m.Component, m.Status, []int{}, struct{ Title,Icon string; Permission []string `json:"permission"` }{m.Name, m.Icon, p}})
			}
			response.Success(c, gin.H{"list": list})
		})
		admin.GET("/roles/list", func(c *gin.Context) {
			var raw []struct{ ID uint64; Name string; Code string; Description string; Status int16 }
			db.Table("roles").Find(&raw)
			type out struct{ ID uint64 `json:"id"`; RoleName string `json:"roleName"`; Code string `json:"code"`; Remark string `json:"remark"`; Status int16 `json:"status"`; MenuIDs []uint64 `json:"menuIds"` }
			list := make([]out, len(raw))
			for i, r := range raw {
				var mids []uint64
				db.Table("role_menus").Where("role_id = ?", r.ID).Pluck("menu_id", &mids)
				list[i] = out{r.ID, r.Name, r.Code, r.Description, r.Status, mids}
			}
			response.Success(c, gin.H{"list": list})
		})
		admin.POST("/roles", func(c *gin.Context) {
			var req struct{ RoleName string `json:"roleName"`; Code string `json:"code"`; Status int16 `json:"status"` }
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			r := model.Role{Name: req.RoleName, Code: req.Code, Status: req.Status}
			db.Create(&r); response.Success(c, r)
		})
		admin.PUT("/roles/:id", func(c *gin.Context) {
			var req struct{ RoleName string `json:"roleName"`; Code string `json:"code"`; Status *int16 `json:"status"`; Remark string `json:"remark"`; Menu []struct{ ID uint64 `json:"id"` } `json:"menu"` }
			if err := c.ShouldBindJSON(&req); err != nil { response.BadRequest(c, "参数错误"); return }
			updates := map[string]any{}
			if req.RoleName != "" { updates["name"] = req.RoleName }
			if req.Code != "" { updates["code"] = req.Code }
			if req.Status != nil { updates["status"] = *req.Status }
			if req.Remark != "" { updates["description"] = req.Remark }
			var rid uint64; fmt.Sscanf(c.Param("id"), "%d", &rid)
			db.Model(&model.Role{}).Where("id = ?", rid).Updates(updates)
			if req.Menu != nil {
				db.Where("role_id = ?", rid).Delete(&model.RoleMenu{})
				for _, m := range req.Menu { db.Create(&model.RoleMenu{RoleID: rid, MenuID: m.ID}) }
			}
			response.Success(c, gin.H{"message": "updated"})
		})
		admin.DELETE("/roles/:id", func(c *gin.Context) {
			db.Delete(&model.Role{}, c.Param("id")); response.Success(c, nil)
		})
	}
	r.Static("/uploads", cfg.Upload.Path)
	return r
}
