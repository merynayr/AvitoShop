package shop

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
)

const (
	tableMerchName     = "merch_prices"
	tableInventoryName = "inventory"
	tableUsersName     = "users"

	nameColumn = "username"

	itemNameColumn = "item_name"
	priceColumn    = "price"

	userIDColumn   = "user_id"
	quantityColumn = "quantity"

	tableTransactionsName = "transactions"

	transactionIDColumn = "id"
	fromUserIDColumn    = "from_user_id"
	toUserIDColumn      = "to_user_id"
	amountColumn        = "amount"
	createdAtColumn     = "created_at"
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

func (r *repo) CheckInventory(ctx context.Context, userID int64, item string) (bool, int64, error) {
	query, args, err := sq.Select(quantityColumn).
		From(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDColumn: userID, itemNameColumn: item}).
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

	return true, Quantity, nil
}

func (r *repo) InsertNewInventory(ctx context.Context, userID int64, item string) error {
	query, args, err := sq.Insert(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, itemNameColumn, quantityColumn).
		Values(userID, item, 1).
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

func (r *repo) UpdateInventory(ctx context.Context, item string, id, newQuantity int64) error {
	query, args, err := sq.Update(tableInventoryName).
		PlaceholderFormat(sq.Dollar).
		Set(quantityColumn, newQuantity).
		Where(sq.Eq{userIDColumn: id, itemNameColumn: item}).
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

func (r *repo) CreateTransaction(ctx context.Context, fromUserID, toUserID, amount int64) error {
	query, args, err := sq.Insert(tableTransactionsName).
		PlaceholderFormat(sq.Dollar).
		Columns(fromUserIDColumn, toUserIDColumn, amountColumn, createdAtColumn).
		Values(fromUserID, toUserID, amount, time.Now()).
		ToSql()

	if err != nil {
		return fmt.Errorf("could not build SQL query: %v", err)
	}

	q := db.Query{
		Name:     "transactions_repository.CreateTransaction",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	return nil
}

// GetUserInventory получает инвентарь пользователя
func (r *repo) GetUserInventory(ctx context.Context, userID int64) ([]model.InventoryItem, error) {
	query, args, err := sq.Select(itemNameColumn, quantityColumn).
		From(tableInventoryName).
		Where(sq.Eq{userIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserInventory",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.InventoryItem
	for rows.Next() {
		var item model.InventoryItem
		if err := rows.Scan(&item.ItemName, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetUserTransactions получает историю транзакций пользователя (полученные и отправленные)
func (r *repo) GetUserTransactions(ctx context.Context, userID int64) ([]model.Received, []model.Sent, error) {
	receivedQuery, receivedArgs, err := sq.Select(
		"t."+transactionIDColumn, "u_from."+nameColumn, "t."+amountColumn, "t."+createdAtColumn).
		From(tableTransactionsName + " AS t").
		Join(tableUsersName + " AS u_from ON t." + fromUserIDColumn + " = u_from.id").
		Where(sq.Eq{"t." + toUserIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, nil, err
	}

	sentQuery, sentArgs, err := sq.Select(
		"t."+transactionIDColumn, "u_to."+nameColumn, "t."+amountColumn, "t."+createdAtColumn).
		From(tableTransactionsName + " AS t").
		Join(tableUsersName + " AS u_to ON t." + toUserIDColumn + " = u_to.id").
		Where(sq.Eq{"t." + fromUserIDColumn: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, nil, err
	}

	qReceived := db.Query{
		Name:     "user_repository.GetUserTransactions_Received",
		QueryRaw: receivedQuery,
	}

	rowsReceived, err := r.db.DB().QueryContext(ctx, qReceived, receivedArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rowsReceived.Close()

	var received []model.Received
	for rowsReceived.Next() {
		var t model.Received
		if err := rowsReceived.Scan(&t.ID, &t.FromUsername, &t.Amount, &t.CreatedAt); err != nil {
			return nil, nil, err
		}
		received = append(received, t)
	}

	qSent := db.Query{
		Name:     "user_repository.GetUserTransactions_Sent",
		QueryRaw: sentQuery,
	}
	rowsSent, err := r.db.DB().QueryContext(ctx, qSent, sentArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rowsSent.Close()

	var sent []model.Sent
	for rowsSent.Next() {
		var t model.Sent
		if err := rowsSent.Scan(&t.ID, &t.ToUsername, &t.Amount, &t.CreatedAt); err != nil {
			return nil, nil, err
		}
		sent = append(sent, t)
	}

	return received, sent, nil
}
