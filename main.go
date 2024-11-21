package main

import (
	"fmt"
	"log"
	"mapreduce-go/src/framework"
	"mapreduce-go/src/mapfuncs"
	"os"
	"path/filepath"
)

func main() {

	// input files
	inputFiles := []string{"input_one.txt", "input_two.txt", "input_three.txt"}

	log.Println("\n=== Running Word Count Example ===")
	countframework := framework.NewMapReduceFramework(inputFiles, 3, 4, mapfuncs.WordCountMapper, mapfuncs.WordCountReducer)

	if err := countframework.Run(); err != nil {
		log.Fatalf("MapReduce job failed: %v", err)
	}

	// results
	log.Println("\nWord Count Results:")
	for i := 0; i < 3; i++ {
		filename := fmt.Sprintf("mr-out-%d", i)
		if content, err := os.ReadFile(filename); err == nil {
			fmt.Printf("=== Output file %s ===\n%s\n", filename, string(content))
		}
	}

	// Cleanup intermediate files
	countframework.Cleanup()

	// Run Grep example
	log.Println("\n=== Running Grep Example (searching for 'Hello') ===")
	grepFramework := framework.NewMapReduceFramework(inputFiles, 2, 3, mapfuncs.GrepMapper("Hello"), mapfuncs.GrepReducer)

	if err := grepFramework.Run(); err != nil {
		log.Fatalf("Grep MapReduce job failed: %v", err)
	}

	// Display grep results
	log.Println("\nGrep Results:")
	for i := 0; i < 2; i++ {
		filename := fmt.Sprintf("mr-out-%d", i)
		if content, err := os.ReadFile(filename); err == nil {
			fmt.Printf("=== Output file %s ===\n%s\n", filename, string(content))
		}
	}

	// Final cleanup
	grepFramework.Cleanup()

	// Clean up sample input files
	for _, file := range inputFiles {
		os.Remove(file)
	}

	// Clean up output files
	outputs, _ := filepath.Glob("mr-out-*")
	for _, file := range outputs {
		os.Remove(file)
	}

	log.Println("\nMapReduce framework demonstration completed successfully!")
}
