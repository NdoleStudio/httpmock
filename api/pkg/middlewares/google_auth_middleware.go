package middlewares

import (
	"context"
	"fmt"
	"strings"

	"github.com/palantir/stacktrace"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

// GoogleAuth authenticates a user based on the bearer token
func GoogleAuth(logger telemetry.Logger, tracer telemetry.Tracer, webhookEndpoint string, authEmail string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span, ctxLogger := tracer.StartFromFiberCtxWithLogger(c, logger, "middlewares.GoogleAuth")
		defer span.End()

		authToken := c.Get(authHeaderBearer)
		if !strings.HasPrefix(authToken, bearerScheme) {
			span.AddEvent(fmt.Sprintf("The request header has no [%s] token", bearerScheme))
			return c.Next()
		}

		if len(authToken) > len(bearerScheme)+1 {
			authToken = authToken[len(bearerScheme)+1:]
		}

		payload, err := idtoken.Validate(context.Background(), authToken, webhookEndpoint)
		if err != nil {
			msg := fmt.Sprintf("invalid google auto token [%s]", authToken)
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
			return c.Next()
		}

		email, ok := payload.Claims["email"].(string)
		if !ok {
			msg := fmt.Sprintf("cannot cast email from [%s] for google auth token claims [%s]", payload.Claims["email"], tracer.Redact(authToken))
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.NewError(msg)))
			return c.Next()
		}

		if email != authEmail {
			msg := fmt.Sprintf("invalid email [%s] for google auth token [%s]", email, tracer.Redact(authToken))
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.NewError(msg)))
			return c.Next()
		}

		span.AddEvent(fmt.Sprintf("[%s] google auth token is valid", bearerScheme))

		authUser := &entities.AuthUser{
			Email: email,
			ID:    entities.UserID(payload.Claims["sub"].(string)),
		}

		c.Locals(ContextKeyAuthUserID, authUser)

		ctxLogger.Info(fmt.Sprintf("[%T] set successfully for user with ID [%s]", authUser, authUser.ID))
		return c.Next()
	}
}
