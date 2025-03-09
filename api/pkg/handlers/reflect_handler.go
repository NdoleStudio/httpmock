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

// ReflectHandler handles requests to the /reflect* endpoints
type ReflectHandler struct {
	handler
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewReflectHandler creates a new ServerHandler
func NewReflectHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (h *ReflectHandler) {
	return &ReflectHandler{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", h)),
		tracer: tracer,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ReflectHandler) RegisterRoutes(app *fiber.App) {
	router := app.Group("/reflect")
	router.All("/*", h.Handle)
}

// Handle handles the request
func (h *ReflectHandler) Handle(c *fiber.Ctx) error {
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

func (h *ReflectHandler) getResponseStatus(ctxLogger telemetry.Logger, c *fiber.Ctx) int {
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

func (h *ReflectHandler) getResponseHeaders(c *fiber.Ctx) []map[string]string {
	var headers []map[string]string
	for key, value := range c.GetReqHeaders() {
		for _, header := range value {
			headers = append(headers, map[string]string{key: header})
		}
	}
	return headers
}
