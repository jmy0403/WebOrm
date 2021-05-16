package log

import (
	"github.com/jmy0403/WebOrm"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(out, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(out, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	panicLog = log.New(out, "\033[36[panic ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Panic  = panicLog.Panicln
	Panicf = panicLog.Panicf
)
var out io.Writer

func Set(cfg *WebOrm.Config) {

	if cfg.LogLervel == 0 {
		setLevel(3)
	} else {
		setLevel(cfg.LogLervel)
	}
	if cfg.LogOutPut == nil {
		out = os.Stdout
	} else {
		out = cfg.LogOutPut
	}
}

// 1:只输出panic日志 2:只输出error和panic日志 3: 所有日志都输出
func setLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(out)
	}
	if level == 1 {
		errorLog.SetOutput(ioutil.Discard)
		infoLog.SetOutput(ioutil.Discard)
	}
	if level == 2 {
		infoLog.SetOutput(ioutil.Discard)
	}

}
