package algo

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//随机值在闭区间[min,max]
func Random(min, max int64) int64 {
	max += 1
	return rand.Int63n(max-min) + min
}
func RemoveDuplicateInt(list []int64) []int64 {
	var x []int64
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}
func RemoveDuplicateStr(list []string) []string {
	var x []string = []string{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

