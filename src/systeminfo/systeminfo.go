package systeminfo

import (
	"fmt"
	"safnari/config"
	"safnari/logger"

	"github.com/shirou/gopsutil/v3/process"
)

type SystemInfo struct {
	OSVersion        string        `json:"os_version"`
	InstalledPatches []string      `json:"installed_patches"`
	RunningProcesses []ProcessInfo `json:"running_processes"`
	StartupPrograms  []string      `json:"startup_programs"`
	InstalledApps    []string      `json:"installed_apps"`
}

type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent,omitempty"`
	MemoryPercent float32 `json:"memory_percent,omitempty"`
	Cmdline       string  `json:"cmdline,omitempty"`
	Username      string  `json:"username,omitempty"`
	Exe           string  `json:"exe,omitempty"`
}

func GetSystemInfo(cfg *config.Config) (*SystemInfo, error) {
	sysInfo := &SystemInfo{}

	if err := gatherOSVersion(sysInfo); err != nil {
		logger.Warnf("Failed to gather OS version: %v", err)
	}

	if err := gatherInstalledPatches(sysInfo); err != nil {
		logger.Warnf("Failed to gather installed patches: %v", err)
	}

	if err := gatherRunningProcesses(sysInfo, cfg.ExtendedProcessInfo); err != nil {
		logger.Warnf("Failed to gather running processes: %v", err)
	}

	if err := gatherStartupPrograms(sysInfo); err != nil {
		logger.Warnf("Failed to gather startup programs: %v", err)
	}

	if err := gatherInstalledApps(sysInfo); err != nil {
		logger.Warnf("Failed to gather installed applications: %v", err)
	}

	return sysInfo, nil
}

func gatherRunningProcesses(sysInfo *SystemInfo, extended bool) error {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get running processes: %v", err)
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		procInfo := ProcessInfo{
			PID:  p.Pid,
			Name: name,
		}

		if extended {
			cpuPercent, err := p.CPUPercent()
			if err == nil {
				procInfo.CPUPercent = cpuPercent
			}

			memPercent, err := p.MemoryPercent()
			if err == nil {
				procInfo.MemoryPercent = memPercent
			}

			cmdline, err := p.Cmdline()
			if err == nil {
				procInfo.Cmdline = cmdline
			}

			username, err := p.Username()
			if err == nil {
				procInfo.Username = username
			}

			exe, err := p.Exe()
			if err == nil {
				procInfo.Exe = exe
			}
		}

		sysInfo.RunningProcesses = append(sysInfo.RunningProcesses, procInfo)
	}

	return nil
}

// Implement gatherOSVersion, gatherInstalledPatches, gatherStartupPrograms, gatherInstalledApps as per previous implementations or stubs
