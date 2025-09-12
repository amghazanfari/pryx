package auth

import "context"

type ContextKey string

const CtxKey ContextKey = "authctx"

type AuthContext struct {
	UserID   uint
	APIKeyID uint
	Scopes   string
}

func WithAuth(ctx context.Context, a AuthContext) context.Context {
	return context.WithValue(ctx, CtxKey, a)
}

func From(ctx context.Context) (AuthContext, bool) {
	v := ctx.Value(CtxKey)
	if v == nil {
		return AuthContext{}, false
	}
	ac, ok := v.(AuthContext)
	return ac, ok
}
