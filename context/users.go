package context

import (
	"context"

	"github.com/rafaelmdurante/lenslocked/models"
)

type key string

const (
    userKey key = "user"
)

// WithUser function stores a User in the Context.
func WithUser(ctx context.Context, user *models.User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

// User function retrieves the User from the Context. If casting fails or there
// is no User, it returns nil.
func User(ctx context.Context) *models.User {
    val := ctx.Value(userKey)
    // cast the context to User model
    user, ok := val.(*models.User)
    if !ok {
        // the most likely is that nothing was ever stored in the context,
        // so it doesn't have a type of *models.User 
        // it is also possible that other code in this package wrote an invalid
        // value using the user key, so it is important to review code changes
        // in this package
        return nil
    }
    return user
}

