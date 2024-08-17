package main

import (
	"fmt"
	"log/slog"
	"os"
	"process-orchestrator/internal/checks"
	"process-orchestrator/internal/config"
	"process-orchestrator/internal/executor"
	"process-orchestrator/internal/host"
	"process-orchestrator/internal/screen"
)

func main() {
	// Change Working Directory to Directory of EXE
	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exePath := filepath.Dir(ex)
	//os.Chdir(exePath)

	// Setup Logger
	logFile, err := os.Create("process-orchestrator.log")
	if err != nil {
		slog.Error("Error opening log file", slog.String("Error", err.Error()))
	}
	defer func() {
		logFile.Close()
	}()
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})))

	// Load the configuration
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		slog.Error("config.Load() failed", slog.String("Error", err.Error()))
		return
	}

	// Print host information
	if err = host.PrintHostInfo(); err != nil {
		slog.Error("host.PrintHostInfo() failed: %s\n", err)
		return
	}

	// Create the GUI
	win := screen.New().SetImage(cfg.Display.BootImage)
	if cfg.Display.FullScreen {
		win.FullScreen()
	}

	win.SetProgress("Starting")
	// Run the background process in a separate goroutine
	go func() {
		slog.Debug("Waiting for processes to complete",
			slog.String("Component", "main"),
		)
		for _, p := range cfg.Processes {
			win.SetProgress(fmt.Sprintf("Starting Process: %s (%s)", slog.String("Process", p.Name), slog.String("Command", p.Command)))
			err = checks.ExecuteConditions(p.PreConditions)
			if err != nil {
				slog.Info("Pre-Condition failed",
					slog.String("Component", "main"),
					slog.String("Process", p.Name),
				)
			}
			err = executor.New(p, p.Name).Execute()
			if err != nil {
				slog.Info("Error executing process", slog.String("Process", p.Name), slog.String("Error", err.Error()))
			}
		}
		win.SetProgress("All processes completed")
		slog.Debug("All processes are complete",
			slog.String("Component", "main"),
		)
		win.Close()
	}()
	win.GetWindow().ShowAndRun()
}
