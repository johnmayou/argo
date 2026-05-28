package main

import (
	"context"
	"encoding/json"
	"io"

	"github.com/johnmayou/argo/internal/argo"
)

const (
	Echo   argo.MessageType = "echo"
	EchoOk argo.MessageType = "echo_ok"
)

type EchoBody struct {
	argo.BaseBody
	Echo string `json:"echo"`
}

type EchoOkBody struct {
	argo.BaseBody
	Echo string `json:"echo"`
}

type EchoNode struct {
	Ctx    context.Context
	NodeID string
	Mux    *argo.Mux
	Out    io.Writer
}

func NewEchoNode(ctx context.Context, init argo.Message[argo.InitBody], out io.Writer) *EchoNode {
	n := &EchoNode{
		Ctx:    ctx,
		NodeID: init.Body.NodeID,
		Mux:    argo.NewMux(),
		Out:    out,
	}

	argo.MuxRegister(n.Mux, Echo, n.handleEcho)
	argo.MuxRegister(n.Mux, EchoOk, argo.NoOpHandler[EchoOkBody])

	return n
}

func (n *EchoNode) Handle(raw []byte) {
	n.Mux.HandleRaw(raw, n.Out)
}

func (n *EchoNode) handleEcho(msg argo.Message[EchoBody]) error {
	json.NewEncoder(n.Out).Encode(argo.Message[EchoOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: EchoOkBody{
			BaseBody: argo.BaseBody{
				Type:      EchoOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
			Echo: msg.Body.Echo,
		},
	})

	return nil
}

func main() {
	argo.MainLoop(NewEchoNode)
}
