package gopay

import (
	"github.com/hero1s/gotools/log"
	"github.com/iGoogle-ink/gopay"
)

var AliPayClient *gopay.AliPayClient

func InitAliPay(isProd bool) {
	AliPayClient = gopay.NewAliPayClient(PayParam.AliPay.AliAppId, string(keyFromFile(PayParam.AliPay.AliAppPrivateKeyFile)), isProd)

	//设置支付宝请求 公共参数
	AliPayClient.SetCharset("utf-8").
		SetSignType("RSA2"). //设置签名类型，不设置默认 RSA2
		SetReturnUrl(PayParam.AliPay.AliReturnUrl). //设置返回URL
		SetNotifyUrl(PayParam.AliPay.AliAppCallbackUrl) //设置异步通知URL
	//.SetAppAuthToken().SetAuthToken()
}

//* 手机网站支付接口2.0（手机网站支付）：client.AliPayTradeWapPay()
func AliPayTradeWapPay(moneyFee int64, describe, orderId, quitUrl,returnUrl string) (string, error) {
	AliPayClient.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("quit_url", quitUrl)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", "QUICK_WAP_WAY")
	//手机网站支付请求
	payUrl, err := AliPayClient.AliPayTradeWapPay(body)
	if err != nil {
		log.Error("page pay err:", err)
		return payUrl, err
	}
	return payUrl, err
}

//* 统一收单下单并支付页面接口（电脑网站支付）：client.AliPayTradePagePay()
func AliPayTradePagePay(moneyFee int64, describe, orderId string,returnUrl string) (string, error) {
	AliPayClient.SetReturnUrl(returnUrl)
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", "FAST_INSTANT_TRADE_PAY")

	//电脑网站支付请求
	payUrl, err := AliPayClient.AliPayTradePagePay(body)
	if err != nil {
		log.Error("page pay err:", err)
		return payUrl, err
	}
	return payUrl, err
}

//* APP支付接口2.0（APP支付）：client.AliPayTradeAppPay()
func AliPayTradeAppPay(moneyFee int64, describe, orderId string) (string, error) {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", describe)
	body.Set("out_trade_no", orderId)
	body.Set("total_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("product_code", PayParam.AliPay.AliProductCode)
	//手机APP支付参数请求
	payParam, err := AliPayClient.AliPayTradeAppPay(body)
	if err != nil {
		log.Error("阿里支付预下单失败err:", err)
	}
	//log.Info("阿里支付预下单返回payParam:", payParam)
	return payParam, err
}

//* 统一收单交易支付接口（商家扫用户付款码）：client.AliPayTradePay()

//* 统一收单交易创建接口（小程序支付）：client.AliPayTradeCreate()

//* 统一收单线下交易查询：client.AliPayTradeQuery()

//* 统一收单交易关闭接口：client.AliPayTradeClose()

//* 统一收单交易撤销接口：client.AliPayTradeCancel()

//* 统一收单交易退款接口：client.AliPayTradeRefund()
func AliPayTradeRefund(orderId string, moneyFee int64) bool {
	//请求参数
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", orderId)
	body.Set("refund_amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("refund_reason", "测试退款")
	//发起退款请求
	aliRsp, err := AliPayClient.AliPayTradeRefund(body)
	if err != nil {
		log.Error("阿里退款失败err:%v", err)
		return false
	}
	log.Info("阿里退款返回aliRsp:%#v", *aliRsp)
	if aliRsp.AlipayTradeRefundResponse.Code == "10000" {
		log.Debug("阿里退款成功:%#v", aliRsp)
		return true
	} else {
		log.Debug("阿里退款失败:%#v", aliRsp)
	}
	return false
}

//* 统一收单退款页面接口：client.AliPayTradePageRefund()

//* 统一收单交易退款查询：client.AliPayTradeFastPayRefundQuery()

//* 统一收单交易结算接口：client.AliPayTradeOrderSettle()

//* 统一收单线下交易预创建（用户扫商品收款码）：client.AliPayTradePrecreate()

//* 单笔转账到支付宝账户接口（商户给支付宝用户转账）：client.AlipayFundTransToaccountTransfer()
func AlipayFundTransToaccountTransfer(account string, moneyFee int64, desc string) {
	body := make(gopay.BodyMap)
	out_biz_no := gopay.GetRandomString(32)
	body.Set("out_biz_no", out_biz_no)
	body.Set("payee_type", "ALIPAY_LOGONID")
	body.Set("payee_account", account)
	body.Set("amount", AliMoneyFeeToString(float64(moneyFee)))
	body.Set("payer_show_name", "发钱人名字")
	body.Set("payee_real_name", "收钱人名字")
	body.Set("remark", desc)
	//创建订单
	aliRsp, err := AliPayClient.AlipayFundTransToaccountTransfer(body)
	if err != nil {
		log.Error("阿里转账支付错误err:", err)
		return
	}
	log.Debug("阿里转账支付返回aliRsp:%#v", *aliRsp)
	if aliRsp.AlipayFundTransToaccountTransferResponse.Code == "10000" {
		log.Debug("支付宝转账成功:%#v", aliRsp)
	} else {
		log.Error("支付宝转账失败:%#v", aliRsp)
	}
}

//* 换取授权访问令牌（获取access_token，user_id等信息）：client.AliPaySystemOauthToken()

//* 支付宝会员授权信息查询接口（App支付宝登录）：client.AlipayUserInfoShare()

//* 换取应用授权令牌（获取app_auth_token，auth_app_id，user_id等信息）：client.AlipayOpenAuthTokenApp()

//* 获取芝麻信用分：client.ZhimaCreditScoreGet()

//验证支付宝回调
func VerifyAliPaySign(notifyReq *gopay.AliPayNotifyRequest) (ok bool, err error) {
	return gopay.VerifyAliPaySign(string(keyFromFile(PayParam.AliPay.AliPayPublicKeyFile)), notifyReq)
}
