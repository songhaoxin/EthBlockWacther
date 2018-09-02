package sortedint64slice

import (
	"testing"
	"fmt"
	"github.com/ethereum/go-ethereum/metrics"
	"sort"
)

func TestGenerateRangeNum(t *testing.T)  {
	for i := 0; i < 3; i++ {
		acaut := GenerateRangeNum(0,10)
		fmt.Println(acaut)

	}
}

var slice []int64 = make([]int64,0)

func TestSwap(t *testing.T) {
	slice = append(slice,8)
	slice = append(slice,9)

	for k,v := range slice {
		fmt.Println(k,v)
	}

	Swap(&slice[0],&slice[1])

	for k,v := range slice {
		fmt.Println(k,v)
	}
}

func TestPartition(t *testing.T) {
	tcase := []int64{5,3,9,6,2,18}

	idex := Partition(tcase,len(tcase),0,len(tcase)-1)
	fmt.Println(idex)
}

func TestQuickSort(t *testing.T) {
	slice := []int64{5,3,9,6,2,18,0}

	QuickSort(slice,len(slice),0,len(slice) - 1)

	for k,v := range slice {
		fmt.Println(k,v)
	}
}

func TestSortedInt64Slice_QuickSort(t *testing.T) {
	//ms := metrics.Int64Slice{}
	//sort.Sort(ms)
	int64Slice := metrics.Int64Slice{8,1,89,0,23,12,100}
	sort.Sort(int64Slice)

	for k,v := range int64Slice {
		fmt.Println(k,v)
	}
}

func BenchmarkSortedInt64Slice_QuickSort(b *testing.B) {
	//int64Slice := SortedInt64Slice{8,1,89,0,23,12,100}
	int64Slice := metrics.Int64Slice{8,1,89,0,23,12,100}
	for i := 0; i < 11; i++ {
		int64Slice = append(int64Slice,int64Slice...)
	}

	//int64Slice.QuickSort()
	//QuickSort(int64Slice,len(int64Slice),0,len(int64Slice) -1)
	//sort.Sort(int64Slice)
	sort.Sort(int64Slice)

}

func TestM(t *testing.T)  {
	m := map[int64]string {
		3:"abc",
		2:"bcd",
		1:"cde",
	}

	var keys metrics.Int64Slice
	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	for i, k := range keys {
		fmt.Println(i)
		fmt.Println("Key:", k, "Value:", m[k])
	}


}
