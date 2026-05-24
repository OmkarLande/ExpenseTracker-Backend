package domain_auth

import (
	"context"
	"errors"
	"log"
	"time"

	"ExpenseTracker-Backend/internal/config"
	domain_user "ExpenseTracker-Backend/internal/domain/user"
	"ExpenseTracker-Backend/internal/utils"
)

type Service interface {
	Register(ctx context.Context, req RegisterUserRequest) error
	Login(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginUserResponse, error)
	Logout(ctx context.Context, userID int64) error
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, otp, newPassword string) error
}

type service struct {
	repo        Repository
	emailSender *EmailSender
	config      *config.Config
}

func NewService(repo Repository, emailSender *EmailSender, cfg *config.Config) Service {
	return &service{
		repo:        repo,
		emailSender: emailSender,
		config:      cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterUserRequest) error {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &domain_user.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// Generate email verification token 
	// token, _, _, err := utils.GenerateTokens(user.ID, s.config.JWTSecret)
	// if err == nil {
		// _ = s.emailSender.SendVerificationEmail(user.Email, token)
	// }

	return nil
}

func (s *service) Login(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, refreshToken, expiresAt, err := utils.GenerateTokens(user.ID, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshHash, _ := HashPassword(refreshToken)
	refreshExpires := time.Now().Add(7 * 24 * time.Hour)
	err = s.repo.CreateSession(ctx, user.ID, refreshHash, refreshExpires)
	if err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginUserResponse, error) {
	claims, err := utils.ValidateAccessToken(req.RefreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
    
	accessToken, refreshToken, expiresAt, err := utils.GenerateTokens(claims.UserID, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}
	
	refreshHash, _ := HashPassword(refreshToken)
	refreshExpires := time.Now().Add(7 * 24 * time.Hour)
	_ = s.repo.CreateSession(ctx, claims.UserID, refreshHash, refreshExpires)

	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *service) Logout(ctx context.Context, userID int64) error {
	return s.repo.DeleteSession(ctx, userID)
}

func (s *service) VerifyEmail(ctx context.Context, token string) error {
	claims, err := utils.ValidateAccessToken(token, s.config.JWTSecret)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	return s.repo.UpdateUserEmailVerified(ctx, claims.UserID)
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil
	}

	otp, err := GenerateNumericOTP(6)
	if err != nil {
		return err
	}

	otpHash, _ := HashOTP(otp)
	expiresAt := time.Now().Add(15 * time.Minute)

	err = s.repo.CreatePasswordResetOTP(ctx, user.ID, otpHash, expiresAt)
	if err != nil {
		return err
	}

	err = s.emailSender.SendPasswordResetEmail(user.Email, otp)
	if err != nil {
		log.Printf("ForgotPassword: failed to send email to %s: %v", user.Email, err)
	}
	return nil
}

// TODO: Complete this function
func (s *service) ResetPassword(ctx context.Context, otp, newPassword string) error {
	return errors.New("reset password logic requires lookup by email and OTP comparison")
}