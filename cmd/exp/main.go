package main

import (
	stdctx "context"
	"fmt"

	"github.com/simon-lentz/webapp/context"
	"github.com/simon-lentz/webapp/models"
)

func main() {
	ctx := stdctx.Background()

	user := models.User{
		Email: "test@testing.com",
	}

	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)

	fmt.Println(retrievedUser.Email)

}
