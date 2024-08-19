package executor

import (
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"process-orchestrator/internal/process"
	"sync"
	"syscall"
	"time"
)

type Watcher struct {
	exe  *exec.Cmd
	done chan bool
	id   string
}

func NewWatcher(exe *exec.Cmd, id string) *Watcher {
	return &Watcher{exe: exe, done: make(chan bool), id: id}
}

// Start starts the network check
func (w *Watcher) Start(wait, restart bool) error {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-w.done:
				return
			case <-exit:
				return
			default:
				for {
					proc, _ := w.startProcess(w.exe, wait || restart)
					if proc != nil {
						slog.Debug("Process started",
							slog.String("ID", w.id),
							slog.String("Component", "executor"),
							slog.Int("Pid", proc.Pid),
						)
					}
					if !restart || (!wait && restart) {
						break
					}
					w.exe.Process = nil
					w.exe.ProcessState = nil
					time.Sleep(1 * time.Second)
				}
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
	return nil
}

// startProcess starts a process
func (w *Watcher) startProcess(cmd *exec.Cmd, wait bool) (*os.Process, error) {
	if err := cmd.Start(); err != nil {
		slog.Error("Error starting process",
			slog.String("ID", w.id),
			slog.String("Component", "executor"),
			slog.String("Command", cmd.Path),
			slog.String("Error", err.Error()))
		return nil, err
	}
	for cmd.Process == nil || (cmd.ProcessState != nil && cmd.ProcessState.Exited()) {
		slog.Debug("Waiting for process to start",
			slog.String("ID", w.id),
			slog.String("Component", "executor"),
		)
		time.Sleep(100 * time.Millisecond)
	}
	err := process.PrintProcessInfo(int32(cmd.Process.Pid), w.id)
	if err != nil {
		slog.Error("Error printing process info",
			slog.String("ID", w.id),
			slog.String("Component", "executor"),
			slog.String("Error", err.Error()))
	}
	if wait {
		err := cmd.Wait()
		if err != nil {
			slog.Error("Error waiting for process",
				slog.String("ID", w.id),
				slog.String("Component", "executor"),
				slog.String("Error", err.Error()))
			return nil, err
		}
	} else {
		slog.Debug("Not waiting for process",
			slog.String("ID", w.id),
			slog.String("Component", "executor"),
			slog.Int("Pid", cmd.Process.Pid),
		)
		cmd.Process.Release()
	}
	return cmd.Process, nil
}

// Stop stops the process thread
func (w *Watcher) Stop() {
	go func() {
		w.done <- true
	}()
}

// Release releases the process check
func (w *Watcher) Release() {
	if w.exe.Process != nil {
		w.exe.Process.Release()
	}
}
