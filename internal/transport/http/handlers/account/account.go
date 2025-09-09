package account_handler

import (
	"api/internal/domain/customer"
	"api/internal/transport/http/helper"

	"github.com/labstack/echo/v4"
)

type handler struct {
	customerUseCase customer.UseCase
}

func New(customerUseCase customer.UseCase) helper.Handler {
	return &handler{
		customerUseCase: customerUseCase,
	}
}

func (h *handler) Register(g *echo.Group, m helper.Middleware) {
	baseUrl := g.Group("/account")

	baseUrl.POST("/login", h.Login)
}

func (h *handler) Login(c echo.Context) error {
	var dto struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := helper.BindAndValidate(c, &dto); err != nil {
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
		return helper.HandleError(c, err)
	}

	return c.JSON(200, output)
}
