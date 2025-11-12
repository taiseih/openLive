package middleware

import (
	"context"
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
	// FirebaseUIDKey はFirebase UIDのコンテキストキーです
	FirebaseUIDKey ContextKey = "firebaseUID"
)

// AuthMiddleware はFirebase認証ミドルウェアを返します
func AuthMiddleware(firebaseAuth *auth.FirebaseAuth, userUseCase usecase.UserUseCase) echo.MiddlewareFunc {
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

			// Firebase IDトークンを検証
			token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// ユーザー情報を取得または作成
			user, err := userUseCase.GetUserByFirebaseUID(c.Request().Context(), token.UID)
			if err != nil {
				// ユーザーが存在しない場合は、Firebaseから情報を取得して作成
				firebaseUser, err := firebaseAuth.GetUser(context.Background(), token.UID)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user info")
				}

				user, err = userUseCase.CreateUser(
					c.Request().Context(),
					token.UID,
					firebaseUser.Email,
					firebaseUser.DisplayName,
					firebaseUser.PhotoURL,
				)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
				}
			}

			// コンテキストにユーザー情報を設定
			c.Set(string(UserIDKey), user.ID)
			c.Set(string(FirebaseUIDKey), token.UID)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware はオプショナルな認証ミドルウェアです（認証エラーでも続行）
func OptionalAuthMiddleware(firebaseAuth *auth.FirebaseAuth, userUseCase usecase.UserUseCase) echo.MiddlewareFunc {
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

			token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				return next(c)
			}

			user, err := userUseCase.GetUserByFirebaseUID(c.Request().Context(), token.UID)
			if err == nil {
				c.Set(string(UserIDKey), user.ID)
				c.Set(string(FirebaseUIDKey), token.UID)
			}

			return next(c)
		}
	}
}

