package telegram

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`
}

func (u Update) FirstText() (chatID int64, text string, ok bool) {
	if u.Message == nil {
		return 0, "", false
	}
	if u.Message.Text == "" {
		return 0, "", false
	}
	return u.Message.Chat.ID, u.Message.Text, true
}
