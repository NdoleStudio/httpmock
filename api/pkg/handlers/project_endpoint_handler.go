package handlers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/requests"
	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/NdoleStudio/httpmock/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointHandler handles entities.ProjectEndpoint requests.
type ProjectEndpointHandler struct {
	handler
	logger         telemetry.Logger
	tracer         telemetry.Tracer
	validator      *validators.ProjectEndpointHandlerValidator
	projectService *services.ProjectService
	service        *services.ProjectEndpointService
}

// NewProjectEndpointHandler creates a new ProjectEndpointHandler
func NewProjectEndpointHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.ProjectEndpointHandlerValidator,
	service *services.ProjectEndpointService,
	projectService *services.ProjectService,
) (h *ProjectEndpointHandler) {
	return &ProjectEndpointHandler{
		logger:         logger.WithCodeNamespace(fmt.Sprintf("%T", h)),
		tracer:         tracer,
		validator:      validator,
		service:        service,
		projectService: projectService,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *ProjectEndpointHandler) RegisterRoutes(app *fiber.App, middlewares []fiber.Handler) {
	router := app.Group("/v1/projects/:projectId/endpoints")
	router.Get("/", h.computeRoute(h.index, middlewares)...)
	router.Post("/", h.computeRoute(h.store, middlewares)...)
	router.Get("/:projectEndpointId", h.computeRoute(h.show, middlewares)...)
	router.Put("/:projectEndpointId", h.computeRoute(h.update, middlewares)...)
	router.Delete("/:projectEndpointId", h.computeRoute(h.delete, middlewares)...)
	router.Get("/:projectEndpointId/traffic", h.computeRoute(h.traffic, middlewares)...)
}

// @Summary      List of project endpoints
// @Description  Fetches the list of all projects endpoints available to the currently authenticated user and project
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Success      200 		{object}	responses.Ok[[]entities.ProjectEndpoint]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /v1/projects/:projectId/endpoints 	[get]
func (h *ProjectEndpointHandler) index(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectId")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], fetching projects with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseNotFound(c, fmt.Sprintf("cannot list endpoints for project with ID [%s]", c.Params("projectId")))
	}

	authUser := h.userFromContext(c)
	projects, err := h.service.Index(ctx, authUser.ID, uuid.MustParse(c.Params("projectId")))
	if err != nil {
		msg := fmt.Sprintf("cannot fetch project endpoints for user with ID [%s] and projectID [%s]", authUser.ID, c.Params("projectId"))
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "projects fetched successfully", projects)
}

