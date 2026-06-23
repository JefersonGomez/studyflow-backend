package note

import "github.com/JefersonGomez/studyflow-backend/pkg/database"

func CreateNote(note *Note) error {

	result := database.DB.Create(note)

	if result.Error != nil {
		return result.Error
	}
	return nil

}

/*
type Note struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"userID" gorm:"not null"`
	CourseID  *string   `json:"courseID" gorm:"not null"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createAt" gorm:"autoCreateTime" `
	UptadeAt  time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}
*/
func GetNotesByCourseID(courseID string) ([]Note, error) {

	var notes []Note

	result := database.DB.Where("course_id = ?", courseID).Find(&notes)

	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil

}

func GetNotesByID(noteID string) (*Note, error) {

	var note Note

	result := database.DB.Where("id = ?", noteID).First(&note)

	if result.Error != nil {
		return nil, result.Error
	}
	return &note, nil

}

func UpdateNote(note *Note) error {

	result := database.DB.Save(&note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteNote(noteID string) error {

	result := database.DB.Delete(&Note{}, "id = ?", noteID)

	if result.Error != nil {
		return result.Error
	}
	return nil

}
