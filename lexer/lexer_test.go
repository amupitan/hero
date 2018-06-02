package lexer

import (
	"reflect"
	"testing"
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
		// TODO: Add test cases.
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
