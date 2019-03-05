package context

import (
	"context"
	"testing"
)

func TestWithRequestID(t *testing.T) {
	reqID := NewUUID()

	ctx := WithRequestID(context.Background(), reqID)
	if reqID != RequestIDFromContext(ctx) {
		t.Errorf("RequestID %s is not equal to reqID from context %s",
			reqID, RequestIDFromContext(ctx))
	}
}

func TestGetRequestID(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx)
	ctx2, _ := context.WithCancel(ctx)

	if RequestIDFromContext(ctx) != RequestIDFromContext(ctx2) {
		t.Errorf("ctx req id %s not equal to ctx2 req id %s",
			RequestIDFromContext(ctx), RequestIDFromContext(ctx2))
	}
}
