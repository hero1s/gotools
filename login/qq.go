package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hero1s/gotools/login/fetch"
)

var AndroidQQAppId = "101814185"
var IosQQAppId = "101814185"

type (
	// QQ qq
	QQ struct {
		Ret             int    `json:"ret"`            // 返回码
		Msg             string `json:"msg"`            // 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
		NickName        string `json:"nickname"`       // 用户在QQ空间的昵称。
		FigureURL       string `json:"figureurl"`      // 大小为30×30像素的QQ空间头像URL。
		FigureURL1      string `json:"figureurl_1"`    // 大小为50×50像素的QQ空间头像URL。
		FigureURL2      string `json:"figureurl_2"`    // 大小为100×100像素的QQ空间头像URL。
		FigureURLQQ1    string `json:"figureurl_qq_1"` // 大小为40×40像素的QQ头像URL。
		FigureURLQQ2    string `json:"figureurl_qq_2"` // 大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
		Gender          string `json:"gender"`         // 性别。 如果获取不到则默认返回"男"
		Sex             int64  // 值为1时是男性，值为2时是女性，值为0时是未知
		ISYellowVip     string `json:"is_yellow_vip"`      // 标识用户是否为黄钻用户（0：不是；1：是）。
		Vip             string `json:"vip"`                // 标识用户是否为黄钻用户（0：不是；1：是）
		YellowVipLevel  string `json:"yellow_vip_level"`   // 黄钻等级
		Level           string `json:"level"`              // 黄钻等级
		IsYellowYearVip string `json:"is_yellow_year_vip"` // 标识是否为年费黄钻用户（0：不是； 1：是）
	}
)

// User user
func GetQQUserInfo(accessToken, openID string, typ int64) (*QQ, error) {
	var appid string
	if typ == 1 {
		appid = AndroidQQAppId
	} else {
		appid = IosQQAppId
	}
	var result QQ
	url := fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s",
		accessToken,
		appid,
		openID,
	)
	body, err := fetch.Cmd(fetch.Request{
		Method: "GET",
		URL:    url,
	})
	if err != nil {
		return &result, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return &result, err
	}
	if result.Gender == "男" {
		result.Sex = 1
	} else {
		result.Sex = 2
	}
	fmt.Printf("QQ返回值:%#v\n", result)
	if result.Ret < 0 {
		return &result, errors.New(result.Msg)
	}
	return &result, err
}
