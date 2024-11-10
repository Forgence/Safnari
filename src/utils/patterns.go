package utils

import (
    "path/filepath"
    "regexp"
)

func ShouldInclude(path string, includePatterns, excludePatterns []string) bool {
    if len(includePatterns) > 0 && !matchesAnyPattern(path, includePatterns) {
        return false
    }
    if len(excludePatterns) > 0 && matchesAnyPattern(path, excludePatterns) {
        return false
    }
    return true
}

func matchesAnyPattern(path string, patterns []string) bool {
    for _, pattern := range patterns {
        matched, _ := filepath.Match(pattern, filepath.Base(path))
        if matched {
            return true
        }
        // Try regex matching
        if re, err := regexp.Compile(pattern); err == nil {
            if re.MatchString(path) {
                return true
            }
        }
    }
    return false
}