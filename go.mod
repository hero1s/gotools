module github.com/hero1s/gotools

go 1.13

require (
	github.com/aliyun/alibaba-cloud-sdk-go v1.60.317
	github.com/aliyun/aliyun-oss-go-sdk v2.0.4+incompatible
	github.com/astaxie/beego v1.12.1
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/bsm/redislock v0.5.0
	github.com/denverdino/aliyungo v0.0.0-20191128015008-acd8035bbb1d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gansidui/geohash v0.0.0-20141019080235-ebe5ba447f34
	github.com/gansidui/nearest v0.0.0-20141019122829-a5d0cde6ef14
	github.com/garyburd/redigo v1.6.0
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/websocket v1.4.2
	github.com/howeyc/fsnotify v0.9.0
	github.com/iGoogle-ink/gopay v1.3.9
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/shopspring/decimal v0.0.0-20191130220710-360f2bc03045
	github.com/stretchr/testify v1.4.0
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/api v0.13.0
	google.golang.org/grpc v1.31.0
	gopkg.in/fatih/set.v0 v0.2.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	xorm.io/xorm v1.0.3
)

replace github.com/astaxie/beego v1.12.2 => github.com/nicle-lin/beego v1.12.3

replace github.com/iGoogle-ink/gopay v1.3.9 => github.com/hero1s/gopay v1.4.3

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
