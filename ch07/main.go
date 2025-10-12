package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cube/manager"
	"cube/task"
	"cube/worker"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

const (
	DEFAULT_HOST = "localhost"
	DEFAULT_PORT = "5555"
)

func main() {
	host, defined := os.LookupEnv("CUBE_HOST")
	if !defined {
		log.Printf("CUBE_HOST not defined, falling back to default: %s", DEFAULT_HOST)
		host = DEFAULT_HOST
	}

	portStr, defined := os.LookupEnv("CUBE_PORT")
	if !defined {
		log.Printf("CUBE_PORT not defined, falling back to default: %s", DEFAULT_PORT)
		portStr = DEFAULT_PORT
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid CUBE_PORT value: %v", err)
	}

	fmt.Println("Starting Cube worker")
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}
	go runTasks(&w)
	go w.CollectStats()
	go api.Start()

	workers := []string{fmt.Sprintf("%s:%d", host, port)}
	m := manager.New(workers)

	for i := 0; i < 3; i++ {
		t := task.Task{
			ID:    uuid.New(),
			Name:  fmt.Sprintf("test-container-%d", i),
			State: task.Scheduled,
			Image: "strm/helloworld-http",
		}
		te := task.TaskEvent{
			ID:    uuid.New(),
			State: task.Running,
			Task:  t,
		}
		m.AddTask(te)
		m.SendWork()
	}

	go func() {
		for {
			fmt.Printf("[Manager] Updating tasks from %d workers\n", len(m.Workers))
			m.UpdateTasks()
			time.Sleep(15 * time.Second)
		}
	}()

	for {
		for _, t := range m.TaskDb {
			fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
			time.Sleep(15 * time.Second)
		}
	}
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
