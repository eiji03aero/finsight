package router

import (
	"backend/internal/infrastructure/http/handler"
	"backend/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter はGinルーターをセットアップする
func SetupRouter(helloWorldHandler *handler.HelloWorldHandler) *gin.Engine {
	// 1. Ginエンジンの初期化
	r := gin.Default()

	// 2. ミドルウェアの適用
	r.Use(middleware.CORS())

	// 3. エンドポイントの登録
	r.GET("/hello_world", helloWorldHandler.Handle)

	// 4. 将来的なエンドポイント追加のための構造
	// api := r.Group("/api")
	// {
	//     v1 := api.Group("/v1")
	//     {
	//         v1.GET("/hello_world", helloWorldHandler.Handle)
	//         // 他のエンドポイントをここに追加
	//     }
	// }

	return r
}
