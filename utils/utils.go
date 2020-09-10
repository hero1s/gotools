package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/hero1s/gotools/log"
	"math/rand"
	"net"
	"runtime"
	"runtime/debug"
	"strings"
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

//生成随机字符串
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
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

// parse params(name=nick&pass=123)
func ParseUrlString(params string) map[string]string {
	paramsMap := map[string]string{}
	for _, param := range strings.Split(params, "&") {
		if ! strings.Contains(param, "=") {
			continue
		}
		paramList := strings.Split(param, "=")
		paramsMap[paramList[0]] = paramList[1]
	}
	return paramsMap
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
func RemoveElementInt64(nums []int64, val int64) []int64 {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}
func RemoveElementUint64(nums []uint64, val uint64) []uint64 {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}
func RemoveElementString(nums []string, val string) []string {
	if len(nums) == 0 {
		return nums
	}
	index := 0
	for ; index < len(nums); {
		if nums[index] == val {
			nums = append(nums[:index], nums[index+1:]...)
			continue
		}
		index++
	}
	return nums
}

//slice去重
func RemoveRepeatedElement(slc []uint64) []uint64 {
	result := []uint64{}         //存放返回的不重复切片
	tempMap := map[uint64]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0 //当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		//如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e) //当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

//安全执行异步函数
func SafeGoroutine(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("child goroutine panic occure,err:%v", r)
				log.Error("stack:%v", debug.Stack())
			}
		}()
		f()
	}()
}

//安全执行函数
func SafeCallFunc(f func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("call func panic occure,err:%v", r)
			log.Error("stack:%v", debug.Stack())
		}
	}()
	f()
}

// 获取正在运行的函数名
func RunFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
func CallerFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}
