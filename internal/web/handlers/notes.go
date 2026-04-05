package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/components"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/page"
	"github.com/sonjek/go-full-stack-example/internal/web/templ/view"
)

const (
	maxTitleLen = 65
	maxBodyLen  = 500
)

// Notes godoc
// @Summary      List notes
// @Description  Get paginated notes with cursor-based pagination
// @Tags         notes
// @Produce      html
// @Param        cursor  query     int  false  "Cursor for pagination"
// @Success      200     {string} html  "Notes page rendered as HTML"
// @Failure      400     {string} string "Error message"
// @Router       /notes [get]
func (h *Handlers) Notes(c fiber.Ctx) error {
	cursor, err := parseCursor(c, 0)
	if err != nil {
		return sendErrorMsg(c, "Invalid cursor parameter")
	}

	notes, err := h.noteService.LoadMore(cursor)
	if err != nil {
		slog.Error("Failed to load notes", "error", err)
		return sendErrorMsg(c, "Failed to load notes")
	}

	return render(c, http.StatusOK, page.Index(view.NotesView(notes), getCSRFToken(c)))
}

// LoadMoreNotes godoc
// @Summary      Load more notes
// @Description  Load additional notes for infinite scroll (HTMX fragment)
// @Tags         notes
// @Produce      html
// @Param        cursor  query     int  false  "Cursor for pagination"
// @Success      200     {string} html  "Notes list fragment"
// @Failure      400     {string} string "Error message"
// @Router       /notes/load-more [get]
func (h *Handlers) LoadMoreNotes(c fiber.Ctx) error {
	cursor, err := parseCursor(c, 0)
	if err != nil {
		return sendErrorMsg(c, "Invalid cursor parameter")
	}

	notes, err := h.noteService.LoadMore(cursor)
	if err != nil {
		slog.Error("Failed to load more notes", "error", err, "cursor", cursor)
		return sendErrorMsg(c, "Failed to load notes")
	}

	return render(c, http.StatusOK, components.NotesList(notes))
}

// CreateNoteModal godoc
// @Summary      Get create note modal
// @Description  Returns the HTML for the create note modal dialog
// @Tags         notes
// @Produce      html
// @Success      200  {string}  html  "Modal HTML fragment"
// @Router       /add [get]
func (h *Handlers) CreateNoteModal(c fiber.Ctx) error {
	return render(c, http.StatusOK, components.ModalAddNote())
}

// CreateNote godoc
// @Summary      Create a note
// @Description  Create a new note from form data
// @Tags         notes
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        title  formData  string  true  "Note title"
// @Param        body   formData  string  true  "Note body content"
// @Success      200    {string}  html    "Created note as HTML card"
// @Failure      400    {string}  string  "Error message"
// @Router       /api/v1/notes [post]
func (h *Handlers) CreateNote(c fiber.Ctx) error {
	title := strings.TrimSpace(c.FormValue("title"))
	body := strings.TrimSpace(c.FormValue("body"))

	if errs := validateNote(title, body); len(errs) > 0 {
		return sendFieldErrors(c, errs)
	}

	note, err := h.noteService.Create(title, body)
	if err != nil {
		if service.IsDuplicateTitle(err) {
			return sendFieldErrors(c, fieldErrors{"title": service.MsgTitleAlreadyExists})
		}
		slog.Error("Failed to create note", "error", err)
		return sendErrorMsg(c, "Failed to create note")
	}

	return render(c, http.StatusOK, components.NoteItem(note))
}

// EditNoteModal godoc
// @Summary      Get edit note modal
// @Description  Returns the HTML for the edit note modal dialog
// @Tags         notes
// @Produce      html
// @Param        id   path      int  true  "Note ID"
// @Success      200  {string}  html  "Modal HTML fragment"
// @Failure      400  {string}  string  "Error message"
// @Router       /edit/{id} [get]
func (h *Handlers) EditNoteModal(c fiber.Ctx) error {
	noteID, err := parseNoteID(c)
	if err != nil {
		return sendErrorMsg(c, "Invalid note ID")
	}

	note, err := h.noteService.Get(noteID)
	if err != nil {
		if service.IsRecordNotFound(err) {
			return sendErrorMsg(c, "Note not found")
		}
		slog.Error("Failed to get note", "error", err)
		return sendErrorMsg(c, "Failed to load note")
	}

	return render(c, http.StatusOK, components.ModalEditNote(note))
}

