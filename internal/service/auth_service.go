package service

import (
	"context"
	"employment-server/internal/config"
	"employment-server/internal/model"
	"employment-server/internal/repository"
	"employment-server/pkg/hash"
	jwtPkg "employment-server/pkg/jwt"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrAccountLocked      = errors.New("账号已被锁定，请稍后再试")
	ErrAccountDisabled    = errors.New("账号已被禁用")
	ErrTokenRevoked       = errors.New("Token已被撤销")
	ErrRefreshFailed      = errors.New("刷新Token失败")
)

type AuthService struct {
	userRepo *repository.UserRepo
	jwtSvc   *jwtPkg.Service
	cfg      *config.Config
	redis    *redis.Client
}

func NewAuthService(userRepo *repository.UserRepo, jwtSvc *jwtPkg.Service, cfg *config.Config, redis *redis.Client) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSvc: jwtSvc, cfg: cfg, redis: redis}
}

func (s *AuthService) Login(ctx context.Context, username, password, ip string) (*jwtPkg.TokenPair, *model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if user.Status == 0 {
		return nil, nil, ErrAccountDisabled
	}
	if user.Status == 2 && user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, nil, ErrAccountLocked
	}

	if !hash.VerifyPassword(password, user.PasswordHash, s.cfg.Encrypt.Pepper) {
		s.userRepo.IncrementLoginAttempts(user.ID)
		// 连续失败5次 → 锁定30分钟
		if user.LoginAttempts >= 4 {
			lockUntil := time.Now().Add(30 * time.Minute)
			s.userRepo.LockAccount(user.ID, lockUntil)
		}
		return nil, nil, ErrInvalidCredentials
	}

	roles := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roles = append(roles, r.Code)
	}

	pair, err := s.jwtSvc.GenerateTokenPair(user.ID, user.Username, roles)
	if err != nil {
		return nil, nil, err
	}

	refreshKey := fmt.Sprintf("refresh:%d:%s", user.ID, "session")
	if err := s.redis.Set(ctx, refreshKey, pair.RefreshToken, s.cfg.JWT.RefreshTTL).Err(); err != nil {
		return nil, nil, fmt.Errorf("存储RefreshToken失败: %w", err)
	}

	s.userRepo.UpdateLoginInfo(user.ID, ip)

	return pair, user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*jwtPkg.TokenPair, error) {
	claims, err := s.jwtSvc.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrTokenRevoked
	}

	refreshKey := fmt.Sprintf("refresh:%d:%s", claims.UserID, "session")
	stored, err := s.redis.Get(ctx, refreshKey).Result()
	if err != nil || stored != refreshToken {
		return nil, ErrTokenRevoked
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roles = append(roles, r.Code)
	}

	newPair, err := s.jwtSvc.GenerateTokenPair(user.ID, user.Username, roles)
	if err != nil {
		return nil, err
	}

	s.redis.Del(ctx, refreshKey)
	if err := s.redis.Set(ctx, refreshKey, newPair.RefreshToken, s.cfg.JWT.RefreshTTL).Err(); err != nil {
		return nil, ErrRefreshFailed
	}

	return newPair, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uint64) error {
	refreshKey := fmt.Sprintf("refresh:%d:%s", userID, "session")
	return s.redis.Del(ctx, refreshKey).Err()
}
