package repositories

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormProjectEndpointRepository is responsible for persisting entities.ProjectEndpoint
type gormProjectEndpointRequestRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormProjectEndpointRequestRepository creates the GORM version of the ProjectEndpointRequestRepository
func NewGormProjectEndpointRequestRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectEndpointRequestRepository {
	return &gormProjectEndpointRequestRepository{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", &gormProjectEndpointRequestRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectEndpointRequestRepository) Store(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Create(request).Error; err != nil {
		msg := fmt.Sprintf("cannot save project endpoint request with ID [%s]", request.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}
