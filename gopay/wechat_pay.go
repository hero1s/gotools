package gopay

import (
	"github.com/iGoogle-ink/gopay"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/utils"
	"strconv"
	"time"
)

var WeChatPayClient *gopay.WeChatClient

func InitWechatPay(isProd bool) {
	WeChatPayClient = gopay.NewWeChatClient(PayParam.WechatPay.WeChatAppId, PayParam.WechatPay.WeChatMchId, PayParam.WechatPay.WeChatKey, isProd)
	WeChatPayClient.SetCountry(gopay.China)
}

//微信预下单
func UnifiedOrder(moneyFee int64, describe, orderId string) (map[string]string, error) {
	//初始化参数Map
	body := make(gopay.BodyMap)
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("body", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_fee", moneyFee) //单位分
	body.Set("spbill_create_ip", utils.LocalIP())
	body.Set("notify_url", PayParam.WechatPay.WeChatCallbackUrl)
	body.Set("trade_type", gopay.TradeType_App)
	body.Set("device_info", "APP")
	body.Set("sign_type", gopay.SignType_MD5)

	//请求支付下单，成功后得到结果
	var c = make(map[string]string)
	wxRsp, err := WeChatPayClient.UnifiedOrder(body)
	if err != nil {
		log.Error("微信预下单:%#v  \n支付失败Error:%v", body, err.Error())
		return c, err
	} else {
		log.Info("微信预下单wxRsp:%#v", *wxRsp)
	}

	c["appid"] = wxRsp.Appid
	c["partnerid"] = wxRsp.MchId
	c["prepayid"] = wxRsp.PrepayId
	c["package"] = "Sign=WXPay"
	c["noncestr"] = wxRsp.NonceStr
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := gopay.GetAppPaySign(wxRsp.Appid, "", wxRsp.NonceStr, wxRsp.PrepayId, gopay.SignType_MD5, timeStamp, PayParam.WechatPay.WeChatKey)
	c["paySign"] = sign
	c["timestamp"] = timeStamp

	return c, err
}

// 提交付款码支付：client.Micropay()

// 查询订单：client.QueryOrder()

// 关闭订单：client.CloseOrder()

// 撤销订单：client.Reverse()

// 申请退款：client.Refund()
func Refund(orderId string, moneyFee int64) bool {
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderId)
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("sign_type", gopay.SignType_MD5)
	body.Set("out_refund_no", orderId)
	body.Set("total_fee", moneyFee)
	body.Set("refund_fee", moneyFee)

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	//    certFilePath：cert证书路径
	//    keyFilePath：Key证书路径
	//    pkcs12FilePath：p12证书路径
	wxRsp, err := WeChatPayClient.Refund(body, PayParam.WechatPay.WeChatCertFile, PayParam.WechatPay.WeChatKeyFile, PayParam.WechatPay.WeChatP12File)
	if err != nil {
		log.Error("微信退款Error:%v", err)
		return false
	}
	log.Debug("微信退款wxRsp：%#v", *wxRsp)
	if wxRsp.ReturnCode == gopay.SUCCESS {
		log.Debug("微信退款成功:%#v", wxRsp)
		return true
	} else {
		log.Error("微信退款失败:%#v", wxRsp)
	}
	return false
}

// 查询退款：client.QueryRefund()

// 下载对账单：client.DownloadBill()

// 下载资金账单：client.DownloadFundFlow()

// 拉取订单评价数据：client.BatchQueryComment()

// 企业向微信用户个人付款：client.Transfer()
func Transfer(orderId, openid, userName, desc string, moneyFee int64) {
	nonceStr := gopay.GetRandomString(32)
	log.Info("partnerTradeNo:%v", orderId)
	//初始化参数结构体
	body := make(gopay.BodyMap)
	body.Set("nonce_str", nonceStr)
	body.Set("partner_trade_no", orderId)
	body.Set("openid", openid)
	body.Set("check_name", "FORCE_CHECK") // NO_CHECK：不校验真实姓名 , FORCE_CHECK：强校验真实姓名
	body.Set("re_user_name", userName)    //收款用户真实姓名。 如果check_name设置为FORCE_CHECK，则必填用户真实姓名
	body.Set("amount", moneyFee)          //企业付款金额，单位为分
	body.Set("desc", desc)                //企业付款备注，必填。注意：备注中的敏感词会被转成字符*
	body.Set("spbill_create_ip", "127.0.0.1")

	//请求申请退款（沙箱环境下，证书路径参数可传空）
	//    body：参数Body
	//    certFilePath：cert证书路径
	//    keyFilePath：Key证书路径
	//    pkcs12FilePath：p12证书路径
	wxRsp, err := WeChatPayClient.Transfer(body, PayParam.WechatPay.WeChatCertFile, PayParam.WechatPay.WeChatKeyFile, PayParam.WechatPay.WeChatP12File)
	if err != nil {
		log.Error("微信付款Error:", err)
		return
	}
	log.Info("wxRsp：%#v", *wxRsp)
	if wxRsp.ReturnCode == gopay.SUCCESS {
		log.Debug("微信转账成功:%#v", wxRsp)
	} else {
		log.Error("微信转账失败:%#v", wxRsp)
	}

}

//验证微信回调
func VerifyWeChatSign(notifyReq *gopay.WeChatNotifyRequest) (ok bool, err error) {
	//验签操作
	return gopay.VerifyWeChatSign(PayParam.WechatPay.WeChatKey, gopay.SignType_MD5, notifyReq)
}
