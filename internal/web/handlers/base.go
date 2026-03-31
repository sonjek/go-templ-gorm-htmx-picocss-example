package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/components"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/page"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/view"
	"gorm.io/gorm"
)

type Handlers struct {
	db          *gorm.DB
	noteService *service.NoteService
}

func NewHandler(db *gorm.DB, ns *service.NoteService) *Handlers {
	return &Handlers{
		noteService: ns,
		db:          db,
	}
}

func handleRenderError(err error) {
	if err != nil {
		fmt.Println("Render error: ", err)
	}
}

// Set header and render error message
func sendErrorMsg(w http.ResponseWriter, r *http.Request, errorMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	handleRenderError(components.ErrorMsg(errorMsg).Render(r.Context(), w))
}

func (h *Handlers) Page404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templ.Handler(
		page.Index(view.NotFoundComponent()),
	).ServeHTTP(w, r)
}
