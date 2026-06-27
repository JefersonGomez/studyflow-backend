package task

import (
	"github.com/JefersonGomez/studyflow-backend/pkg/database"
)

/* type Task struct {
	ID          string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID      string    `json:"userID" gorm:"not null"`
	CourseID    *string   `json:"courseID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"dueDate" gorm:"not null"`
	CreatedAt   time.Time `json:"createAt" gorm:"autoCreateTime" `
	UptadeAt    time.Time `json:"updateAt" gorm:"autoUpdateTime"`
} */

func CreateTask(task *Task) error {

	result := database.DB.Create(task)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func GetTaskByCourseID(courseID string) ([]Task, error) {

	var tasks []Task

	result := database.DB.Where("course_id = ?", courseID).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil

}

func GetAllTasks(userID string) ([]Task, error) {

	var tasks []Task

	result := database.DB.Where("user_id = ? ", userID).Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}
func GetTaskByID(taskID string) (*Task, error) {

	var task Task

	result := database.DB.Where("id = ?", taskID).First(&task)

	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil

}

func UpdateTask(task *Task) error {

	result := database.DB.Save(&task)

	if result.Error != nil {
		return result.Error
	}
	return nil

}

func DeleteTask(taskID string) error {

	result := database.DB.Delete(&Task{}, "id=?", taskID)

	if result.Error != nil {
		return result.Error
	}
	return nil

}
