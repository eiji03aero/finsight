package usecase

import (
	"backend/internal/domain/service"
)

// HelloWorldUsecaseResult は HelloWorldUsecase の実行結果
type HelloWorldUsecaseResult struct {
	Message      string
	FunnyMessage string
}

// HelloWorldUsecase は /hello_world のユースケース
type HelloWorldUsecase interface {
	Execute() *HelloWorldUsecaseResult
}

type helloWorldUsecase struct {
	messageService service.MessageService
}

// NewHelloWorldUsecase は HelloWorldUsecase の実装を返す
func NewHelloWorldUsecase(messageService service.MessageService) HelloWorldUsecase {
	return &helloWorldUsecase{
		messageService: messageService,
	}
}

// Execute はユースケースを実行し、結果を返す
func (u *helloWorldUsecase) Execute() *HelloWorldUsecaseResult {
	// 1. ドメインサービスからランダムメッセージを取得
	funnyMessage := u.messageService.GetRandomFunnyMessage()

	// 2. 結果を構築して返す
	return &HelloWorldUsecaseResult{
		Message:      "hello world",
		FunnyMessage: funnyMessage.GetContent(),
	}
}
