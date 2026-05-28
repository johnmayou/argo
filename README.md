# argo

Go solutions to the [Gossip Glomers](https://fly.io/dist-sys/) distributed systems challenges, built on [Maelstrom](https://github.com/jepsen-io/maelstrom).

## Challenges

| # | Challenge |
|---|-----------|
| [1](https://fly.io/dist-sys/1/) | Echo |
| [2](https://fly.io/dist-sys/2/) | Unique ID Generation |
| [3a](https://fly.io/dist-sys/3a/) | Broadcast — Single Node |
| [3b](https://fly.io/dist-sys/3b/) | Broadcast — Multi Node |
| [3c](https://fly.io/dist-sys/3c/) | Broadcast — Fault Tolerant |
| [3d](https://fly.io/dist-sys/3d/) | Broadcast — Efficient I |
| [3e](https://fly.io/dist-sys/3e/) | Broadcast — Efficient II |
| [4](https://fly.io/dist-sys/4/) | Grow-Only Counter |
| [5a](https://fly.io/dist-sys/5a/) | Kafka-Style Log |
| [5b](https://fly.io/dist-sys/5b/) | Kafka-Style Log — Efficient |
| [5c](https://fly.io/dist-sys/5c/) | Kafka-Style Log — Efficient II |
| [6a](https://fly.io/dist-sys/6a/) | Totally-Available Transactions |
| [6b](https://fly.io/dist-sys/6b/) | Totally-Available Transactions — Read Uncommitted |
| [6c](https://fly.io/dist-sys/6c/) | Totally-Available Transactions — Read Committed |

## Binaries

Each challenge lives under `cmd/<name>/` and compiles to `bin/<name>`.

| Binary | Challenge |
|--------|-----------|
| `echo` | #1 |
| `unique-ids` | #2 |
| `broadcast` | #3a–3e |
| `grow-only-counter` | #4 |
| `kafka-log` | #5a–5c |
| `totally-available` | #6a–6c |

## Maelstrom

The Maelstrom binary is bundled at `./maelstrom/maelstrom`. See the [protocol docs](https://github.com/jepsen-io/maelstrom/blob/main/doc/protocol.md) for message format details.

## Make Targets

```
make serve          # start Maelstrom web UI at localhost:8080
make echo           # build + test echo (1 node, 10s)
make unique-ids     # build + test unique-ids (3 nodes, partitions)
make broadcast-1    # build + test broadcast (1 node)
make broadcast-2    # build + test broadcast (5 nodes)
```
