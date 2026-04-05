//go:build swagger

package web

import (
	"github.com/gofiber/contrib/v3/swaggo"
	_ "github.com/sonjek/go-full-stack-example/docs"
)

func (ws *Server) setupSwagger() {
	ws.app.Get("/swagger/*", swaggo.HandlerDefault)
}
