package xeno

//go:generate mockgen -source=player.go -destination=./player_mock.go -package xeno

import (
	"fmt"
	"log"
	"math/rand"
)

type Hand struct {
	cards []int
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

type PlayerConfig struct {
	Name   string
	Manual bool
}

// PlayerStrategyによりコンピュータや人間などにより判断する部分をPlayerから移譲
type PlayerStrategy interface {
	SelectDiscard(g *Game, p *Player) CardEvent                         // 通常の２枚の手持ちから捨てるカードを選ぶ
	SelectFromWise(g *Game, candidates []int) int                       // 自分の賢者イベント 3枚から1枚選ぶ
	SelectOnPublicExecution(p, target *Player, pair Hand) (discard int) // 相手への公開処刑処理 2枚から1枚選ぶ
	SelectOnPlague(p, target *Player, hand Hand) (discard int)          // 相手への疫病イベント処理 2枚から1枚選ぶ
	KnowByClairvoyance(g *Game, player, target *Player, c int)
	OnOpponentEvent(g *Game, player, opponent *Player, e CardEvent)
}

type PlayerID int

type Player struct {
	id         PlayerID
	name       string
	hand       Hand // 手札
	discarded  []int
	protected  bool
	calledWise bool
	dropped    bool
	manual     bool
	strategy   PlayerStrategy // 戦略
}

var (
	playerCount = 0
)

func NewPlayer(conf PlayerConfig) *Player {
	playerCount++
	id := PlayerID(playerCount)
	name := conf.Name
	if len(name) == 0 {
		name = fmt.Sprintf("プレイヤー%d", id)
	}
	var s PlayerStrategy
	if conf.Manual {
		s = ManualStrategy{}
	} else {
		s = CommStrategy{opponentInfo: map[PlayerID]int{}}
	}
	return &Player{
		id:       id,
		name:     name,
		hand:     Hand{cards: []int{}},
		manual:   conf.Manual,
		strategy: s,
	}
}
func (p *Player) ID() PlayerID {
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

func (p *Player) Discarded() []int {
	return append([]int{}, p.discarded...)
}

func (p *Player) Discard(g *Game) CardEvent {
	if p.hand.Count() < 2 {
		return CardEvent{}
	}
	e := p.strategy.SelectDiscard(g, p)
	p.DiscardSpecified(e.Card)
	return e
}

func (p *Player) TakeFromWise(g *Game, candidates []int) (remains []int) {
	selected := p.strategy.SelectFromWise(g, candidates)
	found := false
	for _, c := range candidates {
		if c == selected && !found {
			// 見つかった番号1枚目は選択したカードとしてスキップ
			found = true
			continue
		}
		// 残った2枚をremainsに入れる
		remains = append(remains, c)
	}
	p.Take(selected)
	return
}

// Targetの捨てカードを選ぶ
func (p *Player) SelectOnPublicExecution(target *Player, hand Hand) (discard int) {
	// 可視
	return p.strategy.SelectOnPublicExecution(p, target, hand)
}

func (p *Player) SelectOnPlague(target *Player, hand Hand) (discard int) {
	// 不可視
	return p.strategy.SelectOnPlague(p, target, hand)
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

// 透視による開示
func (p Player) ShowForClairvoyance() int {
	return p.hand.Get()
}

func (p *Player) KnowByClairvoyance(g *Game, target *Player, c int) {
	p.strategy.KnowByClairvoyance(g, p, target, c)
}

func (p *Player) OnOpponentEvent(g *Game, opponent *Player, e CardEvent) {
	p.strategy.OnOpponentEvent(g, p, opponent, e)
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
