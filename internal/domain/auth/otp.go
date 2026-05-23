package domain_auth

import (
	"crypto/rand"
	"math/big"
)

func GenerateNumericOTP(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

func HashOTP(otp string) (string, error) {
    return HashPassword(otp)
}

func CheckOTPHash(otp, hash string) bool {
    return CheckPasswordHash(otp, hash)
}