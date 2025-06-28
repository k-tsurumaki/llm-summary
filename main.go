package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k-tsurumaki/fuselage"
)

func main() {
	app := fuselage.New()

	// CORS設定
	app.Use(func(next fuselage.HandlerFunc) fuselage.HandlerFunc {
		return func(c *fuselage.Context) error {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")

			if c.Request.Method == "OPTIONS" {
				return c.String(http.StatusOK, "")
			}

			return next(c)
		}
	})

	app.Use(fuselage.RequestID)
	app.Use(fuselage.Logger)
	app.Use(fuselage.Recover)

	// ルート設定
	setupRoutes(app)

	// サーバー起動
	port := ":" + Config.PORT
	server := fuselage.NewServer(port, app)

	// グレースフルシャットダウン
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
