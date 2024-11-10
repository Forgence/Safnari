package config

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "runtime"
    "strconv"
    "strings"
)

type Config struct {
    StartPaths          []string `json:"start_paths"`
    AllDrives           bool     `json:"all_drives"`
    ScanFiles           bool     `json:"scan_files"`
    ScanProcesses       bool     `json:"scan_processes"`
    OutputFormat        string   `json:"output_format"`
    OutputFileName      string   `json:"output_file_name"`
    ConcurrencyLevel    int      `json:"concurrency_level"`
    NiceLevel           string   `json:"nice_level"`
    HashAlgorithms      []string `json:"hash_algorithms"`
    SearchTerms         []string `json:"search_terms"`
    IncludePatterns     []string `json:"include_patterns"`
    ExcludePatterns     []string `json:"exclude_patterns"`
    MaxFileSize         int64    `json:"max_file_size"`
    MaxOutputFileSize   int64    `json:"max_output_file_size"`
    LogLevel            string   `json:"log_level"`
    MaxIOPerSecond      int      `json:"max_io_per_second"`
    ConfigFile          string   `json:"config_file"`
    ExtendedProcessInfo bool     `json:"extended_process_info"`
    SensitiveDataTypes  []string `json:"sensitive_data_types"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        ScanFiles:     true, // Default to scanning files
        ScanProcesses: true, // Default to scanning processes
    }

    // Define command-line flags
    startPath := flag.String("path", "", "Start path(s) for scanning (comma-separated)")
    flag.BoolVar(&cfg.AllDrives, "all-drives", false, "Scan all local drives (Windows only)")
    flag.BoolVar(&cfg.ScanFiles, "scan-files", true, "Enable or disable file scanning")
    flag.BoolVar(&cfg.ScanProcesses, "scan-processes", true, "Enable or disable process scanning")
    flag.StringVar(&cfg.OutputFormat, "format", "json", "Output format: json or csv")
    flag.StringVar(&cfg.OutputFileName, "output", "output.json", "Output file name")
    flag.IntVar(&cfg.ConcurrencyLevel, "concurrency", 4, "Concurrency level")
    flag.StringVar(&cfg.NiceLevel, "nice", "medium", "Nice level: high, medium, low")
    hashes := flag.String("hashes", "md5,sha1,sha256", "Hash algorithms to use (comma-separated)")
    searches := flag.String("search", "", "Search terms (comma-separated)")
    includes := flag.String("include", "", "Include patterns (comma-separated)")
    excludes := flag.String("exclude", "", "Exclude patterns (comma-separated)")
    flag.Int64Var(&cfg.MaxFileSize, "max-file-size", 10485760, "Maximum file size to process (bytes)")
    flag.Int64Var(&cfg.MaxOutputFileSize, "max-output-file-size", 104857600, "Maximum output file size before rotation (bytes)")
    flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level: debug, info, warn, error, fatal, panic")
    flag.IntVar(&cfg.MaxIOPerSecond, "max-io-per-second", 1000, "Maximum disk I/O operations per second")
    flag.StringVar(&cfg.ConfigFile, "config", "", "Path to JSON configuration file")
    flag.BoolVar(&cfg.ExtendedProcessInfo, "extended-process-info", false, "Gather extended process information (requires elevated privileges)")
    sensitiveDataTypes := flag.String("sensitive-data-types", "", "Sensitive data types to scan for (comma-separated)")
    help := flag.Bool("help", false, "Display help message")

    flag.Parse()

    // Display help if requested or if no flags are provided
    if *help || len(os.Args) == 1 {
        displayHelp()
        os.Exit(0)
    }

    // Load configuration from file if specified
    if cfg.ConfigFile != "" {
        err := cfg.loadFromFile(cfg.ConfigFile)
        if err != nil {
            return nil, err
        }
    }

    // Override with command-line flags if set
    cfg.overrideWithFlags()

    // Parse comma-separated values
    cfg.StartPaths = parseCommaSeparated(*startPath)
    cfg.HashAlgorithms = parseCommaSeparated(*hashes)
    cfg.SearchTerms = parseCommaSeparated(*searches)
    cfg.IncludePatterns = parseCommaSeparated(*includes)
    cfg.ExcludePatterns = parseCommaSeparated(*excludes)
    cfg.SensitiveDataTypes = parseCommaSeparated(*sensitiveDataTypes)

    // Validate configuration
    err := cfg.validate()
    if err != nil {
        return nil, err
    }

    return cfg, nil
}

func displayHelp() {
    fmt.Println("Safnari - Advanced Cybersecurity Scanner")
    fmt.Println()
    fmt.Println("Usage:")
    fmt.Println("  safnari.exe [options]")
    fmt.Println()
    fmt.Println("Options:")
    flag.PrintDefaults()
    fmt.Println()
    fmt.Println("Examples:")
    fmt.Println("  safnari.exe --path \"C:\\\"")
    fmt.Println("  safnari.exe --path \"C:\\,D:\\\"")
    fmt.Println("  safnari.exe --all-drives --scan-files=false --scan-processes=true")
}

func (cfg *Config) loadFromFile(path string) error {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return fmt.Errorf("could not read config file: %v", err)
    }
    err = json.Unmarshal(data, cfg)
    if err != nil {
        return fmt.Errorf("invalid config file format: %v", err)
    }
    return nil
}

func (cfg *Config) overrideWithFlags() {
    flag.Visit(func(f *flag.Flag) {
        switch f.Name {
        case "path":
            cfg.StartPaths = parseCommaSeparated(f.Value.String())
        case "all-drives":
            cfg.AllDrives = true
        case "scan-files":
            cfg.ScanFiles = parseBoolFlagValue(f)
        case "scan-processes":
            cfg.ScanProcesses = parseBoolFlagValue(f)
        case "format":
            cfg.OutputFormat = f.Value.String()
        case "output":
            cfg.OutputFileName = f.Value.String()
        case "concurrency":
            cfg.ConcurrencyLevel = getIntFlagValue(f)
        case "nice":
            cfg.NiceLevel = f.Value.String()
        case "log-level":
            cfg.LogLevel = f.Value.String()
        case "max-file-size":
            cfg.MaxFileSize = getInt64FlagValue(f)
        case "max-output-file-size":
            cfg.MaxOutputFileSize = getInt64FlagValue(f)
        case "max-io-per-second":
            cfg.MaxIOPerSecond = getIntFlagValue(f)
        case "extended-process-info":
            cfg.ExtendedProcessInfo = true
        case "sensitive-data-types":
            cfg.SensitiveDataTypes = parseCommaSeparated(f.Value.String())
        }
    })
}

func (cfg *Config) validate() error {
    if !cfg.ScanFiles && !cfg.ScanProcesses {
        return fmt.Errorf("at least one of --scan-files or --scan-processes must be enabled")
    }
    if len(cfg.StartPaths) == 0 && !cfg.AllDrives && cfg.ScanFiles {
        return fmt.Errorf("either start path(s) or --all-drives must be specified for file scanning")
    }
    if cfg.AllDrives && runtime.GOOS != "windows" {
        return fmt.Errorf("--all-drives flag is only supported on Windows")
    }
    if cfg.OutputFormat != "json" && cfg.OutputFormat != "csv" {
        return fmt.Errorf("invalid output format: %s", cfg.OutputFormat)
    }
    if cfg.ConcurrencyLevel <= 0 {
        return fmt.Errorf("concurrency level must be positive")
    }
    if cfg.NiceLevel != "high" && cfg.NiceLevel != "medium" && cfg.NiceLevel != "low" {
        return fmt.Errorf("invalid nice level: %s", cfg.NiceLevel)
    }
    if cfg.LogLevel != "debug" && cfg.LogLevel != "info" && cfg.LogLevel != "warn" &&
        cfg.LogLevel != "error" && cfg.LogLevel != "fatal" && cfg.LogLevel != "panic" {
        return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
    }
    return nil
}

func parseCommaSeparated(input string) []string {
    if input == "" {
        return []string{}
    }
    items := strings.Split(input, ",")
    for i, item := range items {
        items[i] = strings.TrimSpace(item)
    }
    return items
}

func getIntFlagValue(f *flag.Flag) int {
    value, err := strconv.Atoi(f.Value.String())
    if err != nil {
        return 0
    }
    return value
}

func getInt64FlagValue(f *flag.Flag) int64 {
    value, err := strconv.ParseInt(f.Value.String(), 10, 64)
    if err != nil {
        return 0
    }
    return value
}

func parseBoolFlagValue(f *flag.Flag) bool {
    value, err := strconv.ParseBool(f.Value.String())
    if err != nil {
        return false
    }
    return value
}
