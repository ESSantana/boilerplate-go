package routes

import (
	"net/http"
	"time"

	"github.com/ESSantana/boilerplate-backend/internal/controllers"
	"github.com/ESSantana/boilerplate-backend/internal/domain/constants"
	"github.com/ESSantana/boilerplate-backend/internal/routes/middlewares"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/log"
	"github.com/gofiber/fiber/v3"
)

func ConfigRoutes(router *fiber.App, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	router.Get("/metrics", middlewares.PrometheusMetricsHandler())
	router.Get("/health-check", func(c fiber.Ctx) error {
		dbHealthStatus, cacheHealthStatus := serviceManager.HealthCheck()
		payload := map[string]any{
			"checked_at": time.Now().Format(time.RFC3339),
			"mysql":      dbHealthStatus,
			"redis":      cacheHealthStatus,
		}
		c.Status(http.StatusOK)
		return c.JSON(payload)
	})

	configAuth(router, logger, serviceManager, cacheManager)
	configCustomer(router, logger, serviceManager, cacheManager)
	configPayment(router, logger, serviceManager)
}

func configAuth(app *fiber.App, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewAuthController(logger, serviceManager, cacheManager)

	// Anonymous routes
	auth := app.Group("/auth")
	auth.Post("/login", controller.CustomerLogin)
	auth.Post("/customer/recover-password", controller.RecoverPassword)
	auth.Get("/:provider", controller.SSORequest)
	auth.All("/callback/{provider}", controller.SSOCallback)
}

func configCustomer(app *fiber.App, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewCustomerController(logger, serviceManager, cacheManager)

	customer := app.Group("/customer")

	// Anonymous routes
	customer.Post("/", controller.Create)

	// Admin and Manager routes
	customer.Get("/", controller.GetAllCustomers, middlewares.AuthMiddleware([]string{constants.RoleAdmin, constants.RoleManager}))

	// Authenticated routes
	customer.Get("/:id", controller.GetCustomerById, middlewares.AuthMiddleware([]string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	customer.Put("/", controller.Update, middlewares.AuthMiddleware([]string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	customer.Delete("/", controller.SoftDelete, middlewares.AuthMiddleware([]string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
}

func configPayment(app *fiber.App, logger log.Logger, serviceManager svc_interfaces.ServiceManager) {
	controller := controllers.NewPaymentController(logger, serviceManager)

	// Authenticated routes
	payment := app.Group("/payment", middlewares.AuthMiddleware([]string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
	payment.Post("/", controller.ExecutePayment)
	payment.Post("/webhook", controller.PaymentWebhook)
}
