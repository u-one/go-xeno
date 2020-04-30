package xeno

import (
	"reflect"
	"testing"
)

func TestHand_Set(t *testing.T) {
	h := Hand{cards: []int{}}
	h.Add(1)
	h.Add(2)

	got := h.cards
	want := []int{1, 2}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestHand_Get(t *testing.T) {
	h := Hand{cards: []int{1}}

	got := h.Get()
	want := 1
	if want != got {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestHand_Remove(t *testing.T) {
	h := Hand{cards: []int{1}}

	got := h.Remove()
	want := 1
	if want != got {
		t.Errorf("want:%v, got: %v", want, got)
	}

	if len(h.cards) != 0 {
		t.Errorf("len(h.cards) != 0, got: %v", len(h.cards))
	}
}

func TestHand_Clear(t *testing.T) {
	h := Hand{cards: []int{1}}

	h.Clear()
	got := h.cards
	want := []int{}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestHand_Count(t *testing.T) {
	h := Hand{cards: []int{1}}

	got := h.Count()
	want := 1
	if want != got {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestHand_Slice(t *testing.T) {
	h := Hand{cards: []int{1, 2}}

	got := h.Slice()
	want := []int{1, 2}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want:%v, got: %v", want, got)
	}
}

func TestHand_At(t *testing.T) {
	tests := []struct {
		name  string
		index int
		want  int
	}{
		{"0", 0, 1},
		{"1", 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hand{cards: []int{1, 2}}
			got := h.At(tt.index)
			if tt.want != got {
				t.Errorf("name: %s, want:%v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestHand_Another(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		arg  int
		want int
	}{
		{"0", Hand{[]int{1, 2}}, 1, 2},
		{"1", Hand{[]int{1, 2}}, 2, 1},
		{"2", Hand{[]int{1, 1}}, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hand
			got := h.Another(tt.arg)
			if tt.want != got {
				t.Errorf("name: %s, want:%v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestHand_Larger(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		want int
	}{
		{"0", Hand{[]int{1, 2}}, 2},
		{"1", Hand{[]int{2, 1}}, 2},
		{"2", Hand{[]int{2, 2}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hand
			got := h.Larger()
			if tt.want != got {
				t.Errorf("name: %s, want:%v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}

func TestHand_Random(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
	}{
		{"0", Hand{[]int{1, 2}}},
		{"1", Hand{[]int{2, 1}}},
		{"2", Hand{[]int{2, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hand
			got := h.Random()
			if got != 1 && got != 2 {
				t.Errorf("name: %s, want: 1 or 2, got: %v", tt.name, got)
			}
		})
	}
}

func TestHand_Has(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		arg  int
		want bool
	}{
		{"includes", Hand{[]int{9, 10}}, 10, true},
		{"includes", Hand{[]int{10, 9}}, 10, true},
		{"not includes", Hand{[]int{10, 9}}, 8, false},
		{"includes, same number", Hand{[]int{1, 1}}, 1, true},
		{"not includes, same number", Hand{[]int{1, 1}}, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.hand
			got := h.Has(tt.arg)
			if tt.want != got {
				t.Errorf("name: %s, want:%v, got: %v", tt.name, tt.want, got)
			}
		})
	}
}
