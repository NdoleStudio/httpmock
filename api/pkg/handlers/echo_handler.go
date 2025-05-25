package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// EchoHandler handles requests to the /reflect* endpoints
type EchoHandler struct {
	handler
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewEchoHandler creates a new EchoHandler
func NewEchoHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (h *EchoHandler) {
	return &EchoHandler{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", h)),
		tracer: tracer,
	}
}

// RegisterRoutes registers the routes for the EchoHandler
func (h *EchoHandler) RegisterRoutes(app *fiber.App) {
	router := app.Group("/echo")
	router.All("/*", h.Handle)
}

// Handle handles the request
func (h *EchoHandler) Handle(c *fiber.Ctx) error {
	_, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	for _, header := range h.getResponseHeaders(c) {
		for key, value := range header {
			c.Response().Header.Set(key, value)
		}
	}

	c.Response().SetStatusCode(h.getResponseStatus(ctxLogger, c))

	if len(c.Body()) > 0 {
		if _, err := c.Response().BodyWriter().Write(c.Body()); err != nil {
			msg := fmt.Sprintf("error while writing response body [%s] for request [%s] with method [%s]", c.Body(), c.BaseURL()+c.OriginalURL(), c.Method())
			ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		}
	}

	return nil
}

func (h *EchoHandler) getResponseStatus(ctxLogger telemetry.Logger, c *fiber.Ctx) int {
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

func (h *EchoHandler) getResponseHeaders(c *fiber.Ctx) []map[string]string {
	var headers []map[string]string
	for key, value := range c.GetReqHeaders() {
		for _, header := range value {
			headers = append(headers, map[string]string{key: header})
		}
	}
	return headers
}
