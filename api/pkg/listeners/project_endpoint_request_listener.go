package listeners

import (
	"context"
	"fmt"
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/events"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointRequestListener listens for events.ProjectEndpointRequest events
type ProjectEndpointRequestListener struct {
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.ProjectEndpointRequestService
}

// NewProjectEndpointRequestListener creates a new ProjectEndpointRequestListener
func NewProjectEndpointRequestListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.ProjectEndpointRequestService,
) *ProjectEndpointRequestListener {
	return &ProjectEndpointRequestListener{
		logger:  logger.WithCodeNamespace(fmt.Sprintf("%T", &ProjectEndpointRequestListener{})),
		tracer:  tracer,
		service: service,
	}
}

// Register the listener to the dispatcher
func (listener *ProjectEndpointRequestListener) Register(dispatcher *services.EventDispatcher) {
	dispatcher.Subscribe(events.ProjectEndpointRequest, listener.onProjectEndpointRequest)
}

func (listener *ProjectEndpointRequestListener) onProjectEndpointRequest(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.ProjectEndpointRequestPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	request := &entities.ProjectEndpointRequest{
		ID:                          strings.ToLower(payload.ProjectEndpointRequestID.String()),
		ProjectID:                   payload.ProjectID,
		ProjectEndpointID:           payload.ProjectEndpointID,
		UserID:                      payload.UserID,
		RequestMethod:               payload.RequestMethod,
		RequestURL:                  payload.RequestURL,
		RequestHeaders:              payload.RequestHeaders,
		RequestBody:                 payload.RequestBody,
		ResponseCode:                payload.ResponseCode,
		ResponseBody:                payload.ResponseBody,
		ResponseHeaders:             payload.ResponseHeaders,
		ResponseDelayInMilliseconds: payload.ResponseDelayInMilliseconds,
		CreatedAt:                   payload.Timestamp,
	}

	if err := listener.service.Store(ctx, request); err != nil {
		msg := fmt.Sprintf("cannot store [%T] for [%s] event with ID [%s] and user ID [%s]", request, event.Type(), event.ID(), payload.UserID)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return nil
}
