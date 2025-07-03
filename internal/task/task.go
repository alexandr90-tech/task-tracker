package task

import "sync"

type Task struct { // структура одной задачи
	ID int 'json:"id"'
	Tittle string 'json:"title"'
	Done bool 'json:"done"'
}

type Store struct { // хранилище задач в оперативной памяти (без БД)
	sync.Mutex
	tasks []Task
	nextID int
}

func NewStore() *Store {
	return &Store{
	tasks: []Task{},
	nextID: 1,
	}

}

func (s *Store) ListTasks() []Task {
s.Lock()
defer s.Unlock()
return append([]Task(nil), s.tasks...)
}

func (s *Store) AddTask(title string) Task{
s.Lock()
defer s.Unlock()


task := Task{
ID: s.NextID,
Title: title,
Done: false,
}

s.tasks = append(s.tasks, task)
s.nextID++
return task

}