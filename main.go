package main

import (
	"log"

	"github.com/pikachu0310/go-backend-template/core"
	"github.com/pikachu0310/go-backend-template/core/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ras0q/goalie"
)

// @title Go Backend Template API
// @version 1.0
// @description バックエンドテンプレートのAPI仕様書です。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/pikachu0310/go-backend-template
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	if err := run(); err != nil {
		log.Fatalf("runtime error: %+v", err)
	}
}

func run() (err error) {
	g := goalie.New()
	defer g.Collect(&err)

	var config core.Config
	config.Parse()

	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// connect to and migrate database
	db, err := database.Setup(config.MySQLConfig())
	if err != nil {
		return err
	}
	defer g.Guard(db.Close)

	s := core.InjectDeps(db)

	core.SetupRoutes(s.Handler, e)

	if err := e.Start(config.AppAddr); err != nil {
		return err
	}

	return nil
}
