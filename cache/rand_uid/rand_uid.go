package rand_uid

import (
	"github.com/hero1s/gotools/algo"
	"github.com/hero1s/gotools/cache"
	"github.com/hero1s/gotools/log"
	"gopkg.in/fatih/set.v0"
)

var randUidKey = "rand_uid_key"

//初始化
func InitRandUid(key string){
	randUidKey = key
}

//获取一个uid
func PopUid() (int64, error) {
	res := cache.Redis.SPop(randUidKey)
	uid, err := res.Int64()
	if err != nil {
		log.Error("获取uid错误:%v", err)
	}
	reslen := cache.Redis.SCard(randUidKey)
	log.Info("获取用户ID:%v,剩余ID数量:%v",uid,reslen.Val())
	return uid, err
}

//重新生成uid
func ResetNewUid(startId, endId, num int64) {
	a := set.New(set.ThreadSafe)
	for i := int64(0); i < num; i = i + 1 {
		id := algo.Random(startId, endId)
		a.Add(id)
	}
	a.Each(func(id interface{}) bool {
		cache.Redis.SAdd(randUidKey, id)
		return true
	})
	res := cache.Redis.SCard(randUidKey)
	log.Info("重新生成用户ID:%v-%v:生成数量:%v,剩余数量:%v",startId,endId,num,res.Val())
}