package server

import (
	"log"

	"rms/internal/config"
	carsh "rms/internal/handlers"
	pub "rms/internal/handlers/public"
	userh "rms/internal/handlers/user"
	authmw "rms/internal/middleware/auth"
	svc "rms/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouterWithConfig(hh *HealthHandler, cfg config.App, authSvc *svc.AuthService) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	var err error
	if len(cfg.Config.TrustedProxies) == 0 {
		err = r.SetTrustedProxies(nil)
	} else {
		err = r.SetTrustedProxies(cfg.Config.TrustedProxies)
	}
	if err != nil {
		log.Printf("trusted proxies config error: %v", err)
	}
	r.GET("/healthz", hh.Get)

	// public auth routes
	r.POST("/auth/login", pub.LoginHandler(authSvc))
	r.POST("/auth/refresh", pub.RefreshHandler(authSvc))

	// public workshop info
	r.GET("/workshop", pub.WorkshopHandler(cfg.CarWorkshop))

	// protected /me endpoints
	me := r.Group("/me")
	me.Use(authmw.RequireAuth(authSvc))
	me.PATCH("", userh.PatchMeHandler(authSvc))
	me.POST("/password", userh.ChangePasswordHandler(authSvc))

	// cars management - only roles admin, editor, operator
	cars := carsh.NewCarsHandlers(hh.DB)
	adminGroup := r.Group("/cars")
	adminGroup.Use(authmw.RequireAuth(authSvc))
	adminGroup.Use(authmw.RequireRoles("admin", "editor", "operator"))
	adminGroup.GET("", cars.List)
	adminGroup.POST("", cars.Create)
	adminGroup.PATCH(":id", cars.Patch)
	adminGroup.DELETE(":id", cars.Delete)

	return r
}
