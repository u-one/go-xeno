package main

import (
	"fmt"
	"log"
	"math/rand"
)

var (
	CardTypes = []string{
		"",          //
		"少年(革命)",    // 1
		"兵士(捜査)",    // 2
		"占師(透視)",    // 3
		"乙女(守護)",    // 4
		"死神(疫病)",    // 5
		"貴族(対決)",    // 6
		"賢者(選択)",    // 7
		"精霊(交換)",    // 8
		"皇帝(公開処刑)",  // 9
		"英雄(潜伏・転生)", // 10
	}

	AllCards []int = []int{1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8, 8, 9, 10}
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

type Deck struct {
	cards     []int
	reincCard int
}

func newDeck() *Deck {
	cards := AllCards
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	d := Deck{
		cards:     cards[:len(cards)-1],
		reincCard: cards[len(cards)-1],
	}

	return &d
}

func (d Deck) finished() bool {
	return len(d.cards) == 0

}

func (d Deck) count() int {
	return len(d.cards)
}

func (d *Deck) take() int {
	if len(d.cards) < 1 {
		log.Fatal("no remaining card")
	}
	c := d.cards[0]
	d.cards = d.cards[1:]
	return c
}

func (d *Deck) takeN(n int) []int {
	var cards []int
	for i := 0; d.count() > 0 && i < 3; i++ {
		cards = append(cards, d.take())
	}
	return cards
}

func (d *Deck) takeBack(cards []int) {
	d.cards = append(d.cards, cards...)

	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// 転生
func (d *Deck) reincarnate() (bool, int) {
	if d.reincCard == 0 {
		return false, 0
	}

	c := d.reincCard
	d.reincCard = 0
	return true, c
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

func (p *Player) selectOne(g *Game, next int) (result SelectResult) {
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

func (p *Player) selectFromWise(candidates []int) (selected int, remains []int) {
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

func (p *Player) discard(card int) (discarded int) {
	discarded = card
	p.discarded = append(p.discarded, discarded)
	return
}

func (p *Player) dropout() {
	p.current = 0
	p.dropped = true
	fmt.Printf("%s脱落\n", p.name)
}

func (p *Player) reincarnate(newCard int) {
	p.current = newCard
}

func (p Player) String() string {
	return fmt.Sprintf("%s: Current:%d %v", p.name, p.current, p.discarded)
}

func (p Player) has(expect int) bool {
	return (p.current == expect)
}

func debugPrintf(msg string, args ...interface{}) {
	m := "--" + msg
	fmt.Printf(m, args...)
}

// 対決
func (g *Game) confrontation(executor, target *Player) {
	if executor.current > target.current {
		target.discard(target.current)
		target.dropout()
	} else if executor.current < target.current {
		executor.discard(executor.current)
		executor.dropout()
	} else {
		target.discard(target.current)
		target.dropout()
		executor.discard(executor.current)
		executor.dropout()
	}
}

// 公開処刑
func (g *Game) publicExecution(executor, target *Player, fromEmperror bool) {
	fmt.Printf("公開処刑 ターゲット:%s\n", target.name)
	if target.protected {
		fmt.Printf("ターゲット:%sは守護下\n", target.name)
		return
	}
	if g.Deck.finished() {
		fmt.Printf("残り山札なし")
		return
	}

	next := g.Deck.take()
	fmt.Printf("[%d][%d]\n", target.current, next)

	cards := []int{target.current, next}

	// 可視
	var discard, remain int
	if cards[0] < cards[1] {
		discard = cards[1]
		remain = cards[0]
	} else {
		discard = cards[0]
		remain = cards[1]
	}
	target.current = remain
	target.discard(discard)
	fmt.Printf("指定:[%d]\n", discard)

	if discard == 10 {
		for _, c := range cards {
			target.discard(c)
		}

		if fromEmperror {
			fmt.Println("英雄が皇帝に見つかった")
			target.dropout()
		} else {
			fmt.Println("英雄が皇帝以外にやられた")
			ok, c := g.Deck.reincarnate()
			if ok {
				fmt.Println("転生")
				target.reincarnate(c)
			} else {
				fmt.Println("転生不可")
				target.dropout()
			}
		}
	}

}

// 疫病
func (g *Game) plague(executor, target *Player) {
	fmt.Printf("疫病 ターゲット:%s\n", target.name)
	if target.protected {
		fmt.Printf("ターゲット:%sは守護下\n", target.name)
		return
	}
	if g.Deck.finished() {
		fmt.Printf("残り山札なし")
		return
	}

	next := g.Deck.take()
	// 不可視
	fmt.Println("[?][?]")
	debugPrintf("[%d][%d]\n", target.current, next)

	cards := []int{target.current, next}

	var discard, remain int
	if 1 == rand.Intn(2) {
		discard = cards[1]
		remain = cards[0]
	} else {
		discard = cards[0]
		remain = cards[1]
	}
	target.current = remain
	target.discard(discard)
	fmt.Printf("指定:[%d]\n", discard)

	if discard == 10 {
		// 死神・兵士・少年の効果で脱落した場合は、持っている手札を全て捨ててから転生札を引き、ゲームに復帰
		for _, c := range cards {
			target.discard(c)
		}

		fmt.Println("英雄がやられた")
		ok, c := g.Deck.reincarnate()
		if ok {
			fmt.Println("転生")
			target.reincarnate(c)
		} else {
			fmt.Println("転生不可")
			target.dropout()
		}
	}
}

func (g *Game) investigate(executor, target *Player, expect int) {
	fmt.Printf("%sに対する捜査 %d\n", target.name, expect)

	correct := target.has(expect)
	if correct {
		fmt.Printf("正解 %s脱落\n", target.name)
		target.dropout()
	} else {
		fmt.Printf("はずれ\n")
	}
}

type Game struct {
	Deck        *Deck
	PlayerNum   int
	Players     []*Player
	boyAppeared bool
	turn        int
}

func NewGame() *Game {
	deck := newDeck()
	playerNum := 2

	fmt.Println("山札:", deck)
	fmt.Println("プレイヤー数:", playerNum)

	players := make([]*Player, playerNum)
	for i := range players {
		players[i] = &Player{
			id:   i,
			name: fmt.Sprintf("Player%d", i),
		}
	}

	return &Game{
		Deck:      deck,
		PlayerNum: playerNum,
		Players:   players,
	}
}

func (g Game) CurrentPlayer() *Player {
	i := g.turn % g.PlayerNum
	return g.Players[i]
}

func (g Game) AlivePlayerCount() int {
	alive := 0
	for _, p := range g.Players {
		if !p.dropped {
			alive++
		}
	}
	return alive
}

func (g Game) OtherPlayers(p *Player) []*Player {
	others := []*Player{}
	for _, op := range g.Players {
		if op.id != (*p).id {
			debugPrintf("op:%v, p:%v\n", op, p)
			others = append(others, op)
		}
	}
	debugPrintf("others:%v", others)
	return others
}

func (g *Game) Loop() {
	for {
		fmt.Println("----")
		p := g.CurrentPlayer()

		done := func() bool {
			fmt.Printf("ターン%d [%s]の番 残り: %d枚\n", g.turn, p.name, g.Deck.count())
			debugPrintf("山札:%v 転生札:%d\n", g.Deck.cards, g.Deck.reincCard)

			if p.dropped {
				fmt.Printf("%s 脱落 スキップ\n", p.name)
				return false
			}

			var next int
			if p.calledWise {
				candidates := g.Deck.takeN(3)
				var remains []int
				next, remains = p.selectFromWise(candidates)
				g.Deck.takeBack(remains)
			} else {
				next = g.Deck.take()
			}
			p.calledWise = false
			p.protected = false
			res := p.selectOne(g, next)

			if res.Discarded > 0 {
				fmt.Printf("捨てたカード %d: %s\n", res.Discarded, CardTypes[res.Discarded])
			}

			switch res.Discarded {
			case 1:
				if g.boyAppeared {
					// 公開処刑
					g.publicExecution(p, res.Target, false)
				}
				g.boyAppeared = true
			case 2: // 捜査
				g.investigate(p, res.Target, res.Expect)
			case 3: // 透視
			case 4: // 守護
				p.protected = true
			case 5: // 疫病
				g.plague(p, res.Target)
			case 6: // 対決
				g.confrontation(p, res.Target)
			case 7: // 選択
				p.calledWise = true
			case 8: // 交換
				fmt.Printf("target: %s\n", res.Target.name)
				pc, tc := p.current, res.Target.current
				p.current, res.Target.current = tc, pc
			case 9: // 公開処刑
				g.publicExecution(p, res.Target, true)
			case 10:
				// 有り得ない
			}
			return false
		}()

		if g.Deck.finished() {
			fmt.Println("山札なし ゲーム終了")
			done = true
		} else if g.AlivePlayerCount() < 2 {
			done = true
		}

		if done {
			break
		}

		g.turn++
	}
	fmt.Printf("Finish game\n")

	fmt.Printf("%v\n", g.Players)
	for _, p := range g.Players {
		if !p.dropped {
			fmt.Printf("%s の勝ち!\n", p.name)
		}
	}
}

func main() {
	game := NewGame()
	game.Loop()
}
