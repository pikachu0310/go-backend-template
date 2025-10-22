package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping godoc
// @Summary Ping API
// @Description サーバーの死活確認用エンドポイント
// @Tags ping
// @Accept json
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (h *Handler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
