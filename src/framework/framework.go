package framework

import (
	"context"
	"fmt"
	"log"
	coordinatorworker "mapreduce-go/src/coordinator_worker"
	types "mapreduce-go/src/types"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MapReduceFramework struct {
	coordinator *coordinatorworker.Coordinator
	workers     []*coordinatorworker.Worker
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewMapReduceFramework(inputFiles []string, nReduce int, nWorkers int, mapf types.MapFunc, reducef types.ReduceFunc) *MapReduceFramework {
	coordinator := coordinatorworker.NewCoordinator(inputFiles, nReduce)

	workers := make([]*coordinatorworker.Worker, nWorkers)
	for i := 0; i < nWorkers; i++ {
		workerID := fmt.Sprintf("worker-%d", i)
		workers[i] = coordinatorworker.NewWorker(workerID, mapf, reducef)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MapReduceFramework{
		coordinator: coordinator,
		workers:     workers,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (mrf *MapReduceFramework) Run() error {
	log.Println("Starting MapReduce framework")

	// Start all workers
	var wg sync.WaitGroup
	for _, worker := range mrf.workers {
		wg.Add(1)
		go func(w *coordinatorworker.Worker) {
			defer wg.Done()
			w.Run(mrf.ctx, mrf.coordinator)
		}(worker)
	}

	// Wait for completion or timeout
	done := make(chan struct{})
	go func() {
		for !mrf.coordinator.TaskDone() {
			time.Sleep(1 * time.Second)
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Println("MapReduce job completed successfully")
		mrf.cancel() // Signal workers to stop
		wg.Wait()    // Wait for all workers to finish
		return nil
	case <-time.After(5 * time.Minute): // 5 minute timeout
		log.Println("MapReduce job timed out")
		mrf.cancel()
		wg.Wait()
		return fmt.Errorf("job timed out")
	}
}

// remove intermediate files
func (mrf *MapReduceFramework) Cleanup() {
	files, _ := filepath.Glob("mr-*-*")
	for _, file := range files {
		if !strings.HasPrefix(filepath.Base(file), "mr-out-") {
			os.Remove(file)
		}
	}
}
