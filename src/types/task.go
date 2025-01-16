package types

import "time"

// represents one unit of work
type Task struct {
	ID        int
	Type      TaskType
	Filename  string
	MapNum    int
	ReduceNum int
	Status    TaskStatus
	WorkerID  string
	StartTime time.Time
}

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	WaitTask
	ExitTask
)

type TaskStatus int

const (
	Idle TaskStatus = iota
	InProgress
	Completed
	Failed
)
