package xeno

//go:generate mockgen -source=game.go -destination=./game_mock.go -package xeno

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

// CardEvent represents event occured by discard
type CardEvent struct {
	Card   int
	Target *Player
	Expect int // 捜査用。TODO: いずれ分離
}

func pause() {
}

// Define Shuffler interface to make test easier
type Shuffler interface {
	Shuffle([]int) []int
}

type RandomShuffler struct{}

func (s RandomShuffler) Shuffle(cards []int) []int {
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

type Deck struct {
	cards     []int
	reincCard int
	shuffler  Shuffler
}

func newDeck() *Deck {
	cards := AllCards

	shuffler := RandomShuffler{}
	cards = shuffler.Shuffle(cards)

	d := Deck{
		cards:     cards[:len(cards)-1],
		reincCard: cards[len(cards)-1],
		shuffler:  shuffler,
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

	d.cards = d.shuffler.Shuffle(d.cards)
}

// 転生
func (d *Deck) ReincarnateCard() (bool, int) {
	if d.reincCard == 0 {
		return false, 0
	}

	c := d.reincCard
	d.reincCard = 0
	return true, c
}

type GameConfig struct {
	Players []PlayerConfig
}

type Game struct {
	Deck        *Deck
	Players     []*Player
	boyAppeared bool
	turn        int
}

func NewGame(conf GameConfig) *Game {
	deck := newDeck()

	players := make([]*Player, len(conf.Players))
	for i, c := range conf.Players {
		players[i] = NewPlayer(c)
	}

	return &Game{
		Deck:    deck,
		Players: players,
	}
}

func (g Game) CurrentPlayer() *Player {
	i := g.turn % len(g.Players)
	return g.Players[i]
}

func (g Game) AlivePlayerCount() int {
	alive := 0
	for _, p := range g.Players {
		if !p.Dropped() {
			alive++
		}
	}
	return alive
}

func (g Game) OtherPlayers(p *Player) []*Player {
	others := []*Player{}
	for _, op := range g.Players {
		if op.ID() != p.ID() {
			others = append(others, op)
		}
	}
	return others
}

func (g *Game) Loop() {
	fmt.Println("山札:", g.Deck)
	fmt.Println("プレイヤー数:", len(g.Players))

	for {
		g.ProcessTurn()

		if g.Deck.finished() {
			fmt.Println("山札なし")
			fmt.Println("ゲーム終了")
			fmt.Println("_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/")
			var max, maxi int
			for i, p := range g.Players {
				if !p.Dropped() {
					fmt.Printf("%sのカード: %s\n", p.Name(), p.Hand())
					if max < p.Hand().Get() {
						max = p.Hand().Get()
						maxi = i
					}
				}
			}
			for i, p := range g.Players {
				if i != maxi && !p.Dropped() {
					p.Dropout()
					fmt.Printf("%s 脱落\n", p.Name())
				}
			}
			break
		} else if g.AlivePlayerCount() < 2 {
			fmt.Println("ゲーム終了")
			break
		}
		g.turn++
	}

	for _, p := range g.Players {
		if !p.Dropped() {
			fmt.Printf("%s の勝ち!\n", p.Name())
		}
	}
	fmt.Println("_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/")
}

func (g *Game) ProcessTurn() {
	fmt.Println(g)
	defer fmt.Println("======================================")

	p := g.CurrentPlayer()

	fmt.Printf("%s の番 \n", p.Name())

	if p.Dropped() {
		fmt.Printf("%s 脱落 スキップ\n", p.Name())
		return
	}

	var next int
	if p.CalledWise() {
		fmt.Println("賢者からの選択: ")
		candidates := g.Deck.takeN(3)
		{
			str := ""
			for _, c := range candidates {
				str += fmt.Sprintf("[%d]", c)
			}
			debugPrintf("%s\n", str)
		}
		var remains []int
		remains = p.TakeFromWise(g, candidates)
		debugPrintf("[%d]を選択\n", next)
		g.Deck.takeBack(remains)
	} else {
		fmt.Println("山札から引く: ")
		next = g.Deck.take()
		p.Take(next)
	}
	p.SetCalledWise(false)
	p.SetProtected(false)

	debugPrintf("引いたカード: [%d]\n", next)
	fmt.Println("どちらかを捨てる:")
	debugPrintf("%s\n", p.Hand())

	event := p.Discard(g)

	if event.Card > 0 {
		fmt.Printf("捨てたカード: [%d] %s\n", event.Card, CardTypes[event.Card])
	}

	switch event.Card {
	case 1:
		if !g.boyAppeared {
			fmt.Println("少年1枚目。効果発動なし。")
		} else {
			fmt.Println("少年2枚目。革命。公開処刑が発動。")
			// 公開処刑
			g.publicExecution(p, event.Target, false)
		}
		g.boyAppeared = true
	case 2: // 捜査
		fmt.Printf("捜査の効果: %sは%sに手札を言い当てられると脱落。\n", event.Target.Name(), p.Name())
		g.investigation(p, event.Target, event.Expect)
	case 3: // 透視
		fmt.Printf("透視の効果: %sは%sの手札を見ることができる。\n", p.Name(), event.Target.Name())
		// TODO: Implement here
	case 4: // 守護
		fmt.Printf("守護の効果: %sは次の手番まで自分への効果が無効。\n", p.Name())
		p.SetProtected(true)
	case 5: // 疫病
		fmt.Printf("疫病の効果: %sは%sに1枚引かせて、非公開で1枚捨てさせる。\n", p.Name(), event.Target.Name())
		g.plague(p, event.Target)
	case 6: // 対決
		fmt.Printf("対決の効果: %sと%sで手札が小さい方が脱落。\n", p.Name(), event.Target.Name())
		g.confrontation(p, event.Target)
	case 7: // 選択
		p.SetCalledWise(true)
		fmt.Printf("選択の効果: %sは次ターンで3枚引く。\n", p.Name())
	case 8: // 交換
		fmt.Printf("交換の効果: %sと%sはカードを交換。\n", p.Name(), event.Target.Name())
		pc := p.Give()
		tc := event.Target.Give()
		p.Take(tc)
		event.Target.Take(pc)
	case 9: // 公開処刑
		fmt.Printf("公開処刑の効果: %sは%sに1枚引かせて、公開し1枚捨てさせる。\n", p.Name(), event.Target.Name())
		g.publicExecution(p, event.Target, true)
	case 10:
		// 有り得ない
	}
	return
}

// 対決
func (g *Game) confrontation(executor, target *Player) {
	if executor.Hand().Get() > target.Hand().Get() {
		fmt.Printf("%s の勝ち\n", executor.Name())
		target.Dropout()
	} else if executor.Hand().Get() < target.Hand().Get() {
		fmt.Printf("%s の勝ち\n", target.Name())
		executor.Dropout()
	} else {
		fmt.Printf("引き分け\n")
		target.Dropout()
		executor.Dropout()
	}
}

// 公開処刑
func (g *Game) publicExecution(executor, target *Player, fromEmperror bool) {
	fmt.Printf("公開処刑 ターゲット:%s\n", target.Name())
	if target.Protected() {
		fmt.Printf("ターゲット:%sは守護下\n", target.Name())
		return
	}
	if g.Deck.finished() {
		fmt.Printf("残り山札なし")
		return
	}

	// target
	next := g.Deck.take()
	target.Take(next)
	fmt.Printf("%sの手札: %s\n", target.Name(), target.Hand())
	// TODO: 引数でPairを渡すか？なるべくゲームルールをここで表現するため、こうしたい
	discard := executor.SelectOnPublicExecution(target, target.Hand())
	target.DiscardSpecified(discard)
	fmt.Printf("捨てるカードを指定:[%d]\n", discard)

	if discard == 10 {

		if fromEmperror {
			fmt.Println("英雄が皇帝に見つかった")
			fmt.Printf("%s 脱落\n", target.Name())
			target.Dropout()
		} else {
			fmt.Println("英雄が皇帝以外にやられた")
			ok, c := g.Deck.ReincarnateCard()
			if ok {
				fmt.Println("転生")
				target.Reincarnate(c)
			} else {
				fmt.Printf("転生不可 %s 脱落\n", target.Name())
				target.Dropout()
			}
		}
	}
}

// 疫病
func (g *Game) plague(executor, target *Player) {
	fmt.Printf("疫病 ターゲット:%s\n", target.Name())
	if target.Protected() {
		fmt.Printf("ターゲット:%sは守護下\n", target.Name())
		return
	}
	if g.Deck.finished() {
		fmt.Printf("残り山札なし")
		return
	}

	next := g.Deck.take()
	target.Take(next)
	fmt.Println("[?][?]")
	debugPrintf("%s\n", target.Hand())
	discard := executor.SelectOnPlague(target, target.Hand())
	target.DiscardSpecified(discard)
	fmt.Printf("指定:[%d]\n", discard)

	if discard == 10 {
		// 死神・兵士・少年の効果で脱落した場合は、持っている手札を全て捨ててから転生札を引き、ゲームに復帰
		fmt.Println("英雄がやられた")
		ok, c := g.Deck.ReincarnateCard()
		if ok {
			fmt.Println("転生")
			target.Reincarnate(c)
		} else {
			fmt.Printf("転生不可 %s 脱落\n", target.Name())
			target.Dropout()
		}
	}
}

func (g *Game) investigation(executor, target *Player, expect int) {
	fmt.Printf("%sに対する捜査 %d\n", target.Name(), expect)

	correct := target.Has(expect)
	if correct {
		fmt.Printf("正解 %s脱落\n", target.Name())
		target.Dropout()
	} else {
		fmt.Printf("はずれ\n")
	}
}

func (g Game) String() string {
	text := ""
	text += fmt.Sprintf("----- ターン%d ------------------------\n", g.turn)
	text += fmt.Sprintf("= 残り: %d枚\n", g.Deck.count())
	if true {
		text += fmt.Sprintf("= 山札:%v 転生札:%d\n", g.Deck.cards, g.Deck.reincCard)
		for _, lp := range g.Players {
			text += fmt.Sprintf("= %s\n", lp)
		}
	}
	text += fmt.Sprintf("--------------------------------------\n")
	return text
}

func debugPrintf(msg string, args ...interface{}) {
	m := "--" + msg
	fmt.Printf(m, args...)
}
