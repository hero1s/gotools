package cache

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/cache"
	"github.com/go-redis/redis"
	"github.com/hero1s/gotools/log"
	"time"
)

var (
	MemCache   cache.Cache
	RedisCache cache.Cache
	Redis      *redis.Client
)

func InitRedis(host, password string) bool {
	Redis = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	pong, err := Redis.Ping().Result()
	if err != nil {
		log.Error(err.Error())
		return false
	}
	log.Info("redis ping rep:%v", pong)
	return true
}

//发布消息
func PublishMessage(channel string, data interface{}) {
	Redis.Publish(channel, data)
}

//接受消息
func SubscribeMessage(channel string, msg_func func(msg *redis.Message)) {
	pubsub := Redis.Subscribe(channel)
	_, err := pubsub.Receive()
	if err != nil {
		return
	}
	ch := pubsub.Channel()
	for msg := range ch {
		log.Info("接受到消息:%v-->%v", msg.Channel, msg.Payload)
		msg_func(msg)
	}
}

func InitCache(host, password, defaultKey string) error {
	var err error
	MemCache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		return err
	}
	RedisCache, err = cache.NewCache("redis",
		`{"conn":"`+host+`", "password":"`+password+`", "key":"`+defaultKey+`"}`)
	return err
}

func SetCache(cc cache.Cache, key string, value interface{}, timeout time.Duration) error {
	data, err := EncodeJson(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("set cache error caught: %v\n", r)
			cc = nil
		}
	}()
	return cc.Put(key, data, timeout)
}

func GetCache(cc cache.Cache, key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		return errors.New("Cache不存在")
	}
	// log.Pinkln(data)
	return DecodeJson(data.([]byte), to)

}

func DelCache(cc cache.Cache, key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()

	return cc.Delete(key)
}

func IsExist(cc cache.Cache, key string) bool {
	if cc == nil {
		return false
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	return cc.IsExist(key)
}

// increase cached int value by key, as a counter.
func Incr(cc cache.Cache, key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	return cc.Incr(key)
}

// decrease cached int value by key, as a counter.
func Decr(cc cache.Cache, key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	return cc.Decr(key)
}

// clear all cache.
func ClearAll(cc cache.Cache, ) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	return cc.ClearAll()
}

// 用json进行数据编码
//
func EncodeJson(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// -------------------
// Decode
// 用json进行数据解码
//
func DecodeJson(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}
