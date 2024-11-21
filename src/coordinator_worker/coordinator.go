package coordinatorworker

import (
	"mapreduce-go/src/types"
	"sync"
	"time"
)

type Coordinator struct {
	Mu           sync.RWMutex
	Files        []string
	ReduceCount  int
	MapTasks     []types.Task
	ReduceTasks  []types.Task
	Phase        Phase
	Done         bool
	TaskCounter  int
	WorkerHealth map[string]time.Time
}

// represents the current phase of MapReduce execution
type Phase int

const (
	MapPhase Phase = iota
	ReducePhase
	DonePhase
)

// creates a new coordinator instance
func NewCoordinator(files []string, nReduce int) *Coordinator {
	c := &Coordinator{
		Files:        files,
		ReduceCount:  nReduce,
		MapTasks:     make([]types.Task, len(files)),
		ReduceTasks:  make([]types.Task, nReduce),
		Phase:        MapPhase,
		Done:         false,
		TaskCounter:  0,
		WorkerHealth: make(map[string]time.Time),
	}

	// Initialize map tasks
	for i, file := range files {
		c.MapTasks[i] = types.Task{
			ID:       c.TaskCounter,
			Type:     types.MapTask,
			Filename: file,
			MapNum:   i,
			Status:   types.Idle,
		}
		c.TaskCounter++
	}

	// Initialize reduce tasks
	for i := 0; i < nReduce; i++ {
		c.ReduceTasks[i] = types.Task{
			ID:        c.TaskCounter,
			Type:      types.ReduceTask,
			ReduceNum: i,
			Status:    types.Idle,
		}
		c.TaskCounter++
	}

	// Start background goroutine to monitor worker health and task timeouts
	go c.monitorTasks()

	return c
}
