package main

import (
	"context"
	"fmt"
)

func main() {
    // get an empty context
    ctx := context.Background()
    // override it with a new value
    ctx = context.WithValue(ctx, "favourite-colour", "blue")
    // get the value from the context
    v := ctx.Value("favourite-colour")
    fmt.Println(v)
}

