package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

func main() {
	tasks := make(chan string)
	quit := make(chan struct{})
	currentWorkerSize := 3
	wg := &sync.WaitGroup{}
	workerSizeConfigFileName := "concurrency.txt"
	defer os.Remove(workerSizeConfigFileName)

	go func() {
		// gracefully shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		close(tasks)
	}()

	go updateWorkerSizeByConfigFile(workerSizeConfigFileName, wg, currentWorkerSize, tasks, quit)
	go produceTasks(tasks)
	consumeTasks(wg, currentWorkerSize, tasks, quit)
	wg.Wait()
}

func worker(wg *sync.WaitGroup, i int, tasks <-chan string, quit <-chan struct{}) {
	fmt.Println("Start worker ", i)
	defer wg.Done()
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			fmt.Printf("Worker %v, task: %v\n", i, task)
		case <-quit:
			fmt.Printf("quit worker %v\n", i)
			return
		}
	}

}

func consumeTasks(wg *sync.WaitGroup, currentWorkerSize int, tasks chan string, quit chan struct{}) {
	wg.Add(currentWorkerSize)
	for i := 1; i <= currentWorkerSize; i++ {
		go worker(wg, i, tasks, quit)
	}
}

func produceTasks(tasks chan string) {
	for i := 0; i <= 100; i++ {
		tasks <- fmt.Sprintf("task_%v", i)
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("close tasks")
	close(tasks)
}

func updateWorkerSizeByConfigFile(fileName string, wg *sync.WaitGroup, currentWorkerSize int, tasks chan string, quit chan struct{}) {
	for {
		time.Sleep(1 * time.Second)
		newWorkerSize, err := readWorkerSizeFromFile(fileName)
		if err != nil {
			fmt.Println(err)

			// Write current worker size to file
			if err := writeWorkerSizeToFile(fileName, fmt.Sprint(currentWorkerSize)); err != nil {
				fmt.Println(err)
				return
			}
			newWorkerSize = currentWorkerSize
		}

		if newWorkerSize == currentWorkerSize {
			continue
		} else if additionalSize := newWorkerSize - currentWorkerSize; additionalSize > 0 {
			fmt.Println("newWorkerSize:", newWorkerSize)
			wg.Add(additionalSize)
			for i := 0; i < additionalSize; i++ {
				go worker(wg, 100+i, tasks, quit)
			}
		} else {
			fmt.Println("newWorkerSize:", newWorkerSize)
			for i := 0; i < -additionalSize; i++ {
				quit <- struct{}{}
			}
		}
		currentWorkerSize = newWorkerSize
	}
}

func readWorkerSizeFromFile(fileName string) (newWorkerSize int, err error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	b := make([]byte, 1) // read only one byte
	f.Read(b)
	workerSizeStr := string(b)
	if workerSizeStr == "" {
		return 0, errors.New("empty WORKER_SIZE")
	}

	// fmt.Print("workerSizeStr:", workerSizeStr)
	size, err := strconv.Atoi(workerSizeStr)
	if err != nil {
		fmt.Println(err)
		return size, err
	}
	newWorkerSize = size
	// fmt.Println(" newWorkerSize: ", newWorkerSize)
	return
}

func writeWorkerSizeToFile(fileName, numConcurrency string) (err error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	_, err = f.WriteString(fmt.Sprint(numConcurrency))
	return err
}
