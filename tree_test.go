package redblack

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestRBTree_NewRBTree(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
	}{
		{"Test1", []int64{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			tree := NewTree(tt.values)
			fmt.Println(tree.Height())

			fmt.Println(tree)
		})
	}
}

func TestTree_ToSortedSlice(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int64
	}{
		{"Test1", []int64{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		//{"Test3", rand.Perm(32)},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})

			fmt.Println(tt.values)
			t := NewTree(tt.values)
			got := t.ToSortedSlice()
			sort.Slice(tt.values, func(i, j int) bool { return tt.values[i] < tt.values[j] })
			if !reflect.DeepEqual(got, tt.values) {
				t1.Errorf("ToSortedSlice() = %v, want %v", got, tt.values)
				t1.Errorf(t.String())
			}
			fmt.Println(got)

		})
	}
}

func TestTree_checkBlackHeight(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int64
		want   uint
	}{
		{"Test1", []int64{1, 2, 3, 4, 5, 6, 7, 8}, 3},
		{"Test2", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 4},
		//{"Test3", rand.Perm(32), 5},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			//fmt.Println(tt.values)
			t := NewTree(tt.values)
			//t.String()
			blackHeight, checked := t.checkBlackHeight()
			if !checked {
				t1.Errorf("checkBlackHeight() returned that the black-height is not equal for all paths")
				t1.Errorf(t.String())
			}
			if blackHeight != tt.want {
				t1.Errorf("checkBlackHeight() yields = %v, want %v", blackHeight, tt.want)
				t1.Errorf(t.String())
			}
		})
	}
}

func TestTree_checkRedRed(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int64
	}{
		{"Test1", []int64{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		//{"Test3", rand.Perm(32)},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			//fmt.Println(tt.values)
			t := NewTree(tt.values)
			//t.String()
			checked := t.checkRedRed()
			if !checked {
				t1.Errorf("Tree violates no-red-red rule.")
				t1.Errorf(t.String())
			}
		})
	}
}

func TestTree_SearchUpper(t1 *testing.T) {
	tests := []struct {
		name    string
		values  []int64
		q       int64
		want    int64
		wantErr bool
	}{
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 10, 14, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 51, 67, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, -10, 1, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 100, 0, true},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 8, 8, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewTree(tt.values)
			got, err := t.SearchUpper(tt.q)
			if (err != nil) != tt.wantErr {
				t1.Errorf("SearchUpper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("SearchUpper() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_SearchLower(t1 *testing.T) {
	tests := []struct {
		name    string
		values  []int64
		q       int64
		want    int64
		wantErr bool
	}{
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 10, 8, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 51, 50, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, -10, 0, true},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 100, 67, false},
		{"Test1", []int64{1, 2, 5, 8, 14, 23, 44, 50, 67}, 8, 8, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewTree(tt.values)
			got, err := t.SearchLower(tt.q)
			if (err != nil) != tt.wantErr {
				t1.Errorf("SearchUpper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("SearchUpper() got = %v, want %v", got, tt.want)
			}
		})
	}
}
