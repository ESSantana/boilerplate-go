package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ESSantana/boilerplate-backend/internal/config"
	"github.com/ESSantana/boilerplate-backend/internal/repositories"
	repo_interfaces "github.com/ESSantana/boilerplate-backend/internal/repositories/interfaces"
	routes "github.com/ESSantana/boilerplate-backend/internal/router"
	"github.com/ESSantana/boilerplate-backend/internal/services"
	svc_interfaces "github.com/ESSantana/boilerplate-backend/internal/services/interfaces"
	"github.com/ESSantana/boilerplate-backend/packages/cache"
	cache_interfaces "github.com/ESSantana/boilerplate-backend/packages/cache/interfaces"
	"github.com/gofiber/fiber/v3"

	"github.com/ESSantana/boilerplate-backend/packages/log"
)

var (
	cfg            *config.Config
	repoManager    repo_interfaces.RepositoryManager
	serviceManager svc_interfaces.ServiceManager
	cacheManager   cache_interfaces.CacheManager
)

func main() {
	var err error
	cfg, err = config.Load("../../")
	if err != nil {
		shutdownApp(err, "Failed to load configuration")
	}

	log.SetGlobalLevel(log.LogLevel[cfg.Server.LogLevel])

	initCache()
	initRepository(context.Background())
	initService()

	startServer()
}

func startServer() {
	app := fiber.New()

	router := routes.NewRouter(app, cfg, serviceManager, cacheManager)
	router.SetupRoutes()

	err := app.Listen(":" + cfg.Server.Port)
	if err != nil {
		shutdownApp(err, "Failed to start server")
	}
}

func initCache() {
	c := sync.Once{}
	c.Do(func() {
		log.Info("Connecting to Redis...")
		cacheManager = cache.NewCacheManager(cfg)
	})
}

func initRepository(ctx context.Context) {
	r := sync.Once{}
	r.Do(func() {
		log.Info("Connecting to MySQL...")
		repoManager = repositories.NewRepositoryManager(ctx, cfg)
	})
}

func initService() {
	s := sync.Once{}
	s.Do(func() {
		log.Info("Setup service manager...")
		serviceManager = services.NewServiceManager(repoManager, cacheManager)
	})
}

func shutdownApp(err error, message string) {
	if err != nil {
		fmt.Println("shutdown with error" + err.Error() + " - " + message)
		os.Exit(1)
	}
}
