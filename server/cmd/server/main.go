package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taiseihayashizaki/openLive/server/internal/config"
	"github.com/taiseihayashizaki/openLive/server/internal/infrastructure/auth"
	"github.com/taiseihayashizaki/openLive/server/internal/infrastructure/database"
	"github.com/taiseihayashizaki/openLive/server/internal/infrastructure/repository"
	"github.com/taiseihayashizaki/openLive/server/internal/infrastructure/websocket"
	"github.com/taiseihayashizaki/openLive/server/internal/interface/handler"
	appMiddleware "github.com/taiseihayashizaki/openLive/server/internal/interface/middleware"
	"github.com/taiseihayashizaki/openLive/server/internal/usecase"
)

func main() {
	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	db, err := database.NewPostgresDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Firebase認証初期化
	ctx := context.Background()
	firebaseAuth, err := auth.NewFirebaseAuth(
		ctx,
		cfg.Firebase.ProjectID,
		cfg.Firebase.AuthEmulatorHost,
		cfg.Firebase.CredentialsFilePath,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}

	// WebSocketハブ初期化
	hub := websocket.NewHub()
	go hub.Run()

	// リポジトリ初期化
	userRepo := repository.NewUserRepository(db)
	liveAreaRepo := repository.NewLiveAreaRepository(db)
	liveAreaMemberRepo := repository.NewLiveAreaMemberRepository(db)
	liveStreamRepo := repository.NewLiveStreamRepository(db)
	liveStreamViewerRepo := repository.NewLiveStreamViewerRepository(db)

	// ユースケース初期化
	userUseCase := usecase.NewUserUseCase(userRepo)
	liveAreaUseCase := usecase.NewLiveAreaUseCase(liveAreaRepo, liveAreaMemberRepo)
	liveStreamUseCase := usecase.NewLiveStreamUseCase(liveStreamRepo, liveStreamViewerRepo, liveAreaRepo, liveAreaMemberRepo)

	// ハンドラー初期化
	userHandler := handler.NewUserHandler(userUseCase)
	liveAreaHandler := handler.NewLiveAreaHandler(liveAreaUseCase)
	liveStreamHandler := handler.NewLiveStreamHandler(liveStreamUseCase)
	wsHandler := handler.NewWebSocketHandler(hub)

	// Echoサーバー初期化
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(appMiddleware.CORSMiddleware(cfg.CORS.AllowOrigins))

	// ヘルスチェック
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// APIルーティング
	api := e.Group("/api/v1")

	// 認証ミドルウェア
	authMiddleware := appMiddleware.AuthMiddleware(firebaseAuth, userUseCase)

	// ユーザーエンドポイント
	users := api.Group("/users")
	users.GET("/me", userHandler.GetMe, authMiddleware)
	users.PUT("/me", userHandler.UpdateMe, authMiddleware)

	// ライブエリアエンドポイント
	liveareas := api.Group("/liveareas")
	liveareas.POST("", liveAreaHandler.Create, authMiddleware)
	liveareas.GET("", liveAreaHandler.GetAll)
	liveareas.GET("/:id", liveAreaHandler.GetByID)
	liveareas.PUT("/:id", liveAreaHandler.Update, authMiddleware)
	liveareas.DELETE("/:id", liveAreaHandler.Delete, authMiddleware)
	liveareas.POST("/:id/invite", liveAreaHandler.InviteMember, authMiddleware)
	liveareas.DELETE("/:id/members/:userId", liveAreaHandler.RemoveMember, authMiddleware)
	liveareas.GET("/:id/members", liveAreaHandler.GetMembers)

	// ライブ配信エンドポイント
	liveareas.POST("/:id/streams", liveStreamHandler.Create, authMiddleware)
	liveareas.GET("/:id/streams", liveStreamHandler.GetByLiveArea)

	streams := api.Group("/streams")
	streams.GET("/live", liveStreamHandler.GetLiveStreams)
	streams.GET("/:id", liveStreamHandler.GetByID)
	streams.POST("/:id/start", liveStreamHandler.Start, authMiddleware)
	streams.DELETE("/:id", liveStreamHandler.End, authMiddleware)
	streams.GET("/:id/viewers", liveStreamHandler.GetViewers)
	streams.GET("/:id/viewers/count", liveStreamHandler.GetViewerCount)

	// WebSocketエンドポイント
	ws := api.Group("/ws")
	ws.GET("/streams/:id", wsHandler.HandleWebSocket, authMiddleware)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	// グレースフルシャットダウン
	go func() {
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 10秒のタイムアウトでシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

