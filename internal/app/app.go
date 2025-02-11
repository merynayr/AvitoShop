package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/merynayr/AvitoShop/internal/closer"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/rakyll/statik/fs"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"

	_ "github.com/merynayr/AvitoShop/statik" //
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
	swaggerServer   *http.Server
}

// NewApp возвращает объект приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

// initDeps инициализирует все зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	logger.Init(a.serviceProvider.LoggerConfig().Level())

	router := gin.Default()

	a.serviceProvider.AccessMiddleware()
	router.Use(a.serviceProvider.middleware.Check())

	a.serviceProvider.ShopAPI(ctx)
	a.serviceProvider.shopAPI.RegisterRoutes(router)

	a.serviceProvider.UserAPI(ctx)
	a.serviceProvider.userAPI.RegisterRoutes(router)

	a.serviceProvider.AuthAPI(ctx)
	a.serviceProvider.authAPI.RegisterRoutes(router)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log := logger.With("address", a.serviceProvider.HTTPConfig().Address())
	log.Info("HTTP server is running")

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Debug("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			_ = file.Close()
		}()

		logger.Debug("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Served swagger file: %s", path)
	}
}
