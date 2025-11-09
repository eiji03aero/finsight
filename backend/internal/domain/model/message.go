package model

// Message は表示するメッセージを表すドメインモデル
type Message struct {
	Content string
}

// NewMessage はMessageを生成する
func NewMessage(content string) *Message {
	return &Message{
		Content: content,
	}
}

// GetContent はメッセージの内容を取得する
func (m *Message) GetContent() string {
	return m.Content
}
