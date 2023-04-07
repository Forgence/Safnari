package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemInfo struct {
	HostInfo         *host.InfoStat  `json:"host_info"`
	Processes        []ProcessInfo   `json:"processes"`
	CPUUsage         float64         `json:"cpu_usage"`
	FileData         *FileScanResult `json:"file_data"`
	HostInfoProvider HostInfoProvider
	CPUProvider      CPUProvider
	ProcessProvider  ProcessProvider
}

type HostInfoProvider interface {
	Info() (*host.InfoStat, error)
}

type DefaultHostInfoProvider struct{}

func (p DefaultHostInfoProvider) Info() (*host.InfoStat, error) {
	return host.Info()
}

type CPUProvider interface {
	Percent(interval time.Duration, percpu bool) ([]float64, error)
}

type DefaultCPUProvider struct{}

func (p DefaultCPUProvider) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu.Percent(interval, percpu)
}

type ProcessProvider interface {
	Processes() ([]*process.Process, error)
}

type DefaultProcessProvider struct{}

func (p DefaultProcessProvider) Processes() ([]*process.Process, error) {
	return process.Processes()
}

func (sysInfo *SystemInfo) getHostInfo(debug bool, errorLog *log.Logger) {
	hostInfo, err := sysInfo.HostInfoProvider.Info()
	if err != nil {
		if debug {
			errorLog.Printf("Error getting host info: %v\n", err)
		}
		return
	}
	sysInfo.HostInfo = hostInfo
}

func (sysInfo *SystemInfo) getProcesses(debug bool, errorLog *log.Logger) {
	processList, err := sysInfo.ProcessProvider.Processes()
	if err != nil {
		if debug {
			errorLog.Printf("Error getting process list: %v\n", err)
		}
		return
	}

	processes := make([]ProcessInfo, 0)

	for _, proc := range processList {
		pid := proc.Pid
		name, _ := proc.Name()
		cmd, _ := proc.Cmdline()
		createTime, err := proc.CreateTime()
		if err != nil {
			if debug {
				errorLog.Printf("Error getting process create time: %v\n", err)
			}
			continue
		}
		memoryInfo, err := proc.MemoryInfo()
		if err != nil {
			if debug {
				errorLog.Printf("Error getting process memory info: %v\n", err)
			}
			continue
		}

		processInfo := ProcessInfo{
			PID:        int(pid),
			Name:       name,
			Cmd:        cmd,
			CreateTime: time.Unix(0, createTime*int64(time.Millisecond)),
			RSS:        memoryInfo.RSS,
			VMS:        memoryInfo.VMS,
		}
		processes = append(processes, processInfo)
	}

	sysInfo.Processes = processes
}

func (sysInfo *SystemInfo) getCPUUsage(debug bool, errorLog *log.Logger) {
	cpuPercent, err := sysInfo.CPUProvider.Percent(0, false)
	if err != nil {
		if debug {
			errorLog.Printf("Error getting CPU percent: %v\n", err)
		}
		return
	}

	if len(cpuPercent) > 0 {
		sysInfo.CPUUsage = cpuPercent[0]
	}
}

func (sysInfo *SystemInfo) scanFiles(ctx context.Context, startDir string, scanSubDirs, calculateHashes bool, excludePattern string, maxFileSize int64, fileTypes, hashAlgorithms []string, debug bool, errorLog *log.Logger) {
	fileProcessor := make(chan FileData)
	fileScanResult := &FileScanResult{
		TotalFiles: 0,
		FileList:   make([]FileData, 0),
	}

	go func() {
		for fileData := range fileProcessor {
			fileScanResult.TotalFiles++
			fileScanResult.FileList = append(fileScanResult.FileList, fileData)
		}
	}()

	traverseFiles(ctx, startDir, fileProcessor, scanSubDirs, calculateHashes, excludePattern, maxFileSize, fileTypes, hashAlgorithms, debug, errorLog)
	close(fileProcessor)

	sysInfo.FileData = fileScanResult
}

func (sysInfo *SystemInfo) saveToFile(outputFilePath string) error {
	data, err := json.MarshalIndent(sysInfo, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
