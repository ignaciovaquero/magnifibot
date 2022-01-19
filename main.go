package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type from struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type chat struct {
	ID        int64  `json:"id"`
	Title     string `json:"title,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Type      string `json:"type"`
}

type userMessage struct {
	MessageID int64  `json:"message_id"`
	From      from   `json:"from"`
	Chat      chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type senderChat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type entities struct {
	Offset int    `json:"offset"`
	Lenfth int    `json:"length"`
	Type   string `json:"type"`
}

type channelPost struct {
	MessageID  int64      `json:"message_id"`
	SenderChat senderChat `json:"sender_chat"`
	Chat       chat       `json:"chat"`
	Date       int64      `json:"date"`
	Text       string     `json:"text"`
	Entities   []entities `json:"entities"`
}

type update struct {
	UpdateID    int64        `json:"update_id"`
	Message     *userMessage `json:"message,omitempty"`
	ChannelPost *channelPost `json:"channel_post,omitempty"`
}

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		var u update

		if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
			resp.Write([]byte(fmt.Sprintf("error decoding the request: %s", err.Error())))
			return
		}

		defer req.Body.Close()

		if u.Message != nil {
			re := regexp.MustCompile("/(\\w*)@?\\w*")
			if !re.Match([]byte(u.Message.Text)) {
				return
			}
			submatch := re.FindStringSubmatch(u.Message.Text)
			text := ""
			if len(submatch) > 1 {
				text = submatch[1]
			}
			fmt.Printf("Chat ID: %d\nUser ID: %d\nDate: %d\nText: %s\n\n", u.Message.Chat.ID, u.Message.From.ID, u.Message.Date, text)
			return
		}

		if u.ChannelPost != nil {
			re := regexp.MustCompile("/(\\w*)@magnifibot_bot") // TODO: Parameterize bot's name
			if !re.Match([]byte(u.ChannelPost.Text)) {
				return
			}
			submatch := re.FindStringSubmatch(u.ChannelPost.Text)
			text := ""
			if len(submatch) > 1 {
				text = submatch[1]
			}
			if re.Match([]byte(u.ChannelPost.Text)) {
				fmt.Printf("Chat ID: %d\nUser ID: %d\nDate: %d\nText: %s\n\n", u.ChannelPost.Chat.ID, u.ChannelPost.SenderChat.ID, u.ChannelPost.Date, text)
			}
			fmt.Printf("")
		}
	})
	http.ListenAndServe(":80", nil)
}
