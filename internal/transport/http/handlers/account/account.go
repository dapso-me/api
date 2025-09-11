package account_handler

import (
	"api/internal/domain/customer"
	"api/internal/transport/http/helper"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type handler struct {
	logger          *zap.Logger
	customerUseCase customer.UseCase
}

func New(customerUseCase customer.UseCase, logger *zap.Logger) helper.Handler {
	return &handler{
		logger:          logger,
		customerUseCase: customerUseCase,
	}
}

func (h *handler) Register(g *echo.Group, m helper.Middleware) {
	baseUrl := g.Group("/account")
	baseUrl.POST("/login", h.login)
	baseUrl.POST("/register", h.register)
	baseUrl.POST("/register/confirm", h.confirmRegister)

	private := baseUrl.Group("", m.Authenticate)
	private.GET("", h.getMe)
	private.PATCH("", h.update)
	private.PATCH("/password", h.updatePassword)
}

func (h *handler) login(c echo.Context) error {
	var dto loginReq

	if err := helper.BindAndValidate(c, &dto); err != nil {
		h.logger.Error("failed to bind", zap.Error(err))
		return helper.HandleError(c, err)
	}

	loginInput := &customer.LoginInput{
		Email:     dto.Email,
		Password:  dto.Password,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	output, err := h.customerUseCase.Login(c.Request().Context(), loginInput)
	if err != nil {
		h.logger.Error("failed to login", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

func (h *handler) register(c echo.Context) error {
	var dto registerReq

	if err := helper.BindAndValidate(c, &dto); err != nil {
		h.logger.Error("failed to bind", zap.Error(err))
		return helper.HandleError(c, err)
	}

	registerInput := &customer.RegisterInput{
		Email:     dto.Email,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}

	err := h.customerUseCase.Register(c.Request().Context(), registerInput)
	if err != nil {
		h.logger.Error("failed to register", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.NoContent(200)
}

func (h *handler) confirmRegister(c echo.Context) error {
	var dto confirmRegisterReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	confirmRegisterInput := &customer.ConfirmRegisterInput{
		Email:     dto.Email,
		Code:      dto.Code,
		Password:  dto.Password,
		Name:      dto.Name,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	output, err := h.customerUseCase.ConfirmRegister(c.Request().Context(), confirmRegisterInput)
	if err != nil {
		h.logger.Error("failed to confirm register", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

func (h *handler) getMe(c echo.Context) error {
	output, err := h.customerUseCase.Authenticate(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to get me", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}

func (h *handler) update(c echo.Context) error {
	var dto updateReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	updateInput := &customer.UpdateInput{
		Name: dto.Name,
	}
	err := h.customerUseCase.Update(c.Request().Context(), updateInput)
	if err != nil {
		h.logger.Error("failed to update", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.NoContent(200)
}

func (h *handler) updatePassword(c echo.Context) error {
	var dto updagePasswordReq
	if err := helper.BindAndValidate(c, &dto); err != nil {
		return helper.HandleError(c, err)
	}

	input := &customer.UpdatePasswordInput{
		CurrentPassword: dto.CurrentPassword,
		NewPassword:     dto.NewPassword,
	}
	if err := h.customerUseCase.UpdatePassword(c.Request().Context(), input); err != nil {
		h.logger.Error("failed to update password", zap.Error(err))
		return helper.HandleError(c, err)
	}

	return c.NoContent(200)
}
