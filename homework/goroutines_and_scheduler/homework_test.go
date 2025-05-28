package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type TasksQueue struct {
	tasks []Task
	index map[int]int
}

func NewTasksQueue() *TasksQueue {
	return &TasksQueue{index: make(map[int]int)}
}

func (q *TasksQueue) Len() int {
	return len(q.tasks)
}

func (q *TasksQueue) Less(i, j int) bool {
	return q.tasks[i].Priority > q.tasks[j].Priority
}

func (q *TasksQueue) Swap(i, j int) {
	q.tasks[i], q.tasks[j] = q.tasks[j], q.tasks[i]
	q.index[q.tasks[i].Identifier] = i
	q.index[q.tasks[j].Identifier] = j
}

func (q *TasksQueue) Push(x any) {
	task := x.(Task)
	q.tasks = append(q.tasks, task)
	q.index[task.Identifier] = len(q.tasks) - 1
}

func (q *TasksQueue) Pop() any {
	task := q.tasks[len(q.tasks)-1]
	q.tasks = q.tasks[:len(q.tasks)-1]
	delete(q.index, task.Identifier)
	return task
}

func (q *TasksQueue) Update(identifier int, priority int) {
	index := q.index[identifier]
	q.tasks[index].Priority = priority
	heap.Fix(q, index)
}

type Scheduler struct {
	queue *TasksQueue
}

func NewScheduler() Scheduler {
	return Scheduler{queue: NewTasksQueue()}
}

func (s *Scheduler) AddTask(task Task) {
	heap.Push(s.queue, task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.queue.Update(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return heap.Pop(s.queue).(Task)
}

func TestScheduler(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)
	task1.Priority = 100

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
