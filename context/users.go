package context

import (
	"context"

	"github.com/simon-lentz/webapp/models"
)

type key int

const (
	userKey key = iota
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		return nil // nil != nil edge case
	}
	return user
}
