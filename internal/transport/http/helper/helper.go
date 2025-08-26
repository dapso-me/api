package helper

import (
	"api/internal/common"
	"errors"

	"github.com/labstack/echo/v4"
)

type Middleware interface {
	// Authenticate(next echo.HandlerFunc) echo.HandlerFunc
}

type Handler interface {
	Register(g *echo.Group, m Middleware)
}

func BindAndValidate(c echo.Context, dto any) error {
	if err := c.Bind(&dto); err != nil {
		return errors.New("INVALID_INPUT_DATA")
	}

	if err := c.Validate(dto); err != nil {
		return err
	}

	return nil
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func HandleError(c echo.Context, err error) error {
	if errors.Is(err, common.ErrUnauthorized) {
		return c.JSON(401, ResponseError{Code: "UNAUTHORIZED", Message: "Unauthorized"})
	} else if errors.Is(err, common.ErrForbidden) {
		return c.JSON(403, ResponseError{Code: "FORBIDDEN", Message: "Forbidden"})
	} else {
		var appErr *common.AppError
		if errors.As(err, &appErr) {
			return c.JSON(400, ResponseError{Code: appErr.Code(), Message: appErr.Message()})
		} else {
			return c.JSON(500, ResponseError{Code: "INTERNAL_SERVER_ERROR", Message: "Internal Server Error"})
		}
	}
}
