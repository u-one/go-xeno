package xeno

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type CommStrategy struct {
	opponentInfo map[PlayerID]int
}

func (s CommStrategy) SelectDiscard(g *Game, p *Player) CardEvent {
	var discard int
	if p.hand.Has(10) {
		// 10は選べない
		discard = p.hand.Another(10)
	} else {
		discard = p.hand.Random()
	}

	event := CardEvent{Card: discard}
	switch event.Card {
	case 2:
		event.Target, event.Expect = s.estimateOpponentHand(g, p)
	case 1, 3, 5, 9:
		event.Target = s.randomSelectTarget(g, p)
	case 6:
		// TODO: 相手の持っているカードを考慮
		event.Target = s.randomSelectTarget(g, p)
	case 8:
		event.Target = s.randomSelectTarget(g, p)
		s.opponentInfo[event.Target.ID()] = p.hand.Another(discard)
	case 4, 7, 10:
	}
	return event
}

func (s CommStrategy) randomSelectTarget(g *Game, p *Player) (target *Player) {
	others := g.OtherPlayers(p)
	for {
		target = others[rand.Intn(len(others))]
		if !target.Dropped() {
			break
		}
	}
	if target == nil {
		log.Fatal("target not found")
	}
	return target
}

func (s CommStrategy) estimateOpponentHand(g *Game, p *Player) (target *Player, card int) {
	// Decide from opponent info
	if len(s.opponentInfo) > 0 {
		targetIdx := rand.Intn(len(s.opponentInfo))
		idx := 0
		for id, c := range s.opponentInfo {
			if idx == targetIdx {
				target = g.Player(id)
				if target.Dropped() {
					continue
				}
				card = c
				return
			}
		}
	}

	// Then, estimate
	appeared := append([]int{}, p.Hand().Slice()...)
	for _, p := range g.Players {
		d := p.Discarded()
		for _, c := range d {
			appeared = append(appeared, c)
		}
	}

	// put all cards and num of each cards
	hiddens := make(map[int]int, len(AllCards))
	for _, c := range AllCards {
		hiddens[c]++
	}

	// remove already appeared on game
	for _, c := range appeared {
		hiddens[c]--
	}

	// find largest count of each cards
	maxCount := 0
	for _, n := range hiddens {
		if n > maxCount {
			maxCount = n
		}
	}

	debugPrintf("estimateOpponentHand - hidden cards: %v\n", hiddens)

	// select candidates from cards which remains largest count
	candidates := []int{}
	for c, n := range hiddens {
		if n == maxCount {
			candidates = append(candidates, c)
		}
	}

	// finally select randomly
	card = candidates[rand.Intn(len(candidates))]
	target = s.randomSelectTarget(g, p)
	return
}

func (s CommStrategy) SelectFromWise(g *Game, candidates []int) int {
	// TODO: select logic
	return candidates[rand.Intn(len(candidates))]
}

func (s CommStrategy) SelectOnPublicExecution(player, target *Player, hand Hand) int {
	return hand.Larger()
}

func (s CommStrategy) SelectOnPlague(player, target *Player, hand Hand) int {
	return hand.Random()
}

func (s CommStrategy) KnowByClairvoyance(g *Game, player, target *Player, c int) {
	s.opponentInfo[target.ID()] = c
}

func (s CommStrategy) OnOpponentEvent(g *Game, player, opponent *Player, e CardEvent) {
	if c, ok := s.opponentInfo[opponent.ID()]; ok {
		if c == e.Card {
			delete(s.opponentInfo, opponent.ID())
		}
	}
}

func userInput(candidates []int) (num int) {
	for {
		fmt.Printf("Select %v to discard\n", candidates)
		var input string
		fmt.Scan(&input)
		//fmt.Printf("%s\n", input)
		i, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("invalid input: %s\n", input)
			continue
		}
		valid := false
		if len(candidates) == 0 {
			// just waited for hit Enter
			return
		}
		for _, c := range candidates {
			if c == i {
				valid = true
			}
		}
		if !valid {
			fmt.Printf("invalid input: %s\n", input)
			continue
		}
		num = i
		break
	}
	return
}

type ManualStrategy struct{}

func (s ManualStrategy) SelectDiscard(g *Game, p *Player) CardEvent {
	fmt.Println(p.hand)

	var discard int
	if p.hand.Has(10) {
		discard = userInput([]int{p.hand.Another(10)})
	} else {
		discard = userInput(p.hand.Slice())
	}

	others := g.OtherPlayers(p)
	var target *Player
	if len(others) == 1 {
		target = others[0]
	} else {
		var indices []int
		for i, o := range others {
			fmt.Printf("%s: [%d]\n", o.Name(), i)
			indices = append(indices, i)
		}
		fmt.Println("相手は？", indices)
		ti := userInput(indices)
		target = others[ti]
	}

	event := CardEvent{Card: discard}
	switch event.Card {
	case 2:
		event.Target = target
		fmt.Println("捜査: 予想は？[1-10]")
		expect := userInput([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
		event.Expect = expect
	case 1, 3, 5, 6, 8, 9:
		event.Target = target
	case 4, 7, 10:
	}
	return event
}

func (s ManualStrategy) SelectFromWise(g *Game, candidates []int) int {
	selected := userInput(candidates)
	return selected
}

func (s ManualStrategy) SelectOnPublicExecution(player, target *Player, hand Hand) (discard int) {
	// 可視
	fmt.Printf("相手のカード: %s", hand)
	fmt.Println("捨てるカードは？")
	discard = userInput(hand.Slice())
	return discard
}

func (s ManualStrategy) SelectOnPlague(player, target *Player, hand Hand) (discard int) {
	// 不可視
	fmt.Println("捨てるカードは？ 左:[0], 右[1]")
	discardIdx := userInput([]int{0, 1})
	return hand.At(discardIdx)
}

func (s ManualStrategy) KnowByClairvoyance(g *Game, player, target *Player, c int) {
	fmt.Printf("%sの手札: [%d]\n", target.Name(), c)
	fmt.Println("put any char")
	input := make([]byte, 1)
	os.Stdin.Read(input)
}

func (s ManualStrategy) OnOpponentEvent(g *Game, player, opponent *Player, e CardEvent) {
	fmt.Println("put any char")
	input := make([]byte, 1)
	os.Stdin.Read(input)
}
