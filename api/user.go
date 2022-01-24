package api

type From struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Title     string `json:"title,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Type      string `json:"type"`
}

type UserMessage struct {
	MessageID int64  `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type SenderChat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Entities struct {
	Offset int    `json:"offset"`
	Lenfth int    `json:"length"`
	Type   string `json:"type"`
}

type ChannelPost struct {
	MessageID  int64      `json:"message_id"`
	SenderChat SenderChat `json:"sender_chat"`
	Chat       Chat       `json:"chat"`
	Date       int64      `json:"date"`
	Text       string     `json:"text"`
	Entities   []Entities `json:"entities"`
}

type Update struct {
	UpdateID    int64        `json:"update_id"`
	Message     *UserMessage `json:"message,omitempty"`
	ChannelPost *ChannelPost `json:"channel_post,omitempty"`
}
