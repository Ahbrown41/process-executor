package host

import (
	"github.com/shirou/gopsutil/v4/host"
	"log/slog"
)

func PrintHostInfo() error {
	uptime, err := host.Uptime()
	if err != nil {
		return err
	}
	info, err := host.Info()
	if err != nil {
		return err
	}
	slog.Info("Host Info",
		slog.String("Hostname", info.Hostname),
		slog.String("OS", info.OS),
		slog.String("Platform", info.Platform),
		slog.String("PlatformFamily", info.PlatformFamily),
		slog.String("PlatformVersion", info.PlatformVersion),
		slog.String("KernelVersion", info.KernelVersion),
		slog.Uint64("Uptime", uptime),
		slog.String("HostID", info.HostID),
		slog.String("VirtualizationSystem", info.VirtualizationSystem),
		slog.String("VirtualizationRole", info.VirtualizationRole),
		slog.String("HostID", info.HostID))
	return nil
}
