package argo

type Body interface {
	isBody()
}

type Message[B Body] struct {
	Src  string `json:"src"`
	Dst  string `json:"dest"`
	Body B      `json:"body"`
}

type MessageType string

type BaseBody struct {
	Type      MessageType `json:"type"`
	ID        *int        `json:"msg_id"`
	InReplyTo *int        `json:"in_reply_to"`
}

func (b BaseBody) isBody() {}
