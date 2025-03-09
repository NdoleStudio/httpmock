package middlewares

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// RequestRouter handles requests to the server
func RequestRouter(
	tracer telemetry.Tracer,
	logger telemetry.Logger,
	hostname string,
	requestService *services.ProjectEndpointRequestService,
	serverHandler fiber.Handler,
	reflectHandler fiber.Handler,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stopwatch := time.Now()

		ctx, span, ctxLogger := tracer.StartFromFiberCtxWithLogger(c, logger.WithCodeNamespace("middlewares.RequestRouter"), "middlewares.RequestRouter")
		defer span.End()

		if len(c.Subdomains()) > 1 {
			ctxLogger.Info(fmt.Sprintf("redirecting HTTP request [%s] -> [%s://%s] since it has more than 1 subdomains [%#+v]", c.BaseURL()+c.OriginalURL(), c.Protocol(), hostname, c.Subdomains()))
			return c.Redirect(fmt.Sprintf("%s://%s", c.Protocol(), hostname), fiber.StatusMovedPermanently)
		}

		if c.Hostname() == hostname || len(c.Subdomains()) == 0 {
			ctxLogger.Info(fmt.Sprintf("handling request with hostname [%s] using the default router", c.Hostname()))
			return c.Next()
		}

		if len(c.Subdomains()[0]) < 8 {
			return handleNamedSubdomains(c, strings.TrimSpace(c.Subdomains()[0]), serverHandler, reflectHandler)
		}

		endpoint, err := requestService.LoadByRequest(ctx, c.Subdomains()[0], c.Method(), c.Path())
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

func handleNamedSubdomains(c *fiber.Ctx, subdomain string, serverHandler fiber.Handler, reflectHandler fiber.Handler) error {
	switch subdomain {
	case "reflect":
		return reflectHandler(c)
	case "server":
		return serverHandler(c)
	}

	if status, err := strconv.Atoi(subdomain); err == nil && status >= 100 && status <= 599 {
		c.Request().Header.Set("response-status", subdomain)
		return serverHandler(c)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "error",
		"message": fmt.Sprintf("We cannot find a registered mock for URL [%s] and HTTP method [%s]", c.BaseURL()+c.OriginalURL(), c.Method()),
	})
}
