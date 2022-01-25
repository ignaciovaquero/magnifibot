package api

// TelegramWebhookSendMessage is a struct that
type TelegramWebhookSendMessage struct {
	Method                   string `json:"method"`
	ChatID                   int64  `json:"chat_id"`
	Text                     string `json:"text"`
	ParseMode                string `json:"parse_mode,omitempty"`
	DisableWebPagePreview    bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification      bool   `json:"disable_notification,omitempty"`
	ProtectContent           bool   `json:"protect_content,omitempty"`
	ReplyToMessageID         int64  `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool   `json:"allow_sending_without_reply,omitempty"`
}
