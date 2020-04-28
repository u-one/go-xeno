package xeno

import (
	"fmt"
	"math/rand"
)

type Playable interface {
	SelectOne(g *Game, next int) (result SelectResult)
	SelectFromWise(g *Game, candidates []int) (selected int, remains []int)
	Discard(card int) (discarded int)
	Dropout()
	Reincarnate(newCard int)
	Has(expect int) bool
}

type Player struct {
	id         int
	name       string
	current    int
	discarded  []int
	protected  bool
	calledWise bool
	dropped    bool
}

type SelectResult struct {
	Discarded int
	Target    *Player
	Expect    int
}

func (p *Player) SelectOne(g *Game, next int) (result SelectResult) {
	debugPrintf("[%d][%d]\n", p.current, next)

	result = SelectResult{}

	if p.current == 0 {
		p.current = next
		return
	}

	// discard one
	discardNext := false

	// 10は選べない
	if p.current == 10 { // 持っていたのが10
		discardNext = true // 引いた方を捨てる
	} else if next == 10 { // 引いたのが10
		discardNext = false // 持ってた方を捨てる
	} else {
		discardNext = (1 == rand.Intn(2)) // 確率で選ぶ
	}

	var discard int
	if discardNext {
		// 引いた方を捨てる
		discard = next
	} else {
		// 持ってた方を捨てる
		discard = p.current
		p.current = next
	}
	p.discarded = append(p.discarded, discard)

	switch discard {
	case 2:
		others := g.OtherPlayers(p)
		// TODO: select target
		target := others[rand.Intn(len(others))]
		result.Target = target
		// TODO: 出ていないカードを挙げて1枚選ぶ
		result.Expect = rand.Intn(11)
	case 1, 5, 6, 8, 10:
		others := g.OtherPlayers(p)
		debugPrintf("OtherPlayers:%v", others)
		//TODO: select target
		target := others[rand.Intn(len(others))]
		debugPrintf("target:%v", target)
		result.Target = target
	}

	result.Discarded = discard
	return
}

func (p *Player) SelectFromWise(candidates []int) (selected int, remains []int) {
	fmt.Println("賢者からの選択: ")
	for _, c := range candidates {
		debugPrintf("[%d]", c)
	}
	debugPrintf("\n")

	// TODO: select logic
	selectedIdx := rand.Intn(len(candidates))
	for i, c := range candidates {
		if i != selectedIdx {
			remains = append(remains, c)
		}
	}
	selected = candidates[selectedIdx]
	return
}

func (p *Player) Discard(card int) (discarded int) {
	discarded = card
	p.discarded = append(p.discarded, discarded)
	return
}

func (p *Player) Dropout() {
	p.current = 0
	p.dropped = true
	fmt.Printf("%s脱落\n", p.name)
}

func (p *Player) Reincarnate(newCard int) {
	p.current = newCard
}

func (p Player) String() string {
	return fmt.Sprintf("%s: Current:%d %v", p.name, p.current, p.discarded)
}

func (p Player) Has(expect int) bool {
	return (p.current == expect)
}
