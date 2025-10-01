package http

import (
	"api/internal/common"
	"api/internal/domain/session"
	"api/internal/transport/http/helper"
	"errors"

	"github.com/labstack/echo/v4"
)

func (h *httpServer) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := c.Request().Header.Get("Authorization")

		if accessToken == "" {
			token, err := c.Cookie("access_token")
			if err != nil {
				return helper.HandleError(c, common.ErrUnauthorized)
			}

			accessToken = token.Value
		}

		sessionEntity, err := h.sessionRepo.FindOneByAccessToken(c.Request().Context(), accessToken)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return helper.HandleError(c, common.ErrUnauthorized)
			}
			return helper.HandleError(c, err)
		}

		ctx := session.SetToContext(c.Request().Context(), sessionEntity)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
