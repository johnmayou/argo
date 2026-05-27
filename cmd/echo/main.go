package main

import (
	"encoding/json"
	"io"
	"os"

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
	NodeID string
	Mux    *argo.Mux
}

func NewEchoNode(init argo.Message[argo.InitBody]) *EchoNode {
	n := &EchoNode{
		NodeID: init.Body.NodeID,
		Mux:    argo.NewMux(),
	}

	argo.MuxRegister(n.Mux, Echo, n.handleEcho)
	argo.MuxRegister(n.Mux, EchoOk, argo.NoOpHandler[EchoOkBody])

	return n
}

func (n *EchoNode) Handle(raw []byte, out io.Writer) {
	n.Mux.Handle(raw, out)
}

func (n *EchoNode) handleEcho(msg argo.Message[EchoBody], out io.Writer) {
	json.NewEncoder(out).Encode(argo.Message[EchoOkBody]{
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
}

func main() {
	argo.MainLoop(NewEchoNode, os.Stdin, os.Stdout)
}
