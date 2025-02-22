package di

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/listeners"

	"github.com/NdoleStudio/httpmock/docs"

	"github.com/caarlos0/env/v11"
	"github.com/lmittmann/tint"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"

	"github.com/clerk/clerk-sdk-go/v2"

	"github.com/NdoleStudio/go-otelroundtripper"
	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/handlers"
	"github.com/NdoleStudio/httpmock/pkg/middlewares"
	"github.com/NdoleStudio/httpmock/pkg/queue"
	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/NdoleStudio/httpmock/pkg/validators"
	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/hashicorp/go-retryablehttp"

	otelMetric "go.opentelemetry.io/otel/metric"

	"github.com/gofiber/contrib/otelfiber"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/jinzhu/now"

	"github.com/uptrace/uptrace-go/uptrace"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"google.golang.org/api/option"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"

	"gorm.io/driver/postgres"
	gormLogger "gorm.io/gorm/logger"
)

var configuration *Configuration

// Container is used to resolve services at runtime
type Container struct {
	projectID       string
	version         string
	db              *gorm.DB
	app             *fiber.App
	eventDispatcher *services.EventDispatcher
	logger          telemetry.Logger
}

// NewLiteContainer creates a Container without any routes or listeners
func NewLiteContainer() (container *Container) {
	// Set location to UTC
	now.DefaultConfig = &now.Config{
		TimeLocation: time.UTC,
	}

	clerk.SetKey(os.Getenv("CLERK_API_KEY"))

	return &Container{
		logger: logger(3).WithCodeNamespace(fmt.Sprintf("%T", container)),
	}
}

// NewContainer creates a new dependency injection container
func NewContainer(projectID string, version string) (container *Container) {
	// Set location to UTC
	now.DefaultConfig = &now.Config{
		TimeLocation: time.UTC,
	}

	clerk.SetKey(os.Getenv("CLERK_API_KEY"))

	container = &Container{
		projectID: projectID,
		version:   version,
		logger:    logger(3).WithCodeNamespace(fmt.Sprintf("%T", container)),
	}

	container.InitializeTraceProvider()

	return container
}

// App creates a new instance of fiber.App
func (container *Container) App() (app *fiber.App) {
	if container.app != nil {
		return container.app
	}

	container.logger.Debug(fmt.Sprintf("creating %T", app))

	app = fiber.New()

	app.Use(otelfiber.Middleware())
	app.Use(fiberLogger.New(fiberLogger.Config{Output: container.Logger(), Format: "${ip} | ${method} | ${path} | ${status} | ${latency}"}))
	app.Use(cors.New())
	app.Use(middlewares.RequestRouter(container.Tracer(), container.Logger(), os.Getenv("APP_HOSTNAME"), container.ProjectEndpointRequestService()))
	app.Use(middlewares.ClerkBearerAuth(container.Logger(), container.Tracer()))
	app.Use(healthcheck.New())

	container.app = app

	container.RegisterEventRoutes()
	container.RegisterProjectRoutes()
	container.RegisterProjectEndpointRoutes()

	container.RegisterProjectEndpointRequestListeners()

	// UnAuthenticated routes
	container.RegisterLemonsqueezyRoutes()

	// this has to be last since it registers the /* route
	container.RegisterSwaggerRoutes()

	return app
}

// GoogleAuthMiddlewares creates router for authenticated requests
func (container *Container) GoogleAuthMiddlewares(audience string, subject string) []fiber.Handler {
	container.logger.Debug("creating GoogleAuthMiddlewares")
	return []fiber.Handler{
		middlewares.GoogleAuth(container.Logger(), container.Tracer(), audience, subject),
		container.AuthenticatedMiddleware(),
	}
}

// AuthenticatedMiddleware creates a new instance of middlewares.Authenticated
func (container *Container) AuthenticatedMiddleware() fiber.Handler {
	container.logger.Debug("creating middlewares.Authenticated")
	return middlewares.Authenticated(container.Tracer())
}

// Logger creates a new instance of telemetry.Logger
func (container *Container) Logger(skipFrameCount ...int) telemetry.Logger {
	container.logger.Debug("creating telemetry.Logger")
	if len(skipFrameCount) > 0 {
		return logger(skipFrameCount[0])
	}
	return logger(2)
}

// GormLogger creates a new instance of gormLogger.Interface
func (container *Container) GormLogger() gormLogger.Interface {
	container.logger.Debug("creating gormLogger.Interface")
	return telemetry.NewGormLogger(
		container.Tracer(),
		container.Logger(5),
	)
}

