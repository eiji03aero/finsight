package service

import (
	"backend/internal/domain/model"
	"math/rand"
	"time"
)

// MessageService はメッセージ生成のドメインサービス
type MessageService interface {
	GetRandomFunnyMessage() *model.Message
}

type messageService struct {
	funnyMessages []string
	rng           *rand.Rand
}

// NewMessageService は MessageService の実装を返す
func NewMessageService() MessageService {
	return &messageService{
		funnyMessages: []string{
			"Why did the Go programmer quit? Because they didn't get arrays!",
			"Go: where 'nil' is not nothing, but something that is nothing.",
			"I would tell you a UDP joke, but you might not get it.",
			"There are 10 types of people: those who understand binary and those who don't.",
			"Programming is 10% writing code and 90% figuring out why it doesn't work.",
			"Bug? That's not a bug, it's an undocumented feature!",
			"Why do programmers prefer dark mode? Because light attracts bugs!",
			"It works on my machine! ¯\\_(ツ)_/¯",
			"Roses are red, violets are blue, unexpected '{' on line 32.",
			"I'm not lazy, I'm just in energy-saving mode.",
		},
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetRandomFunnyMessage はランダムな面白いメッセージを返す
func (s *messageService) GetRandomFunnyMessage() *model.Message {
	randomIndex := s.rng.Intn(len(s.funnyMessages))
	return model.NewMessage(s.funnyMessages[randomIndex])
}
