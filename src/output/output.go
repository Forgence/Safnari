package output

import (
	"encoding/json"
	"os"
	"sync"

	"safnari/config"
	"safnari/systeminfo"
)

var (
	outputFile   *os.File
	outputWriter *JSONWriter
	cfg          *config.Config
	mu           sync.Mutex
	currentSize  int64
)

type Metrics struct {
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	TotalFiles     int    `json:"total_files"`
	FilesProcessed int    `json:"files_processed"`
	TotalProcesses int    `json:"total_processes"`
}

type OutputData struct {
	SystemInfo *systeminfo.SystemInfo    `json:"system_info,omitempty"`
	Processes  *[]systeminfo.ProcessInfo `json:"processes,omitempty"`
	Files      []map[string]interface{}  `json:"files"`
	Metrics    *Metrics                  `json:"metrics,omitempty"`
}

type JSONWriter struct {
	encoder *json.Encoder
	file    *os.File
	data    OutputData
}

func Init(config *config.Config, sysInfo *systeminfo.SystemInfo, metrics *Metrics) error {
	cfg = config
	var err error
	outputFile, err = os.OpenFile(cfg.OutputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	outputWriter = NewJSONWriter(outputFile)

	// Initialize OutputData
	outputWriter.data.SystemInfo = sysInfo
	outputWriter.data.Processes = &sysInfo.RunningProcesses

	// Update metrics with total process count
	if metrics != nil {
		metrics.TotalProcesses = len(sysInfo.RunningProcesses)
	}

	outputWriter.data.Metrics = metrics

	return nil
}

func WriteData(data map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()

	outputWriter.data.Files = append(outputWriter.data.Files, data)

	// Update file count metric
	if outputWriter.data.Metrics != nil {
		outputWriter.data.Metrics.FilesProcessed++
	}

	// Check for output file size rotation if needed (not implemented in this version)
}

func SetMetrics(metrics Metrics) {
	mu.Lock()
	defer mu.Unlock()

	outputWriter.data.Metrics = &metrics
}

func Close() {
	mu.Lock()
	defer mu.Unlock()

	outputWriter.Flush()
	outputFile.Close()
}

func NewJSONWriter(file *os.File) *JSONWriter {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return &JSONWriter{
		encoder: encoder,
		file:    file,
		data: OutputData{
			Files: []map[string]interface{}{},
		},
	}
}

func (w *JSONWriter) Flush() error {
	return w.encoder.Encode(w.data)
}
