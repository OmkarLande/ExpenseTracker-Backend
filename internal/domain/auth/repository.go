package domain_auth

import (
	"context"
	"errors"
	"time"

	domain_user "ExpenseTracker-Backend/internal/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain_user.User, error)
	CreateUser(ctx context.Context, user *domain_user.User) error
	UpdateUserEmailVerified(ctx context.Context, userID int64) error
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error

	CreateSession(ctx context.Context, userID int64, refreshTokenHash string, expiresAt time.Time) error
	GetSessionByRefreshTokenHash(ctx context.Context, hash string) (int64, error)
	DeleteSession(ctx context.Context, userID int64) error

	CreatePasswordResetOTP(ctx context.Context, userID int64, otpHash string, expiresAt time.Time) error
	GetPasswordResetOTP(ctx context.Context, otpHash string) (int64, error)
	MarkOTPAsUsed(ctx context.Context, otpHash string) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*domain_user.User, error) {
	query := `SELECT id, email, password_hash, name, photo_url, photo_key, photo_size, bio, is_email_verified, is_transactions_enabled, is_reporting_enabled, is_notification_enabled, created_at, updated_at FROM users WHERE email = $1`
	var user domain_user.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.PhotoURL, &user.PhotoKey, &user.PhotoSize, &user.Bio, &user.IsEmailVerified, &user.IsTransactionsEnabled, &user.IsReportingEnabled, &user.IsNotificationEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *postgresRepository) CreateUser(ctx context.Context, user *domain_user.User) error {
	query := `INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.Name).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *postgresRepository) UpdateUserEmailVerified(ctx context.Context, userID int64) error {
	query := `UPDATE users SET is_email_verified = TRUE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *postgresRepository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *postgresRepository) CreateSession(ctx context.Context, userID int64, refreshTokenHash string, expiresAt time.Time) error {
	query := `INSERT INTO auth_sessions (user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, refreshTokenHash, expiresAt)
	return err
}

func (r *postgresRepository) GetSessionByRefreshTokenHash(ctx context.Context, hash string) (int64, error) {
	query := `SELECT user_id FROM auth_sessions WHERE refresh_token_hash = $1 AND expires_at > NOW()`
	var userID int64
	err := r.db.QueryRow(ctx, query, hash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("invalid or expired session")
		}
		return 0, err
	}
	return userID, nil
}

func (r *postgresRepository) DeleteSession(ctx context.Context, userID int64) error {
	query := `DELETE FROM auth_sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *postgresRepository) CreatePasswordResetOTP(ctx context.Context, userID int64, otpHash string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_otps (user_id, otp_hash, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, otpHash, expiresAt)
	return err
}

func (r *postgresRepository) GetPasswordResetOTP(ctx context.Context, otpHash string) (int64, error) {
	query := `SELECT user_id FROM password_reset_otps WHERE otp_hash = $1 AND expires_at > NOW() AND used = FALSE`
	var userID int64
	err := r.db.QueryRow(ctx, query, otpHash).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errors.New("invalid or expired otp")
		}
		return 0, err
	}
	return userID, nil
}

func (r *postgresRepository) MarkOTPAsUsed(ctx context.Context, otpHash string) error {
	query := `UPDATE password_reset_otps SET used = TRUE WHERE otp_hash = $1`
	_, err := r.db.Exec(ctx, query, otpHash)
	return err
}