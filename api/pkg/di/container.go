package di

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

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

	"github.com/hirosassa/zerodriver"
	"github.com/rs/zerolog"
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
		logger: logger(3).WithService(fmt.Sprintf("%T", container)),
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
		logger:    logger(3).WithService(fmt.Sprintf("%T", container)),
	}

	container.InitializeTraceProvider()

	container.RegisterEventRoutes()
	container.RegisterProjectRoutes()

	// UnAuthenticated routes
	container.RegisterLemonsqueezyRoutes()

	// this has to be last since it registers the /* route
	container.RegisterSwaggerRoutes()

	return container
}

// App creates a new instance of fiber.App
func (container *Container) App() (app *fiber.App) {
	if container.app != nil {
		return container.app
	}

	container.logger.Debug(fmt.Sprintf("creating %T", app))

	app = fiber.New()

	if os.Getenv("USE_HTTP_LOGGER") == "true" {
		app.Use(fiberLogger.New())
	}

	app.Use(otelfiber.Middleware())
	app.Use(cors.New())
	app.Use(middlewares.HTTPRequestLogger(container.Tracer(), container.Logger()))

	app.Use(middlewares.ClerkBearerAuth(container.Logger().WithService("middlewares.ClerkBearerAuth"), container.Tracer()))

	container.app = app
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
	return logger(3)
}

// GormLogger creates a new instance of gormLogger.Interface
func (container *Container) GormLogger() gormLogger.Interface {
	container.logger.Debug("creating gormLogger.Interface")
	return telemetry.NewGormLogger(
		container.Tracer(),
		container.Logger(6),
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

	if err = db.Use(tracing.NewPlugin()); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot use GORM tracing plugin"))
	}

	container.logger.Debug(fmt.Sprintf("Running migrations for [%T]", db))

	if err = db.AutoMigrate(&entities.Project{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate [%T]", &entities.Project{})))
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
			container.Logger().WithService(fmt.Sprintf("%T", middlewares.ClerkBearerAuth)),
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

// ProjectHandlerValidator creates a new instance of validators.ProjectHandlerValidator
func (container *Container) ProjectHandlerValidator() (validator *validators.ProjectHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewProjectHandlerValidator(
		container.Logger(),
		container.Tracer(),
		container.ProjectRepository(),
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

// ProjectRepository registers a new instance of repositories.ProjectRepository
func (container *Container) ProjectRepository() repositories.ProjectRepository {
	container.logger.Debug("creating GORM repositories.ProjectRepository")
	return repositories.NewGormProjectRepository(
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
	container.App().Get("/*", swagger.HandlerDefault)
}

// InitializeTraceProvider initializes the open telemetry trace provider
func (container *Container) InitializeTraceProvider() func() {
	return container.initializeUptraceProvider(container.version, container.projectID)
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

func logger(skipFrameCount int) telemetry.Logger {
	fields := map[string]string{
		"pid":      strconv.Itoa(os.Getpid()),
		"hostname": hostName(),
	}

	return telemetry.NewZerologLogger(
		os.Getenv("GCP_PROJECT_ID"),
		fields,
		logDriver(skipFrameCount),
		nil,
	)
}

func logDriver(skipFrameCount int) *zerodriver.Logger {
	if isLocal() {
		return consoleLogger(skipFrameCount)
	}
	return jsonLogger(skipFrameCount)
}

func jsonLogger(skipFrameCount int) *zerodriver.Logger {
	logLevel := zerolog.DebugLevel
	zerolog.SetGlobalLevel(logLevel)

	// See: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
	logLevelSeverity := map[zerolog.Level]string{
		zerolog.TraceLevel: "DEFAULT",
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO",
		zerolog.WarnLevel:  "WARNING",
		zerolog.ErrorLevel: "ERROR",
		zerolog.PanicLevel: "CRITICAL",
		zerolog.FatalLevel: "CRITICAL",
	}

	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return logLevelSeverity[l]
	}
	zerolog.TimestampFieldName = "time"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zl := zerolog.New(os.Stderr).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()
	return &zerodriver.Logger{Logger: &zl}
}

func hostName() string {
	h, err := os.Hostname()
	if err != nil {
		h = strconv.Itoa(os.Getpid())
	}
	return h
}

func consoleLogger(skipFrameCount int) *zerodriver.Logger {
	l := zerolog.New(
		zerolog.ConsoleWriter{
			Out: os.Stderr,
		}).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()
	return &zerodriver.Logger{
		Logger: &l,
	}
}

func isLocal() bool {
	return os.Getenv("APP_ENV") == "local"
}
