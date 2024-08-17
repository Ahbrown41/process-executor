package checks

import (
	"context"
	"log/slog"
	"process-orchestrator/internal/config"
)

func ExecuteConditions(conditions []config.Condition) error {
	if conditions == nil || len(conditions) == 0 {
		return nil
	}
	for _, check := range conditions {
		var cancel context.CancelFunc
		var err error
		ctx := context.Background()
		if check.Timeout != 0 {
			ctx, cancel = context.WithTimeout(context.Background(), check.Timeout)
		}
		switch check.Type {
		case "network":
			err = checkNetwork(ctx, check)
		default:
			slog.Debug("Unknown Condition Type",
				slog.String("Component", "checks"),
				slog.String("Type", check.Type),
			)
			continue
		}
		if err != nil {
			slog.Info("Check Failed",
				slog.String("Component", "checks"),
				slog.String("Type", check.Type),
				slog.String("Name", check.Name),
				slog.String("Error", err.Error()),
			)
		} else {
			slog.Debug("Check Passed",
				slog.String("Component", "checks"),
				slog.String("Type", check.Type),
				slog.String("Name", check.Name),
			)
		}
		if cancel != nil {
			cancel()
		}
	}
	return nil
}
