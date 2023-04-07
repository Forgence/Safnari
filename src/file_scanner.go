package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/glaslos/ssdeep"
)

type FileData struct {
	FilePath    string      `json:"file_path"`
	FileName    string      `json:"file_name"`
	Extension   string      `json:"extension"`
	Size        int64       `json:"size"`
	ModTime     time.Time   `json:"mod_time"`
	IsDir       bool        `json:"is_dir"`
	Permissions os.FileMode `json:"permissions"`
	MD5         string      `json:"md5,omitempty"`
	SHA1        string      `json:"sha1,omitempty"`
	SHA256      string      `json:"sha256,omitempty"`
	SSDEEP      string      `json:"ssdeep,omitempty"`
}

type FileScanResult struct {
	TotalFiles int        `json:"total_files"`
	FileList   []FileData `json:"file_list"`
}

type ProcessInfo struct {
	PID        int       `json:"pid"`
	Name       string    `json:"name"`
	Cmd        string    `json:"cmd"`
	CreateTime time.Time `json:"create_time"`
	RSS        uint64    `json:"rss"`
	VMS        uint64    `json:"vms"`
}

func traverseFiles(ctx context.Context, dirPath string, fileProcessor chan<- FileData, scanSubDirs, calculateHashes bool, excludePattern string, maxFileSize int64, fileTypes []string, hashAlgorithms []string, debug bool, errorLog *log.Logger) {
	var excludePatternRegex *regexp.Regexp
	if excludePattern != "" {
		excludePatternRegex = regexp.MustCompile(excludePattern)
	}

	hashAlgorithmsSet := make(map[string]bool)
	for _, alg := range hashAlgorithms {
		hashAlgorithmsSet[alg] = true
	}

	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			if debug {
				errorLog.Printf("Error accessing path %q: %v\n", filePath, err)
			}
			return err
		}

		// Skip processing directories if scanSubDirs is false and filePath is not the root directory
		if !scanSubDirs && info.IsDir() && filePath != dirPath {
			return filepath.SkipDir
		}

		// Apply exclude pattern, max file size, and file types filters
		if excludePatternRegex != nil && excludePatternRegex.MatchString(filePath) {
			return nil
		}
		if maxFileSize > 0 && info.Size() > maxFileSize {
			return nil
		}
		if len(fileTypes) > 0 && !contains(fileTypes, filepath.Ext(filePath)) {
			return nil
		}

		// Send only regular files to the fileProcessor channel
		if !info.IsDir() {
			fileData := FileData{
				FilePath:    filePath,
				FileName:    info.Name(),
				Extension:   filepath.Ext(filePath),
				Size:        info.Size(),
				ModTime:     info.ModTime(),
				IsDir:       info.IsDir(),
				Permissions: info.Mode(),
			}

			// Compute the file hashes if calculateHashes is true
			if calculateHashes {
				file, err := os.Open(filePath)
				if err != nil {
					if debug {
						errorLog.Printf("Error opening file %q: %v\n", filePath, err)
					}
				} else {
					if hashAlgorithmsSet["md5"] {
						md5Hash := md5.New()
						io.Copy(md5Hash, file)
						fileData.MD5 = hex.EncodeToString(md5Hash.Sum(nil))
					}
					if hashAlgorithmsSet["sha1"] {
						sha1Hash := sha1.New()
						io.Copy(sha1Hash, file)
						fileData.SHA1 = hex.EncodeToString(sha1Hash.Sum(nil))
					}
					if hashAlgorithmsSet["sha256"] {
						sha256Hash := sha256.New()
						io.Copy(sha256Hash, file)
						fileData.SHA256 = hex.EncodeToString(sha256Hash.Sum(nil))
					}
					if hashAlgorithmsSet["ssdeep"] {
						fileData.SSDEEP, _ = ssdeep.FuzzyFilename(filePath)
					}
					file.Close()
				}
			}

			// Send the processed file data to the fileProcessor channel
			select {
			case fileProcessor <- fileData:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	})

	if err != nil {
		if debug {
			errorLog.Printf("Error walking the path %q: %v\n", dirPath, err)
		}
	}
}

// Helper function to check if a slice contains a specific value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
