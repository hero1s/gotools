package utils

import (
	"github.com/go-redis/redis"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/utils/uuid"
	"strconv"
)

var (
	Uid *uuid.UUID
)

func InitUUID(redisHost, password string) error {
	newClient := func() (redis.Cmdable, bool, error) {
		return redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: password,
		}), true, nil
	}
	Uid = uuid.NewUUID("uid")
	err := Uid.LoadH52FromRedis(newClient, "UUID:UID:24")
	if err != nil {
		log.Error("初始化UUID错误:%v",err.Error())
		return err
	}

	return nil
}

func GenUid() uint64 {
	return Uid.Next()
}

func GenStringUUID() string {
	return strconv.FormatUint(GenUid(), 10)
}
