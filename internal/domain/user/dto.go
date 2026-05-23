package domain_user

import (
	"time"
)

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	Bio      *string `json:"bio"`
	PhotoURL *string `json:"photo_url"`
}

type UserResponse struct {
	ID                    int64     `json:"id"`
	Email                 string    `json:"email"`
	Name                  string    `json:"name"`
	PhotoURL              *string   `json:"photo_url"`
	PhotoKey              *string   `json:"photo_key"`
	PhotoSize             *int64    `json:"photo_size"`
	Bio                   *string   `json:"bio"`
	IsEmailVerified       bool      `json:"is_email_verified"`
	IsTransactionsEnabled bool      `json:"is_transactions_enabled"`
	IsReportingEnabled    bool      `json:"is_reporting_enabled"`
	IsNotificationEnabled bool      `json:"is_notification_enabled"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type PresignProfilePictureRequest struct {
	ContentType string `json:"content_type" validate:"required"`
}

type PresignProfilePictureResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

type CompleteProfilePictureRequest struct {
	ObjectKey string `json:"object_key" validate:"required"`
}