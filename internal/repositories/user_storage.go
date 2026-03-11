package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// UserStorage handles database operations for User entities.
type UserStorage struct {
	db *sqlx.DB
}

// NewUserStorage creates a new instance of UserStorage with the provided database connection.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

// GetUserByTgID возвращает пользователя по его Telegram ID.
func (us *UserStorage) GetUserByTgID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User

	const sqlString = `SELECT * FROM users WHERE tg_id = $1`

	// Выполнение запроса на получение данных.
	err := dbutils.Get(ctx, us.db, &user, sqlString, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAllUsers возвращает всех пользователей с пагинацией.
// limit — максимальное количество записей, offset — смещение.
func (us *UserStorage) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := `SELECT id, tg_id, name, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2`
	err := dbutils.Select(ctx, us.db, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %w", err)
	}
	return users, nil
}

// GetUserByID возвращает пользователя по его ID (первичный ключ).
func (us *UserStorage) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	const sqlString = `SELECT id, tg_id, name, created_at, updated_at FROM users WHERE id = $1`
	err := dbutils.Get(ctx, us.db, &user, sqlString, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	return &user, nil
}

// CreateUser inserts a new user into the database. Assumes TgID, Name are provided; CreatedAt and UpdatedAt can be set to now.
func (us *UserStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	// Проверка существования пользователя в БД (используем переданный ctx)
	if existingUser, err := us.GetUserByTgID(ctx, user.TgID); err == nil {
		return existingUser, nil
	}

	const sqlString = `
				INSERT INTO users (tg_id, name, timezone, created_at, updated_at) 
				VALUES (:tg_id, :name, :timezone, :created_at, :updated_at)
    `

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	if user.Timezone == "" {
		user.Timezone = "UTC"
	}

	// Выполняем вставку
	_, err := dbutils.NamedExec(ctx, us.db, sqlString, user)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось создать пользователя")
	}

	// После вставки получаем созданного пользователя
	createdUser, err := us.GetUserByTgID(ctx, user.TgID)
	if err != nil {
		return nil, errors.Wrap(err, "пользователь создан, но не удалось получить данные")
	}

	return createdUser, nil
}

// UpdateUserTimezone обновляет часовой пояс пользователя.
func (us *UserStorage) UpdateUserTimezone(ctx context.Context, userID int, timezone string) error {
	const sqlString = `UPDATE users SET timezone = $1, updated_at = NOW() WHERE id = $2`
	_, err := us.db.ExecContext(ctx, sqlString, timezone, userID)
	if err != nil {
		return fmt.Errorf("UpdateUserTimezone: %w", err)
	}
	return nil
}

/*
// UpdateUser updates an existing user in the database. Updates UpdatedAt to now.
func (us *UserStorage) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := us.db.NamedExec("UPDATE users SET tg_id = :tg_id, name = :name, updated_at = :updated_at WHERE id = :id", user)
	return err
}

// DeleteUser removes a user by their ID.
func (us *UserStorage) DeleteUser(id int) error {
	_, err := us.db.Exec("DELETE FROM users WHERE id = \$1", id)
	return err
}*/
