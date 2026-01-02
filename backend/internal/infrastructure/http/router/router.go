package router

import (
	"backend/internal/infrastructure/http/handler"
	"backend/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter はGinルーターをセットアップする
func SetupRouter(signupHandler *handler.SignupHandler) *gin.Engine {
	// 1. Ginエンジンの初期化
	r := gin.Default()

	// 2. ミドルウェアの適用
	r.Use(middleware.CORS())

	// 3. APIエンドポイントの登録
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", signupHandler.Signup)
			auth.GET("/session", signupHandler.GetSession)
		}
	}

	return r
}
