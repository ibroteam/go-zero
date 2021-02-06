package message

import "encoding/json"

type MarkdownMessageContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MarkdownMessage struct {
	MessageBase
	Markdown MarkdownMessageContent `json:"markdown"`
}

func NewMarkdownMessage(title, content string) *MarkdownMessage {
	return &MarkdownMessage{
		MessageBase: MessageBase{
			MsgType: MsgTypeMarkdown,
		},
		Markdown: MarkdownMessageContent{
			Title: title,
			Text:  content,
		},
	}
}

func (msg *MarkdownMessage) AtPeople(phones ...string) *MarkdownMessage {
	msg.At.AtMobiles = phones
	return msg
}

func (msg *MarkdownMessage) AtAll() *MarkdownMessage {
	msg.At.IsAtAll = true
	return msg
}

func (m *MarkdownMessage) ToJson() json.RawMessage {
	data, _ := json.Marshal(&m)
	return data
}
