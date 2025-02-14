package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/repository/user/converter"
	modelRepo "github.com/merynayr/AvitoShop/internal/repository/user/model"
	"github.com/merynayr/AvitoShop/internal/sys"
	"github.com/merynayr/AvitoShop/internal/utils/hash"
)

const (
	tableName = "users"

	idColumn       = "id"
	nameColumn     = "username"
	passwordColumn = "password"
	coinsColumn    = "coins"
)

// Структура репо с клиентом базы данных (интерфейсом)
type repo struct {
	db db.Client
}

// NewRepository возвращает новый объект репо слоя
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	passHash, err := hash.EncryptPassword(user.Password)
	if err != nil {
		return 0, err
	}

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, passwordColumn, coinsColumn).
		Values(user.Username, passHash, user.Coins).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().ScanOneContext(ctx, &userID, q, args...)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByName получает из БД информацию пользователя
func (r *repo) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	query, args, err := sq.Select(idColumn, nameColumn, passwordColumn, coinsColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserByName",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sys.UserNotFoundError
		}
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser обновляет данные пользователя по id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar)

	if user.Username != "" {
		builderUpdate = builderUpdate.Set(nameColumn, &user.Username)
	}

	if user.Coins >= 0 {
		builderUpdate = builderUpdate.Set(coinsColumn, &user.Coins)
	}

	builderUpdate = builderUpdate.Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
