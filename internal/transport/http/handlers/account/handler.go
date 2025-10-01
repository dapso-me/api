package account_handler

import (
	"api/internal/application"
	"api/internal/transport/http/helper"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	logger   *zap.Logger
	customer application.Customer
}

func New(logger *zap.Logger, customer application.Customer) helper.Handler {
	return &handler{
		logger:   logger,
		customer: customer,
	}
}

func (h *handler) Register(g *echo.Group, m helper.Middleware) {
	customer := g.Group("/account")

	customer.POST("/login", h.signIn)

	private := customer.Group("", m.Authenticate)
	private.GET("", h.getMe)
	private.PATCH("/password", h.changePassword)
}

// signIn godoc
// @Summary      Sign in to Account
// @Description  Authenticate customer with username and password
// @Tags         Account
// @Accept       json
// @Produce      json
// @Param        request body signInReq true "Sign in credentials"
// @Success      200 {object} application.AuthOutput
// @Failure      400 {object} helper.ResponseError
// @Failure      401 {object} helper.ResponseError
// @Failure      500 {object} helper.ResponseError
// @Router       /account/login [post]
func (h *handler) signIn(c echo.Context) error {
	var dto signInReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	signInInput := &application.SignInInput{
		Username:  dto.Username,
		Password:  dto.Password,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	output, err := h.customer.SignIn(c.Request().Context(), signInInput)
	if err != nil {
		h.logger.Error("failed to sign in", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

// getMe godoc
// @Summary      Get current account
// @Description  Get authenticated customer information
// @Tags         Account
// @Accept       json
// @Produce      json
// @Success      200 {object} application.AuthOutput
// @Failure      401 {object} helper.ResponseError
// @Failure      500 {object} helper.ResponseError
// @Router       /account [get]
func (h *handler) getMe(c echo.Context) error {
	output, err := h.customer.Authenticate(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to authenticate", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

// changePassword godoc
// @Summary      Change customer password
// @Description  Change password for authenticated customer
// @Tags         Account
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body changePasswordReq true "Password change data"
// @Success      200
// @Failure      400 {object} helper.ResponseError
// @Failure      401 {object} helper.ResponseError
// @Failure      500 {object} helper.ResponseError
// @Router       /account/password [patch]
func (h *handler) changePassword(c echo.Context) error {
	var dto changePasswordReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	changePasswordInput := &application.ChangePasswordInput{
		CurrentPassword: dto.CurrentPassword,
		NewPassword:     dto.NewPassword,
	}
	err := h.customer.ChangePassword(c.Request().Context(), changePasswordInput)
	if err != nil {
		h.logger.Error("failed to change password", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, map[string]any{"status": "ok"})
}
