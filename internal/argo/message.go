package argo

type Body interface {
	GetType() MessageType
	GetID() *int
	GetInReplyTo() *int
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

func (b BaseBody) GetType() MessageType { return b.Type }
func (b BaseBody) GetID() *int          { return b.ID }
func (b BaseBody) GetInReplyTo() *int   { return b.InReplyTo }
