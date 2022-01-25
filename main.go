package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/igvaquero18/magnifibot/api"
)

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		var u api.Update

		if err := json.NewDecoder(req.Body).Decode(&u); err != nil {
			resp.Write([]byte(fmt.Sprintf("error decoding the request: %s", err.Error())))
			return
		}

		defer req.Body.Close()

		if u.Message != nil {
			re := regexp.MustCompile(`/(\w*)@?\w*`)
			if !re.Match([]byte(u.Message.Text)) {
				return
			}
			submatch := re.FindStringSubmatch(u.Message.Text)
			text := ""
			if len(submatch) > 1 {
				text = submatch[1]
			}
			fmt.Printf("Chat ID: %d\nUser ID: %d\nDate: %d\nText: %s\nType: %s\n\n", u.Message.Chat.ID, u.Message.From.ID, u.Message.Date, text, u.Message.Chat.Type)
			return
		}

		if u.ChannelPost != nil {
			name := "magnifibot_bot"
			re := regexp.MustCompile(fmt.Sprintf(`/(\w*)@%s`, name))
			if !re.Match([]byte(u.ChannelPost.Text)) {
				return
			}
			submatch := re.FindStringSubmatch(u.ChannelPost.Text)
			text := ""
			if len(submatch) > 1 {
				text = submatch[1]
			}
			if re.Match([]byte(u.ChannelPost.Text)) {
				fmt.Printf("Chat ID: %d\nUser ID: %d\nDate: %d\nText: %s\nType: %s\n\n", u.ChannelPost.Chat.ID, u.ChannelPost.SenderChat.ID, u.ChannelPost.Date, text, u.ChannelPost.Chat.Type)
			}
			fmt.Printf("")
		}
	})
	http.ListenAndServe(":80", nil)
}
