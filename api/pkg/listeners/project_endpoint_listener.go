package listeners

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/events"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointListener listens for events.ProjectEndpointRequest events
type ProjectEndpointListener struct {
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.ProjectEndpointService
}

// NewProjectEndpointListener creates a new ProjectEndpointListener
func NewProjectEndpointListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.ProjectEndpointService,
) *ProjectEndpointListener {
	return &ProjectEndpointListener{
		logger:  logger.WithCodeNamespace(fmt.Sprintf("%T", &ProjectEndpointListener{})),
		tracer:  tracer,
		service: service,
	}
}

// Register the listener to the dispatcher
func (listener *ProjectEndpointListener) Register(dispatcher *services.EventDispatcher) {
	dispatcher.Subscribe(events.ProjectEndpointRequest, listener.onProjectUpdatedRequest)
}

func (listener *ProjectEndpointListener) onProjectUpdatedRequest(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.ProjectUpdatedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.UpdateProjectSubdomain(ctx, payload.ProjectID, payload.ProjectSubdomain); err != nil {
		msg := fmt.Sprintf("cannot update subdomain for [%s] event with ID [%s] and project ID [%s]", event.Type(), event.ID(), payload.ProjectID)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return nil
}
