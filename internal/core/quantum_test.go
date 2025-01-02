package core

import (
	"reflect"
	"testing"
)

func TestNewUnsignedQuantum(t *testing.T) {
	type args struct {
		contents   []*QContent
		last       string
		nonce      int
		references []string
	}
	tests := []struct {
		name string
		args args
		want *UnsignedQuantum
	}{
		{
			name: "test",
			args: args{
				contents:   []*QContent{},
				last:       DefaultLastSig,
				nonce:      1,
				references: []string{},
			},
			want: &UnsignedQuantum{
				Contents:   []*QContent{},
				Last:       DefaultLastSig,
				Nonce:      1,
				References: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnsignedQuantum(tt.args.contents, tt.args.last, tt.args.nonce, tt.args.references); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnsignedQuantum() = %v, want %v", got, tt.want)
			}
		})
	}
}
