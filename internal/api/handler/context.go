package handler

import (
	"context"

	"github.com/Facundo-Mourelle/go-gym/internal/api/middleware"
)

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}
