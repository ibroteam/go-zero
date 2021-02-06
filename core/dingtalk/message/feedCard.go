package message

import "encoding/json"

type FeedCardMessageContentLinkItem struct {
	Title      string `json:"title"`
	MessageUrl string `json:"messageURL"`
	PicUrl     string `json:"picURL"`
}

func NewFeedCardMessageContentLinkItem(title, messageUrl, picUrl string) FeedCardMessageContentLinkItem {
	return FeedCardMessageContentLinkItem{
		Title:      title,
		PicUrl:     picUrl,
		MessageUrl: messageUrl,
	}
}

type FeedCardMessageContent struct {
	Links []FeedCardMessageContentLinkItem `json:"links"`
}

type FeedCardMessage struct {
	MessageBase
	FeedCard FeedCardMessageContent `json:"feedCard"`
}

func NewFeedCardMessage() *FeedCardMessage {
	return &FeedCardMessage{
		MessageBase: MessageBase{
			MsgType: MsgTypeLink,
			At: At{
				AtMobiles: []string{},
				IsAtAll:   false,
			},
		},
		FeedCard: FeedCardMessageContent{
			Links: make([]FeedCardMessageContentLinkItem, 0, 5),
		},
	}
}

func (msg *FeedCardMessage) AddLink(title, messageUrl, picUrl string) *FeedCardMessage {
	msg.FeedCard.Links = append(msg.FeedCard.Links, NewFeedCardMessageContentLinkItem(title, messageUrl, picUrl))
	return msg
}

func (msg *FeedCardMessage) AtPeople(phones ...string) *FeedCardMessage {
	msg.At.AtMobiles = phones
	return msg
}

func (msg *FeedCardMessage) AtAll() *FeedCardMessage {
	msg.At.IsAtAll = true
	return msg
}

func (m *FeedCardMessage) ToJson() json.RawMessage {
	data, _ := json.Marshal(&m)
	return data
}
