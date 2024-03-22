package main

import (
	"context"
	"fmt"
)

// make it private so only code within your package can use this
type ctxKey string

const (
    // make it private so only code within your package can use this
    favouriteColour ctxKey = "favourite-colour"
)

func main() {
    // get an empty context
    ctx := context.Background()

    // having a custom type of ctxKey makes it different from the same "key" of
    // pure string
    ctx = context.WithValue(ctx, favouriteColour, "blue")
    ctx = context.WithValue(ctx, "favourite-colour", "red")

    // get the value from the context
    fmt.Println(
        ctx.Value(favouriteColour),
        ctx.Value("favourite-colour"),
    )
}

