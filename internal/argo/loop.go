package argo

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Node interface {
	Handle(raw []byte)
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
	factory func(context.Context, Message[InitBody], io.Writer) N,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := Loop(ctx, factory, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func Loop[N Node](
	ctx context.Context,
	factory func(context.Context, Message[InitBody], io.Writer) N,
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

	node := factory(ctx, init, out)
	for scanner.Scan() {
		node.Handle(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading message: %w", err)
	}
	return nil
}
