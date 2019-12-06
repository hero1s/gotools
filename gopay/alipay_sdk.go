package gopay

import (
	"errors"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/utils"
	"github.com/iGoogle-ink/gopay"
)

//UserCertifyOpenInitialize 身份认证初始化服务 https://docs.open.alipay.com/api_2/alipay.user.certify.open.initialize
func UserCertifyOpenInitialize(realName, cardNo, returnUrl string) (string, error) {
	log.Info("========== UserCertifyOpenInitialize ==========")
	bm := make(gopay.BodyMap)
	bm.Set("outer_order_no", utils.GenStringUUID())
	bm.Set("biz_code", "FACE")
	identity := make(map[string]string)
	identity["identity_type"] = "CERT_INFO"
	identity["cert_type"] = "IDENTITY_CARD"
	identity["cert_name"] = realName
	identity["cert_no"] = cardNo
	bm.Set("identity_param", identity)
	// 商户个性化配置，格式为json
	merchant := make(map[string]string)
	merchant["return_url"] = returnUrl
	bm.Set("merchant_config", merchant)
	//发起请求
	aliRsp, err := AliPayClient.AliPayUserCertifyOpenInit(bm)
	if err != nil {
		log.Error("身份验证初始化失败:%v", err)
		return "", err
	}
	log.Info("初始化身份验证返回:%#v", *aliRsp)
	if aliRsp.Response.Code != "10000" {
		log.Info("初始化认证不成功:%v,%v", aliRsp.Response.Msg, aliRsp.Response.SubMsg)
		return "", errors.New("认证初始化不成功")
	}
	log.Info(aliRsp.Response.CertifyId)
	return aliRsp.Response.CertifyId, nil
}

//UserCertifyOpenCertify 身份认证开始认证 https://docs.open.alipay.com/api_2/alipay.user.certify.open.certify
func UserCertifyOpenCertify(certifyId string) (string, error) {
	log.Info("========== UserCertifyOpenCertify ==========")
	bm := make(gopay.BodyMap)
	// 本次申请操作的唯一标识，由开放认证初始化接口调用后生成，后续的操作都需要用到
	bm.Set("certify_id", certifyId)

	rsp, err := AliPayClient.AliPayUserCertifyOpenCertify(bm)
	if err != nil {
		log.Error("%#v 身份认证错误:%v", bm, err)
	}
	log.Info("身份认证返回值:%#v", rsp)
	return rsp, err
}

// UserCertifyOpenQuery 身份认证记录查询 https://docs.open.alipay.com/api_2/alipay.user.certify.open.query/
func UserCertifyOpenQuery(certifyId string) (rsp *gopay.AliPayUserCertifyOpenQueryResponse, err error) {
	log.Info("========== UserCertifyOpenQuery ==========")
	//请求参数
	bm := make(gopay.BodyMap)
	// 本次申请操作的唯一标识，由开放认证初始化接口调用后生成，后续的操作都需要用到
	bm.Set("certify_id", certifyId)
	rsp, err = AliPayClient.AliPayUserCertifyOpenQuery(bm)
	if err != nil {
		log.Error("身份记录查询错误:%v", err)
		return
	}
	log.Info("身份认证记录查询:%#v", rsp)
	return
}
