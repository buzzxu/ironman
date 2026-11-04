package ironman

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buzzxu/ironman/echomw"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestEchoV5ErrorResponses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"application", NewError(409, "Conflict", map[string]string{"reason": "duplicate"}), 409},
		{"wrapped", fmt.Errorf("wrapped: %w", echo.ErrUnauthorized), 401},
		{"http", echo.NewHTTPError(400, "invalid request"), 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			core, logs := observer.New(zap.WarnLevel)
			e := echo.New()
			e.HTTPErrorHandler = httpErrorHandler
			e.Use(echomw.ZapLogger(zap.New(core)))
			e.GET("/", func(c *echo.Context) error { return tc.err })
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d; want %d: %s", rec.Code, tc.status, rec.Body.String())
			}
			entries := logs.All()
			if len(entries) != 1 || entries[0].ContextMap()["status"] != int64(tc.status) {
				t.Fatalf("incorrect access log: %+v", entries)
			}
		})
	}
}

func TestEchoV5HeadAndCommittedResponse(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.HEAD("/", func(c *echo.Context) error { return echo.ErrNotFound })
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/", nil))
	if rec.Code != 404 || rec.Body.Len() != 0 {
		t.Fatalf("HEAD response: %d %q", rec.Code, rec.Body.String())
	}
	e.GET("/committed", func(c *echo.Context) error {
		if err := c.String(202, "accepted"); err != nil {
			return err
		}
		return echo.ErrInternalServerError
	})
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/committed", nil))
	if rec.Code != 202 || rec.Body.String() != "accepted" {
		t.Fatalf("committed response overwritten: %d %q", rec.Code, rec.Body.String())
	}
}

func TestEchoV5TrailingSlashAndPeerIP(t *testing.T) {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/peer", func(c *echo.Context) error { return c.String(200, c.RealIP()) })
	req := httptest.NewRequest(http.MethodGet, "/peer/", nil)
	req.RemoteAddr = "192.0.2.7:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.9")
	req.Header.Set("X-Real-IP", "203.0.113.10")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != 200 || rec.Body.String() != "192.0.2.7" {
		t.Fatalf("untrusted IP accepted or routing failed: %d %q", rec.Code, rec.Body.String())
	}
}

func TestEchoV5ProtectedStaticFiles(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "secret.txt"), []byte("private-marker"), 0600); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	group := e.Group("/private", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Header.Get("Authorization") != "test-token" {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	})
	group.Static("/", root)
	for _, path := range []string{"/private/secret.txt", "/private%2fsecret.txt", "/%70rivate/secret.txt", "/private/%2e%2e/private/secret.txt", "/private/%252fsecret.txt"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if strings.Contains(rec.Body.String(), "private-marker") {
			t.Fatalf("unauthenticated static disclosure at %s", path)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/private/secret.txt", nil)
	req.Header.Set("Authorization", "test-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != 200 || rec.Body.String() != "private-marker" {
		t.Fatalf("authorized file unavailable: %d %q", rec.Code, rec.Body.String())
	}
}
