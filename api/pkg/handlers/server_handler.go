package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ServerHandler handles requests to the /server* endpoints
type ServerHandler struct {
	handler
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewServerHandler creates a new ServerHandler
func NewServerHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (h *ServerHandler) {
	return &ServerHandler{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", h)),
		tracer: tracer,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ServerHandler) RegisterRoutes(app *fiber.App) {
	router := app.Group("/server")
	router.All("/*", h.Handle)
}

// Handle handles the request
func (h *ServerHandler) Handle(c *fiber.Ctx) error {
	_, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	stopwatch := time.Now()
	headers := h.getResponseHeaders(ctxLogger, c)

	delay := h.getResponseDelay(ctxLogger, c) - time.Since(stopwatch)
	if delay > 0 {
		time.Sleep(delay)
	}

	ctxLogger.Debug(fmt.Sprintf("finished handling request with URL [%s] in [%s]", c.BaseURL()+c.OriginalURL(), time.Since(stopwatch).String()))

	for _, header := range headers {
		for key, value := range header {
			c.Response().Header.Set(key, value)
		}
	}

	c.Response().SetStatusCode(h.getResponseStatus(ctxLogger, c))

	responseBody := strings.TrimSpace(c.Get("response-body"))
	if responseBody != "" {
		if _, err := c.Response().BodyWriter().Write([]byte(responseBody)); err != nil {
			msg := fmt.Sprintf("error while writing response body for request [%s] with method [%s]", c.BaseURL()+c.OriginalURL(), c.Method())
			ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		}
	}

	return nil
}

func (h *ServerHandler) getResponseHeaders(ctxLogger telemetry.Logger, c *fiber.Ctx) []map[string]string {
	var headers []map[string]string

	headersString := strings.TrimSpace(c.Get("response-headers"))
	if headersString == "" {
		return headers
	}

	if err := json.Unmarshal([]byte(headersString), &headers); err != nil {
		msg := fmt.Sprintf("Failed to parse headers [%s] into type [%T]", headersString, headers)
		ctxLogger.Warn(errors.New(msg))
		return headers
	}

	return headers
}

func (h *ServerHandler) getResponseDelay(ctxLogger telemetry.Logger, c *fiber.Ctx) time.Duration {
	delay, err := strconv.Atoi(strings.TrimSpace(c.Get("response-delay")))
	if err != nil || delay < 0 || delay > 100000 {
		msg := fmt.Sprintf("The request has an invalid 'response-delay' header [%s], defaulting to 0", c.Get("response-delay"))
		ctxLogger.Warn(errors.New(msg))
		return 0
	}
	return time.Duration(delay) * time.Millisecond
}

func (h *ServerHandler) getResponseStatus(ctxLogger telemetry.Logger, c *fiber.Ctx) int {
	status := strings.TrimSpace(c.Get("response-status"))
	if status == "" {
		return fiber.StatusOK
	}

	statusCode, err := strconv.Atoi(status)
	if err != nil || statusCode < 100 || statusCode > 599 {
		msg := fmt.Sprintf("The request has an invalid 'response-status' header [%s], defaulting to [%d]", c.Get("response-status"), fiber.StatusOK)
		ctxLogger.Warn(errors.New(msg))
		return fiber.StatusOK
	}
	return statusCode
}
