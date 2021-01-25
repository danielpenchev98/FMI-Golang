package main

import (
	"fmt"
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
	/*var (
		wg  sync.WaitGroup
		res int
		err error
	)

	err = checkTasks(c.Tasks)
	if err != nil {
		return 0, err
	}

	channel := make(chan Answer,len(c.Tasks))

	for _, task := range c.Tasks {
		wg.Add(1)
		go ParallelTaskExecute(task,arg,channel,&wg)
	}

	wg.Wait()
	tempResults := make([]int,len(c.Tasks))
	for _, ans := channel*/
	return 0, nil
}

func ConcurrentMapReduce(reduce func(results []int) int, tasks ...Task) Task {
	return &ConcurrentMapReduceTask{
		Tasks:  tasks,
		Reduce: reduce,
	}
}

func main() {
	var (
		res int
		err error
	)

	if res, err = Pipeline(adder{20}, adder{10}, adder{-50}).Execute(100); err != nil {
		fmt.Printf("The pipeline returned an error\n")
	} else {
		fmt.Printf("The pipeline returned %d\n", res)
	}

	f := Fastest(
		lazyAdder{adder{20}, 500},
		lazyAdder{adder{50}, 300},
		adder{150},
	)
	if res, err = f.Execute(1); err != nil {
		fmt.Printf("The fastest task returned an error\n")
	} else {
		fmt.Printf("The fastest task returned %d\n", res)
	}

	_, err = Timed(lazyAdder{adder{20}, 50}, 2*time.Millisecond).Execute(2)
	if err != nil {
		fmt.Printf("The first time task returned an error. Reason %v\n", err)
	}
	res, err = Timed(lazyAdder{adder{20}, 50}, 300*time.Millisecond).Execute(2)
	if err != nil {
		fmt.Printf("The fastest task returned an error. Reason %v\n", err)
	} else {
		fmt.Printf("The fastest task returned %d\n", res)
	}

}
