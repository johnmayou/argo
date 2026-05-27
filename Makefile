MAELSTROM_BIN = ./maelstrom/maelstrom

.PHONY: serve
serve:
	$(MAELSTROM_BIN) serve

.PHONY: echo
echo:
	go build -o ./bin/echo ./cmd/echo
	$(MAELSTROM_BIN) test -w echo --bin ./bin/echo --node-count 1 --time-limit 10

unique-ids:
	go build -o ./bin/unique-ids ./cmd/unique-ids
	$(MAELSTROM_BIN) test -w unique-ids --bin ./bin/unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition