package xeno

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
)

type Hand struct {
	cards []int
}

func NewHand(first, second int) Hand {
	return Hand{cards: []int{first, second}}
}

func (h *Hand) Add(c int) {
	h.cards = append(h.cards, c)
}

func (h *Hand) Set(c int) {
	h.cards = []int{c}
}

func (h Hand) Get() int {
	if len(h.cards) > 1 {
		log.Fatal("Hand.Get() len(h.cards) > 1")
	}
	return h.cards[0]
}

func (h *Hand) Remove() int {
	if len(h.cards) > 1 {
		log.Fatal("Hand.Remove() len(h.cards) > 1")
	}
	c := h.cards[0]
	h.cards = h.cards[1:]
	return c
}

func (h *Hand) Clear() {
	h.cards = []int{}
}

func (h Hand) Count() int {
	return len(h.cards)
}

func (h Hand) Slice() []int {
	return h.cards
}

func (h Hand) At(i int) int {
	if i >= len(h.cards) {
		log.Fatal("Hand.At() i >= len(h.cards)")
	}
	return h.cards[i]
}

// TODO: 実際にsliceから除去するようにするか？
func (h Hand) Another(card int) int {
	if len(h.cards) != 2 {
		log.Fatal("Hand.Another() len(h.cards) != 2")
	}
	var remain int
	if h.cards[0] == card {
		remain = h.cards[1]
	} else if h.cards[1] == card {
		remain = h.cards[0]
	} else {
		// Error
		log.Panicf("no matching with pair: %v, %d", h, card)
	}
	return remain
}

func (h Hand) Larger() int {
	if len(h.cards) != 2 {
		log.Fatal("Hand.Larger() len(h.cards) != 2")
	}

	if h.cards[0] < h.cards[1] {
		return h.cards[1]
	} else {
		return h.cards[0]
	}
}

func (h Hand) Random() int {
	if len(h.cards) != 2 {
		log.Fatal("Hand.Random() len(h.cards) != 2")
	}

	if 1 == rand.Intn(2) {
		return h.cards[1]
	} else {
		return h.cards[0]
	}
}

func (h Hand) Has(n int) bool {
	for _, c := range h.cards {
		if c == n {
			return true
		}
	}
	return false
}

func (h Hand) String() string {
	str := "手札: "
	for _, c := range h.cards {
		str += fmt.Sprintf("[%d]", c)
	}
	return str
}

type Playable interface {
	Take(card int)                                                    // 1枚引く [1->2枚持っている状態に遷移]
	TakeFromWise(g *Game, candidates []int) (remains []int)           // 自分の賢者イベント 3枚から1枚選ぶ処理 [1->2枚持っている状態に遷移]
	Discard(g *Game) CardEvent                                        // 通常の２枚の手持ちから選ぶ処理 [2->1枚持っている状態に遷移]
	SelectOnPublicExecution(target Playable, pair Hand) (discard int) // 相手への公開処刑処理 2枚から1枚選ぶ
	SelectOnPlague(target Playable, hand Hand) (discard int)          // 相手への疫病イベント処理 2枚から1枚選ぶ
	DiscardSpecified(discard int)                                     // 自分の選択、相手からの公開処刑、疫病で指定された方を捨てる [2->1枚持っている状態に遷移]
	ID() int                                                          // プレイヤー情報 ID
	Name() string                                                     // プレイヤー情報 名前
	Has(expect int) bool                                              // 相手からの捜査で手持ちと比較
	Give() int                                                        // 交換時に自分のカードを渡す
	Reincarnate(newCard int)                                          // 復活札で復活
	Dropout()                                                         // 脱落
	Dropped() bool                                                    // 脱落したか
	SetCalledWise(bool)                                               // 賢者を出したかどうかをセット (要らない?)
	CalledWise() bool                                                 // 前回賢者を出したかどうか
	SetProtected(bool)                                                // 守護を出したかどうかをセット(要らない?)
	Protected() bool                                                  // 前回守護を出したかどうか
	Hand() Hand
}