// DB creates an instance of gorm.DB if it has not been created already
func (container *Container) DB() (db *gorm.DB) {
	if container.db != nil {
		return container.db
	}

	container.logger.Debug(fmt.Sprintf("creating %T", db))

	config := &gorm.Config{TranslateError: true}
	if isLocal() {
		config.Logger = container.GormLogger()
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), config)
	if err != nil {
		container.logger.Fatal(err)
	}
	container.db = db

	sqlDB, err := db.DB()
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot get sql.DB from GORM"))
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err = db.Use(tracing.NewPlugin()); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot use GORM tracing plugin"))
	}

	container.logger.Debug(fmt.Sprintf("Running migrations for [%T]", db))

	if err = db.AutoMigrate(&entities.Project{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate [%T]", &entities.Project{})))
	}
	if err = db.AutoMigrate(&entities.ProjectEndpoint{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate [%T]", &entities.ProjectEndpoint{})))
	}
	if err = db.AutoMigrate(&entities.ProjectEndpointRequest{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate [%T]", &entities.ProjectEndpointRequest{})))
	}
	if err = db.AutoMigrate(&entities.User{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate [%T]", &entities.User{})))
	}

	return container.db
}

// GoogleCredentials returns google credentials as bytes.
func (container *Container) GoogleCredentials() []byte {
	container.logger.Debug("creating google credentials")
	return []byte(os.Getenv("GOOGLE_CREDENTIALS"))
}

// CloudTasksClient creates a new instance of *cloudtasks.Client
func (container *Container) CloudTasksClient() (client *cloudtasks.Client) {
	container.logger.Debug(fmt.Sprintf("creating %T", client))

	client, err := cloudtasks.NewClient(context.Background(), option.WithCredentialsJSON(container.GoogleCredentials()))
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot initialize google cloud tasks client"))
	}

	return client
}

// ClerkBearerAuthMiddlewares creates router for authenticated requests
func (container *Container) ClerkBearerAuthMiddlewares() []fiber.Handler {
	container.logger.Debug("creating ClerkBearerAuthRouter")
	return []fiber.Handler{
		middlewares.ClerkBearerAuth(
			container.Logger().WithCodeNamespace(fmt.Sprintf("%T", middlewares.ClerkBearerAuth)),
			container.Tracer(),
		),
		container.AuthenticatedMiddleware(),
	}
}

// RegisterProjectRoutes registers routes for the /projects prefix
func (container *Container) RegisterProjectRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ProjectHandler{}))
	container.ProjectHandler().RegisterRoutes(container.App(), container.ClerkBearerAuthMiddlewares())
}

// RegisterProjectEndpointRoutes registers routes for the /projects/:projectId/endpoints prefix
func (container *Container) RegisterProjectEndpointRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ProjectEndpointHandler{}))
	container.ProjectEndpointHandler().RegisterRoutes(container.App(), container.ClerkBearerAuthMiddlewares())
}

// ProjectHandler creates a new instance of handlers.ProjectHandler
func (container *Container) ProjectHandler() (handler *handlers.ProjectHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewProjectHandler(
		container.Logger(),
		container.Tracer(),
		container.ProjectHandlerValidator(),
		container.ProjectService(),
	)
}

// ProjectEndpointHandler creates a new instance of handlers.ProjectEndpointHandler
func (container *Container) ProjectEndpointHandler() (handler *handlers.ProjectEndpointHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewProjectEndpointHandler(
		container.Logger(),
		container.Tracer(),
		container.ProjectHandlerEndpointValidator(),
		container.ProjectEndpointService(),
		container.ProjectService(),
	)
}

// RegisterProjectEndpointRequestListeners registers event listeners
func (container *Container) RegisterProjectEndpointRequestListeners() {
	container.logger.Debug(fmt.Sprintf("registering %T", &listeners.ProjectEndpointRequestListener{}))
	container.ProjectEndpointRequestListener().Register(container.EventDispatcher())
}

// ProjectEndpointRequestListener creates a new instance of listeners.ProjectEndpointRequestListener
func (container *Container) ProjectEndpointRequestListener() (handler *listeners.ProjectEndpointRequestListener) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return listeners.NewProjectEndpointRequestListener(
		container.Logger(),
		container.Tracer(),
		container.ProjectEndpointRequestService(),
	)
}

