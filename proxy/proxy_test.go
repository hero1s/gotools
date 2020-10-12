package proxy

import (
	"log"
	"testing"
)

func TestService_Proxy(t *testing.T) {
	// 解开注释测试

	// https://api.igoogle.ink/app/v1/ping
	// test path  /app/v1/ping

	c := &Config{
		ProxySchema: SchemaHTTPS,
		ProxyHost: map[string]string{
			"/baidu": "www.baidu.com",
			"/ali":   "www.aliyun.com",
		},
		ServerPort: ":2233",
		Key:        "123",
	}

	handler := NewHandler(c)

	if err := handler.ListenAndServe(); err != nil {
		log.Fatal("Proxy Start Err：", err)
	}
}
