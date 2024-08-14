package executor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"process-orchestrator/internal/config"
	"process-orchestrator/internal/process"
)

type Process struct {
	cmd     *exec.Cmd
	process config.Process
}

func New(process config.Process) *Process {
	return &Process{process: process}
}

// ToString - Print out the object
func (p *Process) ToString() string {
	return fmt.Sprintf("%v", p)
}

// Start starts the executor
func (p *Process) Start() error {
	var ctx context.Context
	var cancel context.CancelFunc
	var waitTimeStr string
	if p.process.Wait && p.process.WaitTime != 0 {
		ctx, cancel = context.WithTimeout(context.Background(), p.process.WaitTime)
		waitTimeStr = p.process.WaitTime.String()
		defer cancel()
	} else {
		waitTimeStr = "process to complete"
		ctx = context.Background()
	}
	slog.Info("Starting Process",
		slog.String("Name", p.process.Name),
		slog.String("Command", p.process.Command),
		slog.Any("Args", p.process.Arguments),
	)
	p.cmd = exec.CommandContext(ctx, p.process.Command, p.process.Arguments...)
	p.cmd.Stdout = p
	p.cmd.Stderr = p
	p.cmd.Env = os.Environ()
	for _, env := range p.process.Environment {
		p.cmd.Env = append(p.cmd.Env, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	if p.process.Watch {
		slog.Info("Watching Process")
		watcher := NewWatcher(p.cmd)
		err := watcher.Watch()
		if err != nil {
			return err
		}
	} else {
		err := p.cmd.Start()
		if err != nil {
			return err
		}
		err = process.PrintProcessInfo(int32(p.cmd.Process.Pid))
		if err != nil {
			return err
		}
		if p.process.Wait {
			slog.Info("Waiting for Process", slog.String("Wait Time", waitTimeStr))
			err = p.cmd.Wait()
			if err != nil {
				return err
			}
			slog.Info("Process Completed")
		}
	}
	return nil
}

// GetPid - Get Process ID
func (p *Process) GetPid() int {
	return p.cmd.Process.Pid
}

func (e Process) Write(p []byte) (int, error) {
	slog.Info(string(p))
	return len(p), nil
}