// @Summary      store a new project endpoint
// @Description  This endpoint stores a new project endpoint for a user
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Param        payload	body 		requests.ProjectEndpointStoreRequest	true 	"project endpoint store payload"
// @Success      200 		{object}	responses.Ok[entities.ProjectEndpoint]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /v1/projects/:projectId/endpoints 	[post]
func (h *ProjectEndpointHandler) store(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.ProjectEndpointStoreRequest
	if err := c.BodyParser(&request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	authUser := h.userFromContext(c)
	request.ProjectID = c.Params("projectId")

	if errors := h.validator.ValidateStore(ctx, authUser.ID, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while storing project endpoint with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while storing mock endpoint")
	}

	project, err := h.projectService.Load(ctx, authUser.ID, uuid.MustParse(request.ProjectID))
	if err != nil {
		msg := fmt.Sprintf("cannot find project with id [%s] for user [%s]", request.ProjectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	endpoint, err := h.service.Store(ctx, project, request.ToProjectEndpointStorePrams(authUser.ID))
	if err != nil {
		ctxLogger.Error(stacktrace.Propagate(err, fmt.Sprintf("cannot store project endpoint for project ID [%s] for user ID [%s]", request.ProjectID, authUser.ID)))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "endpoint created successfully", endpoint)
}

// @Summary      Update a project endpoint
// @Description  This endpoint updates a project endpoint for a user
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Param 		 projectId			path		string true "Project ID"
// @Param 		 projectEndpointId	path		string true "Project Endpoint ID"
// @Param        payload			body		requests.ProjectEndpointUpdateRequest	true 	"project update payload"
// @Success      200 				{object}	responses.Ok[entities.ProjectEndpoint]
// @Failure      400				{object}	responses.BadRequest
// @Failure 	 401    			{object}	responses.Unauthorized
// @Failure      422				{object}	responses.UnprocessableEntity
// @Failure      500				{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId}	[put]
func (h *ProjectEndpointHandler) update(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	request := new(requests.ProjectEndpointUpdateRequest)
	if err := c.BodyParser(request); err != nil {
		msg := fmt.Sprintf("cannot marshall params [%s] into %T", c.OriginalURL(), request)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseBadRequest(c, err)
	}

	authUser := h.userFromContext(c)
	request.ProjectEndpointID = c.Params("projectEndpointId")
	request.ProjectID = c.Params("projectId")

	if errors := h.validator.ValidateUpdate(ctx, authUser.ID, request.Sanitize()); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while updating project endpoint with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating project endpoint")
	}

	if _, err := h.projectService.Load(ctx, authUser.ID, uuid.MustParse(request.ProjectID)); err != nil {
		msg := fmt.Sprintf("cannot find project with id [%s] for user [%s]", request.ProjectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	project, err := h.service.Update(ctx, request.ToProjectEndpointUpdatePrams(authUser.ID))
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot find project endpoint with ID [%s] and project id [%s] for user [%s]", request.ProjectEndpointID, request.ProjectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot update project endpoint with ID [%s] and project id [%s] for user [%s]", request.ProjectEndpointID, request.ProjectID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "endpoint updated successfully", project)
}

// @Summary      Get a project endpoint
// @Description  This URL gets a project endpoint for a user
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Param 		 projectId			path 		string true "Project ID"
// @Param 		 projectEndpointId	path 		string true "Project Endpoint ID"
// @Success      200 				{object}	responses.Ok[entities.ProjectEndpoint]
// @Failure      400				{object}	responses.BadRequest
// @Failure 	 401    			{object}	responses.Unauthorized
// @Failure      422				{object}	responses.UnprocessableEntity
// @Failure      500				{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId} [get]
func (h *ProjectEndpointHandler) show(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.validator.ValidateUUID(c, "projectId"); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching endpoints with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching project endpoint")
	}

	if errors := h.mergeErrors(h.validateUUID(c, "projectEndpointId")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching endpoints with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while fetching project endpoint")
	}

	projectID := uuid.MustParse(c.Params("projectId"))
	projectEndpointID := uuid.MustParse(c.Params("projectEndpointId"))
	authUser := h.userFromContext(c)

	endpoint, err := h.service.Load(ctx, authUser.ID, projectID, projectEndpointID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("project endpoint not found with id [%s] and project id [%s] for user [%s]", projectEndpointID, projectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint with id [%s] and project id [%s] for user [%s]", projectEndpointID, projectID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project endpoint fetched successfully", endpoint)
}

// @Summary      Delete a project endpoint
// @Description  This API deletes a project endpoint for a user
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Param 		 projectId			path 		string true "Project ID"
// @Param 		 projectEndpointId	path 		string true "Project Endpoint ID"
// @Success      200 				{object}	responses.NoContent
// @Failure      400				{object}	responses.BadRequest
// @Failure 	 401    			{object}	responses.Unauthorized
// @Failure 	 404    			{object}	responses.NotFound
// @Failure      422				{object}	responses.UnprocessableEntity
// @Failure      500				{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId} [delete]
func (h *ProjectEndpointHandler) delete(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.mergeErrors(h.validateUUID(c, "projectId")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while deleting project with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting project endpoint")
	}

	if errors := h.mergeErrors(h.validateUUID(c, "projectEndpointId")); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while fetching endpoints with url [%s]", spew.Sdump(errors), c.OriginalURL())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting project endpoint")
	}

	authUser := h.userFromContext(c)
	projectID := uuid.MustParse(c.Params("projectId"))
	projectEndpointID := uuid.MustParse(c.Params("projectEndpointId"))

	err := h.service.Delete(ctx, authUser.ID, projectID, projectEndpointID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("project endpoint not found with ID [%s] and project ID [%s] for user [%s]", projectEndpointID, projectID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot delete project endpoint with ID [%s] and project ID [%s] for user [%s]", projectEndpointID, projectID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "project endpoint deleted successfully")
}

// @Summary      Get project traffic
// @Description  This endpoint returns the time series traffic for a endpoint in the last 30 days.
// @Security	 BearerAuth
// @Tags         ProjectEndpoints
// @Produce      json
// @Param 		 projectId					path 		string true "Project ID"
// @Param 		 projectEndpointId			path 		string true "Project Endpoint ID"
// @Success      200 		{object}	responses.Ok[[]repositories.TimeSeriesData]
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /v1/projects/{projectId}/endpoints/{projectEndpointId}/traffic [get]
func (h *ProjectEndpointHandler) traffic(c *fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if errors := h.validator.ValidateUUID(c, "projectEndpointId"); len(errors) != 0 {
		msg := fmt.Sprintf("validation errors [%s], while loading project endpoint traffic with request [%s]", spew.Sdump(errors), c.Body())
		ctxLogger.Warn(stacktrace.NewError(msg))
		return h.responseUnprocessableEntity(c, errors, "validation errors while creating project")
	}

	projectEndpointID := uuid.MustParse(c.Params("projectEndpointId"))
	authUser := h.userFromContext(c)

	timeSeries, err := h.service.Traffic(ctx, authUser.ID, projectEndpointID)
	if stacktrace.GetCode(err) == repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot load traffic data for project endpoint with id [%s] for user [%s]", projectEndpointID, authUser.ID)
		ctxLogger.Warn(stacktrace.Propagate(err, msg))
		return h.responseNotFound(c, msg)
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load traffic data for project endpoint [%s] user with ID [%s]", projectEndpointID, authUser.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "project endpoint traffic fetched successfully", timeSeries)
}
