package app

import (
	"context"
	"log"

	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/client/db/pg"
	"github.com/merynayr/AvitoShop/internal/client/db/transaction"
	"github.com/merynayr/AvitoShop/internal/closer"
	"github.com/merynayr/AvitoShop/internal/config"
	"github.com/merynayr/AvitoShop/internal/config/env"
	"github.com/merynayr/AvitoShop/internal/middleware/access"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/service"

	authAPI "github.com/merynayr/AvitoShop/internal/api/auth"
	shopAPI "github.com/merynayr/AvitoShop/internal/api/shop"

	accessService "github.com/merynayr/AvitoShop/internal/service/access"
	authService "github.com/merynayr/AvitoShop/internal/service/auth"
	shopService "github.com/merynayr/AvitoShop/internal/service/shop"

	shopRepository "github.com/merynayr/AvitoShop/internal/repository/shop"
	userRepository "github.com/merynayr/AvitoShop/internal/repository/user"

	"github.com/merynayr/AvitoShop/internal/middleware"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig      config.PGConfig
	httpConfig    config.HTTPConfig
	loggerConfig  config.LoggerConfig
	swaggerConfig config.SwaggerConfig
	authConfig    config.AuthConfig
	accessConfig  config.AccessConfig

	dbClient  db.Client
	txManager db.TxManager

	shopAPI        *shopAPI.API
	shopService    service.ShopService
	shopRepository repository.ShopRepository
	userRepository repository.UserRepository

	authAPI     *authAPI.API
	authService service.AuthService

	middleware    middleware.UserMiddleware
	accessService service.AccessService
}

// NewServiceProvider возвращает новый объект API слоя
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := env.NewLoggerConfig()
		if err != nil {
			log.Fatalf("failed to get logger config:%v", err)
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// AuthConfig инициализирует конфиг auth сервиса
func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := env.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config")
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

// AccessConfig инициализирует конфиг access конфига
func (s *serviceProvider) AccessConfig() config.AccessConfig {
	if s.accessConfig == nil {
		cfg, err := env.NewAccessConfig()
		if err != nil {
			log.Fatalf("failed to get access service")
		}

		s.accessConfig = cfg
	}

	return s.accessConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("shop error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ShopAPI(ctx context.Context) *shopAPI.API {
	if s.shopAPI == nil {
		s.shopAPI = shopAPI.NewAPI(s.ShopService(ctx))
	}

	return s.shopAPI
}

func (s *serviceProvider) ShopService(ctx context.Context) service.ShopService {
	if s.shopService == nil {
		s.shopService = shopService.NewService(
			s.ShopRepository(ctx),
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.shopService
}

func (s *serviceProvider) ShopRepository(ctx context.Context) repository.ShopRepository {
	if s.shopRepository == nil {
		s.shopRepository = shopRepository.NewRepository(s.DBClient(ctx))
	}

	return s.shopRepository
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// AuthAPI инициализирует api слой auth
func (s *serviceProvider) AuthAPI(ctx context.Context) *authAPI.API {
	if s.authAPI == nil {
		s.authAPI = authAPI.NewAPI(s.AuthService(ctx), s.AuthConfig())
	}

	return s.authAPI
}

// AuthService иницилизирует сервисный слой auth
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.UserRepository(ctx),
			s.AuthConfig(),
		)
	}

	return s.authService
}

// AccessMiddleware инициализирует middleware доступа
func (s *serviceProvider) AccessMiddleware(ctx context.Context) middleware.UserMiddleware {
	if s.middleware == nil {
		s.middleware = access.NewMiddleware(
			s.AccessService(ctx),
			s.AuthConfig(),
		)
	}
	return s.middleware
}

// AccessService иницилизирует сервисный слой access
func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		uMap, err := s.AccessConfig().UserAccessesMap()
		if err != nil {
			log.Fatalf("failed to get user access map: %v", err)
		}

		s.accessService = accessService.NewService(s.ShopService(ctx), uMap, s.AuthConfig())
	}

	return s.accessService
}
