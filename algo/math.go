package algo

import (
	crand "crypto/rand"
	"math/big"
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
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := crand.Int(crand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
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
