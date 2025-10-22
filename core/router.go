package core

import (
	"github.com/pikachu0310/go-backend-template/core/internal/handler"
	_ "github.com/pikachu0310/go-backend-template/docs" // Swagger docs
	"github.com/pikachu0310/go-backend-template/frontend"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(h *handler.Handler, e *echo.Echo) {
	e.StaticFS("/", frontend.UI)

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1API := e.Group("/api/v1")

	// ping API
	pingAPI := v1API.Group("/ping")
	{
		pingAPI.GET("", h.Ping)
	}

	// user API
	userAPI := v1API.Group("/users")
	{
		userAPI.GET("", h.GetUsers)
		userAPI.POST("", h.CreateUser)
		userAPI.GET("/:userID", h.GetUser)
	}
}
