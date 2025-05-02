package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sylmark/rpc"
)

func main() {
	logFile := setLogger()
	defer logFile.Close()
	slog.Info("Hey, We're up!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Text()
		handleMessage(msg)
	}
}

func handleMessage(msg any) {
	slog.Info(msg.(string))
}

func setLogger() *os.File {
	logFilePath := filepath.Join("/tmp", "sylmark.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)
	fmt.Println("Logging successfully configured")
	return logFile
}
