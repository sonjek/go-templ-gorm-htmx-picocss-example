package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/sonjek/go-full-stack-example/internal/service"
	"github.com/sonjek/go-full-stack-example/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testHandlers(t *testing.T) *Handlers {
	t.Helper()
	db, err := storage.NewInMemoryDbStorage()
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.DBMigrate(db); err != nil {
		t.Fatal(err)
	}
	return NewHandler(service.NewNoteService(db, 2))
}

func testApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New()
	sessionMiddleware, sessionStore := session.NewWithStore()
	app.Use(sessionMiddleware)
	app.Use(csrf.New(csrf.Config{
		Extractor: extractors.Chain(
			extractors.FromHeader("X-Csrf-Token"),
			extractors.FromForm("_csrf"),
		),
		Session: sessionStore,
	}))
	h := testHandlers(t)
	app.Get("/notes", h.Notes)
	app.Get("/notes/load-more", h.LoadMoreNotes)
	app.Get("/add", h.CreateNoteModal)
	app.Post("/notes", h.CreateNote)
	app.Get("/edit/:id", h.EditNoteModal)
	app.Put("/notes/:id", h.EditNote)
	app.Delete("/notes/:id", h.DeleteNote)
	return app
}

func getCSRFCookie(app *fiber.App) (cookie, token string) {
	req := httptest.NewRequest(http.MethodGet, "/add", nil)
	resp, _ := app.Test(req)
	cookies := []string{}
	for _, c := range resp.Cookies() {
		cookies = append(cookies, c.Name+"="+c.Value)
		if c.Name == "csrf_" {
			token = c.Value
		}
	}
	cookie = strings.Join(cookies, "; ")
	return cookie, token
}

func Test_sendErrorMsg(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return sendErrorMsg(c, "MSG")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "MSG", "error message should be in response")
	assert.Contains(t, bodyStr, "Error:", "should contain Error label")
}

func Test_Page404(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return testHandlers(t).Page404(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func Test_Notes(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		app := testApp(t)

		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid cursor", func(t *testing.T) {
		app := testApp(t)

		req := httptest.NewRequest(http.MethodGet, "/notes?cursor=abc", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Invalid cursor parameter")
	})
}

func Test_LoadMoreNotes(t *testing.T) {
	app := testApp(t)

	req := httptest.NewRequest(http.MethodGet, "/notes/load-more", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateNoteModal(t *testing.T) {
	app := testApp(t)

	req := httptest.NewRequest(http.MethodGet, "/add", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_CreateNote_Success(t *testing.T) {
	app := testApp(t)
	cookie, token := getCSRFCookie(app)

	form := url.Values{}
	form.Set("title", "Test Note")
	form.Set("body", "Test body content")
	form.Set("_csrf", token)

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Test Note")
}

func Test_CreateNote_EmptyTitle(t *testing.T) {
	app := testApp(t)
	cookie, token := getCSRFCookie(app)

	form := url.Values{}
	form.Set("title", "")
	form.Set("body", "Test body")
	form.Set("_csrf", token)

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_CreateNote_EmptyBody(t *testing.T) {
	app := testApp(t)
	cookie, token := getCSRFCookie(app)

	form := url.Values{}
	form.Set("title", "Test")
	form.Set("body", "")
	form.Set("_csrf", token)

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_CreateNote_DuplicateTitle(t *testing.T) {
	app := testApp(t)
	cookie, token := getCSRFCookie(app)

	form := url.Values{}
	form.Set("title", "Duplicate")
	form.Set("body", "First")
	form.Set("_csrf", token)

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cookie, token = getCSRFCookie(app)
	form.Set("_csrf", token)

	req = httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err = app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), service.MsgTitleAlreadyExists)
}

func Test_CreateNote_WhitespaceOnlyTitle(t *testing.T) {
	app := testApp(t)
	cookie, token := getCSRFCookie(app)

	form := url.Values{}
	form.Set("title", "   ")
	form.Set("body", "Test body")
	form.Set("_csrf", token)

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func Test_CreateNote_MissingCSRF(t *testing.T) {
	app := testApp(t)

	form := url.Values{}
	form.Set("title", "Test")
	form.Set("body", "Test body")

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func Test_EditNoteModal(t *testing.T) {
	t.Run("invalid ID", func(t *testing.T) {
		app := testApp(t)

		req := httptest.NewRequest(http.MethodGet, "/edit/abc", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		app := testApp(t)

		req := httptest.NewRequest(http.MethodGet, "/edit/999", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Note not found")
	})
}

func Test_EditNote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := testApp(t)
		cookie, token := getCSRFCookie(app)

		form := url.Values{}
		form.Set("title", "Original")
		form.Set("body", "Original body")
		form.Set("_csrf", token)

		req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", cookie)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookie, token = getCSRFCookie(app)
		form = url.Values{}
		form.Set("title", "Updated")
		form.Set("body", "Updated body")
		form.Set("_csrf", token)

		req = httptest.NewRequest(http.MethodPut, "/notes/1", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", cookie)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Updated")
	})

	t.Run("not found", func(t *testing.T) {
		app := testApp(t)
		cookie, token := getCSRFCookie(app)

		form := url.Values{}
		form.Set("title", "Updated")
		form.Set("body", "Updated body")
		form.Set("_csrf", token)

		req := httptest.NewRequest(http.MethodPut, "/notes/999", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", cookie)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Note not found")
	})

	t.Run("missing CSRF token", func(t *testing.T) {
		app := testApp(t)

		form := url.Values{}
		form.Set("title", "Updated")
		form.Set("body", "Updated body")

		req := httptest.NewRequest(http.MethodPut, "/notes/1", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func Test_DeleteNote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := testApp(t)
		cookie, token := getCSRFCookie(app)

		form := url.Values{}
		form.Set("title", "To Delete")
		form.Set("body", "Delete me")
		form.Set("_csrf", token)

		req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", cookie)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		cookie, token = getCSRFCookie(app)

		req = httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
		req.Header.Set("X-Csrf-Token", token)
		req.Header.Set("Cookie", cookie)
		resp, err = app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		app := testApp(t)
		cookie, token := getCSRFCookie(app)

		req := httptest.NewRequest(http.MethodDelete, "/notes/999", nil)
		req.Header.Set("X-Csrf-Token", token)
		req.Header.Set("Cookie", cookie)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "Note not found")
	})

	t.Run("missing CSRF token", func(t *testing.T) {
		app := testApp(t)

		req := httptest.NewRequest(http.MethodDelete, "/notes/1", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
}

func Test_ParseNoteID(t *testing.T) {
	app := fiber.New()
	app.Get("/notes/:id", func(c fiber.Ctx) error {
		id, err := parseNoteID(c)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(strings.Repeat("0", id))
	})

	t.Run("valid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/5", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("zero ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/0", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("negative ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/-1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
