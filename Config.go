package WebOrm

import (
	"io"
)

type Config struct {
	DatabaseName string
	User         string
	Password     string
	//默认：127.0.0.1
	Ip string
	//默认：3306
	Port string
	//默认：3. 1:只输出panic日志 2:只输出error和panic日志 3: 所有日志都输出
	LogLervel int
	//默认 ：os。Stdout
	LogOutPut io.Writer
	//默认：true
	IsShowSql bool
}
