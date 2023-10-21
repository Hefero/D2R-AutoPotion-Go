package itemfilter

import (
	"testing"

	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/data/item"
	"github.com/Hefero/D2R-AutoPotion-Go/pkg/nip"
	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	type args struct {
		i       data.Item
		nipRule string // not the best test but too lazy to write rules manually
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				i: data.Item{
					Name:     item.Name("lightplatedboots"),
					Quality:  item.QualityUnique,
					Ethereal: false,
				},
				nipRule: "[name] == lightplatedboots && [quality] == unique && [flag] != ethereal # [enhanceddefense] == 60 // goblin toe",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule, err := nip.ParseLine(tt.args.nipRule)
			require.NoError(t, err)

			if got := Evaluate(tt.args.i, []nip.Rule{rule}); got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
