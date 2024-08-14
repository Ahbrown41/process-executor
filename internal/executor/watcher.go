package executor

import (
	"log/slog"
	"os/exec"
	"process-orchestrator/internal/process"
	"sync"
)

type Watcher struct {
	exe      *exec.Cmd
	shutdown chan struct{}
}

func NewWatcher(exe *exec.Cmd) *Watcher {
	return &Watcher{exe: exe, shutdown: make(chan struct{})}
}

func (w *Watcher) Shutdown() {
	w.shutdown <- struct{}{}
	close(w.shutdown)
}

func (w *Watcher) Watch() error {
	defer close(w.shutdown)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		select {
		case <-w.shutdown:
			_ = w.exe.Process.Kill
			wg.Done()
		default:
			for {
				err := w.start()
				if err != nil {
					slog.Error("could not start process", slog.String("Error", err.Error()))
					wg.Done()
					return
				}
			}
		}
	}()
	wg.Wait()
	return nil
}

func (w *Watcher) start() error {
	err := w.exe.Run()
	if err != nil {
		return err
	}
	err = process.PrintProcessInfo(int32(w.exe.Process.Pid))
	if err != nil {
		return err
	}
	slog.Info("Process exited",
		slog.Int("Pid", w.exe.ProcessState.Pid()),
		slog.Bool("Exited", w.exe.ProcessState.Exited()),
		slog.Int("Exited", w.exe.ProcessState.ExitCode()),
	)
	return nil
}
