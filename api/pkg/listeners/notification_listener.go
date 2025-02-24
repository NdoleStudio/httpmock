package listeners

import (
	"context"
	"fmt"

	"github.com/pusher/pusher-http-go/v5"

	"github.com/NdoleStudio/httpmock/pkg/events"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// NotificationListener listens for cloudevents.Events
type NotificationListener struct {
	logger       telemetry.Logger
	tracer       telemetry.Tracer
	pusherClient *pusher.Client
}

// NewNotificationListener creates a new NotificationListener
func NewNotificationListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	pusherClient *pusher.Client,
) *NotificationListener {
	return &NotificationListener{
		logger:       logger.WithCodeNamespace(fmt.Sprintf("%T", &NotificationListener{})),
		tracer:       tracer,
		pusherClient: pusherClient,
	}
}

// Register the listener to the dispatcher
func (listener *NotificationListener) Register(dispatcher *services.EventDispatcher) {
	dispatcher.Subscribe(events.ProjectEndpointRequest, listener.onProjectEndpointRequest)
}

func (listener *NotificationListener) onProjectEndpointRequest(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.ProjectEndpointRequestPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.pusherClient.Trigger(payload.UserID.String(), event.Type(), event); err != nil {
		msg := fmt.Sprintf("cannot send real time notification for [%s] event with ID [%s] and user ID [%s]", event.Type(), event.ID(), payload.UserID)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}
