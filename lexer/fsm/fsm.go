package fsm

import "bytes"

type Transition func(current State, input rune) State

type State struct {
	Value   int
	Accepts bool
}

var NullState = State{-1, false}

type FSM struct {
	states       map[State]struct{}
	initial      State
	getNextState Transition
}

/// New returns a Finite State Machine
func New(states []State, initial State, getNextState Transition) *FSM {
	fsm := &FSM{
		initial:      initial,
		getNextState: getNextState,
		states:       make(map[State]struct{}, len(states)),
	}
	for i := range states {
		fsm.states[states[i]] = struct{}{}
	}

	return fsm
}

/// Run returns a value and whether a valid value was found
func (f *FSM) Run(input []rune) (*bytes.Buffer, bool) {
	currentState := &f.initial

	var output bytes.Buffer

	for _, b := range input {
		nextState := f.getNextState(*currentState, b)

		// Check if the state is a NullState
		if nextState.Value == NullState.Value {
			break
		}

		output.WriteRune(b)
		currentState = &nextState
	}

	isValid := currentState.Accepts
	return &output, isValid

}
