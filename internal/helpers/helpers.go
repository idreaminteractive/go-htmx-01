package helpers

import (
	"context"
	"fmt"
	"net/http"
)

func ContextGet(r *http.Request, key string) (interface{}, error) {
	val := r.Context().Value(key)
	if val == nil {
		return nil, fmt.Errorf("no value exists in the context for key %q", key)
	}
	return val, nil
}

func ContextSave(r *http.Request, key string, val interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, key, val) // nolint:staticcheck
	return r.WithContext(ctx)
}
