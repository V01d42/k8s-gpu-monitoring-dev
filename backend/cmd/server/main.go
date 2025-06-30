package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s-gpu-monitoring/internal/handlers"
	"k8s-gpu-monitoring/internal/middleware"
	"k8s-gpu-monitoring/internal/prometheus"
)

func main() {
	// 環境変数から設定を読み込み
	prometheusURL := getEnv("PROMETHEUS_URL", "http://localhost:9090")
	port := getEnv("PORT", "8080")

	log.Printf("Starting GPU Monitoring API Server...")
	log.Printf("Prometheus URL: %s", prometheusURL)
	log.Printf("Server Port: %s", port)

	// Prometheusクライアントを初期化
	promClient := prometheus.NewClient(prometheusURL)

	// ハンドラーを初期化
	gpuHandler := handlers.NewGPUHandler(promClient)

	// Go 1.22の新しいServeMuxを使用
	mux := http.NewServeMux()

	// ルートを設定
	mux.HandleFunc("GET /api/health", gpuHandler.HealthCheck)
	mux.HandleFunc("GET /api/v1/gpu/metrics", gpuHandler.GetGPUMetrics)
	mux.HandleFunc("GET /api/v1/gpu/nodes", gpuHandler.GetGPUNodes)
	mux.HandleFunc("GET /api/v1/gpu/utilization", gpuHandler.GetGPUUtilization)

	// 静的ファイルサーバー（フロントエンド用）
	mux.Handle("GET /", http.FileServer(http.Dir("./static/")))

	// ミドルウェアを適用
	handler := middleware.Chain(
		mux,
		middleware.Logger,
		middleware.CORS,
		middleware.Recovery,
	)

	// HTTPサーバーを設定
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// サーバーを別ゴルーチンで開始
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdownの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// 30秒のタイムアウトでサーバーを停止
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
