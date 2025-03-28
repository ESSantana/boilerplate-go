package routes

import (
	"github.com/application-ellas/ellas-backend/internal/controllers"
	"github.com/application-ellas/ellas-backend/internal/domain/models"
	"github.com/application-ellas/ellas-backend/internal/routes/middlewares"
	svc_interfaces "github.com/application-ellas/ellas-backend/internal/services/interfaces"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/application-ellas/ellas-backend/packages/log"
	"github.com/go-chi/chi/v5"
)

func ConfigRoutes(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	configAuth(router, logger, serviceManager, cacheManager)
	configPayment(router, logger, serviceManager)
	configServiceProvider(router, logger, serviceManager, cacheManager)
}

func configAuth(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewAuthController(logger, serviceManager, cacheManager)

	router.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", controller.SSORequest)
		r.HandleFunc("/callback/{provider}", controller.SSOCallback)
	})
}

func configServiceProvider(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewServiceProviderController(logger, serviceManager, cacheManager)

	router.Route("/service-provider", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware([]string{models.RoleAdmin}))
		r.Post("/create", controller.PromoteUserToServiceProvider)
	})
}

func configPayment(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager) {
	controller := controllers.NewPaymentController(logger, serviceManager)

	router.Route("/payment", func(r chi.Router) {
		r.Post("/", controller.ExecutePayment)
		r.Post("/webhook", controller.PaymentWebhook)
	})
}
