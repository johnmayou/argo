package argo

import (
	"encoding/json"
	"fmt"
	"io"
)

type Mux struct {
	handlers map[MessageType]func([]byte, io.Writer) error
}

func NewMux() *Mux {
	return &Mux{handlers: make(map[MessageType]func([]byte, io.Writer) error)}
}

func MuxRegister[B Body](m *Mux, t MessageType, h func(Message[B], io.Writer)) {
	m.handlers[t] = func(raw []byte, out io.Writer) error {
		var msg Message[B]
		if err := json.Unmarshal(raw, &msg); err != nil {
			return fmt.Errorf("unmarshal message: %w", err)
		}
		h(msg, out)
		return nil
	}
}

func (m *Mux) Handle(raw []byte, out io.Writer) error {
	var base Message[BaseBody]
	if err := json.Unmarshal(raw, &base); err != nil {
		return fmt.Errorf("unmarshal base: %w", err)
	}
	h, ok := m.handlers[base.Body.Type]
	if !ok {
		return fmt.Errorf("unknown message type: %s", base.Body.Type)
	}
	return h(raw, out)
}

func NoOpHandler[B Body](_ Message[B], _ io.Writer) {}
