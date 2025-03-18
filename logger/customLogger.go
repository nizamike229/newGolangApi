package logger

import (
	"fmt"
	"time"
)

var LogStorage []string

func Info(msg string) {
	green := "\033[32m"
	cyan := "\033[36m"
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	LogStorage = append(LogStorage, cyan+"[INFO]"+green+msg)
	fmt.Println(cyan + "[INFO] " + timestamp + green + " " + msg)
}

func Error(msg string) {
	red := "\033[31m"
	cyan := "\033[36m"

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	LogStorage = append(LogStorage, cyan+"[ERROR]"+red+msg)
	fmt.Println(cyan + "[ERROR] " + timestamp + red + " " + msg)
}

func Warning(msg string) {
	yellow := "\033[33m"
	cyan := "\033[36m"

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	LogStorage = append(LogStorage, cyan+"[WARNING]"+yellow+msg)
	fmt.Println(cyan + "[WARNING] " + timestamp + yellow + " " + msg)
}
