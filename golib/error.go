package golib

import "fmt"

type ErrorCode int

const (
	FGEBase = 10000
)

const (
	EInternalError ErrorCode = 5*FGEBase + iota // value --> 0

)

//客户端请求错误 40000+iota
const (
	ErrClientError ErrorCode = 4*FGEBase + iota // value --> 0
)

var evtDesc = map[ErrorCode]string{
	//系统错误 5xx
	EInternalError: "系统内部错误",

	//客户端请求错误 4xx
	ErrClientError: "客户端错误",
}

func (f ErrorCode) String() string {
	if value, ok := evtDesc[f]; ok {
		return value
	}
	return "Unknow"
}

func (f ErrorCode) Error() string {
	return fmt.Sprintf("%d:%s", int(f), f.String())
}
