package main

import (
	"context"
	"os"

	"github.com/ESSantana/boilerplate-backend/internal/repositories"
	repo_interfaces "github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	"github.com/ESSantana/boilerplate-backend/internal/routes"
	"github.com/ESSantana/boilerplate-backend/internal/services"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/cache"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/gofiber/fiber/v3"
	mw_cors "github.com/gofiber/fiber/v3/middleware/cors"
	mw_logger "github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/ESSantana/boilerplate-backend/packages/log"
)

var logger log.Logger
var repoManager repo_interfaces.RepositoryManager
var serviceManager svc_interfaces.ServiceManager
var cacheManager cache_interfaces.CacheManager
var app *fiber.App

func main() {
	logger = log.NewLogger(log.DEBUG)
	singletonCache()
	singletonRepository(context.Background())
	singletonService(logger)

	startServer()
}

func startServer() {
	app = fiber.New()

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

func singletonRepository(ctx context.Context) {
	logger.Info("Connecting to MySQL...")
	if repoManager != nil {
		return
	}
	repoManager = repositories.NewRepositoryManager(ctx)
}

func singletonService(logger log.Logger) {
	logger.Info("Setup service manager...")
	if serviceManager != nil {
		return
	}
	serviceManager = services.NewServiceManager(logger, repoManager, cacheManager)
}

func singletonCache() {
	logger.Info("Connecting to Redis...")
	if cacheManager != nil {
		return
	}
	cacheManager = cache.NewCacheManager()
}
