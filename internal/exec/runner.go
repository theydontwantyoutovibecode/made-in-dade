package exec

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) error
	Output(ctx context.Context, name string, args ...string) (string, error)
	LookPath(name string) (string, error)
}

type SystemRunner struct{}

func NewSystemRunner() *SystemRunner {
	return &SystemRunner{}
}

func (r *SystemRunner) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return wrapExecError(err, name, args, stderr.String())
	}
	return nil
}

func (r *SystemRunner) Output(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", wrapExecError(err, name, args, stderr.String())
	}
	return strings.TrimSpace(string(out)), nil
}

func (r *SystemRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func CommandAvailable(r Runner, name string) bool {
	_, err := r.LookPath(name)
	return err == nil
}

func wrapExecError(err error, name string, args []string, stderr string) error {
	cmd := strings.Join(append([]string{name}, args...), " ")
	if stderr != "" {
		return errors.New(cmd + ": " + strings.TrimSpace(stderr))
	}
	return errors.New(cmd + ": " + err.Error())
}
