package main

import (
	"fmt"
	"math/rand"

	"github.com/u-one/go-xeno/xeno"
)

func init() {
	// var seed int64 = 3
	// var seed int64 = 1587981835 // 残り0まで
	// var seed int64 = 1587984873 // 英雄が皇帝に見つかる
	var seed int64 = 1587991022
	//var seed int64 = time.Now().Unix()
	fmt.Println("Seed:", seed)
	rand.Seed(seed)
}

func main() {
	game := xeno.NewGame()
	game.Loop()
}
