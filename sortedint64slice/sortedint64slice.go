package sortedint64slice

import (
	"math/rand"
	"time"
)

//type SortedInt64Slice struct {
//	data []int64
//}

type SortedInt64Slice []int64

func (s SortedInt64Slice) QuickSort()  {
	QuickSort(s,len(s),0,len(s) - 1)
}

// GenerateRangeNum 生成一个区间范围的随机数
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	return randNum
}

func Swap(a *int64,b *int64)  {
	*a,*b = *b,*a
}

func Partition(data []int64,length int,start int,end int) int  {
	if nil == data || 0 >= length || 0 > start || end >= length {
		panic("err")
	}

	index := GenerateRangeNum(start,end)
	Swap(&data[index],&data[end])

	small := start - 1

	for index = start; index < end; index++ {
		if data[index] < data[end] {
			small ++
			if small != index {
				Swap(&data[index],&data[small])
			}
		}
	}

	small ++

	Swap(&data[small],&data[end])


	return small
}

func QuickSort(data []int64,length int,start int,end int)  {
	if start == end {
		return
	}

	index := Partition(data,length,start,end)

	if index > start {
		QuickSort(data,length,start,index -1)
	}

	if index < end{
		QuickSort(data,length,index + 1,end)
	}

}