// EditNote godoc
// @Summary      Update a note
// @Description  Update an existing note by ID
// @Tags         notes
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        id     path      int     true  "Note ID"
// @Param        title  formData  string  true  "Note title"
// @Param        body   formData  string  true  "Note body content"
// @Success      200    {string}  html    "Updated note as HTML card"
// @Failure      400    {string}  string  "Error message"
// @Router       /api/v1/notes/{id} [put]
func (h *Handlers) EditNote(c fiber.Ctx) error {
	noteID, err := parseNoteID(c)
	if err != nil {
		return sendErrorMsg(c, "Invalid note ID")
	}

	title := strings.TrimSpace(c.FormValue("title"))
	body := strings.TrimSpace(c.FormValue("body"))

	if errs := validateNote(title, body); len(errs) > 0 {
		return sendFieldErrors(c, errs)
	}

	note, err := h.noteService.FindAndUpdate(noteID, title, body)
	if err != nil {
		if service.IsRecordNotFound(err) {
			return sendErrorMsg(c, "Note not found")
		}
		if service.IsDuplicateTitle(err) {
			return sendFieldErrors(c, fieldErrors{"title": service.MsgTitleAlreadyExists})
		}
		slog.Error("Failed to update note", "error", err)
		return sendErrorMsg(c, "Failed to update note")
	}

	return render(c, http.StatusOK, components.NoteItem(note))
}

// DeleteNote godoc
// @Summary      Delete a note
// @Description  Hard-delete a note by ID
// @Tags         notes
// @Param        id  path  int  true  "Note ID"
// @Success      200  "Note deleted successfully"
// @Failure      400  {string} string "Error message"
// @Router       /api/v1/notes/{id} [delete]
func (h *Handlers) DeleteNote(c fiber.Ctx) error {
	noteID, err := parseNoteID(c)
	if err != nil {
		return sendErrorMsg(c, "Invalid note ID")
	}

	if err := h.noteService.Delete(noteID); err != nil {
		if service.IsRecordNotFound(err) {
			return sendErrorMsg(c, "Note not found")
		}
		slog.Error("Failed to delete note", "error", err)
		return sendErrorMsg(c, "Failed to delete note")
	}

	return c.SendStatus(http.StatusOK)
}

func validateNote(title, body string) fieldErrors {
	errs := fieldErrors{}

	if title == "" {
		errs["title"] = "Title is empty"
	} else if len(title) > maxTitleLen {
		errs["title"] = "Title is too long. Maximum " + strconv.Itoa(maxTitleLen) + " characters."
	}

	if body == "" {
		errs["body"] = "Body is empty"
	} else if len(body) > maxBodyLen {
		errs["body"] = "Body is too long. Maximum " + strconv.Itoa(maxBodyLen) + " characters."
	}

	return errs
}

func parseCursor(c fiber.Ctx, defaultVal int) (int, error) {
	if p := c.Query("cursor"); p != "" {
		parsed, err := strconv.Atoi(p)
		if err != nil {
			slog.Warn("Invalid cursor parameter", "cursor", p)
			return 0, errors.New("invalid cursor parameter")
		}
		return parsed, nil
	}
	return defaultVal, nil
}

func parseNoteID(c fiber.Ctx) (int, error) {
	p := c.Params("id")
	if p == "" {
		return 0, errors.New("note ID is empty")
	}
	noteID, err := strconv.Atoi(p)
	if err != nil {
		slog.Warn("Invalid note ID format", "id", p)
		return 0, errors.New("invalid note ID")
	}
	if noteID < 1 {
		slog.Warn("Note ID out of range", "id", noteID)
		return 0, errors.New("invalid note ID")
	}
	return noteID, nil
}
