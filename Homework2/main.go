package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

//Task - represents a task that needs to be done
type Task interface {
	Execute(int) (int, error)
}

//PipelineTask - task involving the execution of tasks in a pipeline fashion
type PipelineTask struct {
	Tasks []Task
}

type adder struct {
	augend int
}

func (a adder) Execute(addend int) (int, error) {
	result := a.augend + addend
	if result > 127 {
		return 0, fmt.Errorf("Result %d exceeds the adder threshold", a)
	}
	return result, nil
}

//Execute - executes the pipeline
//error is returned if tasks are invalid or a tasks fails during its execution
func (p *PipelineTask) Execute(arg int) (int, error) {
	var (
		result = arg
		err    error
	)

	err = checkTasks(p.Tasks)
	if err != nil {
		return 0, err
	}

	for index, task := range p.Tasks {
		result, err = task.Execute(result)
		if err != nil {
			return 0, fmt.Errorf("Task #%d failed during its execution. Reason : %v", index, err)
		}
	}

	return result, err
}

//Pipeline - creates a pipeline of given tasks
func Pipeline(tasks ...Task) Task {
	return &PipelineTask{
		Tasks: tasks,
	}
}

type FastestTask struct {
	Tasks []Task
}

type Answer struct {
	Result int
	Err    error
}

func (f *FastestTask) Execute(arg int) (int, error) {
	err := checkTasks(f.Tasks)
	if err != nil {
		return 0, err
	}

	var wg sync.WaitGroup
	channel := make(chan Answer, len(f.Tasks))
	for _, task := range f.Tasks {
		wg.Add(1)
		go ParallelTaskExecute(task, arg, channel, &wg)
	}

	wg.Wait()
	answer := <-channel
	close(channel)
	return answer.Result, answer.Err
}

func ParallelTaskExecute(task Task, arg int, c chan<- Answer, wg *sync.WaitGroup) {
	defer wg.Done()
	result, err := task.Execute(arg)
	c <- Answer{
		Result: result,
		Err:    err,
	}
}

func Fastest(tasks ...Task) Task {
	return &FastestTask{
		Tasks: tasks,
	}
}

type lazyAdder struct {
	adder
	delay time.Duration
}

func (la lazyAdder) Execute(addend int) (int, error) {
	time.Sleep(la.delay * time.Millisecond)
	return la.adder.Execute(addend)
}

func checkTasks(tasks []Task) error {
	if len(tasks) == 0 {
		return fmt.Errorf("No tasks were given")
	}

	for index, task := range tasks {
		if task == nil {
			return fmt.Errorf("Task #%d is nil", index)
		}
	}
	return nil
}

type TimedTask struct {
	Task    Task
	Timeout time.Duration
}

func (t *TimedTask) Execute(arg int) (int, error) {
	var (
		result int
		err    error
		wg     sync.WaitGroup
	)

	err = checkTasks([]Task{t.Task})
	if err != nil {
		return 0, err
	}

	c := make(chan Answer, 1)
	wg.Add(1)
	go ParallelTaskExecute(t.Task, arg, c, &wg)

	select {
	case ans := <-c:
		result, err = ans.Result, ans.Err
	case <-time.After(t.Timeout):
		result, err = 0, fmt.Errorf("Task timeout")
	}

	wg.Wait()
	close(c)
	return result, err
}

func Timed(task Task, timeout time.Duration) Task {
	return &TimedTask{
		Task:    task,
		Timeout: timeout,
	}
}

type ConcurrentMapReduceTask struct {
	Tasks  []Task
	Reduce func(results []int) int
}

func (c *ConcurrentMapReduceTask) Execute(arg int) (int, error) {
	if err := checkTasks(c.Tasks); err != nil {
		return 0, err
	}

	prodCh := make(chan Answer)
	var producers sync.WaitGroup
	for _, task := range c.Tasks {
		producers.Add(1)
		go ParallelTaskExecute(task, arg, prodCh, &producers)
	}

	consumerCh := make(chan Answer)
	var consumer sync.WaitGroup
	consumer.Add(1)

	go func(output chan<- Answer, input <-chan Answer, wg *sync.WaitGroup) {
		var (
			result int
			err    error
		)
		tempResults := make([]int, 0, len(c.Tasks))
		for ans := range input {
			if ans.Err != nil {
				err = ans.Err
			}
			tempResults = append(tempResults, ans.Result)
		}
		if err == nil {
			result = c.Reduce(tempResults)
		}

		output <- Answer{
			Result: result,
			Err:    err,
		}
		wg.Done()
	}(consumerCh, prodCh, &consumer)

	producers.Wait()
	close(prodCh)
	ans := <-consumerCh
	close(consumerCh)
	return ans.Result, ans.Err
}

