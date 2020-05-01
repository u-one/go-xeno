# go-xeno

![Go](https://github.com/u-one/go-xeno/workflows/Go/badge.svg)

Golang implementation for playing card game "XENO" with CLI.
(work in progress)

カードゲーム"XENO" をCLIで遊んでみるためのgolang実装。
開発中。

## Usage

main.go

```go
  // 4 com players
	conf := xeno.GameConfig{
		Players: []xeno.PlayerConfig{
			{Name: "Player1"},
			{Name: "Player2"},
			{Name: "Player3"},
			{Name: "Player4"},
		},
	}

```

```go
  // 1 com player, 1 Human
	conf := xeno.GameConfig{
		Players: []xeno.PlayerConfig{
			{Name: "Player1"},
			{Name: "Player2", Manual: true},
		},
	}

```
