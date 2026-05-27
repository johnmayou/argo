package main

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/johnmayou/argo/internal/argo"
)

const (
	Generate   argo.MessageType = "generate"
	GenerateOk argo.MessageType = "generate_ok"
)

type GenerateBody struct {
	argo.BaseBody
}

type GenerateOkBody struct {
	argo.BaseBody
	ID string `json:"id"`
}

type UniqueIdsNode struct {
	NodeID string
	Mux    *argo.Mux

	counter int
}

func NewUniqueIdsNode(init argo.Message[argo.InitBody]) *UniqueIdsNode {
	n := &UniqueIdsNode{
		NodeID: init.Body.NodeID,
		Mux:    argo.NewMux(),
	}

	argo.MuxRegister(n.Mux, Generate, n.handleGenerate)
	argo.MuxRegister(n.Mux, GenerateOk, argo.NoOpHandler[GenerateOkBody])

	return n
}

func (n *UniqueIdsNode) Handle(raw []byte, out io.Writer) {
	n.Mux.Handle(raw, out)
}

func (n *UniqueIdsNode) handleGenerate(msg argo.Message[GenerateBody], out io.Writer) {
	json.NewEncoder(out).Encode(argo.Message[GenerateOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: GenerateOkBody{
			BaseBody: argo.BaseBody{
				Type:      GenerateOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
			ID: n.NodeID + strconv.Itoa(n.counter),
		},
	})
	n.counter += 1
}

func main() {
	argo.MainLoop(NewUniqueIdsNode, os.Stdin, os.Stdout)
}
