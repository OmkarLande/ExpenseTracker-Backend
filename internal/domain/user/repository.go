package domain_user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateProfilePictureMetadata(ctx context.Context, userID int64, photoURL string, photoKey string, photoSize int64) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, name, photo_url, photo_key, photo_size, bio, is_email_verified, 
		       is_transactions_enabled, is_reporting_enabled, is_notification_enabled, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.PhotoURL, &user.PhotoKey, &user.PhotoSize, &user.Bio, &user.IsEmailVerified,
		&user.IsTransactionsEnabled, &user.IsReportingEnabled, &user.IsNotificationEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, name, photo_url, photo_key, photo_size, bio, is_email_verified, 
		       is_transactions_enabled, is_reporting_enabled, is_notification_enabled, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.PhotoURL, &user.PhotoKey, &user.PhotoSize, &user.Bio, &user.IsEmailVerified,
		&user.IsTransactionsEnabled, &user.IsReportingEnabled, &user.IsNotificationEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Name).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *postgresRepository) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET name = $1, photo_url = $2, photo_key = $3, photo_size = $4, bio = $5, updated_at = NOW() WHERE id = $6`
	_, err := r.db.Exec(ctx, query, user.Name, user.PhotoURL, user.PhotoKey, user.PhotoSize, user.Bio, user.ID)
	return err
}

func (r *postgresRepository) UpdateProfilePictureMetadata(ctx context.Context, userID int64, photoURL string, photoKey string, photoSize int64) error {
	query := `UPDATE users SET photo_url = $1, photo_key = $2, photo_size = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(ctx, query, photoURL, photoKey, photoSize, userID)
	return err
}