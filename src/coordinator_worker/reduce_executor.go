package coordinatorworker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mapreduce-go/src/types"
	"os"
	"sort"
)

func (w *Worker) executeReduceTask(task types.Task, coordinator *Coordinator) {
	log.Printf("Worker %s executing reduce task %d", w.Id, task.ID)

	// Read all intermediate files for this reduce partition
	intermediate := []types.KeyValue{}

	for mapNum := 0; mapNum < w.MapCpunt; mapNum++ {
		filename := fmt.Sprintf("mr-%d-%d", mapNum, task.ReduceNum)
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("Error opening intermediate file %s: %v", filename, err)
			continue
		}

		decoder := json.NewDecoder(file)
		for {
			var kv types.KeyValue
			if err := decoder.Decode(&kv); err == io.EOF {
				break
			} else if err != nil {
				log.Printf("Error decoding intermediate file %s: %v", filename, err)
				break
			}
			intermediate = append(intermediate, kv)
		}
		file.Close()
	}

	// Sort intermediate key-value pairs by key
	sort.Slice(intermediate, func(i, j int) bool {
		return intermediate[i].Key < intermediate[j].Key
	})

	// Create output file
	outputFilename := fmt.Sprintf("mr-out-%d", task.ReduceNum)
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Printf("Error creating output file %s: %v", outputFilename, err)
		return
	}
	defer outputFile.Close()

	// Group by key and apply reduce function
	i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}

		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}

		output := w.ReduceFunc(intermediate[i].Key, values)
		fmt.Fprintf(outputFile, "%v %v\n", intermediate[i].Key, output)

		i = j
	}

	// Notify coordinator of completion
	coordinator.TaskCompleted(task)
}
