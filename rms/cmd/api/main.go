// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	appcfg "rms/internal/config"
	appdb "rms/internal/db"
	"rms/internal/server"
	svc "rms/internal/service"
)

func main() {
	cfg, err := appcfg.Load("config/config.yml")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	mode := cfg.Config.GinMode
	if env := os.Getenv("GIN_MODE"); env != "" {
		mode = env
	}
	switch strings.ToLower(mode) {
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	gdb, created, err := appdb.EnsureDatabase(cfg.Config.DBPath, "db/schema.sql", appdb.OpenSQLite)
	if err != nil {
		log.Fatalf("db init error: %v", err)
	}
	defer func() {
		sqlDB, _ := gdb.DB()
		_ = sqlDB.Close()
	}()
	if created {
		log.Printf("database created and initialized from schema: %s", cfg.Config.DBPath)
	}
	// ensure admin user exists when DB was just created
	if created && cfg.Config.AdminEmail != "" {
		// prefer provided hash, otherwise hash the plain password
		hash := cfg.Config.AdminPasswordHash
		if hash == "" && cfg.Config.AdminPassword != "" {
			// generate bcrypt hash
			h, err := svc.HashPassword(cfg.Config.AdminPassword)
			if err != nil {
				log.Fatalf("failed to hash admin password: %v", err)
			}
			hash = h
		}
		if hash != "" {
			// insert user if not exists
			var count int64
			if err := gdb.WithContext(context.Background()).Raw("SELECT COUNT(1) FROM Users WHERE email = ?", cfg.Config.AdminEmail).Scan(&count).Error; err != nil {
				log.Fatalf("db check admin error: %v", err)
			}
			if count == 0 {
				type user struct {
					UserName     *string
					Email        string
					PasswordHash string
					Role         string
				}
				if err := gdb.Create(&user{nil, cfg.Config.AdminEmail, hash, "admin"}).Error; err != nil {
					log.Fatalf("create admin user error: %v", err)
				}
				log.Printf("admin user created: %s", cfg.Config.AdminEmail)
			}
		}
	}
	health := server.NewHealthHandler(gdb)
	// create auth service
	authSvc, err := svc.NewAuthService(gdb, []byte(cfg.Auth.JWTKey), cfg.Auth.JWTTTL, cfg.Auth.RefreshTokenTTL)
	if err != nil {
		log.Fatalf("auth service init error: %v", err)
	}
	r := server.NewRouterWithConfig(health, cfg, authSvc)
	addr := fmt.Sprintf("%s:%d", cfg.Config.IP, cfg.Config.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
