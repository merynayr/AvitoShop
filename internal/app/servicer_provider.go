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
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/service"

	shopAPI "github.com/merynayr/AvitoShop/internal/api/shop"
	userAPI "github.com/merynayr/AvitoShop/internal/api/user"

	shopService "github.com/merynayr/AvitoShop/internal/service/shop"
	userService "github.com/merynayr/AvitoShop/internal/service/user"

	shopRepository "github.com/merynayr/AvitoShop/internal/repository/shop"
	userRepository "github.com/merynayr/AvitoShop/internal/repository/user"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig      config.PGConfig
	httpConfig    config.HTTPConfig
	loggerConfig  config.LoggerConfig
	swaggerConfig config.SwaggerConfig

	dbClient  db.Client
	txManager db.TxManager

	shopAPI        *shopAPI.API
	shopService    service.ShopService
	shopRepository repository.ShopRepository

	userAPI        *userAPI.API
	userService    service.UserService
	userRepository repository.UserRepository
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

func (s *serviceProvider) ShopRepository(ctx context.Context) repository.ShopRepository {
	if s.shopRepository == nil {
		s.shopRepository = shopRepository.NewRepository(s.DBClient(ctx))
	}

	return s.shopRepository
}

func (s *serviceProvider) ShopService(ctx context.Context) service.ShopService {
	if s.shopService == nil {
		s.shopService = shopService.NewService(
			s.ShopRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.shopService
}

func (s *serviceProvider) ShopAPI(ctx context.Context) *shopAPI.API {
	if s.shopAPI == nil {
		s.shopAPI = shopAPI.NewAPI(s.ShopService(ctx))
	}

	return s.shopAPI
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserAPI(ctx context.Context) *userAPI.API {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewAPI(s.UserService(ctx))
	}

	return s.userAPI
}