// ProjectHandlerValidator creates a new instance of validators.ProjectHandlerValidator
func (container *Container) ProjectHandlerValidator() (validator *validators.ProjectHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewProjectHandlerValidator(
		container.Logger(),
		container.Tracer(),
		container.ProjectRepository(),
	)
}

// ProjectHandlerEndpointValidator creates a new instance of validators.ProjectEndpointHandlerValidator
func (container *Container) ProjectHandlerEndpointValidator() (validator *validators.ProjectEndpointHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewProjectEndpointHandlerValidator(
		container.Logger(),
		container.Tracer(),
		container.ProjectEndpointRepository(),
	)
}

// ProjectService creates a new instance of services.ProjectService
func (container *Container) ProjectService() (service *services.ProjectService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
	)
}

// ProjectEndpointService creates a new instance of services.ProjectEndpointService
func (container *Container) ProjectEndpointService() (service *services.ProjectEndpointService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectEndpointService(
		container.Logger(),
		container.Tracer(),
		container.ProjectEndpointRepository(),
	)
}

// ProjectEndpointRequestService creates a new instance of services.ProjectEndpointRequestService
func (container *Container) ProjectEndpointRequestService() (service *services.ProjectEndpointRequestService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectEndpointRequestService(
		container.Logger(),
		container.Tracer(),
		container.ProjectEndpointRepository(),
		container.ProjectEndpointRequestRepository(),
		container.EventDispatcher(),
	)
}

// ProjectRepository registers a new instance of repositories.ProjectRepository
func (container *Container) ProjectRepository() repositories.ProjectRepository {
	container.logger.Debug("creating GORM repositories.ProjectRepository")
	return repositories.NewGormProjectRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// ProjectEndpointRepository registers a new instance of repositories.ProjectEndpointRepository
func (container *Container) ProjectEndpointRepository() repositories.ProjectEndpointRepository {
	container.logger.Debug("creating GORM repositories.ProjectEndpointRepository")
	return repositories.NewGormProjectEndpointRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// ProjectEndpointRequestRepository registers a new instance of repositories.ProjectEndpointRequestRepository
func (container *Container) ProjectEndpointRequestRepository() repositories.ProjectEndpointRequestRepository {
	container.logger.Debug("creating GORM repositories.ProjectEndpointRequestRepository")
	return repositories.NewGormProjectEndpointRequestRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// EventsQueue creates a new instance of services.PushQueue
func (container *Container) EventsQueue() queue.Client {
	container.logger.Debug("creating queue.Client")

	return queue.NewGooglePushQueue(
		container.Logger(),
		container.Tracer(),
		container.CloudTasksClient(),
		os.Getenv("EVENTS_QUEUE_NAME"),
		os.Getenv("EVENTS_QUEUE_AUTH_EMAIL"),
	)
}

// Tracer creates a new instance of telemetry.Tracer
func (container *Container) Tracer() (t telemetry.Tracer) {
	container.logger.Debug("creating telemetry.Tracer")
	return telemetry.NewOtelLogger(container.logger)
}

// EventDispatcher creates a new instance of services.EventDispatcher
func (container *Container) EventDispatcher() (dispatcher *services.EventDispatcher) {
	if container.eventDispatcher != nil {
		return container.eventDispatcher
	}

	container.logger.Debug(fmt.Sprintf("creating %T", dispatcher))
	dispatcher = services.NewEventDispatcher(
		container.Logger(),
		container.Tracer(),
		container.Float64Histogram(
			"event.publisher.duration",
			"ms",
			"measures the duration of processing CloudEvents",
		),
		container.EventsQueue(),
		os.Getenv("EVENTS_QUEUE_WEBHOOK"),
	)

	container.eventDispatcher = dispatcher
	return dispatcher
}

// Float64Histogram creates a new instance of metric.Float64Histogram
func (container *Container) Float64Histogram(name, unit, description string) otelMetric.Float64Histogram {
	container.logger.Debug("creating GORM repositories.MessageRepository")
	meter := otel.GetMeterProvider().Meter(
		container.projectID,
		otelMetric.WithInstrumentationVersion(otel.Version()),
	)
	histogram, err := meter.Float64Histogram(name, otelMetric.WithUnit(unit), otelMetric.WithDescription(description))
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot create float64 histogram"))
	}
	return histogram
}

// HTTPClient creates a new http.Client
func (container *Container) HTTPClient(name string) *http.Client {
	container.logger.Debug(fmt.Sprintf("creating %s %T", name, http.DefaultClient))
	return &http.Client{
		Timeout:   60 * time.Second,
		Transport: container.HTTPRoundTripper(name),
	}
}

// HTTPRoundTripper creates an open telemetry http.RoundTripper
func (container *Container) HTTPRoundTripper(name string) http.RoundTripper {
	container.logger.Debug(fmt.Sprintf("Debug: initializing %s %T", name, http.DefaultTransport))
	return otelroundtripper.New(
		otelroundtripper.WithName(name),
		otelroundtripper.WithParent(container.RetryHTTPRoundTripper()),
		otelroundtripper.WithMeter(otel.GetMeterProvider().Meter(container.projectID)),
		otelroundtripper.WithAttributes(container.OtelResources(container.version, container.projectID).Attributes()...),
	)
}

// OtelResources generates default open telemetry resources
func (container *Container) OtelResources(version string, namespace string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(namespace),
		semconv.ServiceVersionKey.String(version),
		semconv.ServiceInstanceIDKey.String(hostName()),
		semconv.DeploymentEnvironmentKey.String(os.Getenv("APP_ENV")),
	)
}

// RetryHTTPRoundTripper creates a retryable http.RoundTripper
func (container *Container) RetryHTTPRoundTripper() http.RoundTripper {
	container.logger.Debug(fmt.Sprintf("initializing retry %T", http.DefaultTransport))
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = container.Logger()
	return retryClient.StandardClient().Transport
}

// EventsHandler creates a new instance of handlers.EventsHandler
func (container *Container) EventsHandler() (handler *handlers.EventsHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))

	return handlers.NewEventsHandler(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
	)
}

