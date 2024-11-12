package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/events"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"

	"github.com/NdoleStudio/httpmock/pkg/queue"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// EventDispatcher dispatches a new event
type EventDispatcher struct {
	logger      telemetry.Logger
	tracer      telemetry.Tracer
	queue       queue.Client
	meter       metric.Float64Histogram
	consumerURL string
	listeners   map[string][]events.EventListener
}

// NewEventDispatcher creates a new EventDispatcher
func NewEventDispatcher(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	meter metric.Float64Histogram,
	queue queue.Client,
	consumerURL string,
) (dispatcher *EventDispatcher) {
	return &EventDispatcher{
		logger:      logger,
		listeners:   make(map[string][]events.EventListener),
		tracer:      tracer,
		meter:       meter,
		consumerURL: consumerURL,
		queue:       queue,
	}
}

// Dispatch a new event by adding it to the queue to be processed async
func (dispatcher *EventDispatcher) Dispatch(ctx context.Context, event *cloudevents.Event) error {
	ctx, span := dispatcher.tracer.Start(ctx)
	defer span.End()

	ctxLogger := dispatcher.tracer.CtxLogger(dispatcher.logger, span)

	if err := event.Validate(); err != nil {
		msg := fmt.Sprintf("cannot dispatch event with ID [%s] and type [%s] because it is invalid", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	task, err := dispatcher.createTask(event)
	if err != nil {
		msg := fmt.Sprintf("cannot create push queue task for event with ID [%s] and type [%s]", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	taskID, err := dispatcher.queue.Enqueue(ctx, task)
	if err != nil {
		msg := fmt.Sprintf("cannot add event with ID [%s] and type [%s] to producer", event.ID(), event.Type())
		return dispatcher.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	ctxLogger.Info(fmt.Sprintf("push queue task enqueued with ID [%s]", taskID))
	return nil
}

// Subscribe a listener to an event
func (dispatcher *EventDispatcher) Subscribe(eventType string, listener events.EventListener) {
	if _, ok := dispatcher.listeners[eventType]; !ok {
		dispatcher.listeners[eventType] = []events.EventListener{}
	}

	dispatcher.listeners[eventType] = append(dispatcher.listeners[eventType], listener)
}

// Publish an event to subscribers
func (dispatcher *EventDispatcher) Publish(ctx context.Context, event cloudevents.Event) {
	ctx, span := dispatcher.tracer.Start(ctx)
	defer span.End()

	start := time.Now()

	ctxLogger := dispatcher.tracer.CtxLogger(dispatcher.logger, span)

	subscribers, ok := dispatcher.listeners[event.Type()]
	if !ok {
		ctxLogger.Info(fmt.Sprintf("no listener is configured for event type [%s] with id [%s]", event.Type(), event.ID()))
		return
	}

	var wg sync.WaitGroup
	for _, sub := range subscribers {
		wg.Add(1)
		go func(ctx context.Context, sub events.EventListener) {
			if err := sub(ctx, event); err != nil {
				msg := fmt.Sprintf("subscriber [%T] cannot handle event [%s]", sub, event.Type())
				ctxLogger.Error(stacktrace.Propagate(err, msg))
			}
			wg.Done()
		}(ctx, sub)
	}

	wg.Wait()

	dispatcher.meter.Record(
		ctx,
		float64(time.Since(start).Microseconds())/1000,
		metric.WithAttributes(
			semconv.CloudeventsEventType(event.Type()),
			semconv.CloudeventsEventSpecVersion(event.SpecVersion()),
		),
	)
}

func (dispatcher *EventDispatcher) createTask(event *cloudevents.Event) (*queue.Task, error) {
	eventContent, err := json.Marshal(event)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("cannot marshall [%T] with ID [%s]", event, event.ID()))
	}

	return &queue.Task{
		Method: http.MethodPost,
		URL:    dispatcher.consumerURL,
		Body:   eventContent,
	}, nil
}
