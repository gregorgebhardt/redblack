package redblack_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/gregorgebhardt/redblack"
)

func TestTree_NewTree(t1 *testing.T) {
	tests := []struct {
		name              string
		values            []int
		ignore_duplicates bool
		want_err          bool
	}{
		{"Successful", []int{1, 2, 3, 4, 5, 6, 7, 8}, false, false},
		{"Duplicate", []int{1, 2, 3, 4, 5, 6, 7, 8, 8}, false, true},
		{"Duplicate Ignore", []int{1, 2, 3, 4, 5, 6, 7, 8, 8}, true, false},
		{"Many items", rand.Perm(1024), false, false},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			t, err := redblack.NewTree(vals, tt.ignore_duplicates)
			if (err != nil) != tt.want_err || (err != nil && err != redblack.KeyExistsError) {
				t1.Errorf("NewTree() error = %v, wantErr %v", err, tt.want_err)
				return
			}
			// break tests if an error is expected
			if err != nil {
				return
			}

			if !redblack.CheckNoRedRed(t) {
				t1.Errorf("NewTree() resulted in red-red nodes")
			}
			if _, ok := redblack.CheckBlackHeight(t); !ok {
				t1.Errorf("NewTree() resulted in different black-heights")
			}
			if !redblack.CheckLeftLeaning(t) {
				t1.Errorf("NewTree() resulted in a right-leaning tree")
			}
		})
	}
}

func TestTree_Height(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"Empty Tree", []int{}, 0},
		{"One Element", []int{1}, 1},
		{"Two Elements", []int{1, 2}, 2},
		{"Tree Elements", []int{1, 2, 3}, 2},
		{"Four Elements", []int{1, 2, 3, 4}, 3},
		{"Five Elements", []int{1, 2, 3, 4, 5}, 3},
		{"Random Elements", rand.Perm(31), 10},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}
			if got := t.Height(); got > tt.want {
				t1.Errorf("Height() = %v, want %v", got, tt.want)
				t1.Error(t.String())
			}
		})
	}
}

func TestTree_Insert(t1 *testing.T) {
	tests := []struct {
		name    string
		values  []int
		insert  int
		want    []int
		wantErr bool
	}{
		{"Empty Tree", []int{}, 1, []int{1}, false},
		{"One Element", []int{1}, 2, []int{1, 2}, false},
		{"Two Elements", []int{1, 2}, 3, []int{1, 2, 3}, false},
		{"Tree Elements", []int{1, 2, 3}, 4, []int{1, 2, 3, 4}, false},
		{"Existing Element", []int{1, 2, 3}, 2, []int{1, 2, 3}, true},
		{"Negative Elements", []int{-1, -2, -3}, -4, []int{-4, -3, -2, -1}, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}
			err = t.Insert(redblack.Ordered(tt.insert))
			if err != nil && !tt.wantErr {
				t1.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if t.Len() != len(tt.want) {
				t1.Errorf("Insert() = %v, want %v", t.Len(), len(tt.want))
			}

			if !reflect.DeepEqual(t.ToSortedSlice(), tt.want) {
				t1.Errorf("Insert() = %v, want %v", t.ToSortedSlice(), tt.want)
			}
		})
	}
}

func TestTree_Delete(t1 *testing.T) {
	tests := []struct {
		name    string
		values  []int
		delete  int
		want    []int
		wantErr bool
	}{
		{"Empty Tree", []int{}, 1, []int{}, true},
		{"One Element", []int{1}, 1, []int{}, false},
		{"Two Elements", []int{1, 2}, 1, []int{2}, false},
		{"Tree Elements", []int{1, 2, 3}, 2, []int{1, 3}, false},
		{"Four Elements", []int{1, 2, 3, 4}, 2, []int{1, 3, 4}, false},
		{"Five Elements", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 4, 5}, false},
		{"Negative Elements", []int{-1, -2, -3, -4}, -3, []int{-4, -2, -1}, false},
		{"Non-Existing Element", []int{1, 2, 3}, 4, []int{1, 2, 3}, true},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}
			success := t.Delete(tt.delete)
			if success != !tt.wantErr {
				t1.Errorf("Delete() = %v, want %v", success, !tt.wantErr)
			}

			if !reflect.DeepEqual(t.ToSortedSlice(), tt.want) {
				t1.Errorf("Delete() = %v, want %v", t.ToSortedSlice(), tt.want)
			}

			if !redblack.CheckNoRedRed(t) {
				t1.Errorf("Delete() resulted in red-red nodes")
			}
			if _, ok := redblack.CheckBlackHeight(t); !ok {
				t1.Errorf("Delete() resulted in different black-heights")
			}
			if !redblack.CheckLeftLeaning(t) {
				t1.Errorf("Delete() resulted in a right-leaning tree")
			}
		})
	}
}

