package xeno

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

type Pair struct {
	vars []int
}

func NewPair(first, second int) Pair {
	return Pair{vars: []int{first, second}}
}

func (p Pair) Slice() []int {
	return p.vars
}

func (p Pair) At(i int) int {
	if i >= len(p.vars) {
		log.Fatal("i >= len(p.vars)")
	}
	return p.vars[i]
}

func (p Pair) Another(card int) int {
	if len(p.vars) != 2 {
		log.Fatal("len(p.pair.vars) != 2")
	}
	var remain int
	if p.vars[0] == card {
		remain = p.vars[1]
	} else if p.vars[1] == card {
		remain = p.vars[0]
	} else {
		// Error
		log.Panicf("no matching with pair: %v, %d", p, card)
	}
	return remain
}

func (p Pair) Larger() int {
	if p.vars[0] < p.vars[1] {
		return p.vars[1]
	} else {
		return p.vars[0]
	}
}

func (p Pair) Random() int {
	if 1 == rand.Intn(2) {
		return p.vars[1]
	} else {
		return p.vars[0]
	}
}

func (p Pair) Has(n int) bool {
	for _, c := range p.vars {
		if c == n {
			return true
		}
	}
	return false
}

func (p Pair) String() string {
	return fmt.Sprintf("[%d][%d]", p.vars[0], p.vars[1])
}

type Playable interface {
	Take(card int)
	TakeAndSelect(g *Game, next int) CardEvent
	SelectFromWise(g *Game, candidates []int) (selected int, remains []int)
	SelectOnPublicExecution(target Playable, pair Pair) (discard int)
	SelectOnPlague(target Playable, pair Pair) (discard int)
	Discard() (discarded int)
	DiscardFromPair(discard int)
	Dropout()
	Reincarnate(newCard int)
	Has(expect int) bool
	ID() int
	Name() string
	Dropped() bool
	CalledWise() bool
	SetCalledWise(bool)
	SetProtected(bool)
	Protected() bool
	Current() int
	SetCurrent(int)
	Pair() Pair
}

type PlayerConfig struct {
	Name   string
	Manual bool
}

type Player struct {
	id         int
	name       string
	current    int
	discarded  []int
	protected  bool
	calledWise bool
	dropped    bool
	manual     bool
	pair       Pair
}

var (
	playerCount = 0
)

func NewPlayer(conf PlayerConfig) *Player {
	playerCount++
	id := playerCount
	name := conf.Name
	if len(name) == 0 {
		name = fmt.Sprintf("プレイヤー%d", id)
	}
	return &Player{
		id:     id,
		name:   name,
		manual: conf.Manual,
	}
}
func (p *Player) ID() int {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Dropped() bool {
	return p.dropped
}

func (p *Player) CalledWise() bool {
	return p.calledWise
}

func (p *Player) SetCalledWise(b bool) {
	p.calledWise = b
}

func (p *Player) SetProtected(b bool) {
	p.protected = b
}

func (p *Player) Protected() bool {
	return p.protected
}

func (p *Player) Current() int {
	return p.current
}

func (p *Player) SetCurrent(cur int) {
	p.current = cur
}

func (p *Player) Take(card int) {
	p.pair = NewPair(p.current, card)
}

// TODO: pairメンバがイマイチなのでリファクタ
func (p *Player) Pair() Pair {
	return p.pair
}

func (p *Player) TakeAndSelect(g *Game, next int) CardEvent {
	if p.current == 0 {
		p.current = next
		return CardEvent{}
	}

	var discard int
	// TODO: pairメンバがイマイチなのでリファクタ
	p.pair = NewPair(p.current, next)
	if p.pair.Has(10) {
		// 10は選べない
		discard = p.pair.Another(10)
	} else {
		discard = p.pair.Random()
	}
	p.DiscardFromPair(discard)

	selectTarget := func() Playable {
		var target Playable
		others := g.OtherPlayers(p)
		for {
			// TODO: select target
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

	event := CardEvent{Card: discard}
	switch event.Card {
	case 2:
		event.Target = selectTarget()
		// TODO: 出ていないカードを挙げて1枚選ぶ
		event.Expect = rand.Intn(11)
	case 1, 5, 6, 8, 9, 10:
		event.Target = selectTarget()
	}
	return event
}

func (p *Player) SelectFromWise(g *Game, candidates []int) (selected int, remains []int) {
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

// Targetの捨てカードを選ぶ
func (p *Player) SelectOnPublicExecution(target Playable, pair Pair) (discard int) {
	// 可視
	return pair.Larger()
}

func (p *Player) SelectOnPlague(target Playable, pair Pair) (discard int) {
	// 不可視
	return pair.Random()
}

// 二枚持っているカードのうち指定されたカードを捨てる
// TODO: pairメンバがイマイチなのでリファクタ
func (p *Player) DiscardFromPair(discard int) {
	remain := p.pair.Another(discard)
	p.discarded = append(p.discarded, discard)
	p.current = remain
	p.pair = Pair{}
	return
}

// 持っているカードを捨てる
func (p *Player) Discard() (discarded int) {
	p.discarded = append(p.discarded, p.current)
	p.current = 0
	return
}

// 脱落
func (p *Player) Dropout() {
	p.Discard()
	p.dropped = true
}

// 転生
func (p *Player) Reincarnate(newCard int) {
	p.current = newCard
}

func (p Player) String() string {
	alive := ""
	if p.dropped {
		alive = "(脱落)"
	}
	return fmt.Sprintf("%s %s: 手持ち:%d 捨てたカード:%v", p.name, alive, p.current, p.discarded)
}

func (p Player) Has(expect int) bool {
	return (p.current == expect)
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

type ManualPlayer struct {
	Player
}

func NewManualPlayer(conf PlayerConfig) *ManualPlayer {
	playerCount++
	id := playerCount
	name := conf.Name
	if len(name) == 0 {
		name = fmt.Sprintf("プレイヤー%d", id)
	}
	return &ManualPlayer{
		Player{
			id:     id,
			name:   name,
			manual: conf.Manual,
		},
	}
}

func (p *ManualPlayer) TakeAndSelect(g *Game, next int) CardEvent {
	fmt.Printf("[%d][%d]\n", p.current, next)

	if p.current == 0 {
		p.current = next
		return CardEvent{}
	}

	var discard int

	p.Player.pair = NewPair(p.current, next)
	fmt.Println(p.Player.pair)
	if p.Player.pair.Has(10) {
		discard = userInput([]int{p.Player.pair.Another(10)})
	} else {
		discard = userInput(p.Player.pair.Slice())
	}
	p.DiscardFromPair(discard)

	others := g.OtherPlayers(p)
	var target Playable
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
	case 1, 5, 6, 8, 9, 10:
		event.Target = target
	}
	return event
}

func (p *ManualPlayer) SelectFromWise(g *Game, candidates []int) (selected int, remains []int) {
	selected = userInput(candidates)
	for _, c := range candidates {
		if c != selected {
			remains = append(remains, c)
		}
	}
	return
}

func (p *ManualPlayer) SelectOnPublicExecution(target Playable, pair Pair) (discard int) {
	// 可視
	fmt.Printf("相手のカード: %s", pair)
	fmt.Println("捨てるカードは？")
	discard = userInput(pair.Slice())
	return discard
}

func (p *ManualPlayer) SelectOnPlague(target Playable, pair Pair) (discard int) {
	// 不可視
	fmt.Println("捨てるカードは？ 左:[0], 右[1]")
	discardIdx := userInput([]int{0, 1})
	return pair.At(discardIdx)
}
