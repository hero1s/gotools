package green_sdk

import (
	"encoding/json"
	"github.com/hero1s/gotools/alisdk/green/greensdksample"
	"github.com/hero1s/gotools/alisdk/green/uuid"
	"github.com/hero1s/gotools/log"
	"path"
	"regexp"
)

var accessKeyId = "<your access key id>"
var accessKeySecret = "<your access key secret>"

func InitGreenSdk(KeyId, KeySecret string) {
	accessKeyId = KeyId
	accessKeySecret = KeySecret
}

type Detail struct {
	Label      string  `json:"label"`
	Rate       float64 `json:"rate"`
	Scene      string  `json:"scene"`
	Suggestion string  `json:"suggestion"`
}

type ScanDataResp struct {
	Code    int64    `json:"code"`
	Content string   `json:"content"`
	Url     string   `json:"url"`
	DataId  string   `json:"dataId"`
	Msg     string   `json:"msg"`
	TaskId  string   `json:"taskId"`
	Results []Detail `json:"results"`
}

type ScanResp struct {
	Code      int64          `json:"code"`
	Msg       string         `json:"msg"`
	RequestId string         `json:"requestId"`
	Data      []ScanDataResp `json:"data"`
}

//鉴定图片
func CheckImageScan(imageUrl string) bool {
	profile := greensdksample.Profile{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}

	path := "/green/image/scan"
	//是否视频
	if CheckFileNameIsVideo(imageUrl) { //视频文件后台自动冻结
		return true
	}
	clientInfo := greensdksample.ClinetInfo{Ip: "127.0.0.1"}

	// 构造请求数据
	bizType := "Green"
	scenes := []string{"porn"}

	task := greensdksample.Task{DataId: uuid.Rand().Hex(), Url: imageUrl}
	tasks := []greensdksample.Task{task}

	bizData := greensdksample.BizData{bizType, scenes, tasks}

	var client greensdksample.IAliYunClient = greensdksample.DefaultClient{Profile: profile}

	// your biz code
	strResp := client.GetResponse(path, clientInfo, bizData)
	log.Debug("色情图片检测:%v", strResp)
	var resp ScanResp
	if err := json.Unmarshal([]byte(strResp), &resp); err != nil {
		log.Error("解析鉴黄结构体错误:%v", err)
	} else {
		log.Debug("鉴黄结果返回:%+v", resp)
		for _, d := range resp.Data {
			for _, res := range d.Results {
				if res.Suggestion == "block" {
					log.Error("鉴定为屏蔽:%v", d.Url)
					return false
				}
				if res.Rate > 80 && res.Label == "porn" {
					log.Error("鉴定涉黄分数:%v", res)
					return false
				}
			}
		}
	}
	return true

}

//鉴定文本
func CheckTextScan(text string) bool {
	profile := greensdksample.Profile{AccessKeyId: accessKeyId, AccessKeySecret: accessKeySecret}

	pattern := "^[A-Za-z0-9]+$"
	ok, _ := regexp.MatchString(pattern, text)
	if ok {
		log.Debug("纯字母数字文本不检测:%v", text)
		return true
	}
	if len(text) < 4 {
		return true
	}

	path := "/green/text/scan";
	clientInfo := greensdksample.ClinetInfo{Ip: "127.0.0.1"}

	// 构造请求数据
	bizType := "Green"
	scenes := []string{"antispam"}

	task := greensdksample.Task{DataId: uuid.Rand().Hex(), Content: text}
	tasks := []greensdksample.Task{task}

	bizData := greensdksample.BizData{bizType, scenes, tasks}

	var client greensdksample.IAliYunClient = greensdksample.DefaultClient{Profile: profile}

	// your biz code
	strResp := client.GetResponse(path, clientInfo, bizData)
	log.Debug("色情文字检测:%v", strResp)
	var resp ScanResp
	if err := json.Unmarshal([]byte(strResp), &resp); err != nil {
		log.Error("解析鉴黄结构体错误:%v", err)
	} else {
		log.Debug("鉴黄结果返回:%+v", resp)
		for _, d := range resp.Data {
			for _, res := range d.Results {
				if res.Suggestion == "block" {
					log.Error("鉴定为屏蔽:%v", d.Content)
					return false
				}
			}
		}
	}
	return true
}

//判断文件是否视频
func CheckFileNameIsVideo(fileName string) bool {
	filenameWithSuffix := path.Base(fileName)  //获取文件名带后缀
	fileSuffix := path.Ext(filenameWithSuffix) //获取文件后缀
	log.Debug("文件后缀名:%v", fileSuffix)
	if fileSuffix == ".mp4" {
		return true
	}

	return false
}