// UserRepository registers a new instance of repositories.UserRepository
func (container *Container) UserRepository() repositories.UserRepository {
	container.logger.Debug("creating GORM repositories.UserRepository")
	return repositories.NewGormUserRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// LemonsqueezyService creates a new instance of services.LemonsqueezyService
func (container *Container) LemonsqueezyService() (service *services.LemonsqueezyService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewLemonsqueezyService(
		container.Logger(),
		container.Tracer(),
		container.UserRepository(),
		container.EventDispatcher(),
	)
}

// LemonsqueezyHandler creates a new instance of handlers.LemonsqueezyHandler
func (container *Container) LemonsqueezyHandler() (handler *handlers.LemonsqueezyHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))

	return handlers.NewLemonsqueezyHandlerHandler(
		container.Logger(),
		container.Tracer(),
		container.LemonsqueezyService(),
		container.LemonsqueezyHandlerValidator(),
	)
}

// LemonsqueezyHandlerValidator creates a new instance of validators.LemonsqueezyHandlerValidator
func (container *Container) LemonsqueezyHandlerValidator() (validator *validators.LemonsqueezyHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewLemonsqueezyHandlerValidator(
		container.Logger(),
		container.Tracer(),
		container.LemonsqueezyClient(),
	)
}

// LemonsqueezyClient creates a new instance of lemonsqueezy.Client
func (container *Container) LemonsqueezyClient() (client *lemonsqueezy.Client) {
	container.logger.Debug(fmt.Sprintf("creating %T", client))
	return lemonsqueezy.New(
		lemonsqueezy.WithHTTPClient(container.HTTPClient("lemonsqueezy")),
		lemonsqueezy.WithAPIKey(os.Getenv("LEMONSQUEEZY_API_KEY")),
		lemonsqueezy.WithSigningSecret(os.Getenv("LEMONSQUEEZY_SIGNING_SECRET")),
	)
}

// RegisterLemonsqueezyRoutes registers routes for the /lemonsqueezy prefix
func (container *Container) RegisterLemonsqueezyRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.LemonsqueezyHandler{}))
	container.LemonsqueezyHandler().RegisterRoutes(container.App())
}

// RegisterEventRoutes registers routes for the /events prefix
func (container *Container) RegisterEventRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.EventsHandler{}))
	container.EventsHandler().RegisterRoutes(
		container.App(),
		container.GoogleAuthMiddlewares(
			os.Getenv("EVENTS_QUEUE_WEBHOOK"),
			os.Getenv("EVENTS_QUEUE_AUTH_EMAIL"),
		),
	)
}

