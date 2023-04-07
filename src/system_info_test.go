package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

type MockHostInfoProvider struct {
	HostInfo *host.InfoStat
	Err      error
}

func (p MockHostInfoProvider) Info() (*host.InfoStat, error) {
	return p.HostInfo, p.Err
}

type MockCPUProvider struct {
	CPUPercent []float64
	Err        error
}

func (p MockCPUProvider) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return p.CPUPercent, p.Err
}

type MockProcessProvider struct {
	ProcessList []*process.Process
	Err         error
}

func (p MockProcessProvider) Processes() ([]*process.Process, error) {
	return p.ProcessList, p.Err
}

func TestGetHostInfo(t *testing.T) {
	var buf bytes.Buffer
	errorLog := log.New(&buf, "", log.LstdFlags)

	t.Run("Success", func(t *testing.T) {
		sysInfo := &SystemInfo{
			HostInfoProvider: MockHostInfoProvider{
				HostInfo: &host.InfoStat{Hostname: "test-host"},
			},
		}
		sysInfo.getHostInfo(false, errorLog)
		if sysInfo.HostInfo == nil || sysInfo.HostInfo.Hostname != "test-host" {
			t.Error("Expected HostInfo to be populated, but it is not")
		}
	})

	t.Run("Error", func(t *testing.T) {
		sysInfo := &SystemInfo{
			HostInfoProvider: MockHostInfoProvider{
				Err: errors.New("mock error"),
			},
		}
		sysInfo.getHostInfo(true, errorLog)
		if sysInfo.HostInfo != nil {
			t.Error("Expected HostInfo to be nil, but it is not")
		}
		if !strings.Contains(buf.String(), "Error getting host info") {
			t.Error("Expected error message to be logged, but it is not")
		}
	})
}

func TestGetProcesses(t *testing.T) {
	var buf bytes.Buffer
	errorLog := log.New(&buf, "", log.LstdFlags)

	t.Run("Success", func(t *testing.T) {
		mockProcess := &process.Process{Pid: 123}
		sysInfo := &SystemInfo{
			ProcessProvider: MockProcessProvider{
				ProcessList: []*process.Process{mockProcess},
			},
		}
		sysInfo.getProcesses(false, errorLog)
		if len(sysInfo.Processes) == 0 {
			t.Error("Expected Processes to be populated, but it is not")
		}
	})

	t.Run("Error", func(t *testing.T) {
		sysInfo := &SystemInfo{
			ProcessProvider: MockProcessProvider{
				Err: errors.New("mock error"),
			},
		}
		sysInfo.getProcesses(true, errorLog)
		if len(sysInfo.Processes) != 0 {
			t.Error("Expected Processes to be empty, but it is not")
		}
		if !strings.Contains(buf.String(), "Error getting process list") {
			t.Error("Expected error message to be logged, but it is not")
		}
	})
}

func TestGetCPUUsage(t *testing.T) {
	var buf bytes.Buffer
	errorLog := log.New(&buf, "", log.LstdFlags)

	t.Run("Success", func(t *testing.T) {
		sysInfo := &SystemInfo{
			CPUProvider: MockCPUProvider{
				CPUPercent: []float64{25.0},
			},
		}
		sysInfo.getCPUUsage(false, errorLog)
		if sysInfo.CPUUsage != 25.0 {
			t.Errorf("Expected CPUUsage to be 			25.0, but got %f", sysInfo.CPUUsage)
		}
	})

	t.Run("Error", func(t *testing.T) {
		sysInfo := &SystemInfo{
			CPUProvider: MockCPUProvider{
				Err: errors.New("mock error"),
			},
		}
		sysInfo.getCPUUsage(true, errorLog)
		if sysInfo.CPUUsage != 0 {
			t.Errorf("Expected CPUUsage to be 0, but got %f", sysInfo.CPUUsage)
		}
		if !strings.Contains(buf.String(), "Error getting CPU percent") {
			t.Error("Expected error message to be logged, but it is not")
		}
	})
}

func TestScanFiles(t *testing.T) {
	var buf bytes.Buffer
	errorLog := log.New(&buf, "", log.LstdFlags)

	sysInfo := &SystemInfo{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// For the sake of brevity, we're using a simple case that scans the current directory ".".
	// We don't provide a mock implementation for traverseFiles, so this test relies on the actual file system.
	sysInfo.scanFiles(ctx, ".", true, true, "", 0, nil, nil, true, errorLog)

	if sysInfo.FileData == nil {
		t.Error("Expected FileData to be populated, but it is nil")
	}
}

func TestSaveToFile(t *testing.T) {
	sysInfo := &SystemInfo{
		HostInfo: &host.InfoStat{Hostname: "test-host"},
	}
	err := sysInfo.saveToFile("test_output.json")
	if err != nil {
		t.Errorf("Expected saveToFile to succeed, but got error: %v", err)
	}

	// Cleanup: remove the test output file.
	os.Remove("test_output.json")
}
