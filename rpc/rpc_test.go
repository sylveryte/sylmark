package rpc_test

import (
	"sylmark/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	expected := 15
	msg := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content ,_ := rpc.DecodeMessage([]byte(msg))
	if method != "hi" {
		t.Fatalf("Expected: %s, Actual: %s", "hi", method)
	}
	if expected != len(content) {
		t.Fatalf("Expected: %d, Actual: %d", expected, len(content))
	}
}
