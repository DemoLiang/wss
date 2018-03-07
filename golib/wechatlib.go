package golib

import (
	"fmt"
	"net/url"
	"net/http"
	"bytes"
	"io/ioutil"
)

const(
	WECHAT_HOST = "https://api.weixin.qq.com"
)

const(
	_ int =iota
	WeChatProgramAuthReq
)

var WeChatProgramReqPath = map[int]string{
	WeChatProgramAuthReq: "/sns/jscode2session",
}


func WeChatProgramAuth(code,appid string) (rsp string, err error) {
	headData := url.Values{}
	headData.Set("appid", appid)
	headData.Set("secret","SECRET")
	headData.Set("js_code",code)
	headData.Set("grant_type","authorization_code")

	u, _ := url.ParseRequestURI(WECHAT_HOST)
	u.Path = WeChatProgramReqPath[WeChatProgramAuthReq]
	u.RawQuery = headData.Encode()
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	request, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(""))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	res, err := client.Do(request)
	if err != nil {
		Log("Post Msg to WeChat push Error: %s\n", err)
		return "", err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		Log("Read Result Error: %v\n", result)
		return "", err
	}
	Log("Post Msg to WeChat Result: %s\n", result)
	return string(result), nil
}