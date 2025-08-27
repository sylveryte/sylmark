package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sylmark-server/lspserver"

	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	logFile := setLogger()
	defer logFile.Close()
	slog.Info("Hey, We're up!--------------------------------------------------")

	ctx := context.Background()
	stream := jsonrpc2.NewBufferedStream(stdwrc{}, jsonrpc2.VSCodeObjectCodec{})

	handler := lspserver.NewHandler()
	handler.SetupGrammars()
	defer handler.Parser.Close()
	jsonHandler := jsonrpc2.HandlerWithError(handler.Handle)
	<- jsonrpc2.NewConn(ctx, stream, jsonHandler).DisconnectNotify()

	slog.Info("Closing the lsp.")
}

func setLogger() *os.File {
	logFilePath := filepath.Join("/tmp", "sylmark-server.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	logger := slog.New(slog.NewTextHandler(logFile, nil))
	slog.SetDefault(logger)
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
