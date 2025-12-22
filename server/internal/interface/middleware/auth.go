package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/taiseihayashizaki/openLive/server/internal/infrastructure/auth"
	"github.com/taiseihayashizaki/openLive/server/internal/usecase"
)

// ContextKey はコンテキストキーの型です
type ContextKey string

const (
	// UserIDKey はユーザーIDのコンテキストキーです
	UserIDKey ContextKey = "userID"
)

// AuthMiddleware はJWT認証ミドルウェアを返します
func AuthMiddleware(jwtAuth *auth.JWTAuth, userUseCase usecase.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// "Bearer <token>" 形式のトークンを取得
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			idToken := parts[1]

			// JWTトークンを検証
			userID, err := jwtAuth.ParseToken(idToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// ユーザー情報を取得
			user, err := userUseCase.GetUserByID(c.Request().Context(), userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			// コンテキストにユーザー情報を設定
			c.Set(string(UserIDKey), user.ID)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware はオプショナルな認証ミドルウェアです（認証エラーでも続行）
func OptionalAuthMiddleware(jwtAuth *auth.JWTAuth, userUseCase usecase.UserUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return next(c)
			}

			idToken := parts[1]

			userID, err := jwtAuth.ParseToken(idToken)
			if err != nil {
				return next(c)
			}

			user, err := userUseCase.GetUserByID(c.Request().Context(), userID)
			if err == nil {
				c.Set(string(UserIDKey), user.ID)
			}

			return next(c)
		}
	}
}