func TestTree_DeleteMin(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{"Empty Tree", []int{}, []int{}},
		{"One Element", []int{1}, []int{}},
		{"Two Elements", []int{1, 2}, []int{2}},
		{"Tree Elements", []int{1, 2, 3}, []int{2, 3}},
		{"Four Elements", []int{1, 2, 3, 4}, []int{2, 3, 4}},
		{"Five Elements", []int{1, 2, 3, 4, 5}, []int{2, 3, 4, 5}},
		{"Negative Elements", []int{-1, -2, -3, -4}, []int{-3, -2, -1}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}
			t.DeleteMin()

			if !reflect.DeepEqual(t.ToSortedSlice(), tt.want) {
				t1.Errorf("DeleteMin() = %v, want %v", t.ToSortedSlice(), tt.want)
			}
			if !redblack.CheckNoRedRed(t) {
				t1.Errorf("Delete() resulted in red-red nodes")
			}
			if _, ok := redblack.CheckBlackHeight(t); !ok {
				t1.Errorf("Delete() resulted in different black-heights")
			}
			if !redblack.CheckLeftLeaning(t) {
				t1.Errorf("Delete() resulted in a right-leaning tree")
			}
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
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			//fmt.Println(tt.values)
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

			got := make([]int, len(tt.values))
			for i, v := range t.ToSortedSlice() {
				got[i] = v
			}
			sort.Slice(tt.values, func(i, j int) bool { return tt.values[i] < tt.values[j] })
			if !reflect.DeepEqual(got, tt.values) {
				t1.Errorf("ToSortedSlice() = %v, want %v", got, tt.values)
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
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}
			//fmt.Println(tt.values)
			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

			blackHeight, checked := redblack.CheckBlackHeight(t)
			if !checked {
				t1.Errorf("checkBlackHeight() returned that the black-height is not equal for all paths")
			}
			if blackHeight != tt.want {
				t1.Errorf("checkBlackHeight() yields = %v, want %v", blackHeight, tt.want)
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
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}

			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

			noredred := redblack.CheckNoRedRed(t)
			if !noredred {
				t1.Errorf("checkRedRed() returned that there are red-red nodes")
				return
			}
		})
	}
}

func TestTree_SearchUpper(t1 *testing.T) {
	tests := []struct {
		name    string
		values  []int
		q       int
		want    int
		wantErr bool
	}{
		{"Test1", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 10, 14, false},
		{"Test2", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 51, 67, false},
		{"Test3", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, -10, 1, false},
		{"Test4", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 100, 100, true},
		{"Test5", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 8, 8, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}

			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

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
		values  []int
		q       int
		want    int
		wantErr bool
	}{
		{"Test1", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 10, 8, false},
		{"Test2", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 51, 50, false},
		{"Test3", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, -10, -10, true},
		{"Test4", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 100, 67, false},
		{"Test5", []int{1, 2, 5, 8, 14, 23, 44, 50, 67}, 8, 8, false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}

			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}
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

func TestTree_Walk(t1 *testing.T) {
	tests := []struct {
		name   string
		values []int
		order  redblack.WalkOrder
		want   [][]int
	}{
		{"IN order", []int{1, 2, 3, 4, 5}, redblack.INORDER, [][]int{{1, 2, 3, 4, 5}}},
		{"PRE order", []int{1, 2, 3, 4, 5}, redblack.PREORDER, [][]int{{3, 2, 1, 5, 4}, {2, 1, 4, 3, 5}, {4, 2, 1, 3, 5}}},
		{"POST order", []int{1, 2, 3, 4, 5}, redblack.POSTORDER, [][]int{{1, 2, 4, 5, 3}, {1, 3, 2, 5, 4}, {1, 3, 5, 4, 2}}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}

			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

			got := make([]int, 0, len(tt.values))
			t.Walk(func(n *redblack.Node[int, redblack.Orderable[int]]) bool {
				if n != nil {
					got = append(got, n.Value())
				}
				return true
			}, tt.order)

			if !slices.ContainsFunc(tt.want, func(e []int) bool { return reflect.DeepEqual(got, e) }) {
				t1.Errorf("Walk() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTree_Sorted(t1 *testing.T) {
	tests := []struct {
		name       string
		values     []int
		breakAfter int
		continueAt int
		want       []int
	}{
		{"Full Loop", []int{1, 2, 3, 4, 5, 6, 7, 8}, -1, -1, []int{1, 2, 3, 4, 5, 6, 7, 8}},
		{"Break after five", []int{1, 2, 3, 4, 5, 6, 7, 8}, 5, -1, []int{1, 2, 3, 4, 5}},
		{"Continue at five", []int{1, 2, 3, 4, 5, 6, 7, 8}, -1, 5, []int{1, 2, 3, 4, 6, 7, 8}},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			rand.Shuffle(len(tt.values), func(i, j int) {
				tt.values[i], tt.values[j] = tt.values[j], tt.values[i]
			})
			vals := make([]redblack.Orderable[int], 0, len(tt.values))
			for _, v := range tt.values {
				vals = append(vals, redblack.Ordered(v))
			}

			t, err := redblack.NewTree(vals, false)
			if err != nil {
				t1.Errorf("redblack.NewTree() error = %v", err)
				return
			}

			got := make([]int, 0, len(tt.values))
			for v := range t.Sorted() {
				if v == tt.continueAt {
					continue
				}
				got = append(got, v)
				if v == tt.breakAfter {
					break
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Sorted() got = %v, want %v", got, tt.values)
			}
		})
	}
}
