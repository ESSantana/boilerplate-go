package routes

import (
	"net/http"
	"time"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/internal/controllers"
	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/router/middlewares"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/gofiber/fiber/v3"
	mw_cors "github.com/gofiber/fiber/v3/middleware/cors"
	mw_logger "github.com/gofiber/fiber/v3/middleware/logger"
)

type Router struct {
	router         *fiber.App
	cfg            *config.Config
	logger         log.Logger
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
}

func NewRouter(
	router *fiber.App,
	cfg *config.Config,
	logger log.Logger,
	serviceManager svc_interfaces.ServiceManager,
	cacheManager cache_interfaces.CacheManager,
) *Router {
	return &Router{
		router:         router,
		cfg:            cfg,
		logger:         logger,
		serviceManager: serviceManager,
		cacheManager:   cacheManager,
	}
}

func (r *Router) SetupRoutes() {
	middlewares.PrometheusInit()

	r.router.Use(middlewares.TrackMetricsMiddleware())
	r.router.Use(mw_logger.New())
	r.router.Use(mw_cors.New(mw_cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.router.Get("/metrics", middlewares.PrometheusMetricsHandler())
	r.router.Get("/health-check", func(c fiber.Ctx) error {
		dbHealthStatus, cacheHealthStatus := r.serviceManager.HealthCheck()
		payload := map[string]any{
			"checked_at": time.Now().Format(time.RFC3339),
			"mysql":      dbHealthStatus,
			"redis":      cacheHealthStatus,
		}
		c.Status(http.StatusOK)
		return c.JSON(payload)
	})

	r.configAuth()
	r.configCustomer()
	r.configPayment()
}

func (r *Router) configAuth() {
	controller := controllers.NewAuthController(r.cfg, r.logger, r.serviceManager, r.cacheManager)

	// Anonymous routes
	auth := r.router.Group("/auth")
	auth.Post("/login", controller.CustomerLogin)
	auth.Post("/customer/recover-password", controller.RecoverPassword)
	auth.Get("/:provider", controller.SSORequest)
	auth.All("/callback/{provider}", controller.SSOCallback)
}

func (r *Router) configCustomer() {
	controller := controllers.NewCustomerController(r.cfg, r.logger, r.serviceManager, r.cacheManager)

	customer := r.router.Group("/customer")

	// Anonymous routes
	customer.Post("/", controller.Create)

	// Admin and Manager routes
	customer.Get("/", controller.GetAllCustomers, middlewares.AuthMiddleware(r.cfg, []string{constants.RoleAdmin, constants.RoleManager}))

	// Authenticated routes
	customer.Get("/:id", controller.GetCustomerById, middlewares.AuthMiddleware(r.cfg, []string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	customer.Put("/", controller.Update, middlewares.AuthMiddleware(r.cfg, []string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	customer.Delete("/", controller.SoftDelete, middlewares.AuthMiddleware(r.cfg, []string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
}

func (r *Router) configPayment() {
	controller := controllers.NewPaymentController(r.cfg, r.logger, r.serviceManager)

	// Authenticated routes
	payment := r.router.Group("/payment", middlewares.AuthMiddleware(r.cfg, []string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	payment.Post("/", controller.ExecutePayment)
	payment.Post("/webhook", controller.PaymentWebhook)
}