// RegisterSwaggerRoutes registers routes for swagger
func (container *Container) RegisterSwaggerRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", swagger.HandlerDefault))
	container.App().Get("/*", swagger.New(swagger.Config{
		Title: docs.SwaggerInfo.Title,
		CustomScript: `
		document.addEventListener("DOMContentLoaded", function(event) {
			var links = document.querySelectorAll("link[rel~='icon']");
			links.forEach(function (link) {
				link.href = 'https://cloud.httpmock.dev/favicon.ico';
			});
		});`,
	}))
}

// InitializeTraceProvider initializes the open telemetry trace provider
func (container *Container) InitializeTraceProvider() func() {
	return container.initializeAxiomProvider(container.version, container.projectID)
}

func (container *Container) initializeUptraceProvider(version string, namespace string) (flush func()) {
	container.logger.Debug("initializing uptrace provider")
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(os.Getenv("UPTRACE_DSN")),
		uptrace.WithServiceName(namespace),
		uptrace.WithServiceVersion(version),
		uptrace.WithDeploymentEnvironment(os.Getenv("APP_ENV")),
	)

	// Send buffered spans and free resources.
	return func() {
		err := uptrace.Shutdown(context.Background())
		if err != nil {
			container.logger.Error(err)
		}
	}
}

// Config loads the configuration from .env
func Config() *Configuration {
	if configuration != nil {
		return configuration
	}

	configuration = new(Configuration)
	if err := env.Parse(configuration); err != nil {
		panic(stacktrace.Propagate(err, fmt.Sprintf("cannot parse [%T]", configuration)))
	}

	return configuration
}

// Resource creates a new instance of resource.Resource
func (container *Container) Resource(version string, namespace string) *resource.Resource {
	// Defines resource with service name, version, and environment.
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(namespace),
		semconv.ServiceVersionKey.String(version),
		semconv.DeploymentEnvironmentKey.String(os.Getenv("APP_ENV")),
	)
}

func (container *Container) initializeAxiomProvider(version string, namespace string) func() {
	// Sets up OTLP HTTP exporter with endpoint, headers, and TLS config.
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("api.axiom.co"),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization":   "Bearer " + os.Getenv("AXIOM_API_KEY"),
			"X-AXIOM-DATASET": os.Getenv("APP_ENV"),
		}),
		otlptracehttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot initialize axiom OTLP HTTP exporter"))
	}

	// Configures the tracer provider with the exporter and resource.
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(container.Resource(version, namespace)),
	)
	otel.SetTracerProvider(tracerProvider)

	// Sets global propagator to W3C Trace Context and Baggage.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create the OTLP log exporter that sends logs to configured destination
	logExporter, err := otlploghttp.New(
		context.Background(),
		otlploghttp.WithEndpoint("api.axiom.co"),
		otlploghttp.WithHeaders(map[string]string{
			"Authorization":   "Bearer " + os.Getenv("AXIOM_API_KEY"),
			"X-AXIOM-DATASET": os.Getenv("APP_ENV"),
		}),
		otlploghttp.WithTLSClientConfig(&tls.Config{}),
	)
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot initialize axiom OTLP HTTP exporter"))
	}

	// Create the logger provider
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(container.Resource(version, namespace)),
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
	)
	global.SetLoggerProvider(loggerProvider)

	return func() {
		err = tracerProvider.Shutdown(context.Background())
		if err != nil {
			container.logger.Error(stacktrace.Propagate(err, "cannot shutdown axiom trace exporter"))
		}
		err = loggerProvider.Shutdown(context.Background())
		if err != nil {
			container.logger.Error(stacktrace.Propagate(err, "cannot shutdown axiom log exporter"))
		}
	}
}

func logger(skipFrameCount int) telemetry.Logger {
	return telemetry.NewSlogLogger(
		context.Background(),
		skipFrameCount,
		getSlogHandler(),
		[]any{string(semconv.ProcessPIDKey), os.Getpid(), string(semconv.HostNameKey), hostName()},
	)
}

func getSlogHandler() slog.Handler {
	// Create a new Slog handler
	if Config().UseOpenTelemetryLogger {
		return otelslog.NewHandler(os.Getenv("GCP_PROJECT_ID"))
	}
	return tint.NewHandler(os.Stderr, &tint.Options{Level: slog.LevelDebug})
}

func hostName() string {
	h, err := os.Hostname()
	if err != nil {
		h = strconv.Itoa(os.Getpid())
	}
	return h
}

func isLocal() bool {
	return os.Getenv("APP_ENV") == "local"
}
