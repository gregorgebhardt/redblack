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
		values []int
	}{
		{"Test1", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
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
		values []int
	}{
		{"Test1", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		{"Test3", rand.Perm(32)},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			fmt.Println(tt.values)
			t := NewTree(tt.values)
			got := t.ToSortedSlice()
			sort.Ints(tt.values)
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
		values []int
		want   uint
	}{
		{"Test1", []int{1, 2, 3, 4, 5, 6, 7, 8}, 3},
		{"Test2", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 4},
		{"Test3", rand.Perm(32), 5},
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
		values []int
	}{
		{"Test1", []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Test2", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		{"Test3", rand.Perm(32)},
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
