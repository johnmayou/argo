package argo

import (
	"encoding/json"
	"fmt"
	"io"
)

type Mux struct {
	handlers map[MessageType]func([]byte) error
}

func NewMux() *Mux {
	return &Mux{handlers: make(map[MessageType]func([]byte) error)}
}

func MuxRegister[B Body](m *Mux, t MessageType, h func(Message[B]) error) {
	m.handlers[t] = func(raw []byte) error {
		var msg Message[B]
		if err := json.Unmarshal(raw, &msg); err != nil {
			return fmt.Errorf("unmarshal message: %w", err)
		}
		return h(msg)
	}
}

func MuxHandle[B Body](m *Mux, msg Message[B], out io.Writer) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}
	return m.HandleRaw(raw, out)
}

func (m *Mux) HandleRaw(raw []byte, out io.Writer) error {
	var env struct {
		Body struct {
			Type MessageType `json:"type"`
		} `json:"body"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return fmt.Errorf("unmarshal message type: %w", err)
	}
	h, ok := m.handlers[env.Body.Type]
	if !ok {
		return fmt.Errorf("unknown message type: %s", env.Body.Type)
	}
	return h(raw)
}

func NoOpHandler[B Body](_ Message[B]) error {
	return nil
}
