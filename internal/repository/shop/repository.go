package shop

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableMerchName     = "merch_prices"
	tableInventoryName = "inventory"

	itemNameColumn = "item_name"
	priceColumn    = "price"

	userIDColumn   = "user_id"
	quantityColumn = "quantity"
)

// Репозиторий
type repo struct {
	db db.Client
}

// NewRepository возвращает объект репозитория
func NewRepository(db db.Client) repository.ShopRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) GetMerchPrice(ctx context.Context, item string) (int64, error) {
	query, args, err := sq.Select(priceColumn).
		From(tableMerchName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{itemNameColumn: item}).
		Limit(1).ToSql()

	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "shop_repository.GetMerchPrice",
		QueryRaw: query,
	}

	var price int64
	err = r.db.DB().ScanOneContext(ctx, &price, q, args...)
	if err != nil {
		logger.Debug("%s: failed to get price on merch: %v", q.Name, err)
		return 0, err
	}
	return price, nil
}

func (r *repo) CheckInventory(ctx context.Context, user *model.User, item string) (bool, int64, error) {
	query, args, err := sq.Select("id", quantityColumn).
		From(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDColumn: user.ID, itemNameColumn: item}).
		Suffix("FOR UPDATE").
		ToSql()

	if err != nil {
		return false, 0, fmt.Errorf("could not build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "shop_repository.CheckInventory",
		QueryRaw: query,
	}

	var Quantity int64
	err = r.db.DB().ScanOneContext(ctx, &Quantity, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, 0, nil
		}

		return false, 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return true, Quantity + 1, nil
}

func (r *repo) InsertNewInventory(ctx context.Context, user *model.User, item string) error {
	query, args, err := sq.Insert(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, itemNameColumn, quantityColumn).
		Values(user.ID, item, 1).
		ToSql()

	if err != nil {
		return fmt.Errorf("could not build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "shop_repository.InsertNewInventory",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to insert new inventory record: %v", err)
	}

	return nil
}
func (r *repo) UpdateInventory(ctx context.Context, id, newQuantity int64) error {
	query, args, err := sq.Update(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Set(quantityColumn, newQuantity).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("could not build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "shop_repository.UpdateInventory",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	return nil
}
