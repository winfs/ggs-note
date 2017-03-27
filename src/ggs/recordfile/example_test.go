package recordfile_test

import (
	"fmt"
	"ggs/recordfile"
)

func Example() {
	type Record struct {
		IndexInt int    "index"
		IndexStr string "index"
		_Number  int32
		Str      string
		Arr1     [2]int
		Arr2     [3][2]int
		Arr3     []int
		St       struct {
			Name string "name"
			Num  int    "num"
		}
		M map[string]int
	}

	rf, err := recordfile.New(Record{})
	if err != nil {
		return
	}

	err = rf.Read("test.txt")
	if err != nil {
		return
	}

	for i, n := 0, rf.NumRecord(); i < n; i++ {
		r := rf.Record(i).(*Record) // 返回第i条记录
		fmt.Println(r.IndexInt)     // 1 2 3
	}

	r := rf.Index(2).(*Record) // 获取索引为2的那条记录
	fmt.Println(r.Str)         // cat

	r = rf.Indexes(0)[2].(*Record) // 获取第0列(即数字索引列)索引为2的那条记录
	fmt.Println(r.Str)             // cat

	r = rf.Indexes(1)["three"].(*Record) // 获取第1列(即字符串索引列)索引为"three"的那条记录
	fmt.Println(r.Str)                   // book
	fmt.Println(r.Arr1[1])               // 6
	fmt.Println(r.Arr2[2][0])            // 4
	fmt.Println(r.Arr3[0])               // 6
	fmt.Println(r.St.Name)               // name5566
	fmt.Println(r.M["key6"])             // 6

	// Output:
	// 1
	// 2
	// 3
	// cat
	// cat
	// book
	// 6
	// 4
	// 6
	// name5566
	// 6
}
