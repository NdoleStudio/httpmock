package middlewares

import (
	"fmt"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/palantir/stacktrace"
)

// ClerkBearerAuth authenticates a user based on the bearer token
func ClerkBearerAuth(logger telemetry.Logger, tracer telemetry.Tracer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span, ctxLogger := tracer.StartFromFiberCtxWithLogger(c, logger, "middlewares.ClerkBearerAuth")
		defer span.End()

		authToken := c.Get(authHeaderBearer)
		if !strings.HasPrefix(authToken, bearerScheme) {
			span.AddEvent(fmt.Sprintf("the request header has no [%s] token", bearerScheme))
			return c.Next()
		}

		if len(authToken) > len(bearerScheme)+1 {
			authToken = authToken[len(bearerScheme)+1:]
		}

		claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{Token: authToken})
		if err != nil {
			msg := fmt.Sprintf("invalid clerk id token [%s] and error message [%s]", tracer.Redact(authToken), err.Error())
			span.AddEvent(msg)
			return c.Next()
		}

		u, err := user.Get(c.Context(), claims.Subject)
		if err != nil {
			msg := fmt.Sprintf("cannot fetch user with ID [%s]", claims.Subject)
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
			return c.Next()
		}

		span.AddEvent(fmt.Sprintf("the clerk [%s] token [%s] is valid", bearerScheme, tracer.Redact(authToken)))

		authUser := &entities.AuthUser{
			Email:     u.EmailAddresses[0].EmailAddress,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			ID:        entities.UserID(u.ID),
		}

		c.Locals(ContextKeyAuthUserID, authUser)

		ctxLogger.Info(fmt.Sprintf("[%T] set successfully for user with ID [%s]", authUser, authUser.ID))
		return c.Next()
	}
}
