package wechat

type WeChatUserInfoRsp struct {
	Openid      string `json:"openid"`
	Session_key string `json:"session_key"`
}
