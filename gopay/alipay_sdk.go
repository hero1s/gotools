package gopay

import (
	"encoding/json"
	"errors"
	"github.com/hero1s/gotools/common"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/utils"
	"github.com/iGoogle-ink/gopay"
)

//UserCertifyOpenInitialize 身份认证初始化服务 https://docs.open.alipay.com/api_2/alipay.user.certify.open.initialize
func UserCertifyOpenInitialize(realName, identity, returnUrl string) (string, error) {
	log.Debug("========== UserCertifyOpenInitialize ==========")
	var p = gopay.UserCertifyOpenInitialize{}
	p.OuterOrderNo = utils.GenStringUUID()
	p.BizCode = gopay.K_CERTIFY_BIZ_CODE_FACE
	p.IdentityParam.IdentityType = "CERT_INFO"
	p.IdentityParam.CertType = "IDENTITY_CARD"
	p.IdentityParam.CertName = realName
	p.IdentityParam.CertNo = identity
	p.MerchantConfig.ReturnURL = returnUrl

	pxx, _ := json.Marshal(p)
	log.Debug("身份初始化验证参数:%v", string(pxx))

	rsp, err := AliPayClient.UserCertifyOpenInitialize(common.ChangeStructToJsonMap(p))
	if err != nil {
		log.Error("身份验证初始化失败:%v", err)
		return "", err
	}
	log.Debug("初始化身份验证返回:%#v", rsp)
	if rsp.Content.Code != "10000" {
		log.Error("初始化认证不成功:%v,%v", rsp.Content.Msg, rsp.Content.SubMsg)
		return "", errors.New("认证初始化不成功")
	}
	log.Debug(rsp.Content.CertifyId)
	return rsp.Content.CertifyId, nil
}

//UserCertifyOpenCertify 身份认证开始认证 https://docs.open.alipay.com/api_2/alipay.user.certify.open.certify
func UserCertifyOpenCertify(certifyId string) (string, error) {
	log.Debug("========== UserCertifyOpenCertify ==========")
	body := map[string]interface{}{
		"certify_id": certifyId,
	}
	rsp, err := AliPayClient.UserCertifyOpenCertify(body)
	if err != nil {
		log.Error("%#v 身份认证错误:%v", body, err)
	}
	log.Debug("身份认证返回值:%#v", string(rsp))
	return string(rsp), err
}

// UserCertifyOpenQuery 身份认证记录查询 https://docs.open.alipay.com/api_2/alipay.user.certify.open.query/
func UserCertifyOpenQuery(certifyId string) (rsp gopay.UserCertifyOpenQueryRsp, err error) {
	log.Debug("========== UserCertifyOpenQuery ==========")
	var p = gopay.UserCertifyOpenQuery{}
	p.CertifyId = certifyId
	rsp, err = AliPayClient.UserCertifyOpenQuery(common.ChangeStructToJsonMap(p))
	if err != nil {
		log.Error("身份记录查询错误:%v", err)
		return
	}
	if rsp.Content.Code != "10000" {
		log.Error("验证信息不成功:", rsp.Content.Msg, rsp.Content.SubMsg)
	}
	log.Debug("身份认证记录查询:%#v", rsp.Content)
	return
}
