package jwtt

import (
	"errors"
	"github.com/buzzxu/ironman"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// 预定义错误类型
var (
	ErrJWTMissing      = errors.New("missing or malformed jwt")
	ErrJWTInvalid      = errors.New("invalid or expired jwt")
	ErrJWTUnauthorized = errors.New("unauthorized")
)

// JWTConfig JWT配置
type JWTConfig struct {
	Claims        jwt.Claims
	SigningMethod string
	ContextKey    string
	SigningKey    []byte
	TokenLookup   string
	KeyFunc       jwt.Keyfunc
	ExpiresIn     time.Duration // 默认令牌过期时间
	ErrorConfig   ErrorConfig   // 错误处理配置
}

// ErrorConfig 错误处理配置
type ErrorConfig struct {
	// 错误处理函数，返回正确的响应
	ErrorHandler func(c echo.Context, err error) error
	// 各类错误对应的HTTP状态码
	ErrorStatusCodes map[string]int
	// 自定义错误消息
	ErrorMessages map[string]string
}

// 默认错误状态码和消息
var (
	defaultErrorStatusCodes = map[string]int{
		ErrJWTMissing.Error():      http.StatusBadRequest,
		ErrJWTInvalid.Error():      http.StatusUnauthorized,
		ErrJWTUnauthorized.Error(): http.StatusForbidden,
	}

	defaultErrorMessages = map[string]string{
		ErrJWTMissing.Error():      "Missing or malformed JWT",
		ErrJWTInvalid.Error():      "Invalid or expired JWT",
		ErrJWTUnauthorized.Error(): "Unauthorized access",
	}
)

// DefaultErrorConfig 默认错误处理配置
var DefaultErrorConfig = ErrorConfig{
	ErrorHandler:     defaultErrorHandler,
	ErrorStatusCodes: defaultErrorStatusCodes,
	ErrorMessages:    defaultErrorMessages,
}

// defaultErrorHandler 默认错误处理函数
func defaultErrorHandler(c echo.Context, err error) error {
	var message string
	var status int

	// 使用局部变量而不是引用 DefaultJWTConfig
	switch err.Error() {
	case ErrJWTMissing.Error():
		message = defaultErrorMessages[ErrJWTMissing.Error()]
		status = defaultErrorStatusCodes[ErrJWTMissing.Error()]
	case ErrJWTInvalid.Error():
		message = defaultErrorMessages[ErrJWTInvalid.Error()]
		status = defaultErrorStatusCodes[ErrJWTInvalid.Error()]
	default:
		message = defaultErrorMessages[ErrJWTUnauthorized.Error()]
		status = defaultErrorStatusCodes[ErrJWTUnauthorized.Error()]
	}

	return c.JSON(status, map[string]string{
		"error": message,
	})
}

// DefaultJWTConfig 默认JWT配置
var DefaultJWTConfig = JWTConfig{
	SigningMethod: "HS256", // jwt v5 uses string constants
	ContextKey:    "user",
	SigningKey:    []byte("ironman"),
	TokenLookup:   "header:Authorization,query:Authorization,param:Authorization",
	ExpiresIn:     time.Hour * 168, // 默认168小时过期
	Claims:        nil,
	KeyFunc:       nil,
	ErrorConfig:   DefaultErrorConfig,
}

// CustomClaims 自定义JWT Claims结构体
type CustomClaims struct {
	jwt.RegisteredClaims
	// 可以添加自定义字段
	UserID   string `json:"userId,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}

// NewClaims 创建一个新的空Claims对象
func NewClaims() jwt.Claims {
	return jwt.MapClaims{}
}

// NewCustomClaims 创建一个带有自定义字段的Claims对象
func NewCustomClaims() *CustomClaims {
	if customClaims, ok := DefaultJWTConfig.Claims.(*CustomClaims); ok {
		return customClaims
	}
	return &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
	}
}

// SetErrorConfig 设置错误处理配置
func SetErrorConfig(config ErrorConfig) {
	DefaultJWTConfig.ErrorConfig = config
}

// SetErrorHandler 设置错误处理函数
func SetErrorHandler(handler func(c echo.Context, err error) error) {
	DefaultJWTConfig.ErrorConfig.ErrorHandler = handler
}

// SetErrorStatusCode 设置特定错误的HTTP状态码
func SetErrorStatusCode(errType error, statusCode int) {
	DefaultJWTConfig.ErrorConfig.ErrorStatusCodes[errType.Error()] = statusCode
}

// SetErrorMessage 设置特定错误的消息
func SetErrorMessage(errType error, message string) {
	DefaultJWTConfig.ErrorConfig.ErrorMessages[errType.Error()] = message
}

// JwtConfig 设置JWT 配置
func JwtConfig(skipper middleware.Skipper) echojwt.Config {
	return JwtConfigWithClaims(skipper, func(c echo.Context) jwt.Claims {
		return NewCustomClaims()
	})
}

// JwtConfigWithClaims 设置JWT 配置，并使用自定义的Claims
func JwtConfigWithClaims(skipper middleware.Skipper, claimsFunc func(c echo.Context) jwt.Claims) echojwt.Config {
	return echojwt.Config{
		Skipper:       skipper,
		SigningMethod: DefaultJWTConfig.SigningMethod,
		ContextKey:    DefaultJWTConfig.ContextKey,
		SigningKey:    DefaultJWTConfig.SigningKey,
		TokenLookup:   DefaultJWTConfig.TokenLookup,
		KeyFunc:       DefaultJWTConfig.KeyFunc,
		NewClaimsFunc: claimsFunc,
		ErrorHandler:  DefaultJWTConfig.ErrorConfig.ErrorHandler,
	}
}

// GenerateToken 生成token，使用默认过期时间
func GenerateToken(claims jwt.Claims) (string, error) {
	// 如果是自定义Claims并且没有设置过期时间，设置默认过期时间
	if customClaims, ok := claims.(*CustomClaims); ok {
		if customClaims.ExpiresAt == nil || customClaims.ExpiresAt.Time.IsZero() {
			customClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(DefaultJWTConfig.ExpiresIn))
		}
	} else if mapClaims, ok := claims.(jwt.MapClaims); ok {
		// 检查MapClaims是否有过期时间
		if _, exists := mapClaims["exp"]; !exists {
			mapClaims["exp"] = time.Now().Add(DefaultJWTConfig.ExpiresIn).Unix()
		}
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(DefaultJWTConfig.SigningMethod), claims)
	return token.SignedString(DefaultJWTConfig.SigningKey)
}

// GenerateTokenWithExpiration 生成有指定过期时间的token
func GenerateTokenWithExpiration(claims jwt.Claims, expiresIn time.Duration) (string, error) {
	// 设置过期时间
	if customClaims, ok := claims.(*CustomClaims); ok {
		customClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))
	} else if mapClaims, ok := claims.(jwt.MapClaims); ok {
		mapClaims["exp"] = time.Now().Add(expiresIn).Unix()
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(DefaultJWTConfig.SigningMethod), claims)
	return token.SignedString(DefaultJWTConfig.SigningKey)
}

// ParserToken 解析token
func ParserToken(c echo.Context) *ironman.Optional {
	// Get the token from the context - it will be placed there by echo-jwt middleware
	token, ok := c.Get(DefaultJWTConfig.ContextKey).(*jwt.Token)
	//claims := reflect.ValueOf(DefaultJWTConfig.Claims).Interface().(jwt.Claims)
	if !ok || token == nil {
		return ironman.OptionalOfNil()
	}

	// We don't need to check token.Valid as the middleware has already verified it
	return ironman.OptionalOf(token)
}

// ParserTokenUnverified 解析token但不验证签名
func ParserTokenUnverified(c echo.Context, name string) (*ironman.Optional, error) {
	auth, err := GetToken(c, name)
	if err != nil {
		return nil, err
	}
	return TokenOf(auth)
}

// TokenOf 从字符串创建token对象
func TokenOf(auth string) (*ironman.Optional, error) {
	if len(auth) == 0 {
		return ironman.OptionalOfNil(), nil
	}

	// Create parser with default claims
	parser := jwt.NewParser()
	// Parse the token without verification
	token, _, err := parser.ParseUnverified(auth, DefaultJWTConfig.Claims)
	if err == nil {
		return ironman.OptionalOf(token), nil
	}
	return ironman.OptionalOfNil(), err
}

// GetToken 从请求的不同位置获取token
func GetToken(c echo.Context, name string) (string, error) {
	// Check header
	token := c.Request().Header.Get(name)
	if token != "" {
		return token, nil
	}

	// Check query param
	token = c.QueryParam(name)
	if token != "" {
		return token, nil
	}

	// Check cookie
	cookie, err := c.Cookie(name)
	if err == nil && cookie != nil {
		return cookie.Value, nil
	}

	return "", nil
}

// IsTokenExpired 检查token是否过期
func IsTokenExpired(token *jwt.Token) bool {
	claims, ok := token.Claims.(jwt.Claims)
	if !ok {
		return true
	}

	// 检查过期时间
	if expTime, err := claims.GetExpirationTime(); err == nil {
		return expTime.Before(time.Now())
	}

	return true // 如果无法获取过期时间，默认为已过期
}

// GetTokenRemainingTime 获取token剩余有效时间
func GetTokenRemainingTime(token *jwt.Token) (time.Duration, error) {
	claims, ok := token.Claims.(jwt.Claims)
	if !ok {
		return 0, jwt.ErrTokenInvalidClaims
	}

	expTime, err := claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}

	return time.Until(expTime.Time), nil
}

// Optional 为了保持与原代码兼容，保留Optional类型的支持
// 注意：你需要确保你有Optional的相关实现（OptionalOf和OptionalOfNil函数）
