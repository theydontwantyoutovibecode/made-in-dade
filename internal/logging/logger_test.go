package logging

import (
	"bytes"
	"strings"
	"testing"
)

type bufferSink struct {
	out bytes.Buffer
	err bytes.Buffer
}

func TestLoggerPlainOutput(t *testing.T) {
	sink := bufferSink{}
	logger := New(&sink.out, &sink.err, false)

	logger.Info("hello")
	logger.Success("done")
	logger.Warn("careful")
	logger.Error("bad")

	out := sink.out.String()
	err := sink.err.String()

	if !strings.Contains(out, "hello") {
		t.Fatalf("expected info output")
	}
	if !strings.Contains(out, "✓ done") {
		t.Fatalf("expected success output")
	}
	if !strings.Contains(out, "⚠ careful") {
		t.Fatalf("expected warn output")
	}
	if !strings.Contains(err, "✗ bad") {
		t.Fatalf("expected error output")
	}
}

func TestLoggerSilent(t *testing.T) {
	sink := bufferSink{}
	logger := New(&sink.out, &sink.err, false)
	logger.SetSilent(true)

	logger.Info("hello")
	logger.Error("bad")

	if sink.out.Len() != 0 || sink.err.Len() != 0 {
		t.Fatalf("expected silent logger to write nothing")
	}
}

