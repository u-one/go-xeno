package xeno

import (
	reflect "reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestPlayer_ShowForClairvoyance(t *testing.T) {
	p := Player{
		hand: Hand{cards: []int{10}},
	}

	got := p.ShowForClairvoyance()
	want := 10
	if got != want {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestPlayer_KnowByClairvoyance(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStrategy := NewMockPlayerStrategy(ctrl)

	g := Game{}

	p := Player{
		id:       0,
		hand:     Hand{cards: []int{}},
		strategy: mockStrategy,
	}

	o := Player{
		id: 1,
	}

	mockStrategy.EXPECT().KnowByClairvoyance(gomock.Any(), &p, &o, 10)

	p.KnowByClairvoyance(&g, &o, 10)
}

func TestComStrategy_KnowByClairvoyance(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStrategy := NewMockPlayerStrategy(ctrl)

	g := Game{}

	p := Player{
		id:       0,
		hand:     Hand{cards: []int{}},
		strategy: mockStrategy,
	}

	o := Player{
		id: 1,
	}

	s := CommStrategy{
		opponentInfo: map[PlayerID]int{},
	}

	s.KnowByClairvoyance(&g, &p, &o, 10)

	got := s
	want := CommStrategy{
		opponentInfo: map[PlayerID]int{1: 10},
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestComStrategy_OnOpponentEvent(t *testing.T) {
	tests := []struct {
		name  string
		state CommStrategy
		event CardEvent
		want  CommStrategy
	}{
		{
			name: "delete from opponent info",
			state: CommStrategy{
				opponentInfo: map[PlayerID]int{1: 1},
			},
			event: CardEvent{Card: 1},
			want: CommStrategy{
				opponentInfo: map[PlayerID]int{},
			},
		},
		{
			name: "should not delete from opponent info",
			state: CommStrategy{
				opponentInfo: map[PlayerID]int{1: 1},
			},
			event: CardEvent{Card: 2},
			want: CommStrategy{
				opponentInfo: map[PlayerID]int{1: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Game{}

			p := Player{
				id:       0,
				hand:     Hand{cards: []int{}},
				strategy: tt.state,
			}

			o := Player{
				id: 1,
			}

			s := tt.state

			s.OnOpponentEvent(&g, &p, &o, tt.event)

			got := s

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want:%v, got: %v", tt.want, got)
			}

		})
	}

}

func TestComStrategy_SelectDiscard(t *testing.T) {
	tests := []struct {
		name  string
		state CommStrategy
		hand  Hand
		want  CommStrategy
	}{
		{
			name: "add opponent info exchanged",
			state: CommStrategy{
				opponentInfo: map[PlayerID]int{},
			},
			hand: Hand{cards: []int{10, 8}},
			want: CommStrategy{
				opponentInfo: map[PlayerID]int{1: 10},
			},
		},
		{
			name: "no change in opponent info",
			state: CommStrategy{
				opponentInfo: map[PlayerID]int{},
			},
			hand: Hand{cards: []int{10, 9}},
			want: CommStrategy{
				opponentInfo: map[PlayerID]int{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{
				id:       0,
				hand:     tt.hand,
				strategy: tt.state,
			}
			o := Player{
				id: 1,
			}

			g := Game{
				Players: []*Player{&p, &o},
			}

			s := tt.state

			s.SelectDiscard(&g, &p)

			got := s

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want:%v, got: %v", tt.want, got)
			}

		})
	}

}
