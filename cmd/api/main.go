package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/application-ellas/ellas-backend/internal/repositories"
	repo_interfaces "github.com/application-ellas/ellas-backend/internal/repositories/interfaces"
	"github.com/application-ellas/ellas-backend/internal/routes"
	"github.com/application-ellas/ellas-backend/internal/services"
	svc_interfaces "github.com/application-ellas/ellas-backend/internal/services/interfaces"
	"github.com/application-ellas/ellas-backend/packages/cache"
	cache_interfaces "github.com/application-ellas/ellas-backend/packages/cache/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/application-ellas/ellas-backend/packages/log"
)

var logger log.Logger
var repoManager repo_interfaces.RepositoryManager
var serviceManager svc_interfaces.ServiceManager
var cacheManager cache_interfaces.CacheManager
var router *chi.Mux

func main() {
	logger = log.NewLogger(log.DEBUG)
	singletonRepository(context.Background())
	singletonService(logger)
	singletonCache()
	setupRoute()

	defer startServer(router)
	fmt.Printf("Server listening on port :%s\n", os.Getenv("SERVER_PORT"))
}

func setupRoute() {
	router = chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routes.ConfigRoutes(router, logger, serviceManager, cacheManager)
}

func startServer(router *chi.Mux) {
	port := ":" + os.Getenv("SERVER_PORT")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	if err := http.Serve(listen, router); err != nil {
		panic(err)
	}
}

func singletonRepository(ctx context.Context) {
	if repoManager != nil {
		return
	}
	repoManager = repositories.NewRepositoryManager(ctx)
}

func singletonService(logger log.Logger) {
	if serviceManager != nil {
		return
	}
	serviceManager = services.NewServiceManager(logger, repoManager)
}

func singletonCache() {
	if cacheManager != nil {
		return
	}
	cacheManager = cache.NewCacheManager()
}
