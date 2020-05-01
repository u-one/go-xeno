package xeno

import (
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestGame_ProcessTurn(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrategyH := NewMockPlayerStrategy(ctrl)
	mockStrategyN := NewMockPlayerStrategy(ctrl)
	playerH := &Player{
		id:        0,
		name:      "ヒカル",
		hand:      Hand{cards: []int{7}},
		discarded: []int{},
		strategy:  mockStrategyH,
	}
	playerN := &Player{
		id:        1,
		name:      "中田",
		hand:      Hand{cards: []int{8}},
		discarded: []int{},
		strategy:  mockStrategyN,
	}

	g := Game{
		Deck: &Deck{cards: []int{5, 4, 1, 2, 2, 3, 3, 4, 5, 6, 6, 7, 8, 9, 10}, reincCard: 1},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: false,
		turn:        2,
	}

	mockStrategyH.EXPECT().SelectDiscard(gomock.Any(), gomock.Any()).Return(CardEvent{Card: 5, Target: playerN})
	mockStrategyH.EXPECT().SelectOnPlague(gomock.Any(), gomock.Any(), Hand{cards: []int{8, 4}}).Return(4)

	g.ProcessTurn()

	gwant := Game{
		Deck: &Deck{cards: []int{1, 2, 2, 3, 3, 4, 5, 6, 6, 7, 8, 9, 10}, reincCard: 1},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: false,
		turn:        2,
	}

	if !reflect.DeepEqual(g, gwant) {
		t.Errorf("want: %v, got: %v\n", gwant, g)
	}

}
