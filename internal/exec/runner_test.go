package exec

import (
	"context"
	"errors"
	"testing"
)

type fakeRunner struct {
	lookPathErr error
	runErr      error
	output      string
	outputErr   error
}

func (f *fakeRunner) Run(_ context.Context, _ string, _ ...string) error {
	return f.runErr
}

func (f *fakeRunner) Output(_ context.Context, _ string, _ ...string) (string, error) {
	return f.output, f.outputErr
}

func (f *fakeRunner) LookPath(_ string) (string, error) {
	if f.lookPathErr != nil {
		return "", f.lookPathErr
	}
	return "/usr/bin/fake", nil
}

func TestCommandAvailable(t *testing.T) {
	if CommandAvailable(&fakeRunner{lookPathErr: errors.New("missing")}, "git") {
		t.Fatalf("expected command unavailable")
	}
	if !CommandAvailable(&fakeRunner{}, "git") {
		t.Fatalf("expected command available")
	}
}

func TestWrapExecError(t *testing.T) {
	err := wrapExecError(errors.New("boom"), "git", []string{"clone"}, "fatal: nope")
	if err == nil || err.Error() != "git clone: fatal: nope" {
		t.Fatalf("unexpected error: %v", err)
	}
}
