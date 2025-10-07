package menu_handler

import (
	"api/internal/application"
	"api/internal/transport/http/helper"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	logger   *zap.Logger
	category application.Category
	dish     application.Dish
}

func New(logger *zap.Logger, category application.Category, dish application.Dish) helper.Handler {
	return &handler{
		logger:   logger,
		category: category,
		dish:     dish,
	}
}

func (h *handler) Register(g *echo.Group, m helper.Middleware) {
	project := g.Group("/project/:projectID")
	project.GET("/menu", h.getMenu)

	private := project.Group("", m.Authenticate)
	private.POST("/category", h.addCategory)
	private.POST("/dish", h.addDish)
}

// getMenu godoc
// @Summary      Get Menu Categories
// @Description  Get all menu categories for a specific project
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        projectID path string true "Project ID" format(uuid)
// @Success      200 {array} menu.Category
// @Failure      400 {object} helper.ResponseError
// @Failure      500 {object} helper.ResponseError
// @Router       /project/{projectID} [get]
func (h *handler) getMenu(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		return helper.HandleError(c, err)
	}

	menu, err := h.category.FindByProjectID(c.Request().Context(), projectID)
	if err != nil {
		h.logger.Error("failed to get menu", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, menu)
}

// addCategory godoc
// @Summary      Add Menu Category
// @Description  Add a new category to project menu
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        projectID path string true "Project ID" format(uuid)
// @Param        request body addCategoryReq true "Category data"
// @Success      200 {object} menu.Category
// @Failure      400 {object} helper.ResponseError
// @Failure      500 {object} helper.ResponseError
// @Router       /project/{projectID}/category [post]
func (h *handler) addCategory(c echo.Context) error {
	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		return helper.HandleError(c, err)
	}

	var dto addCategoryReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	input := &application.AddCategoryInput{
		ProjectID:    projectID,
		Position:     dto.Position,
		Translations: dto.Translations,
	}
	output, err := h.category.Add(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("failed to add category", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

func (h *handler) addDish(c echo.Context) error {
	var dto addDishReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		h.logger.Error("failed to add dish", zap.Error(err))
		return helper.HandleError(c, err)
	}

	projectID, err := uuid.Parse(c.Param("projectID"))
	if err != nil {
		h.logger.Error("failed to add dish", zap.Error(err))
		return helper.HandleError(c, err)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(dto.Photo)
	if err != nil {
		h.logger.Error("failed to add dish", zap.Error(err))
		return helper.HandleError(c, err)
	}

	input := &application.AddDishInput{
		ProjectID:    projectID,
		CategoryID:   dto.CategoryID,
		Position:     dto.Position,
		Photo:        imageBytes,
		Price:        dto.Price,
		Translations: dto.Translations,
	}
	output, err := h.dish.Add(c.Request().Context(), input)
	if err != nil {
		h.logger.Error("failed to add dish", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}
