package coordinatorworker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"mapreduce-go/src/types"
	"os"
	"strings"
)

func (w *Worker) executeMapTask(task types.Task, coordinator *Coordinator) {
	log.Printf("Worker %s executing map task %d on file %s", w.Id, task.ID, task.Filename)

	// Read input file
	content, err := w.readFile(task.Filename)
	if err != nil {
		log.Printf("Error reading file %s: %v", task.Filename, err)
		return
	}

	// Execute map function
	kvs := w.MapFunc(task.Filename, content)

	// Create intermediate files
	intermediateFiles := make([]*os.File, w.ReduceCount)
	encoders := make([]*json.Encoder, w.ReduceCount)

	for i := 0; i < w.ReduceCount; i++ {
		filename := fmt.Sprintf("mr-%d-%d", task.MapNum, i)
		file, err := os.Create(filename)
		if err != nil {
			log.Printf("Error creating intermediate file %s: %v", filename, err)
			return
		}
		intermediateFiles[i] = file
		encoders[i] = json.NewEncoder(file)
	}

	// Partition and write key-value pairs
	for _, kv := range kvs {
		reduceNum := w.ihash(kv.Key) % w.ReduceCount
		err := encoders[reduceNum].Encode(&kv)
		if err != nil {
			log.Printf("Error encoding key-value pair: %v", err)
			return
		}
	}

	// Close all files
	for _, file := range intermediateFiles {
		file.Close()
	}

	// Notify coordinator of completion
	coordinator.TaskCompleted(task)
}

func (w *Worker) readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}

// generates a hash for key partitioning
func (w *Worker) ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}
