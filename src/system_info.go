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
	HostInfo  *host.InfoStat  `json:"host_info"`
	Processes []ProcessInfo   `json:"processes"`
	CPUUsage  float64         `json:"cpu_usage"`
	FileData  *FileScanResult `json:"file_data"`
}

func (sysInfo *SystemInfo) getHostInfo(debug bool, errorLog *log.Logger) {
	hostInfo, err := host.Info()
	if err != nil {
		if debug {
			errorLog.Printf("Error getting host info: %v\n", err)
		}
		return
	}
	sysInfo.HostInfo = hostInfo
}

// system_info.go
func (sysInfo *SystemInfo) getProcesses(debug bool, errorLog *log.Logger) {
	processList, err := process.Processes()
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
		createTime, err := proc.CreateTime() // Handle error from CreateTime
		if err != nil {
			if debug {
				errorLog.Printf("Error getting process create time: %v\n", err)
			}
			continue
		}
		memoryInfo, err := proc.MemoryInfo() // Handle error from MemoryInfo
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
	cpuPercent, err := cpu.Percent(0, false)
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

// system_info.go
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

	// Updated call to traverseFiles
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
