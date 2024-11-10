package scanner

import "regexp"

var sensitivePatterns = map[string]*regexp.Regexp{
    "email":        regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`),
    "credit_card":  regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
    "ssn":          regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),
    "ip_address":   regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`),
    "api_key":      regexp.MustCompile(`(?i)(api_key|api-secret|access-token)[\s:=]+"?[\w\-]+"?`),
    "phone_number": regexp.MustCompile(`\b\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}\b`),
    // Add more patterns as needed
}

func GetPatterns(types []string) map[string]*regexp.Regexp {
    patterns := make(map[string]*regexp.Regexp)
    for _, t := range types {
        if pattern, exists := sensitivePatterns[t]; exists {
            patterns[t] = pattern
        }
    }
    return patterns
}
