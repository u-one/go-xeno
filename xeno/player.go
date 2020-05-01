package xeno

//go:generate mockgen -source=player.go -destination=./player_mock.go -package xeno

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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
		s = CommStrategy{}
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
