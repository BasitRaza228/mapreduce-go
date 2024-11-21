package coordinatorworker

import (
	"log"
	"mapreduce-go/src/types"
	"time"
)

func (c *Coordinator) RequestTask(req TaskRequest) TaskResponse {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.WorkerHealth[req.WorkerID] = time.Now()

	var task types.Task
	found := false

	if c.Phase == MapPhase {
		// Look for available map tasks
		for i := range c.MapTasks {
			if c.MapTasks[i].Status == types.Idle {
				c.MapTasks[i].Status = types.InProgress
				c.MapTasks[i].WorkerID = req.WorkerID
				c.MapTasks[i].StartTime = time.Now()
				task = c.MapTasks[i]
				found = true
				break
			}
		}

		// Check if all map tasks are completed
		if !found && c.allMapTasksCompleted() {
			c.Phase = ReducePhase
			log.Println("All map tasks completed. Moving to reduce phase.")
		}
	}

	if c.Phase == ReducePhase {
		// Look for available reduce tasks
		for i := range c.ReduceTasks {
			if c.ReduceTasks[i].Status == types.Idle {
				c.ReduceTasks[i].Status = types.InProgress
				c.ReduceTasks[i].WorkerID = req.WorkerID
				c.ReduceTasks[i].StartTime = time.Now()
				task = c.ReduceTasks[i]
				found = true
				break
			}
		}

		// Check if all reduce tasks are completed
		if !found && c.allReduceTasksCompleted() {
			c.Phase = DonePhase
			c.Done = true
			log.Println("All reduce tasks completed. MapReduce job finished.")
		}
	}

	if !found {
		if c.Done {
			task.Type = types.ExitTask
		} else {
			task.Type = types.WaitTask
		}
	}

	return TaskResponse{
		Task:        task,
		ReduceCount: c.ReduceCount,
		MapCount:    len(c.Files),
	}
}

func (c *Coordinator) TaskCompleted(task types.Task) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if task.Type == types.MapTask {
		for i := range c.MapTasks {
			if c.MapTasks[i].ID == task.ID && c.MapTasks[i].WorkerID == task.WorkerID {
				c.MapTasks[i].Status = types.Completed
				log.Printf("Map task %d completed by worker %s", task.ID, task.WorkerID)
				break
			}
		}
	} else if task.Type == types.ReduceTask {
		for i := range c.ReduceTasks {
			if c.ReduceTasks[i].ID == task.ID && c.ReduceTasks[i].WorkerID == task.WorkerID {
				c.ReduceTasks[i].Status = types.Completed
				log.Printf("Reduce task %d completed by worker %s", task.ID, task.WorkerID)
				break
			}
		}
	}
}

func (c *Coordinator) allMapTasksCompleted() bool {
	for _, task := range c.MapTasks {
		if task.Status != types.Completed {
			return false
		}
	}
	return true
}

func (c *Coordinator) allReduceTasksCompleted() bool {
	for _, task := range c.ReduceTasks {
		if task.Status != types.Completed {
			return false
		}
	}
	return true
}

func (c *Coordinator) monitorTasks() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.Mu.Lock()

		now := time.Now()

		// Check for timed out map tasks
		for i := range c.MapTasks {
			if c.MapTasks[i].Status == types.InProgress {
				if now.Sub(c.MapTasks[i].StartTime) > 30*time.Second {
					log.Printf("Map task %d timed out, reassigning", c.MapTasks[i].ID)
					c.MapTasks[i].Status = types.Idle
					c.MapTasks[i].WorkerID = ""
				}
			}
		}

		// Check for timed out reduce tasks
		for i := range c.ReduceTasks {
			if c.ReduceTasks[i].Status == types.InProgress {
				if now.Sub(c.ReduceTasks[i].StartTime) > 30*time.Second {
					log.Printf("Reduce task %d timed out, reassigning", c.ReduceTasks[i].ID)
					c.ReduceTasks[i].Status = types.Idle
					c.ReduceTasks[i].WorkerID = ""
				}
			}
		}

		c.Mu.Unlock()

		if c.Done {
			break
		}
	}
}

// returns whether the MapReduce job is complete
func (c *Coordinator) TaskDone() bool {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Done
}
