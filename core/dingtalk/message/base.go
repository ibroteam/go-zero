package message

import "encoding/json"

const (
	MsgTypeText       = "text"
	MsgTypeLink       = "link"
	MsgTypeMarkdown   = "markdown"
	MsgTypeActionCard = "actionCard"
	MsgTypeFeedCard   = "feedCard"
)

type Message interface {
	ToJson() json.RawMessage
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type MessageBase struct {
	MsgType string `json:"msgtype"`
	At      At     `json:"at"l`
}
