package domain_user

import "time"

type User struct {
	ID                    int64
	Email                 string
	PasswordHash          string
	Name                  string
	PhotoURL              *string
	PhotoKey              *string
	PhotoSize             *int64
	Bio                   *string
	IsEmailVerified       bool
	IsTransactionsEnabled bool
	IsReportingEnabled    bool
	IsNotificationEnabled bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}