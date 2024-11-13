package main

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/di"
	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/davecgh/go-spew/spew"
	"github.com/palantir/stacktrace"
)

func main() {
	di.LoadEnv("../../.env")

	container := di.NewLiteContainer()
	logger := container.Logger()
	tracer := container.Tracer()

	token := ""
	claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{Token: token})
	if err != nil {
		msg := fmt.Sprintf("invalid clerk id token [%s]", tracer.Redact(token))
		logger.Fatal(stacktrace.Propagate(err, msg))
	}

	u, err := user.Get(context.Background(), claims.Subject)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch user with ID [%s]", claims.Subject)
		logger.Fatal(stacktrace.Propagate(err, msg))
	}

	spew.Dump(u)

	authUser := entities.AuthUser{
		Email:     u.EmailAddresses[0].EmailAddress,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ID:        entities.UserID(u.ID),
	}

	spew.Dump(authUser)
}
