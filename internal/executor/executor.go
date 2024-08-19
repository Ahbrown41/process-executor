package executor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"process-orchestrator/internal/config"
)

type Process struct {
	cmd     *exec.Cmd
	process config.Process
	id      string
}

func New(process config.Process, id string) *Process {
	return &Process{process: process, id: id}
}

// ToString - Print out the object
func (p *Process) ToString() string {
	return fmt.Sprintf("%v", p)
}

// Execute starts the executor
func (p *Process) Execute() error {
	var ctx context.Context
	var cancel context.CancelFunc
	if p.process.Wait && p.process.WaitMax != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), p.process.WaitMax)
		defer cancel()
	} else {
		ctx = context.Background()
	}
	attrs := []any{
		slog.String("ID", p.id),
		slog.String("Component", "executor"),
		slog.String("Name", p.process.Name),
		slog.String("WorkDir", p.process.WorkDir),
		slog.String("Command", p.process.Command),
	}
	if p.process.Arguments == nil || len(p.process.Arguments) == 0 {
		p.cmd = exec.CommandContext(ctx, p.process.Command)
		attrs = append(attrs, slog.Any("Args", p.process.Arguments))
	} else {
		p.cmd = exec.CommandContext(ctx, p.process.Command, p.process.Arguments...)
	}
	if p.process.WorkDir != "" {
		p.cmd.Dir = p.process.WorkDir
	}
	p.cmd.Stdout = p
	p.cmd.Stderr = p
	p.cmd.Env = os.Environ()
	for _, env := range p.process.Environment {
		env := fmt.Sprintf("%s=%s", env.Key, env.Value)
		slog.Debug("Setting Environment Variable",
			slog.String("Component", "executor"),
			slog.String("ID", p.id),
			slog.String("Env", env),
		)
		p.cmd.Env = append(p.cmd.Env, env)
	}

	cmd := NewWatcher(p.cmd, p.id)
	slog.Info("Starting Process", attrs...)
	err := cmd.Start(p.process.Wait, p.process.Restart)
	if err != nil {
		return err
	}

	defer func() {
		cmd.Release()
		cmd.Stop()
	}()
	return nil
}

// GetPid - Get Process ID
func (p *Process) GetPid() int {
	return p.cmd.Process.Pid
}

// Write - Write to the log
func (p *Process) Write(b []byte) (int, error) {
	if len(b) > 0 {
		slog.Info(string(b))
	}
	return len(b), nil
}
