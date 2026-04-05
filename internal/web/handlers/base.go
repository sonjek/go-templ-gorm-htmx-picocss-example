package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/components"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/page"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/view"
)

type Handlers struct {
	noteService *service.NoteService
}

func NewHandler(ns *service.NoteService) *Handlers {
	return &Handlers{
		noteService: ns,
	}
}

func getCSRFToken(c fiber.Ctx) string {
	return csrf.TokenFromContext(c)
}

func render(c fiber.Ctx, statusCode int, component templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(c.Context(), buf); err != nil {
		slog.Error("Template render error", "error", err)
		return err
	}
	c.Status(statusCode)
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(buf.Bytes())
}

func sendErrorMsg(c fiber.Ctx, errorMsg string) error {
	return render(c, http.StatusBadRequest, components.ErrorMsg(errorMsg))
}

type fieldErrors map[string]string

func sendFieldErrors(c fiber.Ctx, errors fieldErrors) error {
	c.Status(http.StatusBadRequest)
	c.Type("json")
	return c.JSON(fieldErrorsResponse{Errors: errors})
}

type fieldErrorsResponse struct {
	Errors fieldErrors `json:"errors"`
}

// Health godoc
// @Summary      Health check
// @Description  Returns 200 OK for container orchestrators
// @Tags         health
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /health [get]
func (h *Handlers) Health(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func (h *Handlers) Page404(c fiber.Ctx) error {
	return render(c, http.StatusNotFound, page.Index(view.NotFoundComponent(), getCSRFToken(c)))
}
