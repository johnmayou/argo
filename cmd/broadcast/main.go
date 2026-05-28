package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/johnmayou/argo/internal/argo"
)

const (
	Broadcast   argo.MessageType = "broadcast"
	BroadcastOk argo.MessageType = "broadcast_ok"
	Read        argo.MessageType = "read"
	ReadOk      argo.MessageType = "read_ok"
	Topology    argo.MessageType = "topology"
	TopologyOk  argo.MessageType = "topology_ok"
)

type BroadcastBody struct {
	argo.BaseBody
	Message int `json:"message"`
}

type BroadcastOkBody struct {
	argo.BaseBody
}

type ReadBody struct {
	argo.BaseBody
}

type ReadOkBody struct {
	argo.BaseBody
	Messages []int `json:"messages"`
}

type TopologyBody struct {
	argo.BaseBody
}

type TopologyOkBody struct {
	argo.BaseBody
}

type BroadcastNode struct {
	NodeID   string
	Mux      *argo.Mux
	Messages []int
}

func NewBroadcastNode(init argo.Message[argo.InitBody]) *BroadcastNode {
	n := &BroadcastNode{
		NodeID:   init.Body.NodeID,
		Mux:      argo.NewMux(),
		Messages: make([]int, 0),
	}

	argo.MuxRegister(n.Mux, Broadcast, n.handleBroadcast)
	argo.MuxRegister(n.Mux, BroadcastOk, argo.NoOpHandler[BroadcastOkBody])
	argo.MuxRegister(n.Mux, Read, n.handleRead)
	argo.MuxRegister(n.Mux, ReadOk, argo.NoOpHandler[ReadOkBody])
	argo.MuxRegister(n.Mux, Topology, n.handleTopology)
	argo.MuxRegister(n.Mux, TopologyOk, argo.NoOpHandler[TopologyOkBody])

	return n
}

func (n *BroadcastNode) Handle(raw []byte, out io.Writer) {
	n.Mux.Handle(raw, out)
}

func (n *BroadcastNode) handleBroadcast(msg argo.Message[BroadcastBody], out io.Writer) {
	n.Messages = append(n.Messages, msg.Body.Message)

	json.NewEncoder(out).Encode(argo.Message[BroadcastOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: BroadcastOkBody{
			BaseBody: argo.BaseBody{
				Type:      BroadcastOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
		},
	})
}

func (n *BroadcastNode) handleRead(msg argo.Message[ReadBody], out io.Writer) {
	json.NewEncoder(out).Encode(argo.Message[ReadOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: ReadOkBody{
			BaseBody: argo.BaseBody{
				Type:      ReadOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
			Messages: n.Messages,
		},
	})
}

func (n *BroadcastNode) handleTopology(msg argo.Message[TopologyBody], out io.Writer) {
	json.NewEncoder(out).Encode(argo.Message[TopologyOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: TopologyOkBody{
			BaseBody: argo.BaseBody{
				Type:      TopologyOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
		},
	})
}

func main() {
	argo.MainLoop(NewBroadcastNode, os.Stdin, os.Stdout)
}
