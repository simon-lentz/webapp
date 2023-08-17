package main

import (
	"context"
	"fmt"
)

type ctxKey string

const (
	favoriteColorKey ctxKey = "favorite-color"
)

func main() {
	// Using an unexported custom type rather than a
	// builtin as a context key is important because
	// if multiple packages are accessing the context
	// the key could be overwritten. The custom type
	// helps to prevent this.
	ctx := context.Background()
	ctx = context.WithValue(ctx, favoriteColorKey, "blue")
	value := ctx.Value(favoriteColorKey)
	fmt.Println(value)

	// Without the custom type...
	// We set the color
	ctx = context.WithValue(ctx, "favorite-color", "blue")
	// Another package sets a color, overwriting our favorite color.
	ctx = context.WithValue(ctx, "favorite-color", "red")
	fmt.Println(ctx)
}
