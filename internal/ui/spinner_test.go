package ui

import (
	"errors"
	"strings"
	"testing"
)

func TestSpinnerFallbackSuccess(t *testing.T) {
	buffer := &strings.Builder{}
	spinner := NewSpinner(buffer, false)
	err := spinner.Run("Doing work", func() error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buffer.String(), "done") {
		t.Fatalf("expected done output")
	}
}

func TestSpinnerFallbackFailure(t *testing.T) {
	buffer := &strings.Builder{}
	spinner := NewSpinner(buffer, false)
	err := spinner.Run("Doing work", func() error { return errors.New("fail") })
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(buffer.String(), "failed") {
		t.Fatalf("expected failed output")
	}
}
