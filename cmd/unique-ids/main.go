package main

import (
	"context"
	"encoding/json"
	"io"
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
	Ctx    context.Context
	NodeID string
	Mux    *argo.Mux
	Out    io.Writer

	counter int
}

func NewUniqueIdsNode(ctx context.Context, init argo.Message[argo.InitBody], out io.Writer) *UniqueIdsNode {
	n := &UniqueIdsNode{
		Ctx:    ctx,
		NodeID: init.Body.NodeID,
		Mux:    argo.NewMux(),
		Out:    out,
	}

	argo.MuxRegister(n.Mux, Generate, n.handleGenerate)
	argo.MuxRegister(n.Mux, GenerateOk, argo.NoOpHandler[GenerateOkBody])

	return n
}

func (n *UniqueIdsNode) Handle(raw []byte) {
	n.Mux.HandleRaw(raw, n.Out)
}

func (n *UniqueIdsNode) handleGenerate(msg argo.Message[GenerateBody]) error {
	json.NewEncoder(n.Out).Encode(argo.Message[GenerateOkBody]{
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

	return nil
}

func main() {
	argo.MainLoop(NewUniqueIdsNode)
}
