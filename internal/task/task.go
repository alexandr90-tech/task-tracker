package task

import (
	"encoding/json"
	"fmt"
	"os"
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
	file       string //файл для сохранения задач
}

func NewStore(filename string) *Store {
	s := &Store{
		tasks:  []Task{},
		nextID: 1,
		file:   filename,
	}
	err := s.Load()
	if err != nil {
		//если файл не существует, то ошибки можно игнорировать
		if !os.IsNotExist(err) {
			fmt.Printf("Ошибка загрузки задач: %v\n", err)
		}
	}
	return s

}

// Load загружает задачи из JSON файла
func (s *Store) Load() error {
	s.Lock()
	defer s.Unlock()

	file, err := os.Open(s.file)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&s.tasks); err != nil {
		return err
	}

	//выставляем nextID на max+1
	maxID := 0
	for _, t := range s.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	s.nextID = maxID + 1

	return nil
}

// Save сохраняет текущие задачи в JSON файл
func (s *Store) Save() error {
	s.Lock()
	defer s.Unlock()

	file, err := os.Create(s.file)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(s.tasks)
}

func (s *Store) ListTasks() []Task { // возвращает список всех задач
	s.Lock()                               //блокирует доступ к данным пока мы читаем, чтобы никто не мог изменить в это время
	defer s.Unlock()                       // разблокируем автоматически после выхода из функции
	return append([]Task(nil), s.tasks...) // создаём копию среза и возвращаем, чтобы никто не мог изменить внутренний срез
}

func (s *Store) AddTask(title string) Task { // добавление новой задачи
	s.Lock()
	task := Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}

	s.tasks = append(s.tasks, task) // добавляем новую задачу в срез
	s.nextID++                      // увеличиваем ID для следующей задачи

	s.Unlock()
	//сохраняем изменения в файл
	_ = s.Save() //вызов Save без блокировки
	return task  // возвращаем созданную задачу

}

func (s *Store) UpdateTaskDone(id int, done bool) (*Task, error) { // ищем задачу по ID, меняем статус на Done и возвращаем обновлённую задачу и флаг успешности
	s.Lock()
	var updatedTask *Task
	var err error

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Done = done
			updatedTask = &s.tasks[i]
			break
		}
	}
	if updatedTask == nil {
		err = fmt.Errorf("task with ID %d not found", id)
	}

	s.Unlock()

	if err == nil {
		_ = s.Save() // Save вне блокировки
	}
	return updatedTask, err

}

// метод ищет задачу по ID и удаляет её из среза, если находит
func (s *Store) DeleteTask(id int) bool {
	s.Lock()
	var found bool
	for i, t := range s.tasks {
		if t.ID == id {
			// удаляем задачу из среза tasks
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			found = true
			break
		}
	}
	s.Unlock()

	if found {
		_ = s.Save() // Save вне блокировки
	}

	return found
}
