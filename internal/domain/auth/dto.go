package domain_auth

import (
	"errors"
	"time"
    "regexp"
)

// Request
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *RegisterUserRequest) Validate() error {
    if r.Email == "" || r.Name == "" || r.Password == "" {
        return errors.New("missing required fields")
    }
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    if !emailRegex.MatchString(r.Email) {
        return errors.New("invalid email format")
    }
    if len(r.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginUserRequest) Validate() error {
    if r.Email == "" || r.Password == "" {
        return errors.New("missing required fields")
    }
    return nil
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email"`
}

type ResetPasswordRequest struct {
    OTP string `json:"otp"`
    NewPassword string `json:"new_password"`
}

// Response
type LoginUserResponse struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshTokenRequest) Validate() error {
    if r.RefreshToken == "" {
        return errors.New("refresh token is required")
    }
    return nil
}
