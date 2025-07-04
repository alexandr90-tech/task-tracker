package task

import (
	"fmt"
	"sync"
)

type Task struct { // структура одной задачи
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Store struct { // хранилище задач в оперативной памяти (без БД)
	sync.Mutex        // чтобы защититься от одновременного доступа
	tasks      []Task // срез задач - само хранилище
	nextID     int    // следующий id для новой задачи
}

func NewStore() *Store { // конструктор создаёт новый пустой Store с пустым срезом задач и стартовым nextID = 1, возвращает указатель на Store
	return &Store{
		tasks:  []Task{},
		nextID: 1,
	}

}

func (s *Store) ListTasks() []Task { // возвращает список всех задач
	s.Lock()                               //блокирует доступ к данным пока мы читаем, чтобы никто не мог изменить в это время
	defer s.Unlock()                       // разблокируем автоматически после выхода из функции
	return append([]Task(nil), s.tasks...) // создаём копию среза и возвращаем, чтобы никто не мог изменить внутренний срез
}

func (s *Store) AddTask(title string) Task { // добавление новой задачи
	s.Lock()
	defer s.Unlock()

	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}

	s.tasks = append(s.tasks, task) // добавляем новую задачу в срез
	s.nextID++                      // увеличиваем ID для следующей задачи
	return task                     // возвращаем созданную задачу

}

func (s *Store) UpdateTaskDone(id int, done bool) (*Task, error) { // ищем задачу по ID, меняем статус на Done и возвращаем обновлённую задачу и флаг успешности
	s.Lock()
	defer s.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Done = done
			return &s.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// метод ищет задачу по ID и удаляет её из среза, если находит
func (s *Store) DeleteTask(id int) bool {
	s.Lock()
	defer s.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			// удаляем задачу из среза tasks
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}
