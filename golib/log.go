package golib

import "fmt"

func Log(format string, a ...interface{}) (n int, err error) {
	//if a != nil{
		return fmt.Printf(format,a...)
	//}else {
	//	return fmt.Printf(format,nil)
	//}
}
