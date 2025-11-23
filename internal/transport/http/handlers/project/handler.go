package project_handler

import (
	"api/internal/application"
	"api/internal/transport/http/helper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	logger  *zap.Logger
	project application.Project
}

func New(logger *zap.Logger, project application.Project) helper.Handler {
	return &handler{
		logger:  logger,
		project: project,
	}
}

func (h *handler) Register(g *echo.Group, m helper.Middleware) {
	customer := g.Group("/project")

	customer.GET("/:slug", h.findBySlug)

	private := customer.Group("", m.Authenticate)
	private.GET("/my", h.findMy)
	private.POST("", h.add)
	private.DELETE("/:ID", h.remove)
}

// @Tags         Projects
// @Param        slug  path      string  true  "Project slug"
// @Router       /project/{slug} [get]
func (h *handler) findBySlug(c echo.Context) error {
	project, err := h.project.FindOneBySlug(c.Request().Context(), c.Param("slug"))
	if err != nil {
		h.logger.Error("failed to find by slug", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, project)
}

// @Tags         Projects
// @Param        request  body      addProjectReq  true  "Project data"
// @Success      200      {object}  project.Project
// @Router       /project [post]
func (h *handler) add(c echo.Context) error {
	var dto addProjectReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	input := &application.AddProjectInput{
		CustomerID: dto.CustomerID,
		Slug:       dto.Slug,
		Name:       dto.Name,
	}
	output, err := h.project.Add(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("failed to add project", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

// @Tags         Projects
// @Success      200  {array}   project.Project
// @Router       /project/my [get]
func (h *handler) findMy(c echo.Context) error {
	projects, err := h.project.FindByOwner(c.Request().Context())
	if err != nil {
		return helper.HandleError(c, err)
	}

	return c.JSON(200, projects)
}

// @Summary      Delete project
// @Description  Delete project by ID
// @Tags         projects
// @Param        ID   path      string  true  "Project ID" format(uuid)
// @Success      200  {object}  helper.StatusResponse
// @Router       /project/{ID} [delete]
func (h *handler) remove(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("ID"))
	if err != nil {
		return helper.HandleError(c, err)
	}

	if err := h.project.Remove(c.Request().Context(), projectID); err != nil {
		h.logger.Error("failed to remove project", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return helper.StatusOk(c)
}
