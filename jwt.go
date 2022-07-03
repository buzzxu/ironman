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
	extractor := jwtGet(parts[1], jwtConfig.AuthScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	case "header":
		extractor = jwtFromHeader(parts[1], jwtConfig.AuthScheme)

	}
	DefaultJWTConfig.jwtExtractor = extractor
	DefaultJWTConfig.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != jwtConfig.SigningMethod {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
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

func ParserTokenUnverified(c echo.Context, name string, authScheme string) (*Optional, error) {
	auth, err := GetToken(c, name, authScheme)
	if err != nil {
		return nil, err
	}
	return TokenOf(auth)
}

func TokenOf(auth string) (*Optional, error) {
	if len(auth) == 0 {
		return OptionalOfNil(), nil
	}
	token := new(jwt.Token)
	claims := reflect.ValueOf(DefaultJWTConfig.Claims).Interface().(jwt.Claims)
	var err error
	token, _, err = new(jwt.Parser).ParseUnverified(auth, claims)
	if err == nil {
		return OptionalOf(token), nil
	}
	return OptionalOfNil(), err
}
func GetToken(c echo.Context, name string, authScheme string) (string, error) {
	token := c.Request().Header.Get(name)
	if token == "" {
		token = c.QueryParam(name)
		if token == "" {
			cookie, err := c.Cookie(name)
			if err != nil {
				return "", nil
			}
			token = cookie.Value
		}
	}
	if token != "" {
		l := len(authScheme)
		if len(token) > l+1 && token[:l] == authScheme {
			return token[l+1:], nil
		} else {
			return token, nil
		}
	}
	return "", nil
}
func jwtGet(name string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		return GetToken(c, name, authScheme)
	}
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if auth != "" {
			if len(auth) > l+1 && auth[:l] == authScheme {
				return auth[l+1:], nil
			} else {
				return auth, nil
			}
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
