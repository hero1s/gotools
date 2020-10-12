package retry

import (
	"errors"
	"github.com/siddontang/go/log"
	"testing"
	"time"

	"github.com/iGoogle-ink/gotil/xlog"
)

func TestRetry(t *testing.T) {
	err := Retry(func() error {
		log.Debugf("retry func")
		return errors.New("please retry")
	}, 3, 2*time.Second)
	if err != nil {
		xlog.Error(err)
	}
}
