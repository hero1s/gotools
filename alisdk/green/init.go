package green

import "github.com/hero1s/gotools/alisdk/green/api"

var GreenClient *api.Client

func InitGreen(accessKeyId,secretKey,regionId string){
	cfg := api.Config{
		Url:         "",
		RegionId:    regionId,
		AccessKeyId: accessKeyId,
		SecretKey:   secretKey,
		BodyLimit:   0,
	}
	GreenClient = api.New(cfg)
}


