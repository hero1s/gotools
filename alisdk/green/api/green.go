package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sluu99/uuid"
)

const (
	ActionImageDetection string = "ImageDetection"
	DefaultUrl           string = "http://green.cn-hangzhou.aliyuncs.com"
)

const (
	CodeSuccess string = "Success"
)

type ImageDetectionResult struct {
	ImageName string `json:"ImageName"`
	TaskId    string `json:"TaskId"`
}

type ImageDetectionResults struct {
	Items []ImageDetectionResult `json:"ImageResult"`
}

type ImageDetectionRsp struct {
	Code    string                `json:"Code"`
	Msg     string                `json:"Msg"`
	Results ImageDetectionResults `json:"ImageResults"`
}

type Config struct {
	Url         string
	RegionId    string
	AccessKeyId string
	SecretKey   string
	BodyLimit   int64
}

type Client struct {
	cfg *Config
}

func New(cfg Config) (cli *Client) {
	if cfg.Url == "" {
		cfg.Url = DefaultUrl
	}
	if cfg.BodyLimit <= 0 {
		cfg.BodyLimit = 10 * 1024 * 1024
	}
	return &Client{
		cfg: &cfg,
	}
}

func (cli *Client) fillFormData(imageUrls []string, action string, data url.Values) {

	ts := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	uuidx := uuid.Rand().Hex()
	data.Set("Format", "json")
	data.Set("Version", "2016-10-18")
	data.Set("AccessKeyId", cli.cfg.AccessKeyId)
	data.Set("SignatureMethod", "HMAC-SHA1")
	data.Set("Timestamp", ts)
	data.Set("SignatureVersion", "1.0")
	data.Set("SignatureNonce", uuidx)
	data.Set("Action", action)
	data.Set("RegionId", cli.cfg.RegionId)
	data.Set("Async", "true")
	for k, v := range imageUrls {
		data.Set("ImageUrl."+strconv.Itoa(k+1), v)
	}
	data.Set("Scene.1", "porn")
	data.Set("NotifyUrl", "input your notify url")
	data.Set("Signature", cli.sign(data))

}

func (cli *Client) ImageDetection(imageUrls []string) (*ImageDetectionRsp, error) {

	data := make(url.Values)

	cli.fillFormData(imageUrls, ActionImageDetection, data)
	htprsp, err := http.PostForm(cli.cfg.Url, data)
	if err != nil {
		return nil, err
	}
	defer htprsp.Body.Close()
	bs, err := ioutil.ReadAll(io.LimitReader(htprsp.Body, cli.cfg.BodyLimit))
	if err != nil {
		return nil, err
	}
	var rsp ImageDetectionRsp
	err = json.Unmarshal(bs, &rsp)
	return &rsp, err
}

func (cli *Client) sign(data url.Values) string {
	canonicalizedQueryString := percentReplace(data.Encode())
	stringToSign := "POST" + "&%2F&" + url.QueryEscape(canonicalizedQueryString)
	hmacSha1 := hmac.New(sha1.New, []byte(cli.cfg.SecretKey+"&"))
	hmacSha1.Write([]byte(stringToSign))
	sign := hmacSha1.Sum(nil)
	return base64.StdEncoding.EncodeToString(sign)
}

func percentReplace(str string) string {
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "*", "%2A", -1)
	str = strings.Replace(str, "%7E", "~", -1)

	return str
}
