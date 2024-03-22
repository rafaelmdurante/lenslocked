package main

import (
	"context"
	"fmt"
	"strings"
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

    // get the value from the context
    value := ctx.Value(favouriteColour)


    // example of casting the value to its original type - ie string
    str, ok := value.(string)
    if !ok {
        fmt.Println("not a string")
    } else {
        // it won't compile strings.HasPrefix(value, "b") because 'value' is
        // of type 'any'
        // this can only be performed on strings
        fmt.Println(strings.HasPrefix(str, "b")) // output: true
    }

    // print both variables with the same value
    fmt.Println(value, str) // output: blue blue
}

