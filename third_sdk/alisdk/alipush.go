package alisdk

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/denverdino/aliyungo/push"
)

var (
	accessKeyId_push     = ""
	accessKeySecret_push = ""
	pushAppKey           int64
)

func InitAliPush(accessKey, accessKeySecret string, appKey int64) {
	accessKeyId_push = accessKey
	accessKeySecret_push = accessKeySecret
	pushAppKey = appKey
}

// 推送消息
func PushMsg(title, body, target, targetValue string) {
	args := push.PushArgs{}
	args.AppKey = pushAppKey
	args.Target = target
	args.TargetValue = targetValue
	args.DeviceType = push.PushDeviceTypeAll
	args.PushType = push.PushTypeMessage
	args.Title = title
	args.Body = body
	args.Summary = body
	clien := push.NewClient(accessKeyId_push, accessKeySecret_push)
	rep, err := clien.Push(&args)
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Info(fmt.Sprintf("push msg %s", rep.MessageId))
}

// 推送通知
func PushNotice(title, body, target, targetValue string) {
	args := push.PushArgs{}
	args.AppKey = pushAppKey
	args.Target = target
	args.TargetValue = targetValue
	args.DeviceType = push.PushDeviceTypeAll
	args.PushType = push.PushTypeNotice
	args.Title = title
	args.Body = body
	args.Summary = body
	clien := push.NewClient(accessKeyId_push, accessKeySecret_push)
	rep, err := clien.Push(&args)
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Info(fmt.Sprintf("push notice %s", rep.MessageId))
}
