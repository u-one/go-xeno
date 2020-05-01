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
		name:      "Hikaru",
		hand:      Hand{cards: []int{7}},
		discarded: []int{},
		strategy:  mockStrategyH,
	}
	playerN := &Player{
		id:        1,
		name:      "Nakata",
		hand:      Hand{cards: []int{8}},
		discarded: []int{},
		strategy:  mockStrategyN,
	}

	g := Game{
		Deck: &Deck{cards: []int{5, 4}, reincCard: 1, shuffler: RandomShuffler{}},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: false,
		turn:        2,
	}

	mockStrategyH.EXPECT().SelectDiscard(gomock.Any(), gomock.Any()).Return(CardEvent{Card: 5, Target: playerN})
	mockStrategyH.EXPECT().SelectOnPlague(gomock.Any(), gomock.Any(), Hand{cards: []int{8, 4}}).Return(4)

	mockStrategyN.EXPECT().OnOpponentEvent(gomock.Any(), gomock.Any(), playerH, CardEvent{Card: 5, Target: playerN})

	g.ProcessTurn()

	gwant := Game{
		Deck: &Deck{cards: []int{}, reincCard: 1, shuffler: RandomShuffler{}},
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

func TestGame_ProcessTurn_Dropped(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStrategyH := NewMockPlayerStrategy(ctrl)
	mockStrategyN := NewMockPlayerStrategy(ctrl)
	playerH := &Player{
		id:        0,
		name:      "Hikaru",
		hand:      Hand{cards: []int{7}},
		discarded: []int{},
		strategy:  mockStrategyH,
		dropped:   true,
	}
	playerN := &Player{
		id:        1,
		name:      "Nakata",
		hand:      Hand{cards: []int{8}},
		discarded: []int{},
		strategy:  mockStrategyN,
	}

	g := Game{
		Deck: &Deck{cards: []int{5, 4}, reincCard: 1, shuffler: RandomShuffler{}},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: false,
		turn:        2,
	}

	g.ProcessTurn()

	gwant := Game{
		Deck: &Deck{cards: []int{5, 4}, reincCard: 1, shuffler: RandomShuffler{}},
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

func TestGame_ProcessTurn_Wise(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockShuffler := NewMockShuffler(ctrl)

	mockStrategyH := NewMockPlayerStrategy(ctrl)
	mockStrategyN := NewMockPlayerStrategy(ctrl)
	playerH := &Player{
		id:        0,
		name:      "Hikaru",
		hand:      Hand{cards: []int{6}},
		discarded: []int{4, 1, 2},
		strategy:  mockStrategyH,
	}
	playerN := &Player{
		id:         1,
		name:       "Nakata",
		hand:       Hand{cards: []int{6}},
		discarded:  []int{5, 7},
		strategy:   mockStrategyN,
		calledWise: true,
	}

	g := Game{
		Deck: &Deck{cards: []int{8, 1, 4, 7}, reincCard: 1, shuffler: mockShuffler},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: true,
		turn:        5,
	}

	mockStrategyN.EXPECT().SelectFromWise(gomock.Any(), []int{8, 1, 4}).Return(1)
	mockShuffler.EXPECT().Shuffle([]int{7, 8, 4}).Return([]int{7, 8, 4})
	mockStrategyN.EXPECT().SelectDiscard(gomock.Any(), gomock.Any()).Return(CardEvent{Card: 1, Target: playerH})
	mockStrategyN.EXPECT().SelectOnPublicExecution(gomock.Any(), gomock.Any(), Hand{cards: []int{6, 7}}).Return(7)

	mockStrategyH.EXPECT().OnOpponentEvent(gomock.Any(), gomock.Any(), playerN, CardEvent{Card: 1, Target: playerH})

	g.ProcessTurn()

	gwant := Game{
		Deck: &Deck{cards: []int{8, 4}, reincCard: 1, shuffler: mockShuffler},
		Players: []*Player{
			playerH,
			playerN,
		},
		boyAppeared: true,
		turn:        5,
	}

	if !reflect.DeepEqual(g, gwant) {
		t.Errorf("want: \n%v, got: \n%v\n", gwant, g)
	}

}
