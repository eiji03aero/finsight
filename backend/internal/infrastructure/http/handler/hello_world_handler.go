package handler

import (
	"net/http"

	"backend/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

// GetHelloWorldResponse は GET /hello_world エンドポイントのレスポンス
type GetHelloWorldResponse struct {
	Message      string `json:"message"`
	FunnyMessage string `json:"funny_message"`
}

// HelloWorldHandler は /hello_world エンドポイントのハンドラー
type HelloWorldHandler struct {
	usecase usecase.HelloWorldUsecase
}

// NewHelloWorldHandler は HelloWorldHandler を生成する
func NewHelloWorldHandler(usecase usecase.HelloWorldUsecase) *HelloWorldHandler {
	return &HelloWorldHandler{
		usecase: usecase,
	}
}

// Handle は GET /hello_world リクエストを処理する
func (h *HelloWorldHandler) Handle(c *gin.Context) {
	// 1. ユースケースを実行
	result := h.usecase.Execute()

	// 2. レスポンスを構築
	response := &GetHelloWorldResponse{
		Message:      result.Message,
		FunnyMessage: result.FunnyMessage,
	}

	// 3. JSONレスポンスを返す
	c.JSON(http.StatusOK, response)
}
