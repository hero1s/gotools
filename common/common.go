package common

import (
	"encoding/json"
	"github.com/hero1s/gotools/log"
	"math/rand"
	"reflect"
	"strconv"
	"time"

)

func RandomString(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	var result string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64
	for i = 0; i < length; i++ {
		result = result + string(str[r.Intn(len(str))])
	}
	return result
}

func RandomNum() int64 {
	num1 := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(10000)*5))
	return num1[r.Intn(len(num1))]
}

//根据时间戳,获取星座
func Constellation(tt int64) string {
	t := time.Unix(tt, 0).Format("0102")
	d, _ := strconv.ParseInt(t, 10, 64)
	if d >= 321 && d <= 419 {
		return "白羊座"
	}
	if d >= 420 && d <= 520 {
		return "金牛座"
	}
	if d >= 521 && d <= 621 {
		return "双子座"
	}
	if d >= 622 && d <= 722 {
		return "巨蟹座"
	}
	if d >= 723 && d <= 822 {
		return "狮子座"
	}
	if d >= 823 && d <= 922 {
		return "处女座"
	}
	if d >= 923 && d <= 1023 {
		return "天秤座"
	}
	if d >= 1024 && d <= 1122 {
		return "天蝎座"
	}

	if d >= 1123 && d <= 1221 {
		return "射手座"
	}
	if d >= 1222 || d <= 119 {
		return "魔羯座"
	}
	if d >= 120 && d <= 218 {
		return "水平座"
	}
	if d >= 219 && d <= 320 {
		return "双鱼座"
	}

	return "水平座"
}

// 过来结构体空字段，转换json字段的map
func ChangeStructPointToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if !f.IsNil() {
			data[t.Field(i).Tag.Get("json")] = f.Interface()
		}
	}
	return data
}

func ChangeStructToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		data[t.Field(i).Tag.Get("json")] = f.Interface()
	}
	return data
}

//json数据转换
func ChangeJsonStruct(from, to interface{}) error {
	str, err := json.Marshal(from)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = json.Unmarshal(str, &to)
	if err != nil {
		return err
	}
	return nil
}