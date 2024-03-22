package main

import (
	stdctx "context"
	"fmt"

	"github.com/rafaelmdurante/lenslocked/context"
	"github.com/rafaelmdurante/lenslocked/models"
)

func main() {
    ctx := stdctx.Background()

    user := models.User{
        Email: "raf@ael.com",
    }

    ctx = context.WithUser(ctx, &user)

    retrievedUser := context.User(ctx)

    fmt.Println(retrievedUser.Email)
}

