package blob

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseIntPair(t *testing.T) {
	type args struct {
		pair string
	}
	tests := []struct {
		name     string
		args     args
		wantData [2]int
	}{
		// TODO: Add test cases.
		{
			name: "ParseIntPair",
			args: args{
				pair: "[1,2]",
			},
			wantData: [2]int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotData := ParseIntPair(tt.args.pair); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ParseIntPair() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestTagSort(t *testing.T) {
	str := `
1: [1,2] [3,4]
2: [1,4]
3: [0,1] [2,3] [2,3]
4: [0,1] [4,6] [2,3] [2,3]
5: [0,3] [5,6] [0,3] [0,3]
6: [2,3] [4,6] [2,3] [2,3]
7: [2,5]
9: [4,5]
8: [4,5]
8: [2,3]`
	fmt.Println(tagsort(str))
}
