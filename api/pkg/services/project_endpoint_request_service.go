package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/events"
	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointRequestService is responsible for managing entities.ProjectEndpointRequest
type ProjectEndpointRequestService struct {
	service
	logger                           telemetry.Logger
	tracer                           telemetry.Tracer
	projectEndpointRequestRepository repositories.ProjectEndpointRequestRepository
	projectEndpointRepository        repositories.ProjectEndpointRepository
	eventDispatcher                  *EventDispatcher
}

// NewProjectEndpointRequestService creates a new ProjectEndpointRequestService
func NewProjectEndpointRequestService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	projectEndpointRepository repositories.ProjectEndpointRepository,
	projectEndpointRequestRepository repositories.ProjectEndpointRequestRepository,
	eventDispatcher *EventDispatcher,
) (s *ProjectEndpointRequestService) {
	return &ProjectEndpointRequestService{
		logger:                           logger.WithCodeNamespace(fmt.Sprintf("%T", s)),
		tracer:                           tracer,
		projectEndpointRepository:        projectEndpointRepository,
		eventDispatcher:                  eventDispatcher,
		projectEndpointRequestRepository: projectEndpointRequestRepository,
	}
}

// Delete an entities.ProjectEndpointRequest
func (service *ProjectEndpointRequestService) Delete(ctx context.Context, userID entities.UserID, requestID ulid.ULID) error {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	request, err := service.projectEndpointRequestRepository.Load(ctx, userID, requestID)
	if err != nil {
		msg := fmt.Sprintf("cannot load endpoint request with ID [%s] for user ID [%s]", requestID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err = service.projectEndpointRequestRepository.Delete(ctx, request); err != nil {
		msg := fmt.Sprintf("cannot delete endpoint with ID [%s] and project ID [%s] for user ID [%s]", request.ID, request.ProjectID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	return nil
}

// Index fetches the list of all project endpoint requests available to the currently authenticated user
func (service *ProjectEndpointRequestService) Index(ctx context.Context, userID entities.UserID, endpointID uuid.UUID, limit uint, previousID *ulid.ULID) ([]*entities.ProjectEndpointRequest, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	requests, err := service.projectEndpointRequestRepository.Index(ctx, userID, endpointID, limit, previousID)
	if err != nil {
		msg := fmt.Sprintf("cannot fetch project endpoint requests for user with ID [%s] and project endpoint ID [%s]", userID, endpointID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return requests, nil
}

// HandleHTTPRequest registers a new HTTP request for an endpoint
func (service *ProjectEndpointRequestService) HandleHTTPRequest(ctx context.Context, c *fiber.Ctx, stopwatch time.Time, endpoint *entities.ProjectEndpoint) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	requestID := ulid.Make()

	service.storeProjectEndpointRequestEvent(ctx, requestID, stopwatch, c, endpoint)
	headers := service.getHTTPHeaders(ctxLogger, c, endpoint)

	delay := endpoint.ResponseDelayInMilliseconds - uint(time.Since(stopwatch).Milliseconds())
	if endpoint.ResponseDelayInMilliseconds > 0 && delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	ctxLogger.Debug(fmt.Sprintf("finished handling request with URL [%s] in [%s] and request ID [%s]", c.BaseURL()+c.OriginalURL(), time.Since(stopwatch).String(), requestID))
	for _, header := range headers {
		for key, value := range header {
			c.Response().Header.Set(key, value)
		}
	}

	c.Response().SetStatusCode(int(endpoint.ResponseCode))

	if endpoint.ResponseBody != nil {
		if _, err := c.Response().BodyWriter().Write([]byte(*endpoint.ResponseBody)); err != nil {
			msg := fmt.Sprintf("error while writing response body for request [%s] with method [%s] and request ID [%s]", c.BaseURL()+c.OriginalURL(), c.Method(), requestID)
			ctxLogger.Error(service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		}
	}
}

// LoadByRequest a project endpoint by request method and path
func (service *ProjectEndpointRequestService) LoadByRequest(ctx context.Context, subdomain string, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	return service.projectEndpointRepository.LoadByRequest(ctx, subdomain, requestMethod, requestPath)
}

// Store a project endpoint request
func (service *ProjectEndpointRequestService) Store(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	if err := service.projectEndpointRequestRepository.Store(ctx, request); err != nil {
		msg := fmt.Sprintf("cannot save [%T] with ID [%s] for project with ID [%s] and user with ID [%s]", request, request.ID, request.ProjectID, request.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := service.projectEndpointRepository.RegisterRequest(ctx, request.ProjectEndpointID); err != nil {
		msg := fmt.Sprintf("cannot register request for [%T] with ID [%s] for project with ID [%s] and user with ID [%s]", &entities.ProjectEndpoint{}, request.ProjectEndpointID, request.ProjectID, request.UserID)
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (service *ProjectEndpointRequestService) getHTTPHeaders(ctxLogger telemetry.Logger, c *fiber.Ctx, endpoint *entities.ProjectEndpoint) []map[string]string {
	var headers []map[string]string

	if endpoint.ResponseHeaders == nil || *endpoint.ResponseHeaders == "" {
		return headers
	}

	if err := json.Unmarshal([]byte(*endpoint.ResponseHeaders), &headers); err != nil {
		msg := fmt.Sprintf("error while unmarshalling response headers [%s] for request [%s] with method [%s]", *endpoint.ResponseHeaders, c.BaseURL()+c.OriginalURL(), c.Method())
		ctxLogger.Error(stacktrace.Propagate(err, msg))
	}

	return headers
}

func (service *ProjectEndpointRequestService) storeProjectEndpointRequestEvent(
	ctx context.Context,
	requestID ulid.ULID,
	stopwatch time.Time,
	c *fiber.Ctx,
	endpoint *entities.ProjectEndpoint,
) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	requestBody := c.Body()
	if len(requestBody) == 0 {
		requestBody = nil
	}

	source := c.BaseURL() + c.OriginalURL()
	event, err := service.createEvent(events.ProjectEndpointRequest, source, &events.ProjectEndpointRequestPayload{
		UserID:                      endpoint.UserID,
		ProjectID:                   endpoint.ProjectID,
		ProjectEndpointID:           endpoint.ID,
		ProjectEndpointRequestID:    requestID,
		RequestURL:                  source,
		RequestMethod:               c.Method(),
		RequestBody:                 service.getRequestBody(c),
		RequestHeaders:              service.getRequestHeaders(ctxLogger, c),
		ResponseCode:                endpoint.ResponseCode,
		ResponseBody:                endpoint.ResponseBody,
		ResponseHeaders:             endpoint.ResponseHeaders,
		ResponseDelayInMilliseconds: endpoint.ResponseDelayInMilliseconds,
		IPAddress:                   c.IP(),
		Timestamp:                   stopwatch,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for  project endpiont request with ID [%s]", events.ProjectEndpointRequest, requestID)
		ctxLogger.Error(service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		return
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project endpiont request with ID [%s]", event.Type(), requestID)
		ctxLogger.Error(service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
	}

	return
}

func (service *ProjectEndpointRequestService) getRequestHeaders(ctxLogger telemetry.Logger, c *fiber.Ctx) *string {
	if len(c.GetReqHeaders()) == 0 {
		return nil
	}

	var headers []map[string]string
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			headers = append(headers, map[string]string{key: value})
		}
	}
	result, err := json.Marshal(headers)
	if err != nil {
		ctxLogger.Error(stacktrace.Propagate(err, fmt.Sprintf("error while marshalling request headers [%s]", headers)))
	}

	resultString := string(result)
	return &resultString
}

func (service *ProjectEndpointRequestService) getRequestBody(c *fiber.Ctx) *string {
	if len(c.Body()) == 0 {
		return nil
	}
	requestBody := string(c.Body())
	return &requestBody
}
