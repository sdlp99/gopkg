package msg

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/sdlp99/sdpkg/utils/logger"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 发送钉钉文本信息
func SendDingTalkText(accId string, secret string, title string, text string) (string, error) {

	value, _ := sjson.Set("", "msgtype", "markdown")
	value, _ = sjson.Set(value, "markdown.title", title)
	value, _ = sjson.Set(value, "markdown.text", text)

	timeStampNow := time.Now().UnixNano() / 1000000
	signStr := fmt.Sprintf("%d\n%s", timeStampNow, secret)

	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(signStr))
	sum := hash.Sum(nil)

	encode := base64.StdEncoding.EncodeToString(sum)
	urlEncode := url.QueryEscape(encode)

	// 构建 请求 url
	UrlAddress := fmt.Sprintf("%s?access_token=%s&timestamp=%d&sign=%s",
		"https://oapi.dingtalk.com/robot/send", accId, timeStampNow, urlEncode)

	// 构建 请求体
	request, err := http.NewRequest("POST", UrlAddress, strings.NewReader(value))
	if err != nil {
		return "", err
	}
	// 设置库端口
	client := &http.Client{}

	// 请求头添加内容
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	//fmt.Println("response: ", response)

	// 关闭 读取 reader
	defer response.Body.Close()

	// 读取内容
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	return string(all), err
}

// https://api.day.app/eL8azgbJhmbthyNK7VVRrj/
func SendBarkText(url string, secret string, title string, text string) (string, error) {

	value := "" //sjson.Set("", "device_key", secret)
	value, _ = sjson.Set(value, "title", title)
	value, _ = sjson.Set(value, "body", text)

	// 构建 请求体
	request, err := http.NewRequest("POST", url+"/"+secret+"/"+title+"/"+text, strings.NewReader(value))
	if err != nil {
		return "", err
	}
	// 设置库端口
	client := &http.Client{}

	// 请求头添加内容
	request.Header.Set("Content-Type", "application/json")

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	//fmt.Println("response: ", response)

	// 关闭 读取 reader
	defer response.Body.Close()

	// 读取内容
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
	return string(all), err
}

func SendQQEmail(from, secret, to string, title string, text string) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, from)
	m.SetAddressHeader("To", to, to)

	m.SetHeader("Subject", title)
	m.SetBody("text/plain", text)
	d := gomail.NewDialer("smtp.qq.com", 587, from, secret)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.GetLogger().Error(err.Error())
		return err
	}
	return nil
}
