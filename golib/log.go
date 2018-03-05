package golib

import "fmt"

func Log(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format, a...)
}
