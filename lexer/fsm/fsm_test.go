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
	mockNextState := func(current State, input byte) State {
		value := int(input - '0')
		isAccepts := value%2 == 0 //isAccepts if even
		return State{value, isAccepts}
	}
	_ = mockNextState //TODO
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FSM{
				states:       tt.fields.states,
				initial:      tt.fields.initial,
				getNextState: tt.fields.getNextState,
			}
			got, got1 := f.Run(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FSM.Run() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FSM.Run() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