func ConcurrentMapReduce(reduce func(results []int) int, tasks ...Task) Task {
	return &ConcurrentMapReduceTask{
		Tasks:  tasks,
		Reduce: reduce,
	}
}

type GreatSearcherTask struct {
	MaxErrors     int
	IncomingTasks <-chan Task
}

func (g *GreatSearcherTask) Execute(arg int) (int, error) {
	incomingResults := make(chan Answer)
	resultChan := make(chan Answer)
	var reducer sync.WaitGroup
	reducer.Add(1)
	go func(output chan<- Answer, input <-chan Answer, wg *sync.WaitGroup) {
		var (
			err       error
			errors    int
			maxResult = math.MinInt64
		)

		for ans := range input {
			if ans.Err != nil {
				errors++
			} else if maxResult < ans.Result {
				maxResult = ans.Result
			}
		}
		if errors > g.MaxErrors {
			err = fmt.Errorf("Max allowed errors exceeded")
		}

		output <- Answer{
			Result: maxResult,
			Err:    err,
		}
		wg.Done()
	}(resultChan, incomingResults, &reducer)

	var producer sync.WaitGroup
	for task := range g.IncomingTasks {
		if task == nil {
			continue //or also count it as an error
		}
		producer.Add(1)
		go ParallelTaskExecute(task, arg, incomingResults, &producer)
	}

	producer.Wait()
	close(incomingResults)
	ans := <-resultChan
	reducer.Wait()
	close(resultChan)
	return ans.Result, ans.Err
}

func GreatestSearcher(errorLimit int, tasks <-chan Task) Task {
	return &GreatSearcherTask{
		MaxErrors:     errorLimit,
		IncomingTasks: tasks,
	}
}

func main() {
	var (
		res int
		err error
	)

	//Task1
	if res, err = Pipeline(adder{50}, adder{60}).Execute(10); err != nil {
		fmt.Println("Error: Task1 successfull scenario failure")
	} else {
		fmt.Printf("Task1 returned %d\n", res)
	}

	if res, err = Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100); err != nil {
		fmt.Printf("Task1 returned an error\n")
	} else {
		fmt.Println("Error: Task1 unsuccessfull scenario failure")
	}

	//Task2
	f := Fastest(
		lazyAdder{adder{20}, 500},
		lazyAdder{adder{50}, 300},
		adder{41},
	)

	if res, err = f.Execute(1); err != nil {
		fmt.Println("Error: Task2 unsuccessfull scenario failure")
	} else {
		fmt.Printf("Task2 returned %d\n", res)
	}

	//Task3
	_, err = Timed(lazyAdder{adder{20}, 50}, 2*time.Millisecond).Execute(2)
	if err != nil {
		fmt.Printf("Task3 returned an error. Reason %v\n", err)
	} else {
		fmt.Printf("The unsuccessful scenario of task3 failed")
	}

	res, err = Timed(lazyAdder{adder{20}, 50}, 300*time.Millisecond).Execute(2)
	if err != nil {
		fmt.Printf("The successfull scenario of task3 failed")
	} else {
		fmt.Printf("The fastest task returned %d\n", res)
	}

	//Task4
	reduce := func(results []int) int {
		smallest := 128
		for _, v := range results {
			if v < smallest {
				smallest = v
			}
		}
		return smallest
	}

	mr := ConcurrentMapReduce(reduce, adder{30}, adder{50}, adder{20})
	if res, err := mr.Execute(5); err != nil {
		fmt.Printf("The successfull sceanrio of task4 failed\n")
	} else {
		fmt.Printf("The ConcurrentMapReduce returned %d\n", res)
	}

	//Task5
	tasks := make(chan Task)
	gs := GreatestSearcher(2, tasks)

	go func() {
		tasks <- adder{4}
		tasks <- nil
		tasks <- lazyAdder{adder{22}, 20}
		tasks <- adder{125}
		time.Sleep(50 * time.Millisecond)
		tasks <- adder{32}

		tasks <- Timed(lazyAdder{adder{100}, 2000}, 20*time.Millisecond)
		//tasks <- adder{127}
		close(tasks)
	}()

	res, err = gs.Execute(10)
	if err != nil {
		fmt.Printf("The fastest task returned an error. Reason %v\n", err)
	} else {
		fmt.Printf("The fastest task returned %d\n", res)
	}
}
