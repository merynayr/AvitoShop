package integrationtests

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/client/db/pg"
	"github.com/merynayr/AvitoShop/internal/client/db/transaction"
	"github.com/merynayr/AvitoShop/internal/logger"
	shopRepo "github.com/merynayr/AvitoShop/internal/repository/shop"
	userRepo "github.com/merynayr/AvitoShop/internal/repository/user"
	"github.com/merynayr/AvitoShop/internal/service"
	shopService "github.com/merynayr/AvitoShop/internal/service/shop"
)

var (
	testDSN = "postgres://test_user:test_password@localhost:5433/avito_test?sslmode=disable"
)

type MyNewIntegrationSuite struct {
	suite.Suite
	pool *pgxpool.Pool
	r    db.Client
}

func TestMyNewIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MyNewIntegrationSuite))
}

func (s *MyNewIntegrationSuite) SetupSuite() {
	ctx := context.Background()
	logger.Init("Error")
	pool, err := pgxpool.New(ctx, testDSN)
	if err != nil {
		log.Fatal(err)
	}
	s.pool = pool

	s.r, err = pg.New(ctx, testDSN)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *MyNewIntegrationSuite) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.r != nil {
		err := s.r.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func (s *MyNewIntegrationSuite) TearDownTest() {
	_, err := s.pool.Exec(context.Background(), "TRUNCATE TABLE users, inventory, transactions RESTART IDENTITY CASCADE")
	if err != nil {
		log.Printf("failed to clean DB: %v", err)
	}
}

func (s *MyNewIntegrationSuite) SetupTest() service.ShopService {
	txManager := transaction.NewTransactionManager(s.r.DB())
	userRepo := userRepo.NewRepository(s.r)
	shopRepo := shopRepo.NewRepository(s.r)
	userSrv := shopService.NewService(shopRepo, userRepo, txManager)
	return userSrv
}

// Цены на товары
var merchPrices = map[string]int64{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}
