package logger

import (
	"fmt"
	"time"
)

var DefaultLogger = new(printfLogger)

type Logger interface {
	InfoF(format string, args ...interface{})
	ErrorF(format string, args ...interface{})
}

type printfLogger struct{}

func (printfLogger) InfoF(format string, args ...interface{}) {
	fmt.Printf("[INFO] - %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}

func (printfLogger) ErrorF(format string, args ...interface{}) {
	fmt.Printf("[\033[31mERROR\033[0m] - %s - [\033[31m%s\033[0m]\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}
