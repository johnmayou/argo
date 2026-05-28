package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/johnmayou/argo/internal/argo"
	"github.com/johnmayou/argo/internal/hashset"
)

const (
	Broadcast   argo.MessageType = "broadcast"
	BroadcastOk argo.MessageType = "broadcast_ok"
	Read        argo.MessageType = "read"
	ReadOk      argo.MessageType = "read_ok"
	Topology    argo.MessageType = "topology"
	TopologyOk  argo.MessageType = "topology_ok"
	Gossip      argo.MessageType = "gossip"
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
	Topology map[string][]string `json:"topology"`
}

type TopologyOkBody struct {
	argo.BaseBody
}

type GossipBody struct {
	argo.BaseBody
	Seen []int `json:"seen"`
}

type BroadcastNode struct {
	Ctx          context.Context
	NodeID       string
	Mux          *argo.Mux
	Neighborhood []string
	Messages     *hashset.HashSet[int]
	Known        map[string]*hashset.HashSet[int]
	Out          io.Writer
}

func NewBroadcastNode(ctx context.Context, init argo.Message[argo.InitBody], out io.Writer) *BroadcastNode {
	n := &BroadcastNode{
		Ctx:          ctx,
		NodeID:       init.Body.NodeID,
		Mux:          argo.NewMux(),
		Neighborhood: make([]string, 0),
		Messages:     hashset.NewHashSet[int](),
		Known:        make(map[string]*hashset.HashSet[int]),
		Out:          out,
	}

	argo.MuxRegister(n.Mux, Broadcast, n.handleBroadcast)
	argo.MuxRegister(n.Mux, BroadcastOk, argo.NoOpHandler[BroadcastOkBody])
	argo.MuxRegister(n.Mux, Read, n.handleRead)
	argo.MuxRegister(n.Mux, ReadOk, argo.NoOpHandler[ReadOkBody])
	argo.MuxRegister(n.Mux, Topology, n.handleTopology)
	argo.MuxRegister(n.Mux, TopologyOk, argo.NoOpHandler[TopologyOkBody])
	argo.MuxRegister(n.Mux, Gossip, n.handleGossip)

	n.InitNeighborGossip()

	return n
}

func (n *BroadcastNode) InitNeighborGossip() {
	ticker := time.NewTicker(300 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-n.Ctx.Done():
				return
			case <-ticker.C:
				for _, neiID := range n.Neighborhood {
					var seen []int
					known, ok := n.Known[neiID]
					if ok {
						for msg := range n.Messages.Values() {
							if !known.Contains(msg) {
								seen = append(seen, msg)
							}
						}
					} else {
						seen = slices.Collect(n.Messages.Values())
					}
					msg := argo.Message[GossipBody]{
						Src: n.NodeID,
						Dst: neiID,
						Body: GossipBody{
							BaseBody: argo.BaseBody{
								Type:      Gossip,
								ID:        nil,
								InReplyTo: nil,
							},
							Seen: seen,
						},
					}
					argo.MuxHandle(n.Mux, msg, n.Out)
				}
			}
		}
	}()
}

func (n *BroadcastNode) Handle(raw []byte) {
	n.Mux.HandleRaw(raw, n.Out)
}

func (n *BroadcastNode) handleBroadcast(msg argo.Message[BroadcastBody]) error {
	n.Messages.Add(msg.Body.Message)

	json.NewEncoder(n.Out).Encode(argo.Message[BroadcastOkBody]{
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

	return nil
}

func (n *BroadcastNode) handleRead(msg argo.Message[ReadBody]) error {
	json.NewEncoder(n.Out).Encode(argo.Message[ReadOkBody]{
		Src: msg.Dst,
		Dst: msg.Src,
		Body: ReadOkBody{
			BaseBody: argo.BaseBody{
				Type:      ReadOk,
				ID:        msg.Body.ID,
				InReplyTo: msg.Body.ID,
			},
			Messages: slices.Collect(n.Messages.Values()),
		},
	})

	return nil
}

func (n *BroadcastNode) handleTopology(msg argo.Message[TopologyBody]) error {
	topology, ok := msg.Body.Topology[n.NodeID]
	if !ok {
		return fmt.Errorf("no topology found for node")
	}
	n.Neighborhood = topology

	json.NewEncoder(n.Out).Encode(argo.Message[TopologyOkBody]{
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

	return nil
}

func (n *BroadcastNode) handleGossip(msg argo.Message[GossipBody]) error {
	if _, ok := n.Known[msg.Src]; !ok {
		n.Known[msg.Src] = hashset.NewHashSet[int]()
	}
	known := n.Known[msg.Src]

	for _, message := range msg.Body.Seen {
		n.Messages.Add(message)
		known.Add(message)
	}

	return nil
}

func main() {
	argo.MainLoop(NewBroadcastNode)
}
