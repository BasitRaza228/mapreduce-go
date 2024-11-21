package coordinatorworker

import (
	"context"
	"log"
	"mapreduce-go/src/types"
	"time"
)

type Worker struct {
	Id          string
	MapFunc     types.MapFunc
	ReduceFunc  types.ReduceFunc
	ReduceCount int
	MapCpunt    int
}

type TaskRequest struct {
	WorkerID string
}

type TaskResponse struct {
	Task        types.Task
	ReduceCount int
	MapCount    int
}

func NewWorker(id string, mapf types.MapFunc, reducef types.ReduceFunc) *Worker {
	return &Worker{
		Id:         id,
		MapFunc:    mapf,
		ReduceFunc: reducef,
	}
}

func (w *Worker) Run(ctx context.Context, coordinator *Coordinator) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Request a task from coordinator
			req := TaskRequest{WorkerID: w.Id}
			resp := coordinator.RequestTask(req)

			w.ReduceCount = resp.ReduceCount
			w.MapCpunt = resp.MapCount

			switch resp.Task.Type {
			case types.MapTask:
				w.executeMapTask(resp.Task, coordinator)
			case types.ReduceTask:
				w.executeReduceTask(resp.Task, coordinator)
			case types.WaitTask:
				time.Sleep(1 * time.Second)
			case types.ExitTask:
				log.Printf("Worker %s exiting", w.Id)
				return
			}
		}
	}
}
