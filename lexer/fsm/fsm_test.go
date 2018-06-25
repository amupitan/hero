package fsm

import (
	"reflect"
	"testing"
)

func TestFSM_Run(t *testing.T) {
	type fields struct {
		states       []State
		initial      State
		getNextState Transition
	}
	type args struct {
		input []rune
	}
	state1, state2 := State{1, true}, State{2, true}
	mockNextState := func(current State, input rune) State {
		value := int(input - '0')
		isAccepts := value%2 == 0 //isAccepts if even
		return State{value, isAccepts}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []rune
		want1  bool
	}{
		{
			"one state",
			fields{
				states:       []State{state1, state2},
				initial:      state1,
				getNextState: mockNextState,
			},
			args{[]rune("1 + 2")},
			[]rune("1 + 2"),
			true,
		},
		{
			"null state",
			fields{
				states:       []State{state1, NullState},
				initial:      state1,
				getNextState: func(current State, input rune) State { return NullState },
			},
			args{[]rune(`doesn't matter`)},
			[]rune(``),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(
				tt.fields.states,
				tt.fields.initial,
				tt.fields.getNextState,
			)
			got, got1 := f.Run(tt.args.input)
			if !reflect.DeepEqual(got.String(), string(tt.want)) {
				t.Errorf("FSM.Run() got = %v, want %v", got.String(), string(tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("FSM.Run() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
