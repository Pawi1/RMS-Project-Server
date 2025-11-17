package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"rms/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB         *gorm.DB
	JWTSecret  []byte
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type AccessClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(db *gorm.DB, jwtSecret []byte, accessTTL, refreshTTL time.Duration) (*AuthService, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("auth: jwt secret is empty")
	}
	if accessTTL <= 0 || refreshTTL <= 0 {
		return nil, errors.New("auth: invalid TTLs")
	}
	return &AuthService{
		DB:         db,
		JWTSecret:  jwtSecret,
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}, nil
}

func genJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	var u models.User
	// only select fields we need
	if err := s.DB.WithContext(ctx).Select("id", "password_hash", "role").Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid credentials")
		}
		return "", "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", errors.New("invalid credentials")
	}
	if len(s.JWTSecret) == 0 {
		return "", "", errors.New("server misconfiguration")
	}

	now := time.Now()
	jti, err := genJTI()
	if err != nil {
		return "", "", fmt.Errorf("generating jti: %w", err)
	}

	accessClaims := AccessClaims{
		Role: string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(u.ID),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString(s.JWTSecret)
	if err != nil {
		return "", "", err
	}

	// refresh token: keep minimal claims but include jti for possible revocation
	refreshClaims := jwt.RegisteredClaims{
		Subject:   fmt.Sprint(u.ID),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        jti,
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString(s.JWTSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// IssueAccessForSubject issues an access token for the given subject (user id string).
// jti can be provided to keep token linkage to refresh; if empty, a new jti is generated.
// role, if non-empty, will be embedded into the access token's claims.
func (s *AuthService) IssueAccessForSubject(ctx context.Context, subject, jti, role string) (string, error) {
	if len(s.JWTSecret) == 0 {
		return "", errors.New("server misconfiguration")
	}
	now := time.Now()
	if jti == "" {
		var err error
		jti, err = genJTI()
		if err != nil {
			return "", err
		}
	}
	claims := AccessClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.JWTSecret)
}

// HashPassword returns bcrypt hash of the provided password.
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}
