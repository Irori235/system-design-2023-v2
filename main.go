package main

import (
	"github.com/Irori235/system-design-2023-v2/internal/handler"
	"github.com/Irori235/system-design-2023-v2/internal/migration"
	"github.com/Irori235/system-design-2023-v2/internal/pkg/config"
	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/middleware"
	"github.com/jmoiron/sqlx"
)

func main() {
	// setup gin
	r := gin.Default()

	// middlewares
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery())

	// connect to database
	db, err := sqlx.Connect("mysql", config.MySQL().FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		e.Logger.Fatal(err)
	}

	// setup repository
	repo := repository.New(db)

	// setup routes
	h := handler.New(repo)
	v1API := e.Group("/api/v1")
	h.SetupRoutes(v1API)

	e.Logger.Fatal(e.Start(config.AppAddr()))
}
