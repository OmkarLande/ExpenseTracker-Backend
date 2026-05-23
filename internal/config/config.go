package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string

	DBUrl string

	JWTSecret string

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	AWSRegion             string
	AWSAccessKeyID         string
	AWSSecretAccessKey     string
	AWSS3Bucket           string
	ProfileImageMaxSizeMB int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
	}

	maxSizeMB := 10 // Default 10 MB
	if sizeStr := os.Getenv("PROFILE_IMAGE_MAX_SIZE_MB"); sizeStr != "" {
		if val, err := strconv.Atoi(sizeStr); err == nil {
			maxSizeMB = val
		}
	}

	cfg := &Config{
		AppEnv:                os.Getenv("APP_ENV"),
		Port:                  os.Getenv("PORT"),
		DBUrl:                 os.Getenv("DB_URL"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		SMTPHost:              os.Getenv("SMTP_HOST"),
		SMTPPort:              os.Getenv("SMTP_PORT"),
		SMTPUser:              os.Getenv("SMTP_USER"),
		SMTPPassword:          os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:              os.Getenv("SMTP_FROM"),
		AWSRegion:             os.Getenv("AWS_REGION"),
		AWSAccessKeyID:         os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:     os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSS3Bucket:           os.Getenv("AWS_S3_BUCKET"),
		ProfileImageMaxSizeMB: maxSizeMB,
	}

	return cfg, nil
}