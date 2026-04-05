package web

import (
	"context"
	"embed"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/sonjek/go-full-stack-example/internal/web/handlers"
	"github.com/sonjek/go-full-stack-example/internal/web/middleware"
)

const (
	defaultPort      = "3000"
	requestBodyLimit = 1 << 20 // 1 MiB
	maxCacheAge      = 86400   // 24 hours
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	app      *fiber.App
	handlers *handlers.Handlers
}

func NewServer(h *handlers.Handlers) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit: requestBodyLimit,
	})

	return &Server{
		app:      app,
		handlers: h,
	}
}

func (ws *Server) SetupMiddleware() {
	ws.app.Use(middleware.LoggingMiddleware)

	// Add Global Security Headers
	ws.app.Use(helmet.New(helmet.Config{
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		// Use CSP middleware instead
		ContentSecurityPolicy: "",
	}))

	ws.app.Use(func(c fiber.Ctx) error {
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data:; " +
			"font-src 'self';"
		c.Set("Content-Security-Policy", csp)
		return c.Next()
	})

	sessionMiddleware, sessionStore := session.NewWithStore(session.Config{
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		IdleTimeout:    30 * time.Minute,
	})
	ws.app.Use(sessionMiddleware)

	ws.app.Use(csrf.New(csrf.Config{
		CookieName:        "csrf_",
		CookieHTTPOnly:    true, // Server-side only cookie (no JavaScript access)
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		SingleUseToken:    false,
		Extractor:         extractors.FromHeader("X-Csrf-Token"),
		Session:           sessionStore,
	}))

	// CSRF Token Refresher Middleware (for htmx requests)
	ws.app.Use(func(c fiber.Ctx) error {
		// Let request finish (CSRF and Handlers run)
		err := c.Next()

		// No refresh if the request is not successful
		if c.Response().StatusCode() >= 400 {
			return err
		}

		// Set CSRF token in the response header for HTML and HTMX requests only
		isHTML := strings.Contains(c.Get(fiber.HeaderAccept), "text/html")
		isHTMX := c.Get("HX-Request") != ""
		if isHTML || isHTMX {
			// Extract updated token from the context AFTER handler execution
			if token := csrf.TokenFromContext(c); token != "" {
				c.Set("X-Csrf-Token", token)
			}
		}

		return err
	})
}

func (ws *Server) SetupRoutes() error {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	ws.app.Use("/static", static.New("", static.Config{
		FS:     subFS,
		MaxAge: maxCacheAge,
	}))

	ws.app.Get("/favicon.ico", func(c fiber.Ctx) error {
		return c.Redirect().To("/static/favicon.ico")
	})

	ws.app.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().To("/notes")
	})

	api := ws.app.Group("/api/v1")
	api.Post("/notes", ws.handlers.CreateNote)
	api.Put("/notes/:id", ws.handlers.EditNote)
	api.Delete("/notes/:id", ws.handlers.DeleteNote)

	ws.app.Get("/notes", ws.handlers.Notes)
	ws.app.Get("/notes/load-more", ws.handlers.LoadMoreNotes)
	ws.app.Get("/add", ws.handlers.CreateNoteModal)
	ws.app.Get("/edit/:id", ws.handlers.EditNoteModal)

	ws.setupSwagger()
	ws.setupMonitor()

	ws.app.Get("/health", ws.handlers.Health)

	ws.app.Use(func(c fiber.Ctx) error {
		return ws.handlers.Page404(c)
	})

	return nil
}

func (ws *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return ws.app.Listen(":" + port)
}

func (ws *Server) Shutdown(ctx context.Context) error {
	return ws.app.ShutdownWithContext(ctx)
}
