package argo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Node interface {
	Handle(raw []byte, out io.Writer)
}

const (
	Init   MessageType = "init"
	InitOk MessageType = "init_ok"
)

type InitBody struct {
	BaseBody
	NodeID  string   `json:"node_id"`
	NodeIDs []string `json:"node_ids"`
}

type InitOkBody struct {
	BaseBody
}

func MainLoop[N Node](
	factory func(Message[InitBody]) N,
	in io.Reader,
	out io.Writer,
) {
	if err := Loop(factory, in, out); err != nil {
		log.Fatal(err)
	}
}

func Loop[N Node](
	factory func(Message[InitBody]) N,
	in io.Reader,
	out io.Writer,
) error {
	scanner := bufio.NewScanner(in)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("reading init: %w", err)
		}
		return fmt.Errorf("no init message")
	}

	var init Message[InitBody]
	if err := json.Unmarshal(scanner.Bytes(), &init); err != nil {
		return fmt.Errorf("unmarshal init: %w", err)
	}
	json.NewEncoder(out).Encode(Message[InitOkBody]{
		Src: init.Dst,
		Dst: init.Src,
		Body: InitOkBody{
			BaseBody: BaseBody{
				Type:      InitOk,
				InReplyTo: init.Body.ID,
			},
		},
	})

	node := factory(init)
	for scanner.Scan() {
		node.Handle(scanner.Bytes(), out)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading message: %w", err)
	}
	return nil
}
