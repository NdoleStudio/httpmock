package handlers

import (
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/requests"
	"github.com/oklog/ulid/v2"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/NdoleStudio/httpmock/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointRequestHandler handles entities.ProjectEndpointRequest.
type ProjectEndpointRequestHandler struct {
	handler
	logger                 telemetry.Logger
	tracer                 telemetry.Tracer
	validator              *validators.ProjectEndpointRequestHandlerValidator
	projectEndpointService *services.ProjectEndpointService
	service                *services.ProjectEndpointRequestService
}

// NewProjectEndpointRequestHandler creates a new ProjectEndpointRequestHandler
func NewProjectEndpointRequestHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.ProjectEndpointRequestHandlerValidator,
	service *services.ProjectEndpointRequestService,
	projectEndpointService *services.ProjectEndpointService,
) (h *ProjectEndpointRequestHandler) {
	return &ProjectEndpointRequestHandler{
		logger:                 logger.WithCodeNamespace(fmt.Sprintf("%T", h)),
		tracer:                 tracer,
		validator:              validator,
		service:                service,
		projectEndpointService: projectEndpointService,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ProjectEndpointRequestHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectId/endpoints/:projectEndpointId/requests")
	router.Get("/", h.computeRoute(h.index, middlewares)...)
	router.Delete("/:projectEndpointRequestId", h.computeRoute(h.delete, middlewares)...)
}

// @Summary      List of project endpoint requests
// @Description  Fetches the list of all projects endpoint requests available to the currently authenticated user
// @Security	 BearerAuth
// @Tags         ProjectEndpointRequests
// @Produce      json
// @Param        prev		query  string  	false	"ID of the last request returned in the previous page"
// @Param        limit		query  int  	false	"number of messages to return"			minimum(1)	maximum(100)
// @Success      200 		{object}	responses.Ok[[]entities.ProjectEndpointRequest]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId}/requests 	[get]
func (h *ProjectEndpointRequestHandler) index(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ProjectEndpointRequestIndexRequest
	if err := c.QueryParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params in [%s] into [%T]", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	if errors := h.validator.ValidateIndex(request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], fetching project endpoint requests with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseNotFound(c, fmt.Sprintf("cannot list requests project endpoint with ID [%s]", request.ProjectEndpointID))
	}

	endpointRequests, err := h.service.Index(ctx, h.userIDFomContext(c), uuid.MustParse(request.ProjectEndpointID), request.Limit, request.PrevID())
	if err != nil {
		msg := fmt.Sprintf("cannot fetch project endpoints for user with ID [%s] and projectID [%s]", h.userIDFomContext(c), request.ProjectID)
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project endpoint requests fetched successfully", endpointRequests)
}

// @Summary      Delete a project endpoint request
// @Description  This API deletes a project endpoint request for a user
// @Security	 BearerAuth
// @Tags         ProjectEndpointRequests
// @Produce      json
// @Param 		 projectId					path 		string true "Project ID"
// @Param 		 projectEndpointId			path 		string true "Project Endpoint ID"
// @Param 		 projectEndpointRequestId	path 		string true "Project Endpoint Request ID"
// @Success      204 						{object}	responses.NoContent
// @Failure      400						{object}	responses.BadRequest
// @Failure 	 401    					{object}	responses.Unauthorized
// @Failure 	 404    					{object}	responses.NotFound
// @Failure      422						{object}	responses.UnprocessableEntity
// @Failure      500						{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId}/requests/{projectEndpointRequestId} [delete]
func (h *ProjectEndpointRequestHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if validationErrors := h.mergeErrors(h.validateUUID(c, "projectEndpointRequestId")); len(validationErrors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting project with url [%s]", spew.Sdump(validationErrors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, validationErrors, "validation errors while deleting project endpoint request")
	}

	requestID := ulid.MustParse(c.Params("projectEndpointRequestId"))

	err := h.service.Delete(ctx, h.userIDFomContext(c), requestID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("project endpoint request not found with ID [%s] and for user [%s]", requestID, h.userIDFomContext(c))
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot delete project endpoint request with ID [%s] for user [%s]", requestID, h.userIDFomContext(c))
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "project endpoint request deleted successfully")
}