type PlayerConfig struct {
	Name   string
	Manual bool
}

type Player struct {
	id         int
	name       string
	hand       Hand // 手札
	discarded  []int
	protected  bool
	calledWise bool
	dropped    bool
	manual     bool
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
		hand:   Hand{cards: []int{}},
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

func (p *Player) Take(next int) {
	p.hand.Add(next)
}

func (p *Player) Give() int {
	return p.hand.Remove()
}

func (p *Player) Hand() Hand {
	return p.hand
}

func (p *Player) Discard(g *Game) CardEvent {
	if p.hand.Count() < 2 {
		return CardEvent{}
	}

	var discard int
	// TODO* Pair
	if p.hand.Has(10) {
		// 10は選べない
		discard = p.hand.Another(10)
	} else {
		discard = p.hand.Random()
	}
	p.DiscardSpecified(discard)

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

func (p *Player) TakeFromWise(g *Game, candidates []int) (remains []int) {
	// TODO: select logic
	selectedIdx := rand.Intn(len(candidates))
	for i, c := range candidates {
		if i != selectedIdx {
			remains = append(remains, c)
		}
	}
	selected := candidates[selectedIdx]
	p.Take(selected)
	return
}

// Targetの捨てカードを選ぶ
func (p *Player) SelectOnPublicExecution(target Playable, hand Hand) (discard int) {
	// 可視
	return hand.Larger()
}

func (p *Player) SelectOnPlague(target Playable, hand Hand) (discard int) {
	// 不可視
	return hand.Random()
}

// 二枚持っているカードのうち指定されたカードを捨てる
// TODO: pairメンバがイマイチなのでリファクタ
func (p *Player) DiscardSpecified(discard int) {
	remain := p.hand.Another(discard)
	p.discarded = append(p.discarded, discard)
	p.hand.Set(remain)
	return
}

// 脱落
func (p *Player) Dropout() {
	for _, c := range p.hand.Slice() {
		p.discarded = append(p.discarded, c)
	}
	p.hand.Clear()
	p.dropped = true
}

// 転生
func (p *Player) Reincarnate(newCard int) {
	for _, c := range p.hand.Slice() {
		p.discarded = append(p.discarded, c)
	}
	p.hand.Set(newCard)
}

func (p Player) String() string {
	alive := ""
	if p.dropped {
		alive = "(脱落)"
	}
	return fmt.Sprintf("%s %s: %s 捨てたカード:%v", p.name, alive, p.hand, p.discarded)
}

func (p Player) Has(expect int) bool {
	if p.hand.Count() != 1 {
		log.Fatal("p.hand.Count() != 1")
	}
	return p.hand.Has(expect)
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

func (p *ManualPlayer) Discard(g *Game) CardEvent {
	fmt.Println(p.hand)

	if p.hand.Count() == 1 {
		return CardEvent{}
	}

	var discard int

	if p.Player.hand.Has(10) {
		discard = userInput([]int{p.Player.hand.Another(10)})
	} else {
		discard = userInput(p.Player.hand.Slice())
	}
	p.DiscardSpecified(discard)

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

func (p *ManualPlayer) TakeFromWise(g *Game, candidates []int) (remains []int) {
	selected := userInput(candidates)
	for _, c := range candidates {
		if c != selected {
			remains = append(remains, c)
		}
	}
	p.Take(selected)
	return
}

func (p *ManualPlayer) SelectOnPublicExecution(target Playable, hand Hand) (discard int) {
	// 可視
	fmt.Printf("相手のカード: %s", hand)
	fmt.Println("捨てるカードは？")
	discard = userInput(hand.Slice())
	return discard
}

func (p *ManualPlayer) SelectOnPlague(target Playable, hand Hand) (discard int) {
	// 不可視
	fmt.Println("捨てるカードは？ 左:[0], 右[1]")
	discardIdx := userInput([]int{0, 1})
	return hand.At(discardIdx)
}
