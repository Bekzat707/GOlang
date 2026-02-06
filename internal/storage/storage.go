package storage

import (
	"github.com/bekzatsaparbekov/task-api/internal/models"
)

type TaskStorage struct {
	tasks  []models.Task
	nextID int
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  []models.Task{},
		nextID: 1,
	}
}

func (s *TaskStorage) Create(title string) models.Task {
	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks = append(s.tasks, task)
	s.nextID++
	return task
}

func (s *TaskStorage) GetByID(id int) (models.Task, bool) {
	for _, task := range s.tasks {
		if task.ID == id {
			return task, true
		}
	}
	return models.Task{}, false
}

func (s *TaskStorage) GetAll() []models.Task {
	return s.tasks
}

func (s *TaskStorage) Update(id int, done bool) (models.Task, bool) {
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Done = done
			return s.tasks[i], true
		}
	}
	return models.Task{}, false
}
