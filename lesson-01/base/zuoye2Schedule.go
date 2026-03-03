package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func() error

type TaskResult struct {
	Name     string
	Duration time.Duration
	Error    error
}

func runTask(name string, task Task) *TaskResult {
	start := time.Now()
	err := task()
	duration := time.Since(start)
	return &TaskResult{
		Name:     name,
		Duration: duration,
		Error:    err,
	}
}
func ScheduleTasks(tasks map[string]Task) []*TaskResult {
	var wg sync.WaitGroup
	results := make([]*TaskResult, 0, len(tasks))
	resultChin := make(chan *TaskResult, len(tasks))

	for name, task := range tasks {
		wg.Add(1)
		go func(n string, t Task) {
			defer wg.Done()
			resultChin <- runTask(n, t)
		}(name, task)
	}
	wg.Wait()
	close(resultChin)
	for result := range resultChin {
		results = append(results, result)
	}
	return results
}
func task1() error {
	time.Sleep(2 * time.Second)
	fmt.Println("Task 1 completed")
	return nil
}
func task2() error {
	time.Sleep(1 * time.Second)
	fmt.Println("Task 2 completed")
	return nil
}
func task3() error {
	time.Sleep(3 * time.Second)
	fmt.Println("Task 3 completed")
	return nil
}
func main() {
	tasks := map[string]Task{
		"Task 1": task1,
		"Task 2": task2,
		"Task 3": task3,
	}
	fmt.Println("Starting task execution...")
	results := ScheduleTasks(tasks)
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("Task %s failed: %v\n", res.Name, res.Error)
		} else {
			fmt.Printf("Task %s completed in %v\n", res.Name, res.Duration)
		}
	}
}
