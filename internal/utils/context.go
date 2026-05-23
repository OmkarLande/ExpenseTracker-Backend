package utils

import (
	"context"

	"ExpenseTracker-Backend/internal/types"
)

func GetUserIDFromContext(
	ctx context.Context,
) (int64, bool) {

	userID, ok := ctx.Value(
		types.UserIDContextKey,
	).(int64)

	return userID, ok
}