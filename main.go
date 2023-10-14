package main

import (
	"log"

	"github.com/Irori235/system-design-2023-v2/internal/handler"
	"github.com/Irori235/system-design-2023-v2/internal/migration"
	"github.com/Irori235/system-design-2023-v2/internal/pkg/config"
	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
	// setup gin
	r := gin.Default()

	// connect to database
	db, err := sqlx.Connect("mysql", config.MySQL().FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		log.Fatal(err)
	}

	// setup repository
	repo := repository.New(db)

	// setup routes
	h := handler.New(repo)
	v1API := r.Group("/api/v1")
	h.SetupRoutes(v1API)

	log.Fatal(r.Run(config.AppAddr()))
}
