package helper

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/hero1s/gotools/cache"
	"github.com/hero1s/gotools/log"
	"time"
)

func TryLock(lockName string, ttl,waitTime time.Duration) (*redislock.Lock, error) {
	locker := redislock.New(cache.Redis)
	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()
	opt := redislock.Options{
		RetryStrategy: redislock.ExponentialBackoff(100*time.Millisecond, waitTime),
		Context:       ctx,
	}
	lock, err := locker.Obtain(lockName, ttl, &opt)
	if err == redislock.ErrNotObtained {
		log.Error("Could not obtain lock! lockerName:%v",lockName)
	} else if err != nil {
		log.Error(err.Error())
	}
	return lock, err
}
