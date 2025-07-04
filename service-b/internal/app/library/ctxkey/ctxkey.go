package ctxkey

import (
	"context"
)

type ctxKey int

const (
	keyAuthorizationHeader ctxKey = iota
	keyIdempotencyKey
	keyRequestID
)

func GetIdempotencyKey(ctx context.Context) (string, bool) {
	if s, ok := ctx.Value(keyIdempotencyKey).(string); ok {
		return s, true
	}

	return "", false
}

func PutIdempotencyKey(ctx context.Context, idempotencyKey string) context.Context {
	return context.WithValue(ctx, keyIdempotencyKey, idempotencyKey)
}

func GetRequestID(ctx context.Context) (string, bool) {
	if s, ok := ctx.Value(keyRequestID).(string); ok {
		return s, true
	}

	return "", false
}

func PutRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}
