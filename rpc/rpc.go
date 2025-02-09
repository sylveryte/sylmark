package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Content-Length: ...\r\n
// \r\n
//
//	{
//		"jsonrpc": "2.0",
//		"id": 1,
//		"method": "textDocument/completion",
//		"params": {
//			...
//		}
//	}
func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if found == false {
		return "", nil, errors.New("Syntax error")
	}

	contentLengthString := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthString))

	if err != nil {
		return "", nil, err
	}
	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "", nil, err
	}

	return baseMessage.Method, content[:contentLength], nil

}
