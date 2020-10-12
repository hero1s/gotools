package xhttp

import (
	"github.com/hero1s/gotools/log"
	"testing"
	"time"

)

type HttpGet struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func TestHttpGet(t *testing.T) {
	client := NewClient()
	client.Timeout = 10 * time.Second

	rsp := new(HttpGet)
	_, errs := client.Type(TypeJSON).Get("https://api.igoogle.ink/app/v1/ping").EndStruct(rsp)
	if len(errs) > 0 {
		log.Error(errs[0].Error())
		return
	}
	log.Debug("%v",rsp)

	// test
	_, bs, errs := client.Get("http://www.baidu.com").EndBytes()
	if len(errs) > 0 {
		log.Error(errs[0].Error())
		return
	}
	log.Debug(string(bs))
}
