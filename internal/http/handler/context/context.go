package context

import (
	"context"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("userID")

func ContextSetUserID(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, userID)
	return r.WithContext(ctx)
}

func ContextGetUserID(r *http.Request) int64 {
	userID, ok := r.Context().Value(userContextKey).(int64)
	if !ok {
		panic("missing user value in request context")
	}

	return userID
}
