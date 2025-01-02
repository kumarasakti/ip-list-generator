# IP List Generator

A high-performance command-line tool for generating IP addresses from CIDR notation and saving them to a file.

## Features

- Fast IP address generation from CIDR ranges
- Buffered file writing for optimal performance
- Progress tracking for large IP ranges
- Customizable output directory and filename
- Detailed execution summary with performance metrics
- Built-in path validation and error handling

## Usage
```bash
ip-list-generator -cidr 192.168.1.0/24 -filename list.txt -output <file directory>
```
### Available Flag
```bash  
  -cidr string
        CIDR range (e.g., 192.168.1.0/24)
  -filename string
        Custom filename (optional)
  -output string
        Output directory path
```
## Installation
Build from source code  
```bash
git clone https://github.com/kumarasakti/ip-list-generator.git
go build ip-list-generator.go
ip-list-generator --help
```
Use go get
```bash
go get github.com/yourusername/ip-list-generator
```