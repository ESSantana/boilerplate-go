package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/application-ellas/ella-backend/internal/controllers"
	"github.com/application-ellas/ella-backend/internal/domain/constants"
	"github.com/application-ellas/ella-backend/internal/routes/middlewares"
	svc_interfaces "github.com/application-ellas/ella-backend/internal/services/interfaces"
	cache_interfaces "github.com/application-ellas/ella-backend/packages/cache/interfaces"
	"github.com/application-ellas/ella-backend/packages/log"
	"github.com/go-chi/chi/v5"
)

func ConfigRoutes(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		dbHealthStatus, cacheHealthStatus := serviceManager.HealthCheck()
		payload := map[string]any{
			"checked_at": time.Now().Format(time.RFC3339),
			"mysql":      dbHealthStatus,
			"redis":      cacheHealthStatus,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	})

	configAuth(router, logger, serviceManager, cacheManager)
	configPayment(router, logger, serviceManager)
	configCustomer(router, logger, serviceManager, cacheManager)
}

func configAuth(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewAuthController(logger, serviceManager, cacheManager)

	router.Route("/auth", func(r chi.Router) {
		r.Post("/customer", controller.CustomerLogin)
		r.Post("/customer/recover-password", controller.RecoverPassword)
		// r.Get("/{provider}", controller.SSORequest)
		// r.HandleFunc("/callback/{provider}", controller.SSOCallback)
	})
}

func configCustomer(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewCustomerController(logger, serviceManager, cacheManager)

	// Anonymous routes
	router.Group(func(r chi.Router) {
		r.Post("/customer", controller.Create)
	})

	// Admin and Manager routes
	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware([]string{constants.RoleAdmin, constants.RoleManager}))
		r.Get("/customer", controller.GetAllCustomers)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware([]string{constants.RoleCustomer, constants.RoleAdmin, constants.RoleManager}))
		r.Get("/customer/{id}", controller.GetCustomerById)
		r.Put("/customer", controller.Update)
		r.Delete("/customer", controller.SoftDelete)
	})
}

func configPayment(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager) {
	controller := controllers.NewPaymentController(logger, serviceManager)

	router.Route("/payment", func(r chi.Router) {
		r.Post("/", controller.ExecutePayment)
		r.Post("/webhook", controller.PaymentWebhook)
	})
}
