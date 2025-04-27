package server

import "context"

// PreHookFunc is called before executing a storage operation.
// If it returns an error, the operation is aborted.
type PreHookFunc func(ctx context.Context, method string, req interface{}) error

// PostHookFunc is called after successfully executing a storage operation.
// It receives both the request and the response.
// If it returns an error, it does not affect the original operation.
type PostHookFunc func(ctx context.Context, method string, req interface{}, resp interface{}) error
