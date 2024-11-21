# MapReduce Framework in Go

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/yourusername/mapreduce-framework)

A comprehensive, production-ready MapReduce framework implementation in Go, designed for distributed processing of large datasets with fault tolerance and scalability.

## 🚀 Features

- **Complete MapReduce Implementation**: Full coordinator-worker architecture
- **Fault Tolerance**: Automatic task reassignment and worker failure handling
- **Scalable Design**: Support for multiple workers
- **Thread-Safe Operations**: Concurrent processing with proper synchronization
- **Extensible Framework**: Easy to implement custom map and reduce functions
- **Built-in Examples**: Word count and grep, implementations
- **Production Ready**: Comprehensive error handling and monitoring
- **Clean Architecture**: Modular design following Go best practices


## 🏃 Quick Start

```bash
# Clone the repository
git clone https://github.com/BasitRaza228/MapReduce-golang
cd MapReduce-golang

# Run the demo
go run main.go

```

### Output Example
```
=== Word Count Results ===
hello 15
world 8
mapreduce 12
framework 6
golang 4
```

## 🏗️ Architecture

The framework implements the classic MapReduce paradigm with the following components:

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│ Coordinator │◄──►│   Worker 1   │    │   Input     │
│             │    │              │◄───┤   Files     │
│  - Task Mgmt│    │ - Map Tasks  │    │             │
│  - Scheduling│   │ - Reduce Tasks│   └─────────────┘
│  - Fault Tol.│   └──────────────┘
└─────────────┘           │
       │                  │
       ▼                  ▼
┌──────────────┐    ┌─────────────┐
│   Worker 2   │    │ Intermediate│
│              │    │   Files     │
│ - Map Tasks  │───►│ (Partitioned)│
│ - Reduce Tasks│   └─────────────┘
└──────────────┘           │
                          ▼
                   ┌─────────────┐
                   │   Output    │
                   │   Files     │
                   └─────────────┘
```

### Key Components

- **Coordinator**: Manages task distribution, worker health, and job lifecycle
- **Workers**: Execute map and reduce tasks independently
- **Task Management**: Automatic task assignment, timeout handling, and reassignment
- **Fault Tolerance**: Worker failure detection and task redistribution
- **Data Partitioning**: Intelligent key-based partitioning for load balancing

## 💻 Installation

### Prerequisites
- Go 1.22 or later
- Git

### Install from Source
```bash
git clone https://github.com/BasitRaza228/MapReduce-golang
cd MapReduce-golang
go mod download
go run main.go
```

### Using Go Modules
```go
// go.mod
module your-project

require github.com/yourusername/mapreduce-framework v1.0.0
```

## 📖 Usage

### Basic Usage

```go
package main

import (
    "github.com/yourusername/mapreduce-framework/pkg/framework"
    "github.com/yourusername/mapreduce-framework/pkg/mapfuncs"
)

func main() {
    // Define input files
    inputFiles := []string{"input1.txt", "input2.txt", "input3.txt"}
    
    // Create framework instance
    framework := framework.NewMapReduceFramework(
        inputFiles,
        3, // number of reduce tasks
        4, // number of workers
        mapfuncs.WordCountMapper,
        mapfuncs.WordCountReducer,
    )
    
    // Execute MapReduce job
    if err := framework.Run(); err != nil {
        log.Fatal("MapReduce job failed:", err)
    }
    
    // Cleanup intermediate files
    framework.Cleanup()
}
```

### Custom Map and Reduce Functions

```go
// Custom mapper function
func MyMapper(filename string, contents string) []types.KeyValue {
    // Your custom mapping logic here
    var kvs []types.KeyValue
    // ... process contents ...
    return kvs
}

// Custom reducer function
func MyReducer(key string, values []string) string {
    // Your custom reduction logic here
    // ... aggregate values ...
    return result
}

// Use custom functions
framework := framework.NewMapReduceFramework(
    inputFiles, 3, 4, MyMapper, MyReducer,
)
```

## 🎯 Examples

### 1. Word Count
```bash
go run main.go
```

Count occurrences of each word in text files.

### 2. Grep Pattern Matching
```bash
go run main.go"
```

## 📚 API Reference

### Core Types

```go
// KeyValue represents a key-value pair
type KeyValue struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// MapFunc defines the map function signature
type MapFunc func(filename string, contents string) []KeyValue

// ReduceFunc defines the reduce function signature
type ReduceFunc func(key string, values []string) string
```

### Framework Methods

```go
// NewMapReduceFramework creates a new framework instance
func NewMapReduceFramework(
    inputFiles []string,
    nReduce int,
    nWorkers int,
    mapf MapFunc,
    reducef ReduceFunc,
) *MapReduceFramework

// Run executes the MapReduce job
func (mrf *MapReduceFramework) Run() error

// Cleanup removes intermediate files
func (mrf *MapReduceFramework) Cleanup()
```

### Built-in Map/Reduce Functions

| Function | Description | Use Case |
|----------|-------------|----------|
| `WordCountMapper/Reducer` | Count word occurrences | Text analysis |
| `GrepMapper/Reducer` | Pattern matching | Log searching |

### Scaling Guidelines

- **Small files (<100MB)**: 2-4 workers, 2-3 reduce tasks
- **Medium files (100MB-1GB)**: 4-8 workers, 3-5 reduce tasks  
- **Large files (>1GB)**: 8-16 workers, 5-10 reduce tasks
- **Very large datasets**: Consider distributed deployment

## 🔧 Advanced Features

### Fault Tolerance
- Automatic task reassignment on worker failures
- Health monitoring and recovery
- Graceful shutdown handling


### Contribution Areas
- 🐛 Bug fixes and improvements
- ✨ New map/reduce function implementations
- 📚 Documentation enhancements
- ⚡ Performance optimizations
- 🧪 test cases
- 🔧 DevOps and deployment

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Google's original MapReduce paper (Dean & Ghemawat, 2004)
- Built following Go community best practices

## 📞 Support

- 📖 [Documentation](https://github.com/BasitRaza228/MapReduce-golang/)
- 🐛 [Issue Tracker](https://github.com/BasitRaza228/MapReduce-golang/issues)
- 📧 Email: basitraza228@gmail.com

⭐ **Star this repository if you find it helpful!**

Made with ❤️ for the distributed computing community.