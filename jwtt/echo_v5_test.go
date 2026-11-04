package jwtt

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func TestEchoV5JWTVerification(t *testing.T) {
	original := DefaultJWTConfig
	t.Cleanup(func() { DefaultJWTConfig = original })
	DefaultJWTConfig = JWTConfig{
		SigningMethod: "HS256", SigningKey: []byte("test-only-signing-key"),
		ContextKey: "user", TokenLookup: "header:Authorization:Bearer ",
		ErrorConfig: ErrorConfig{ErrorHandler: func(c *echo.Context, err error) error { return c.NoContent(401) }},
	}
	sign := func(method jwt.SigningMethod, key []byte, expires time.Time) string {
		t.Helper()
		token, err := jwt.NewWithClaims(method, jwt.MapClaims{"exp": expires.Unix()}).SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
	valid := sign(jwt.SigningMethodHS256, DefaultJWTConfig.SigningKey, time.Now().Add(time.Hour))
	for _, tc := range []struct {
		name, token string
		want        int
	}{
		{"valid", valid, 200},
		{"missing", "", 401},
		{"malformed", "invalid", 401},
		{"expired", sign(jwt.SigningMethodHS256, DefaultJWTConfig.SigningKey, time.Now().Add(-time.Hour)), 401},
		{"wrong-key", sign(jwt.SigningMethodHS256, []byte("other-key"), time.Now().Add(time.Hour)), 401},
		{"wrong-algorithm", sign(jwt.SigningMethodHS384, DefaultJWTConfig.SigningKey, time.Now().Add(time.Hour)), 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Use(echojwt.WithConfig(JwtConfig(nil)))
			e.GET("/", func(c *echo.Context) error {
				if !ParserToken(c).IsPresent() {
					t.Error("verified token absent from context")
				}
				return c.NoContent(200)
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status = %d; want %d", rec.Code, tc.want)
			}
		})
	}
}
