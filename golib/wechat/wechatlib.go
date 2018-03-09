package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DemoLiang/wss/golib"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	WECHAT_HOST = "https://api.weixin.qq.com"
	APPID       = ""
	SECRET      = ""
)

const (
	_ int = iota
	WeChatProgramUserInfoReqPath
)

var WeChatProgramReqPath = map[int]string{
	WeChatProgramUserInfoReqPath: "/sns/jscode2session",
}

func WeChatProgramUserInfoReq(code, appid string) (rsp []byte, err error) {
	headData := url.Values{}
	headData.Set("appid", appid)
	headData.Set("secret", SECRET)
	headData.Set("js_code", code)
	headData.Set("grant_type", "authorization_code")

	u, _ := url.ParseRequestURI(WECHAT_HOST)
	u.Path = WeChatProgramReqPath[WeChatProgramUserInfoReqPath]
	u.RawQuery = headData.Encode()
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	request, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(""))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	res, err := client.Do(request)
	if err != nil {
		golib.Log("Post Msg to WeChat push Error: %s\n", err)
		return []byte(""), err
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		golib.Log("Read Result Error: %v\n", result)
		return []byte(""), err
	}
	golib.Log("Post Msg to WeChat Result: %s\n", result)
	return result, nil
}

func GetWeChatOpenIdByCode(code string) (openid, session_key string) {
	rsp, err := WeChatProgramUserInfoReq(code, APPID)
	if err != nil {
		return "", ""
	}
	var weChatUserInforsp WeChatUserInfoRsp
	json.Unmarshal([]byte(rsp), &weChatUserInforsp)

	return weChatUserInforsp.Openid, weChatUserInforsp.Session_key
}
