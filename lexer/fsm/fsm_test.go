package fsm

import (
	"reflect"
	"testing"
)

func TestFSM_Run(t *testing.T) {
	type fields struct {
		states       map[State]struct{}
		initial      State
		getNextState Transition
	}
	type args struct {
		input []byte
	}
	empty := struct{}{}
	state1, state2 := State{1, true}, State{2, true}
	mockNextState := func(current State, input byte) State {
		value := int(input - '0')
		isAccepts := value%2 == 0 //isAccepts if even
		return State{value, isAccepts}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  bool
	}{
		{
			"one state",
			fields{
				states:       map[State]struct{}{state1: empty, state2: empty},
				initial:      state1,
				getNextState: mockNextState,
			},
			args{[]byte("1 + 2")},
			[]byte("1 + 2"),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				states:       tt.fields.states,
				initial:      tt.fields.initial,
				getNextState: tt.fields.getNextState,
			}
			got, got1 := f.Run(tt.args.input)
			if !reflect.DeepEqual(got.Bytes(), tt.want) {
				t.Errorf("FSM.Run() got = %v, want %v", got.String(), string(tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("FSM.Run() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
