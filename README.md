# Safnari: File and System Information Gatherer

Safnari is a versatile and comprehensive tool designed to gather information about files and system information from a host machine. The tool scans the specified directory, collects metadata about the files, and retrieves system information such as host details, running processes, and CPU usage. Safnari provides several command-line flags to configure the scanning process, such as filtering file types, excluding patterns, and calculating file hashes.

## Features

- Gather host information such as OS details, uptime, and hostname
- List running processes and their details (PID, name, memory usage, etc.)
- Measure CPU usage percentage
- Scan files in a specified directory with optional recursion into subdirectories
- Filter files based on file types, maximum file size, and exclusion patterns
- Calculate file hashes (MD5, SHA1, SHA256, and SSDeep)
- Output results in JSON format to a file

## Installation

To install Safnari, you can either clone the repository and build from source or download the latest pre-compiled binary from the GitHub releases page.

### Build from Source

\`\`\`sh
git clone https://github.com/Forgence/Safnari.git
cd Safnari/src
go build -o safnari
\`\`\`

### Download Pre-Compiled Binary

Check the [releases page](https://github.com/Forgence/Safnari/releases) for the latest pre-compiled binaries for your operating system.

## Usage

\`\`\`
Usage of safnari:
  -start-dir string
        Starting directory for file scanning (default ".")
  -debug
        Print debug logs
  -scan-subdirs
        Scan subdirectories
  -hashes
        Calculate MD5, SHA1, SHA256, and SSDeep hashes for files
  -exclude-pattern string
        Pattern to exclude files/directories from scanning
  -max-file-size int
        Maximum file size in bytes to be included in the scan (default 0)
  -file-types string
        Comma-separated list of file extensions to include in the scan
  -hash-algorithms string
        Comma-separated list of hash algorithms to use (default "md5,sha1,sha256,ssdeep")
  -timeout duration
        Timeout duration for the scanning process (e.g., 30s, 5m) (default 0)
  -output string
        Output file for the JSON data (default "file_data.json")
\`\`\`

## Example

To scan the current directory and calculate file hashes, run the following command:

\`\`\`sh
./safnari -hashes -start-dir ./
\`\`\`

This will output the scan results to `file_data.json` by default.

## Contributing

Contributions to Safnari are always welcome! Feel free to open issues or submit pull requests to help improve the project.

## License

Safnari is released under the [MIT License](LICENSE).