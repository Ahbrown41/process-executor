package main

import (
	"fmt"
	"log/slog"
	"os"
	"process-orchestrator/internal/config"
	"process-orchestrator/internal/executor"
	"process-orchestrator/internal/host"
	"process-orchestrator/internal/screen"
	"sync"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})))

	// Load the configuration
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		slog.Error("config.Load() failed: %s\n", err)
		return
	}

	// Print host information
	if err = host.PrintHostInfo(); err != nil {
		slog.Error("host.PrintHostInfo() failed: %s\n", err)
		return
	}

	// Create the GUI
	win := screen.New().SetImage(cfg.Screen.BootImage)
	if cfg.Screen.FullScreen {
		win.FullScreen()
	}

	win.SetProgress("Starting")
	var wg sync.WaitGroup

	// Run the background process in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Waiting for processes to complete")
		for _, p := range cfg.Processes {
			win.SetProgress(fmt.Sprintf("Starting Process: %s (%s)", slog.String("Process", p.Name), slog.String("Command", p.Command)))
			err = executor.New(p).Start()
			if err != nil {
				slog.Info("Error executing process", slog.String("Process", p.Name), slog.String("Error", err.Error()))
			}
		}
		win.SetProgress("All processes completed")
		win.Close()
	}()

	win.GetWindow().ShowAndRun()

	// Wait for the background process to complete
	wg.Wait()
}
