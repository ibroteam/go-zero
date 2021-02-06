package message

import "encoding/json"

type LinkMessageContent struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl,omitempty"`
}

type LinkMessage struct {
	MessageBase
	Link LinkMessageContent `json:"link"`
}

func NewLinkMessage(title, text, url string) *LinkMessage {
	return &LinkMessage{
		MessageBase: MessageBase{
			MsgType: MsgTypeLink,
			At: At{
				AtMobiles: []string{},
				IsAtAll:   false,
			},
		},
		Link: LinkMessageContent{
			Title:      title,
			Text:       text,
			MessageUrl: url,
		},
	}
}

func (msg *LinkMessage) PicUrl(picUrl string) *LinkMessage {
	msg.Link.PicUrl = picUrl
	return msg
}

func (msg *LinkMessage) AtPeople(phones ...string) *LinkMessage {
	msg.At.AtMobiles = phones
	return msg
}

func (msg *LinkMessage) AtAll() *LinkMessage {
	msg.At.IsAtAll = true
	return msg
}

func (m *LinkMessage) ToJson() json.RawMessage {
	data, _ := json.Marshal(&m)
	return data
}
