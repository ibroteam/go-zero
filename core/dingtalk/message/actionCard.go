package message

import "encoding/json"

type ActionButtonItem struct {
	Title     string `json:"title"`
	ActionUrl string `json:"actionURL"`
}

func NewActionButtonItem(title, url string) ActionButtonItem {
	return ActionButtonItem{
		Title:     title,
		ActionUrl: url,
	}
}

type ActionCardMessageContent struct {
	Title          string             `json:"title"`
	Text           string             `json:"text"`
	BtnOrientation string             `json:"btnOrientation"`
	SingleTitle    string             `json:"singleTitle"`
	SingleURL      string             `json:"singleURL"`
	Btns           []ActionButtonItem `json:"btns,omitempty"`
}

type ActionCardMessage struct {
	MessageBase
	ActionCard ActionCardMessageContent `json:"actionCard"`
}

func NewActionCardMessage(title, text, singleTitle, singleURL string) *ActionCardMessage {
	return &ActionCardMessage{
		MessageBase: MessageBase{
			MsgType: MsgTypeLink,
			At: At{
				AtMobiles: []string{},
				IsAtAll:   false,
			},
		},
		ActionCard: ActionCardMessageContent{
			Title:          title,
			Text:           text,
			BtnOrientation: "0",
			SingleTitle:    singleTitle,
			SingleURL:      singleURL,
		},
	}
}

func (msg *ActionCardMessage) AddButton(title, url string) *ActionCardMessage {
	if msg.ActionCard.Btns == nil {
		msg.ActionCard.Btns = make([]ActionButtonItem, 0, 10)
	}
	msg.ActionCard.Btns = append(msg.ActionCard.Btns, NewActionButtonItem(title, url))
	return msg
}

func (msg *ActionCardMessage) BtnOrientation(btnOrientation string) *ActionCardMessage {
	msg.ActionCard.BtnOrientation = btnOrientation
	return msg
}

func (msg *ActionCardMessage) AtPeople(phones ...string) *ActionCardMessage {
	msg.At.AtMobiles = phones
	return msg
}

func (msg *ActionCardMessage) AtAll() *ActionCardMessage {
	msg.At.IsAtAll = true
	return msg
}

func (m *ActionCardMessage) ToJson() json.RawMessage {
	data, _ := json.Marshal(&m)
	return data
}
