package router

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// 正常なリクエストかどうかの確認
func (r *Router) requestVerification(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ヘッダーを取得
		h := c.Request().Header

		// 認証トークンの照合
		if h.Get(botTokenHeader) != r.Config.Verification_Token {
			return c.JSON(http.StatusForbidden, errorMessage{Message: "invalid token"})
		}

		// ボディ形式の確認
		if !strings.HasPrefix(h.Get(contentTypeHeader), "application/json") {
			return c.JSON(http.StatusBadRequest, errorMessage{Message: "invalid content type"})
		}

		return next(c)
	}
}
