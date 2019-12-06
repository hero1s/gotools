package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func Md5Sum(plaintext string) string {
	h := md5.New()
	h.Write([]byte(plaintext))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由 addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

// LocalIP 获取机器的IP
func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}

//RandomStr 获取一个随机字符串
func RandomStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func MapStringToStruct(m map[string]interface{}, i interface{}) error {
	bin, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, &i)
	if err != nil {
		return err
	}
	return nil
}

func StructToMapString(i interface{}, m map[string]interface{}) error {
	bin, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bin, &m)
	if err != nil {
		return err
	}
	return nil
}

// 如果有相同的key,会被覆盖
func MergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

func ValueFindinInt(value int64, values ...int64) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

func ValueFindinString(value string, values ...string) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

func ValueFindinUint(value uint64, values ...uint64) bool {
	for i := 0; i < len(values); i++ {
		if value == values[i] {
			return true
		}
	}
	return false
}

//移除切片元素
func RemoveElementInt64(nums []int64, val int64) int {
	if len(nums) == 0 {
		return 0
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return len(nums)
}
func RemoveElementUint64(nums []uint64, val uint64) int {
	if len(nums) == 0 {
		return 0
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return len(nums)
}