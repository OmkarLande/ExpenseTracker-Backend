package domain_auth

import "time"

type TokenStatus int
const (
	Active TokenStatus = 1
	Inactive TokenStatus = 2
	Expired TokenStatus = 3
)

type VerificationTokenType string
const (
	EmailVerification    VerificationTokenType = "email_verification"
	PasswordReset        VerificationTokenType = "password_reset"
)

type VerificationToken struct {
	ID        string    `db:"id"`
	UserID    int64     `db:"user_id"`
	Token     string    `db:"token"`
	Type      VerificationTokenType    `db:"type"`
	Status    TokenStatus `db:"status"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
