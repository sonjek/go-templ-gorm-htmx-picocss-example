package web

import "github.com/gofiber/contrib/v3/monitor"

func (ws *Server) setupMonitor() {
	ws.app.Get("/monitor/api", monitor.New(monitor.Config{
		APIOnly: true,
	}))
}
