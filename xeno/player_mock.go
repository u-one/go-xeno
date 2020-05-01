// Code generated by MockGen. DO NOT EDIT.
// Source: player.go

// Package xeno is a generated GoMock package.
package xeno

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPlayerStrategy is a mock of PlayerStrategy interface
type MockPlayerStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockPlayerStrategyMockRecorder
}

// MockPlayerStrategyMockRecorder is the mock recorder for MockPlayerStrategy
type MockPlayerStrategyMockRecorder struct {
	mock *MockPlayerStrategy
}

// NewMockPlayerStrategy creates a new mock instance
func NewMockPlayerStrategy(ctrl *gomock.Controller) *MockPlayerStrategy {
	mock := &MockPlayerStrategy{ctrl: ctrl}
	mock.recorder = &MockPlayerStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPlayerStrategy) EXPECT() *MockPlayerStrategyMockRecorder {
	return m.recorder
}

// SelectDiscard mocks base method
func (m *MockPlayerStrategy) SelectDiscard(g *Game, p *Player) CardEvent {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectDiscard", g, p)
	ret0, _ := ret[0].(CardEvent)
	return ret0
}

// SelectDiscard indicates an expected call of SelectDiscard
func (mr *MockPlayerStrategyMockRecorder) SelectDiscard(g, p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectDiscard", reflect.TypeOf((*MockPlayerStrategy)(nil).SelectDiscard), g, p)
}

// SelectFromWise mocks base method
func (m *MockPlayerStrategy) SelectFromWise(g *Game, candidates []int) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectFromWise", g, candidates)
	ret0, _ := ret[0].(int)
	return ret0
}

// SelectFromWise indicates an expected call of SelectFromWise
func (mr *MockPlayerStrategyMockRecorder) SelectFromWise(g, candidates interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectFromWise", reflect.TypeOf((*MockPlayerStrategy)(nil).SelectFromWise), g, candidates)
}

// SelectOnPublicExecution mocks base method
func (m *MockPlayerStrategy) SelectOnPublicExecution(p, target *Player, pair Hand) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectOnPublicExecution", p, target, pair)
	ret0, _ := ret[0].(int)
	return ret0
}

// SelectOnPublicExecution indicates an expected call of SelectOnPublicExecution
func (mr *MockPlayerStrategyMockRecorder) SelectOnPublicExecution(p, target, pair interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectOnPublicExecution", reflect.TypeOf((*MockPlayerStrategy)(nil).SelectOnPublicExecution), p, target, pair)
}

// SelectOnPlague mocks base method
func (m *MockPlayerStrategy) SelectOnPlague(p, target *Player, hand Hand) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectOnPlague", p, target, hand)
	ret0, _ := ret[0].(int)
	return ret0
}

// SelectOnPlague indicates an expected call of SelectOnPlague
func (mr *MockPlayerStrategyMockRecorder) SelectOnPlague(p, target, hand interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectOnPlague", reflect.TypeOf((*MockPlayerStrategy)(nil).SelectOnPlague), p, target, hand)
}
