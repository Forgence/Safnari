package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	startDir := flag.String("start-dir", ".", "starting directory for file scanning")
	debug := flag.Bool("debug", false, "print debug logs")
	scanSubDirs := flag.Bool("scan-subdirs", false, "scan subdirectories")
	calculateHashes := flag.Bool("hashes", false, "calculate MD5, SHA1, SHA256 and SSDeep hashes for files")
	excludePattern := flag.String("exclude-pattern", "", "pattern to exclude files/directories from scanning")
	maxFileSize := flag.Int64("max-file-size", 0, "maximum file size in bytes to be included in the scan")
	fileTypes := flag.String("file-types", "", "comma-separated list of file extensions to include in the scan")
	hashAlgorithms := flag.String("hash-algorithms", "md5,sha1,sha256,ssdeep", "comma-separated list of hash algorithms to use")
	timeout := flag.Duration("timeout", 0, "timeout duration for the scanning process (e.g., 30s, 5m)")
	outputFile := flag.String("output", "file_data.json", "output file for the JSON data")
	flag.Parse()

	// Convert comma-separated strings to slices
	var fileTypesSlice []string
	if *fileTypes != "" {
		fileTypesSlice = strings.Split(*fileTypes, ",")
	}
	hashAlgorithmsSlice := strings.Split(*hashAlgorithms, ",")

	if *startDir == "" {
		fmt.Println("Please provide a starting directory with the -start-dir flag.")
		os.Exit(1)
	}

	absStartDir, err := filepath.Abs(*startDir)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	logFile, err := os.Create("error.log")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}
	defer logFile.Close()
	errorLog := log.New(logFile, "", log.LstdFlags)

	ctx, cancel := context.WithCancel(context.Background())
	if *timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, *timeout)
	}
	defer cancel()

	sysInfo := &SystemInfo{
		FileData: &FileScanResult{
			TotalFiles: 0,
			FileList:   make([]FileData, 0),
		},
	}

	fmt.Println("Gathering host info...")
	sysInfo.getHostInfo(*debug, errorLog)

	fmt.Println("Gathering process info...")
	sysInfo.getProcesses(*debug, errorLog)

	fmt.Println("Gathering CPU usage info...")
	sysInfo.getCPUUsage(*debug, errorLog)

	fmt.Println("Scanning files...")
	sysInfo.scanFiles(ctx, absStartDir, *scanSubDirs, *calculateHashes, *excludePattern, *maxFileSize, fileTypesSlice, hashAlgorithmsSlice, *debug, errorLog)

	fmt.Println("Saving data to file...")
	sysInfo.saveToFile(*outputFile)

	fmt.Println("Scan completed!")
}
