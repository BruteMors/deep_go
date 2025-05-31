package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type wrappedTask struct {
	task     Task
	priority int
	index    int
}

type taskHeap []*wrappedTask

func (h taskHeap) Len() int           { return len(h) }
func (h taskHeap) Less(i, j int) bool { return h[i].priority > h[j].priority }
func (h taskHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i]; h[i].index, h[j].index = i, j }
func (h *taskHeap) Push(x any)        { w := x.(*wrappedTask); w.index = len(*h); *h = append(*h, w) }
func (h *taskHeap) Pop() any          { old := *h; n := len(old); w := old[n-1]; *h = old[:n-1]; return w }

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	h        taskHeap
	idToTask map[int]*wrappedTask
}

func NewScheduler() Scheduler {
	return Scheduler{
		h:        make(taskHeap, 0),
		idToTask: make(map[int]*wrappedTask),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, exists := s.idToTask[task.Identifier]; exists {
		return
	}
	w := &wrappedTask{
		task:     task,
		priority: task.Priority,
	}
	heap.Push(&s.h, w)
	s.idToTask[task.Identifier] = w
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if w, ok := s.idToTask[taskID]; ok {
		w.priority = newPriority
		heap.Fix(&s.h, w.index)
	}
}

func (s *Scheduler) GetTask() Task {
	if len(s.h) == 0 {
		return Task{}
	}
	w := heap.Pop(&s.h).(*wrappedTask)
	delete(s.idToTask, w.task.Identifier)
	return w.task
}

func TestTrace(t *testing.T) {
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

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
