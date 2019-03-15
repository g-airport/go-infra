package context

import (
	"context"

	. "github.com/g-airport/go-infra/auth"

	"github.com/google/uuid"
	"github.com/micro/go-micro/metadata"
)

const (
	RequestID = "RequestID"
	Token     = "Token"
)

func SetUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, "user", u)
}

func GetUser(ctx context.Context) *User {
	return ctx.Value("user").(*User)
}

func NewContext(token string) context.Context {
	ctx := metadata.NewContext(context.Background(), map[string]string{
		Token: token,
	})
	return ctx
}

func NewUUID() string {
	u, _ := uuid.NewUUID()
	return u.String()
}

func RequestIDFromContext(ctx context.Context) string {
	if s, ok := ctx.Value(RequestID).(string); ok {
		return s
	}
	return ""
}

func WithRequestID(ctx context.Context, requestID ...string) context.Context {
	if len(requestID) != 0 {
		return context.WithValue(ctx, RequestID, requestID[0])
	}
	return context.WithValue(ctx, RequestID, NewUUID())
}
