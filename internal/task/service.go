package task

import (
	"errors"
	"time"
)

func CreateTaskService(userID, courseID, title, description, status string, dueDate time.Time) (*Task, error) {

	newTask := &Task{
		UserID:      userID,
		CourseID:    &courseID,
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     dueDate,
	}

	err := CreateTask(newTask)

	if err != nil {
		return nil, err
	}
	return newTask, nil

}

func GetTasksByUserService(userID string) ([]Task, error) {

	task, err := GetAllTasks(userID)

	if err != nil {
		return nil, err
	}

	return task, nil

}

func GetCourseTask(courseID string) ([]Task, error) {

	tasks, err := GetTaskByCourseID(courseID)

	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func UpdateTaskService(taskID, userID, courseID, title, description, status string, dueDate time.Time) (*Task, error) {

	task, err := GetTaskByID(taskID)

	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, errors.New("no tienes permiso para editar este curso")
	}

	task.Title = title
	task.Description = description
	task.Status = status
	task.DueDate = dueDate

	err = UpdateTask(task)
	if err != nil {
		return nil, err
	}

	return task, nil

}

func DeleteTaskService(taskID string, userID string) error {

	task, err := GetTaskByID(taskID)

	if err != nil {
		return err
	}

	if task.UserID != userID {
		return errors.New("no tienes permiso para editar este curso")
	}

	err = DeleteTask(task.ID)

	if err != nil {
		return err
	}

	return nil

}
