package logger

import (
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	WARN    = "warn"
	ERROR   = "error"
	FATAL   = "fatal"
)

func LogMessage(logType, message string, a ...interface{}) {
	switch logType {
	case FATAL:
		Fatal(message, a)
		break
	case WARN:
		Warn(message, a)
		break
	case ERROR:
		Error(message, a)
		break
	case SUCCESS:
		Success(message, a)
		break
	case INFO:
	default:
		Info(message, a)
		break
	}
}

func Fatal(message string, a ...interface{}) {
	t := time.Now()
	color.Red("["+t.Format("15:04:05")+" - FATAL]: "+message, a)
	os.Exit(1)
}

func Warn(message string, a ...interface{}) {
	t := time.Now()
	color.Yellow("["+t.Format("15:04:05")+" - WARNING]: "+message, a)
}

func Error(message string, a ...interface{}) {
	t := time.Now()
	color.Red("["+t.Format("15:04:05")+" - ERROR]: "+message, a)
}

func Success(message string, a ...interface{}) {
	t := time.Now()
	color.Green("["+t.Format("15:04:05")+" - SUCCESS]: "+message, a)
}

func Info(message string, a ...interface{}) {
	t := time.Now()
	color.White("["+t.Format("15:04:05")+" - INFO]: "+message, a)
}
