package main

import (
	"context"
	"os"
	"sync"

	"github.com/ESSantana/boilerplate-backend/internal/repositories"
	repo_interfaces "github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/routes"
	"github.com/ESSantana/boilerplate-backend/internal/routes/middlewares"
	"github.com/ESSantana/boilerplate-backend/internal/services"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/cache"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/gofiber/fiber/v3"
	mw_cors "github.com/gofiber/fiber/v3/middleware/cors"
	mw_logger "github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/ESSantana/boilerplate-backend/packages/log"
)

var (
	logger         log.Logger
	repoManager    repo_interfaces.RepositoryManager
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
	app            *fiber.App
)


func main() {
	logger = log.NewLogger(log.DEBUG)

	initCache()
	initRepository(context.Background())
	initService(logger)

	startServer()
}

func startServer() {
	app = fiber.New()
	middlewares.PrometheusInit()

	app.Use(middlewares.TrackMetricsMiddleware())
	app.Use(mw_logger.New())
	app.Use(mw_cors.New(mw_cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routes.ConfigRoutes(app, logger, serviceManager, cacheManager)

	err := app.Listen(":" + os.Getenv("SERVER_PORT"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func initCache() {
	c := sync.Once{}
	c.Do(func() {
		logger.Info("Connecting to Redis...")
		cacheManager = cache.NewCacheManager()
	})
}

func initRepository(ctx context.Context) {
	r := sync.Once{}
	r.Do(func() {
		logger.Info("Connecting to MySQL...")
		repoManager = repositories.NewRepositoryManager(ctx)
	})
}

func initService(logger log.Logger) {
	s := sync.Once{}
	s.Do(func() {
		logger.Info("Setup service manager...")
		serviceManager = services.NewServiceManager(logger, repoManager, cacheManager)
	})
}
