package proxy

import (
	"github.com/hero1s/gotools/log"
	"testing"
)

func TestService_Proxy(t *testing.T) {
	cfg, err := LoadConfig("./config.json")
	if err != nil {
		log.Alert("reload config Err：", err)
	} else {
		if err := ListenAndServe(&cfg); err != nil {
			log.Alert("Proxy Start Err：", err)
		}
	}
}
