package main

import (
	"log"
	"os"

	"backend/internal/application/usecase"
	"backend/internal/domain/service"
	"backend/internal/infrastructure/http/handler"
	"backend/internal/infrastructure/http/router"
)

func main() {
	// 1. 依存関係の初期化 (依存性注入)
	//    Domain Layer → Application Layer → Infrastructure Layer の順

	// Domain Layer
	messageService := service.NewMessageService()

	// Application Layer
	helloWorldUsecase := usecase.NewHelloWorldUsecase(messageService)

	// Infrastructure Layer
	helloWorldHandler := handler.NewHelloWorldHandler(helloWorldUsecase)

	// 2. ルーターのセットアップ
	r := router.SetupRouter(helloWorldHandler)

	// 3. サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
