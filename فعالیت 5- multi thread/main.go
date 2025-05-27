package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numWorkers = 1000
	numTasks   = 1000
)

// SimpleThreadPool represents a basic worker pool implementation
type SimpleThreadPool struct {
	wg            sync.WaitGroup
	workerChan    chan struct{}
	completedTasks int64
}

// NewSimpleThreadPool creates a new simple thread pool
func NewSimpleThreadPool(numWorkers int) *SimpleThreadPool {
	return &SimpleThreadPool{
		workerChan: make(chan struct{}, numWorkers),
	}
}

// ExecuteTasks runs the specified number of tasks
func (p *SimpleThreadPool) ExecuteTasks() {
	for i := 0; i < numTasks; i++ {
		p.wg.Add(1)
		p.workerChan <- struct{}{} // Acquire worker slot
		go func(taskID int) {
			defer p.wg.Done()
			defer func() { <-p.workerChan }() // Release worker slot
			
			// Simulate work
			time.Sleep(100 * time.Millisecond)
			p.completedTasks++
		}(i)
	}
}

// WaitForCompletion waits for all tasks to complete
func (p *SimpleThreadPool) WaitForCompletion() {
	p.wg.Wait()
}

// GetCompletedTasks returns the number of completed tasks
func (p *SimpleThreadPool) GetCompletedTasks() int64 {
	return p.completedTasks
}

// ApacheThreadPool represents a more sophisticated worker pool implementation
type ApacheThreadPool struct {
	wg            sync.WaitGroup
	workerPool    chan *Worker
	completedTasks int64
}

// Worker represents a worker in the pool
type Worker struct {
	ID int
}

// NewApacheThreadPool creates a new Apache-style thread pool
func NewApacheThreadPool(numWorkers int) *ApacheThreadPool {
	pool := &ApacheThreadPool{
		workerPool: make(chan *Worker, numWorkers),
	}
	
	// Initialize worker pool
	for i := 0; i < numWorkers; i++ {
		pool.workerPool <- &Worker{ID: i}
	}
	
	return pool
}

// ExecuteTasks runs the specified number of tasks
func (p *ApacheThreadPool) ExecuteTasks() {
	for i := 0; i < numTasks; i++ {
		p.wg.Add(1)
		go func(taskID int) {
			defer p.wg.Done()
			
			// Get worker from pool
			worker := <-p.workerPool
			defer func() { p.workerPool <- worker }() // Return worker to pool
			
			// Simulate work
			time.Sleep(100 * time.Millisecond)
			p.completedTasks++
		}(i)
	}
}

// WaitForCompletion waits for all tasks to complete
func (p *ApacheThreadPool) WaitForCompletion() {
	p.wg.Wait()
}

// GetCompletedTasks returns the number of completed tasks
func (p *ApacheThreadPool) GetCompletedTasks() int64 {
	return p.completedTasks
}

func main() {
	// Benchmark Simple Thread Pool
	start := time.Now()
	simplePool := NewSimpleThreadPool(numWorkers)
	simplePool.ExecuteTasks()
	simplePool.WaitForCompletion()
	simpleDuration := time.Since(start)
	
	fmt.Printf("Simple Thread Pool Results:\n")
	fmt.Printf("Completed Tasks: %d\n", simplePool.GetCompletedTasks())
	fmt.Printf("Total Time: %v\n\n", simpleDuration)
	
	// Benchmark Apache Thread Pool
	start = time.Now()
	apachePool := NewApacheThreadPool(numWorkers)
	apachePool.ExecuteTasks()
	apachePool.WaitForCompletion()
	apacheDuration := time.Since(start)
	
	fmt.Printf("Apache Thread Pool Results:\n")
	fmt.Printf("Completed Tasks: %d\n", apachePool.GetCompletedTasks())
	fmt.Printf("Total Time: %v\n\n", apacheDuration)
	
	// Calculate and display performance difference
	diff := float64(apacheDuration-simpleDuration) / float64(simpleDuration) * 100
	fmt.Printf("Performance Difference: %.2f%%\n", diff)
}