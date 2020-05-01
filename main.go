package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/u-one/go-xeno/xeno"
)

func init() {
	var seed int64 = time.Now().Unix()
	fmt.Println("Seed:", seed)
	rand.Seed(seed)
}

func main() {

	conf := xeno.GameConfig{
		Players: []xeno.PlayerConfig{
			{Name: "Player1"},
			{Name: "Player2"},
		},
	}

	game := xeno.NewGame(conf)
	game.Loop()
}
