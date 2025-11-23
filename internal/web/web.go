package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/handlers"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/middleware"
)

type Server struct {
	mux      *http.ServeMux
	handlers *handlers.Handlers
}

func NewServer(h *handlers.Handlers) *Server {
	mux := http.NewServeMux()

	return &Server{
		mux:      mux,
		handlers: h,
	}
}

//go:embed static/*
var staticFiles embed.FS

func (ws Server) Start() {
	ws.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested URL is one of the defined handlers
		// If not, redirect to the custom 404 page
		_, pattern := ws.mux.Handler(r)
		if r.URL.Path != pattern {
			ws.handlers.Page404(w, r)
		}

		// Redirect from site index to notes page
		http.Redirect(w, r, "/notes", http.StatusSeeOther)
	})

	ws.mux.HandleFunc("/notes", ws.handlers.Notes)
	ws.mux.HandleFunc("/notes/load-more", ws.handlers.LoadMoreNotes)
	ws.mux.HandleFunc("/add", ws.handlers.CreateNoteModal)
	ws.mux.HandleFunc("POST /notes", ws.handlers.CreateNote)
	ws.mux.HandleFunc("/edit/{id}", ws.handlers.EditNoteModal)
	ws.mux.HandleFunc("PUT /note/{id}", ws.handlers.EditNote)
	ws.mux.HandleFunc("DELETE /note/{id}", ws.handlers.DeleteNote)

	// Use embed.FS to create a file system from the embedded files
	ws.mux.Handle("/static/", http.FileServerFS(staticFiles))

	fileSystem, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	ws.mux.Handle("/favicon.ico", http.StripPrefix("/", http.FileServerFS(fileSystem)))

	fmt.Println("Starting web interface on port: 8089")

	// Create stack for handle multiple middlewares
	middlewares := middleware.CreateMiddlewareStack(
		middleware.LoggingMiddleware,
		middleware.CacheStaticFiles,
		middleware.SlowdownMiddleware,
	)

	server := &http.Server{
		Addr:         ":8089",
		Handler:      middlewares(ws.mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
