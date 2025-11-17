package server

import (
	"log"

	"rms/internal/config"
	admin "rms/internal/handlers/admin"
	operatorh "rms/internal/handlers/operator"
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

	// operator routes: create users and manage invitations
	op := r.Group("/operator")
	op.Use(authmw.RequireAuth(authSvc))
	op.Use(authmw.RequireRoles("admin", "operator", "editor"))
	// operator/editor create clients only
	op.POST("/clients", operatorh.CreateClientHandler(authSvc))
	op.PATCH("/clients/:id", operatorh.PatchClientHandler(authSvc))
	op.POST("/orders", operatorh.CreateOrderHandler(hh.DB))
	op.POST("/tasks", operatorh.CreateTaskHandler(hh.DB))
	// operator cars handlers: create an instance and register its methods
	opCars := operatorh.NewCarsHandlers(hh.DB)
	op.GET("/cars", opCars.ListCars)
	op.POST("/cars", opCars.CreateCar)
	op.PATCH("/cars/:id", opCars.PatchCar)
	op.DELETE("/cars/:id", opCars.DeleteCar)

	// admin routes (admins can also create users)
	adm := r.Group("/admin")
	adm.Use(authmw.RequireAuth(authSvc))
	adm.Use(authmw.RequireRoles("admin"))
	adm.POST("/users", admin.CreateUserHandler(authSvc))
	adm.PATCH("/users/:id", admin.PatchUserHandler(authSvc))
	adm.POST("/orders", admin.CreateOrderHandler(hh.DB))
	adm.POST("/tasks", admin.CreateTaskHandler(hh.DB))
	admCars := operatorh.NewCarsHandlers(hh.DB)
	adm.GET("/cars", admCars.ListCars)
	adm.POST("/cars", admCars.CreateCar)
	adm.PATCH("/cars/:id", admCars.PatchCar)
	adm.DELETE("/cars/:id", admCars.DeleteCar)

	return r
}
