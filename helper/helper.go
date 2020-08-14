package helper

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method string
	URL    string
	Body   string
	Header http.Header
}

func HttpsRequest(args Request) ([]byte, error) {
	client := &http.Client{Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			deadline := time.Now().Add(20 * time.Second)
			c, err := net.DialTimeout(network, addr, 18*time.Second)
			if err != nil {
				return nil, err
			}
			c.SetDeadline(deadline)
			return c, nil
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	}
	req, err := http.NewRequest(args.Method, args.URL, strings.NewReader(args.Body))
	if err != nil {
		return nil, nil
	}
	req.Close = true
	req.Header = args.Header

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//MD5WithRsa
func Md5WithRsa(params string, privateKey []byte) (string, error) {
	data := []byte(params)
	hashMd5 := md5.Sum(data)
	hashed := hashMd5[:]

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Parse private key  error:%v", err))
	}
	p := priv.(*rsa.PrivateKey)
	signature, err := rsa.SignPKCS1v15(rand.Reader, p, crypto.MD5, hashed)

	return base64.StdEncoding.EncodeToString(signature), err
}

//发送钉钉消息
func SendMsgToDingDing(msg string, tokenUrl string) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json;charset=utf-8")
	type Msg struct {
		Content string `json:"content"`
	}
	var m Msg
	m.Content = msg
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    m,
	}
	buff, _ := json.Marshal(data)
	HttpsRequest(Request{
		Method: "POST",
		URL:    tokenUrl,
		Header: headers,
		Body:   string(buff),
	})
}
