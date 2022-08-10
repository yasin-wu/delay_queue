package logger

import (
	"fmt"
	"time"
)

var DefaultLogger = new(printfLogger)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type printfLogger struct{}

func (printfLogger) Infof(format string, args ...any) {
	fmt.Printf("[INFO] - %s - %s\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}

func (printfLogger) Errorf(format string, args ...any) {
	fmt.Printf("[\033[31mERROR\033[0m] - %s - [\033[31m%s\033[0m]\n",
		time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
}
