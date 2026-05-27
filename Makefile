MAELSTROM_BIN = ./maelstrom/maelstrom

.PHONY: serve
serve:
	$(MAELSTROM_BIN) serve

.PHONY: echo
echo:
	go build -o ./bin/echo ./cmd/echo
	$(MAELSTROM_BIN) test -w echo --bin ./bin/echo --node-count 1 --time-limit 10
