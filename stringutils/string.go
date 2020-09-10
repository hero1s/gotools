package stringutils

import (
	"encoding/json"
	"github.com/hero1s/gotools/log"
	"github.com/zheng-ji/goSnowFlake"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

const (
	regular = "^1[3|4|5|7|8][0-9]{9}$"
)

/*
func AjaxMsg(code int, message string) {
	this.Data["json"] = map[string]interface{}{"code": code, "message": message}
	this.ServeJSON()
	return
}
*/
//字串截取
func SubString(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetFileSuffix(s string) string {
	re, _ := regexp.Compile(".(jpg|jpeg|png|gif|exe|doc|docx|ppt|pptx|xls|xlsx)")
	suffix := re.ReplaceAllString(s, "")
	return suffix
}

func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func Unicode(rs string) string {
	json := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

/**
 *  to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail
 */
func SendMail(to, subject, body, user, password, host string) error {

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	content_type = "Content-type:text/html;charset=utf-8"

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func SnowFlakeId() int64 {
	iw, _ := goSnowFlake.NewIdWorker(1)
	if id, err := iw.NextId(); err != nil {
		return 0
	} else {
		return id
	}
}

//数组去重来源网络
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 1 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

//判断是否手机号
func CheckPhoneNum(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//json数据转换
func ChangeJsonStruct(from, to interface{}) error {
	str, err := json.Marshal(from)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = json.Unmarshal(str, &to)
	if err != nil {
		return err
	}
	return nil
}

