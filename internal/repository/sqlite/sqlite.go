package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tamper000/hysteria-auth/internal/repository"
	_ "modernc.org/sqlite"
)

type SqliteRepo struct {
	db *sql.DB
}

func New(dbPath string) (*SqliteRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Инициализация таблицы
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return &SqliteRepo{db: db}, nil
}

func initSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		auth TEXT PRIMARY KEY,
		id TEXT NOT NULL,
		optional TEXT
	);`

	_, err := db.Exec(query)
	return err
}

func (r *SqliteRepo) Close() error {
	return r.db.Close()
}

func (r *SqliteRepo) CreateUser(ctx context.Context, user repository.User) error {
	query := `INSERT INTO user (id, auth, optional) VALUES (?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Auth, user.Optional)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *SqliteRepo) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM user WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}

func (r *SqliteRepo) UserExists(ctx context.Context, auth string) (bool, error) {
	query := `SELECT COUNT(*) FROM user WHERE auth = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, auth).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

func (r *SqliteRepo) GetIDByAuth(ctx context.Context, auth string) (string, bool, error) {
	query := `SELECT id FROM user WHERE auth = ?`

	var id string
	err := r.db.QueryRowContext(ctx, query, auth).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to get user id by auth: %w", err)
	}

	return id, true, nil
}
