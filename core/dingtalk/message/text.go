package message

import "encoding/json"

type TextMessageContent struct {
	Content string `json:"content"`
}

type TextMessage struct {
	MessageBase
	Text TextMessageContent `json:"text"`
}

func NewTextMessage(content string) *TextMessage {
	return &TextMessage{
		MessageBase: MessageBase{
			MsgType: MsgTypeText,
			At: At{
				AtMobiles: []string{},
				IsAtAll:   false,
			},
		},
		Text: TextMessageContent{
			Content: content,
		},
	}
}

func (msg *TextMessage) AtPeople(phones ...string) *TextMessage {
	msg.At.AtMobiles = phones
	return msg
}

func (msg *TextMessage) AtAll() *TextMessage {
	msg.At.IsAtAll = true
	return msg
}

func (m *TextMessage) ToJson() json.RawMessage {
	data, _ := json.Marshal(&m)
	return data
}
