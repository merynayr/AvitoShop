package user

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/merynayr/AvitoShop/internal/client/db"
	"github.com/merynayr/AvitoShop/internal/logger"
	"github.com/merynayr/AvitoShop/internal/model"
	"github.com/merynayr/AvitoShop/internal/repository"
	"github.com/merynayr/AvitoShop/internal/repository/user/converter"
	modelRepo "github.com/merynayr/AvitoShop/internal/repository/user/model"
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
	op := "CreateUser"

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
		logger.Debug("%s: failed to create builder: %v", op, err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err = r.db.DB().ScanOneContext(ctx, &userID, q, args...)
	if err != nil {
		logger.Debug("%s: failed to insert user: %v", op, err)
		return 0, err
	}

	logger.Debug("%s: inserted user with id: %d", op, userID)
	return userID, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("user with id %d doesn't exist", userID)
	}

	query, args, err := sq.Select(idColumn, nameColumn, coinsColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserByID",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		logger.Debug("%s: failed to select user: %v", q.Name, err)
		return nil, err
	}

	logger.Debug("%s: selected user %d", q.Name, userID)
	return converter.ToUserFromRepo(&user), nil
}

// GetUserByName получает из БД информацию пользователя
func (r *repo) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	exist, err := r.IsNameExist(ctx, name)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, fmt.Errorf("user with name %s doesn't exist", name)
	}

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
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		logger.Debug("%s: failed to select user: %v", q.Name, err)
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser обновляет данные пользователя по id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdate) error {
	op := "UpdateUser"

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
		logger.Debug("%s: failed to create builder: %v", op, err)
		return err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Debug("%s: failed to update user: %v", op, err)
		return err
	}

	logger.Debug("%s: updated user %d", op, user.ID)
	return nil
}

// DeleteUser удаляет пользователя по id
func (r *repo) DeleteUser(ctx context.Context, userID int64) error {
	exist, err := r.IsExistByID(ctx, userID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("user with id %d doesn't exist", userID)
	}

	query, args, err := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Debug("%s: failed to delete user: %v", q.Name, err)
		return err
	}

	return nil
}

// IsExistById проверяет, существует ли в БД пользователь с указанным ID
func (r *repo) IsExistByID(ctx context.Context, userID int64) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userID}).
		ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsExistByID",
		QueryRaw: query,
	}

	var user int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// IsNameExist проверяет, существует ли в БД указанный name
func (r *repo) IsNameExist(ctx context.Context, name string) (bool, error) {
	query, args, err := sq.Select("1").
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{nameColumn: name}).
		Limit(1).ToSql()

	if err != nil {
		return false, err
	}

	q := db.Query{
		Name:     "user_repository.IsNameExist",
		QueryRaw: query,
	}

	var one int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
