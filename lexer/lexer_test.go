package lexer

import (
	"reflect"
	"testing"

	"./fsm"
)

func TestNew(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want *Lexer
	}{
		{
			"create lexer",
			args{"this is code"},
			&Lexer{input: []byte("this is code")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextState(t *testing.T) {
	type args struct {
		currentState fsm.State
		input        byte
	}
	tests := []struct {
		name string
		args args
		want fsm.State
	}{
		{
			"Initial with digit",
			args{InitialState, '8'},
			IntegerState,
		},
		{
			"Initial with decimal point",
			args{InitialState, '.'},
			NullState,
		},
		{
			"Initial with exponent,e",
			args{InitialState, 'e'},
			NullState,
		},
		{
			"Initial with exponent,E",
			args{InitialState, 'E'},
			NullState,
		},
		{
			"Initial with sign,+",
			args{InitialState, '+'},
			NullState,
		},
		{
			"Initial with sign,-",
			args{InitialState, '-'},
			NullState,
		},
		{
			"Integer with digit",
			args{IntegerState, '8'},
			IntegerState,
		},
		{
			"Integer with decimal point",
			args{IntegerState, '.'},
			BeginsFloatState,
		},
		{
			"Integer with exponent,e",
			args{IntegerState, 'e'},
			BeginExpState,
		},
		{
			"Integer with exponent,E",
			args{IntegerState, 'E'},
			BeginExpState,
		},
		{
			"Integer with sign,+",
			args{IntegerState, '+'},
			NullState,
		},
		{
			"Integer with sign,-",
			args{IntegerState, '-'},
			NullState,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextState(tt.args.currentState, tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextState() = %v, want %v", got, tt.want)
			}
		})
	}
}
