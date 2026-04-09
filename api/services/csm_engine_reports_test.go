package services

import (
	"reflect"
	"testing"
)

func TestGenerateBelBuildupOneGroup(t *testing.T) {
	type args struct {
		ifrs17Group string
		productCode string
		runDate     string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateBelBuildupOneGroup(tt.args.ifrs17Group, tt.args.productCode, tt.args.runDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateBelBuildupOneGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
