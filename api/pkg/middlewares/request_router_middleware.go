package middlewares

import (
	"fmt"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// RequestRouter handles requests to the server
func RequestRouter(tracer telemetry.Tracer, logger telemetry.Logger, hostname string, requestService *services.ProjectEndpointRequestService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stopwatch := time.Now()

		ctx, span, ctxLogger := tracer.StartFromFiberCtxWithLogger(c, logger.WithCodeNamespace("middlewares.RequestRouter"), "middlewares.RequestRouter")
		defer span.End()

		ctxLogger.Info(fmt.Sprintf("handling request with hostname [%#+v]", c.Subdomains()))

		if c.Hostname() == hostname || len(c.Subdomains()) == 0 {
			ctxLogger.Info(fmt.Sprintf("handling request with hostname [%s] using the default router", c.Hostname()))
			return c.Next()
		}

		endpoint, err := requestService.FetchEndpoint(ctx, c.Subdomains()[0], c.Method(), c.Path())
		if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("We cannot find a registered mock for URL [%s] and HTTP method [%s]", c.BaseURL()+c.OriginalURL(), c.Method()),
			})
		}

		if err != nil {
			msg := fmt.Sprintf("error while fetching endpoint [%s] with method [%s]", c.BaseURL()+c.OriginalURL(), c.Method())
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "We ran into an internal server error occurred while processing your request. We have been notified about it it already.",
			})
		}

		requestService.HandleHTTPRequest(ctx, c, stopwatch, endpoint)
		return nil
	}
}
