package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sylmark/lsp"

	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	logFile := setLogger()
	defer logFile.Close()
	slog.Info("Hey, We're up!--------------------------------------------------")

	ctx := context.Background()
	stream := jsonrpc2.NewBufferedStream(stdwrc{}, jsonrpc2.VSCodeObjectCodec{})

	handler := lsp.NewHandler()
	handler.SetupGrammars()
	defer handler.Parser.Close()
	jsonHandler := jsonrpc2.HandlerWithError(handler.Handle)
	<- jsonrpc2.NewConn(ctx, stream, jsonHandler).DisconnectNotify()

	slog.Info("Closing the lsp.")
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

type stdwrc struct{}

func (stdwrc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdwrc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdwrc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
