package process

import (
	"fmt"
	processLib "github.com/shirou/gopsutil/v4/process"
	"log/slog"
)

// FindByName finds a process by name
func FindByName(name string) int32 {
	procs, err := processLib.Processes()
	if err != nil {
		slog.Error("Error getting processes", slog.String("Error", err.Error()))
	}
	for _, proc := range procs {
		pname, err := proc.Name()
		if err != nil {
			slog.Error("Error getting process name", slog.String("Error", err.Error()))
		}
		if pname == name {
			slog.Debug("Process",
				slog.Int("Pid", int(proc.Pid)),
				slog.String("Name", pname),
			)
		}
		return proc.Pid
	}
	return -1
}

// KillByName kills a process by name
func KillByName(name string) {
	pid := FindByName(name)
	if pid != -1 {
		proc, err := processLib.NewProcess(pid)
		if err != nil {
			slog.Error("Error getting process", slog.String("Error", err.Error()))
		}
		err = proc.Kill()
		if err != nil {
			slog.Error("Error killing process",
				slog.String("Name", name),
				slog.String("Error", err.Error()),
			)
		}
	}
}

// Children finds the children of a process by pid
func Children(pid int32) error {
	proc, err := processLib.NewProcess(pid)
	if err != nil {
		return err
	}
	attrs := []slog.Attr{
		slog.Int("Pid", int(proc.Pid)),
	}
	if children, err := proc.Children(); err == nil {
		for _, child := range children {
			name, err := child.Name()
			if err != nil {
				name = "Unknown"
			}
			attrs = append(attrs, slog.String("Child", fmt.Sprintf("%s (%d)", name, child.Pid)))
		}
	}
	slog.Debug("Children", attrs)
	return nil
}

func PrintProcesses() error {
	procs, err := processLib.Processes()
	if err != nil {
		return err
	}
	for _, proc := range procs {
		err := PrintProcessInfo(proc.Pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintProcessInfo prints information about a process
func PrintProcessInfo(pid int32) error {
	proc, err := processLib.NewProcess(pid)
	if err != nil {
		return err
	}
	attrs := []any{
		slog.Int("Pid", int(proc.Pid)),
	}

	name, err := proc.Name()
	if err != nil {
		name = "Unknown"
	}
	attrs = append(attrs, slog.String("Name", name))

	if parent, err := proc.Parent(); err == nil && parent != nil {
		parent, err := proc.Parent()
		if err != nil {
			return err
		}
		parentName, err := proc.Name()
		if err != nil {
			name = "Unknown"
		}
		attrs = append(attrs, slog.String("Parent", fmt.Sprintf("%s (%d)", parentName, parent.Pid)))
	}
	if children, err := proc.Children(); err == nil {
		for _, child := range children {
			name, err := child.Name()
			if err != nil {
				name = "Unknown"
			}
			attrs = append(attrs, slog.String("Child", fmt.Sprintf("%s (%d)", name, child.Pid)))
		}
	}
	slog.Debug("Process", attrs...)
	return nil
}
