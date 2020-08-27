package ironman

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Claims        jwt.Claims
	SigningMethod string
	ContextKey    string
	SigningKey    []byte
	TokenLookup   string
	AuthScheme    string
	jwtExtractor
	keyFunc jwt.Keyfunc
}
type jwtExtractor func(echo.Context) (string, error)

// DefaultJWTConfig 默认JWT配置
var DefaultJWTConfig = JWTConfig{
	Claims:        nil,
	SigningMethod: jwt.SigningMethodHS256.Name,
	ContextKey:    "user",
	SigningKey:    []byte("ironman"),
	TokenLookup:   "header:" + echo.HeaderAuthorization,
	AuthScheme:    "Bearer",
}

//JwtConfig 设置JWT 配置
func JwtConfig(skiper middleware.Skipper) (jwtConfig middleware.JWTConfig) {
	jwtConfig = middleware.JWTConfig{
		Skipper:       skiper,
		Claims:        DefaultJWTConfig.Claims,
		SigningMethod: DefaultJWTConfig.SigningMethod,
		ContextKey:    DefaultJWTConfig.ContextKey,
		SigningKey:    DefaultJWTConfig.SigningKey,
		TokenLookup:   DefaultJWTConfig.TokenLookup,
		AuthScheme:    DefaultJWTConfig.AuthScheme,
	}
	parts := strings.Split(jwtConfig.TokenLookup, ":")
	extractor := jwtFromHeader(parts[1], jwtConfig.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}
	DefaultJWTConfig.jwtExtractor = extractor
	DefaultJWTConfig.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != jwtConfig.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		return jwtConfig.SigningKey, nil
	}
	return jwtConfig
}

//GenerateToken 生成token
func GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.GetSigningMethod(DefaultJWTConfig.SigningMethod), claims)
	t, error := token.SignedString(DefaultJWTConfig.SigningKey)
	if error != nil {
		return "", error
	}
	return t, nil
}

//ParserToken 解析token
func ParserToken(c echo.Context) *Optional {
	auth, err := DefaultJWTConfig.jwtExtractor(c)
	if len(auth) == 0 {
		return OptionalOfNil()
	}
	token := new(jwt.Token)
	claims := reflect.ValueOf(DefaultJWTConfig.Claims).Interface().(jwt.Claims)
	token, err = jwt.ParseWithClaims(auth, claims, DefaultJWTConfig.keyFunc)
	if err == nil && token.Valid {
		return OptionalOf(token)
	}
	return OptionalOfNil()
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", nil
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		if token == "" {
			return "", nil
		}
		return token, nil
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from the named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", nil
		}
		return cookie.Value, nil
	}
}
