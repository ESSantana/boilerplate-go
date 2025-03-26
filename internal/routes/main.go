package routes

import (
	"github.com/ESSantana/boilerplate-go/internal/controllers"
	svc_interfaces "github.com/ESSantana/boilerplate-go/internal/services/interfaces"
	cache_interfaces "github.com/ESSantana/boilerplate-go/packages/cache/interfaces"
	"github.com/ESSantana/boilerplate-go/packages/log"
	"github.com/go-chi/chi/v5"
)

func ConfigRoutes(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	configAuth(router, logger, serviceManager, cacheManager)
	configPayment(router, logger, serviceManager)
}

func configAuth(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager, cacheManager cache_interfaces.CacheManager) {
	controller := controllers.NewAuthController(logger, serviceManager, cacheManager)

	router.Route("/auth", func(r chi.Router) {
		r.Get("/{provider}", controller.SSORequest)
		r.HandleFunc("/callback/{provider}", controller.SSOCallback)
	})
}

func configPayment(router *chi.Mux, logger log.Logger, serviceManager svc_interfaces.ServiceManager) {
	controller := controllers.NewPaymentController(logger, serviceManager)

	router.Route("/payment", func(r chi.Router) {
		r.Post("/", controller.ExecutePayment)
		r.Post("/webhook", controller.PaymentWebhook)
	})
}
