package domain_user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"ExpenseTracker-Backend/internal/config"
	"ExpenseTracker-Backend/internal/storage"
)

type Service interface {
	GetProfile(ctx context.Context, userID int64) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID int64, req UpdateUserRequest) (*UserResponse, error)
	PresignProfilePicture(ctx context.Context, userID int64, req PresignProfilePictureRequest) (*PresignProfilePictureResponse, error)
	CompleteProfilePicture(ctx context.Context, userID int64, req CompleteProfilePictureRequest) (*UserResponse, error)
}

type service struct {
	repo    Repository
	storage storage.Storage
	config  *config.Config
}

func NewService(repo Repository, storage storage.Storage, cfg *config.Config) Service {
	return &service{
		repo:    repo,
		storage: storage,
		config:  cfg,
	}
}

func (s *service) GetProfile(ctx context.Context, userID int64) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.PhotoKey != nil && *user.PhotoKey != "" {
		downloadURL, err := s.storage.GenerateDownloadURL(ctx, *user.PhotoKey)
		if err == nil {
			user.PhotoURL = &downloadURL
		}
	}

	return toUserResponse(user), nil
}

func (s *service) UpdateProfile(ctx context.Context, userID int64, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = strings.TrimSpace(req.Name)
	if user.Name == "" {
		user.Name = "User"
	}

	if req.Bio != nil {
		trimmedBio := strings.TrimSpace(*req.Bio)
		user.Bio = &trimmedBio
	}

	if req.PhotoURL != nil {
		trimmedPhoto := strings.TrimSpace(*req.PhotoURL)
		user.PhotoURL = &trimmedPhoto
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if user.PhotoKey != nil && *user.PhotoKey != "" {
		downloadURL, err := s.storage.GenerateDownloadURL(ctx, *user.PhotoKey)
		if err == nil {
			user.PhotoURL = &downloadURL
		}
	}

	return toUserResponse(user), nil
}

func (s *service) PresignProfilePicture(ctx context.Context, userID int64, req PresignProfilePictureRequest) (*PresignProfilePictureResponse, error) {
	contentType := strings.ToLower(strings.TrimSpace(req.ContentType))
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" && contentType != "image/gif" {
		return nil, errors.New("invalid file type, allowed: image/jpeg, image/png, image/webp, image/gif")
	}

	ext := ".jpg"
	if contentType == "image/png" {
		ext = ".png"
	} else if contentType == "image/webp" {
		ext = ".webp"
	} else if contentType == "image/gif" {
		ext = ".gif"
	}

	key := fmt.Sprintf("ExpenseTracker/Users/%d/Profile/profile_%d%s", userID, time.Now().Unix(), ext)

	sizeLimitBytes := int64(s.config.ProfileImageMaxSizeMB) * 1024 * 1024
	if sizeLimitBytes <= 0 {
		sizeLimitBytes = 10 * 1024 * 1024 // Default 10MB
	}

	uploadURL, err := s.storage.GenerateUploadURL(ctx, key, contentType)
	if err != nil {
		return nil, fmt.Errorf("could not generate pre-signed upload URL: %w", err)
	}

	return &PresignProfilePictureResponse{
		UploadURL: uploadURL,
		ObjectKey: key,
	}, nil
}

func (s *service) CompleteProfilePicture(ctx context.Context, userID int64, req CompleteProfilePictureRequest) (*UserResponse, error) {
	objectKey := strings.TrimSpace(req.ObjectKey)
	if objectKey == "" {
		return nil, errors.New("object_key is required")
	}

	expectedPrefix := fmt.Sprintf("ExpenseTracker/Users/%d/Profile/", userID)
	if !strings.HasPrefix(objectKey, expectedPrefix) {
		return nil, errors.New("unauthorized: invalid object key destination")
	}
	exists, size, err := s.storage.ObjectExists(ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("error verifying object on storage: %w", err)
	}
	if !exists {
		return nil, errors.New("uploaded profile picture not found on storage, upload may have failed")
	}

	maxSizeBytes := int64(s.config.ProfileImageMaxSizeMB) * 1024 * 1024
	if size > maxSizeBytes {
		_ = s.storage.DeleteObject(ctx, objectKey)
		return nil, fmt.Errorf("uploaded file size exceeds maximum limit of %d MB", s.config.ProfileImageMaxSizeMB)
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	oldPhotoKey := ""
	if user.PhotoKey != nil {
		oldPhotoKey = *user.PhotoKey
	}

	newPhotoURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.AWSS3Bucket, s.config.AWSRegion, objectKey)

	err = s.repo.UpdateProfilePictureMetadata(ctx, userID, newPhotoURL, objectKey, size)
	if err != nil {
		return nil, err
	}

	if oldPhotoKey != "" && oldPhotoKey != objectKey {
		go func(key string) {
			deleteCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if delErr := s.storage.DeleteObject(deleteCtx, key); delErr != nil {
				log.Printf("[WARNING] Could not delete old S3 object key %s: %v", key, delErr)
			} else {
				log.Printf("[INFO] Successfully deleted old S3 profile object key %s", key)
			}
		}(oldPhotoKey)
	}

	user.PhotoURL = &newPhotoURL
	user.PhotoKey = &objectKey
	user.PhotoSize = &size

	if user.PhotoKey != nil && *user.PhotoKey != "" {
		downloadURL, err := s.storage.GenerateDownloadURL(ctx, *user.PhotoKey)
		if err == nil {
			user.PhotoURL = &downloadURL
		}
	}

	return toUserResponse(user), nil
}

func toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:                    user.ID,
		Email:                 user.Email,
		Name:                  user.Name,
		PhotoURL:              user.PhotoURL,
		PhotoKey:              user.PhotoKey,
		PhotoSize:             user.PhotoSize,
		Bio:                   user.Bio,
		IsEmailVerified:       user.IsEmailVerified,
		IsTransactionsEnabled: user.IsTransactionsEnabled,
		IsReportingEnabled:    user.IsReportingEnabled,
		IsNotificationEnabled: user.IsNotificationEnabled,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
	}
}