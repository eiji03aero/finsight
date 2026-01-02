package main

import (
	"fmt"
	"log"

	"backend/internal/application/usecase"
	"backend/internal/config"
	"backend/internal/infrastructure/ent"
	"backend/internal/infrastructure/http/handler"
	"backend/internal/infrastructure/http/router"
	"backend/internal/infrastructure/repositories"
	"backend/internal/infrastructure/session"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Database connection
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Name, cfg.Database.Password)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	// Run auto migration (optional - we use manual migrations)
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("Failed to create schema: %v", err)
	// }

	log.Println("Database connection established")

	// 2. Session initialization
	session.InitStore(cfg.Session.Secret)
	log.Println("Session store initialized")

	// 3. Repository layer
	userRepo := repositories.NewUserRepository(client)
	workspaceRepo := repositories.NewWorkspaceRepository(client)

	// 4. Use case layer
	signupUseCase := usecase.NewSignupUseCase(userRepo, workspaceRepo, client)

	// 5. Handler layer
	signupHandler := handler.NewSignupHandler(signupUseCase)

	// 6. Router setup
	r := router.SetupRouter(signupHandler)

	// 7. Server startup
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
