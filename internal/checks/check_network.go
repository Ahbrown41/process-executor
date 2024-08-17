package checks

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"process-orchestrator/internal/config"
	"syscall"
	"time"
)

// checkNetwork checks the network
func checkNetwork(ctx context.Context, check config.Condition) error {
	attributes := make(map[string]string)
	for _, val := range check.Attributes {
		attributes[val.Key] = val.Value
	}
	if attributes["hostPort"] == "" {
		slog.Error("Host Port not provided",
			slog.String("Component", "network"),
		)
		return errors.New("host Port not provided")
	}
	hostPort := attributes["hostPort"]

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	slog.Debug("Awaiting Network",
		slog.String("Component", "network"),
		slog.String("Name", check.Name),
		slog.String("HostPort", hostPort),
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-sig:
			return nil
		default:
			err := checkHostPort(hostPort)
			if err == nil || check.Wait == false {
				return err
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
	return nil
}

// checkNetwork checks the network
func checkHostPort(hostPort string) error {
	conn, err := net.DialTimeout("tcp", hostPort, time.Second*3)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
