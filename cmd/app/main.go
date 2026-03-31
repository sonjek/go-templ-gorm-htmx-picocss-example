package main

import (
	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/storage"
	"github.com/sonjek/go-full-stack-example/internal/web"
	"github.com/sonjek/go-full-stack-example/internal/web/handlers"
)

func main() {
	db := storage.NewDbStorage()
	storage.DBMigrate(db)
	storage.SeedData(db)

	noteService := service.NewNoteService(db)
	appHandlers := handlers.NewHandler(db, noteService)
	webServer := web.NewServer(appHandlers)
	webServer.Start()
}
