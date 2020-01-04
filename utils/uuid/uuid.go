package uuid

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hero1s/gotools/log"
	"io"
	"sync"
	"sync/atomic"
)

// 思想就是低位和高位都是跳动变化，这样达到不是一个完全顺序的数
// 低位满了之后，高位重新生成一位,如果高位也满了，则再向高位生成一位，直到64位都用完
//  ---高32位---|---中20位---|---低12位---
const (
	// RenewInterval indicates how often renew retries are performed
	// 低12位
	Renew12Interval uint64 = 0x3FF
	// 低32位
	Renew32Interval uint64 = 0xFFFFFFFF
	// 64位都用完了
	PanicValue uint64 = 0xFFFFFFFFFFFFFFFF // 16个F
)

// UUID is for internal use only.
type UUID struct {
	sync.Mutex
	N           uint64
	Tag         string
	Renew32     func() error
	Renew20     func() error
	H20Verifier func(h20 uint64) error // 高20位
	H32Verifier func(h32 uint64) error // 更高位32位
}

type Option func(*UUID)

func NewUUID(tag string, opts ...Option) *UUID {
	u := &UUID{Tag: tag}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

// Next is for internal use only.
func (u *UUID) Next() uint64 {
	x := atomic.AddUint64(&u.N, 1)
	// 值已经大于64位的最大值了
	if x >= PanicValue {
		log.Error("<uuid> 已经达到极限值, tag: %s", u.Tag)
		panic(fmt.Errorf("<uuid> 已经达到极限值, tag: %s", u.Tag))
	}
	// 低12位已经满了
	if x&Renew12Interval == Renew12Interval {
		err := u.Renew20Now()
		if err != nil {
			log.Error("<uuid> renew 20 failed. tag: %s, reason: %+v", u.Tag, err)
		} else {
			log.Debug("<uuid> renew 20 succeeded. tag: %s", u.Tag)
		}
	} else if x&Renew32Interval == Renew32Interval { // 中20位已经满了
		err := u.Renew32Now()
		if err != nil {
			log.Error("<uuid> renew 32 failed. tag: %s, reason: %+v", u.Tag, err)
		} else {
			log.Debug("<uuid> renew 32 succeeded. tag: %s", u.Tag)
		}
	}
	return x
}

// RenewNow reacquires the high 32 bits from your data store immediately
func (u *UUID) Renew32Now() error {
	u.Lock()
	renew := u.Renew32
	u.Unlock()

	return renew()
}

// RenewNow reacquires the high 20 bits from your data store immediately
func (u *UUID) Renew20Now() error {
	u.Lock()
	renew := u.Renew20
	u.Unlock()

	return renew()
}

// Reset is for internal use only.
func (u *UUID) Reset(n uint64) {
	atomic.StoreUint64(&u.N, n)
}

// VerifyH20 is for internal use only.
func (u *UUID) VerifyH20(h20 uint64) error {
	if h20 == 0 {
		return errors.New("the high 20 bits should not be 0. tag: " + u.Tag)
	}

	if h20 > 0xFFFFF {
		return errors.New("the high 20 bits should not exceed 0x0FFFFFFF. tag: " + u.Tag)
	}

	if u.H20Verifier != nil {
		if err := u.H20Verifier(h20); err != nil {
			return err
		}
	}
	return nil
}

// 校验更高位
// VerifyH32 is for internal use only.
func (u *UUID) VerifyH32(h32 uint64) error {
	// 高32位可以是为0
	if h32 > 0xFFFFFFFF {
		return errors.New("the higher 32 bits should not exceed 0xFFFFFFFF. tag: " + u.Tag)
	}

	if u.H32Verifier != nil {
		if err := u.H32Verifier(h32); err != nil {
			return err
		}
	}
	return nil
}

type NewClient func() (client redis.Cmdable, autoDisconnect bool, err error)

func (u *UUID) LoadH24FromRedis(newClient NewClient, key string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty. tag: " + u.Tag)
	}

	client, autoDisconnect, err := newClient()
	if err != nil {
		return err
	}
	if autoDisconnect {
		defer func() {
			closer := client.(io.Closer)
			_ = closer.Close()
		}()
	}

	n, err := client.Incr(key).Result()
	if err != nil {
		return err
	}
	h20 := uint64(n)
	if err = u.VerifyH20(h20); err != nil {
		return err
	}

	u.Reset(h20 << 12)
	log.Info("<uuid> new h20: %d. tag: %s\n", h20, u.Tag)

	u.Lock()
	defer u.Unlock()

	if u.Renew20 != nil {
		return nil
	}
	u.Renew20 = func() error {
		return u.LoadH24FromRedis(newClient, key)
	}
	return nil
}

func (u *UUID) LoadH32FromRedis(newClient NewClient, key string) error {
	if len(key) == 0 {
		return errors.New("key cannot be empty. tag: " + u.Tag)
	}

	client, autoDisconnect, err := newClient()
	if err != nil {
		return err
	}
	if autoDisconnect {
		defer func() {
			closer := client.(io.Closer)
			_ = closer.Close()
		}()
	}

	n, err := client.Incr(key).Result()
	if err != nil {
		return err
	}
	h32 := uint64(n)
	if err = u.VerifyH32(h32); err != nil {
		return err
	}

	u.Reset(h32 << 32)
	log.Info("<uuid> new h32: %d. tag: %s\n", h32, u.Tag)

	u.Lock()
	defer u.Unlock()

	if u.Renew32 != nil {
		return nil
	}
	u.Renew32 = func() error {
		return u.LoadH32FromRedis(newClient, key)
	}
	return nil
}
