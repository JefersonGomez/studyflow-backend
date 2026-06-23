package note

import "errors"

/* type Note struct {
	ID        string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"userID" gorm:"not null"`
	CourseID  *string   `json:"courseID" gorm:"not null"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createAt" gorm:"autoCreateTime" `
	UptadeAt  time.Time `json:"updateAt" gorm:"autoUpdateTime"`
}
*/

func CreateNoteService(userID, courseID, title, content string) (*Note, error) {

	newNote := &Note{

		UserID:   userID,
		CourseID: &courseID,
		Title:    title,
		Content:  content,
	}

	err := CreateNote(newNote)

	if err != nil {

		return nil, err
	}
	return newNote, nil

}

func GetNotesCourseService(courseID string) ([]Note, error) {

	notes, err := GetNotesByCourseID(courseID)

	if err != nil {
		return nil, err
	}
	return notes, nil

}

func UpdateNotesService(noteID, userID, courseID, title, content string) (*Note, error) {

	note, err := GetNotesByID(noteID)

	if err != nil {
		return nil, err
	}

	if note.UserID != userID {
		return nil, errors.New("no tienes permiso para editar esta nota")
	}

	note.Title = title
	note.Content = content

	err = UpdateNote(note)

	if err != nil {
		return nil, err
	}
	return note, nil

}

func DeleteNoteService(noteID, userID string) error {

	note, err := GetNotesByID(noteID)
	if err != nil {
		return err
	}
	if note.UserID != userID {
		return errors.New("no tienes permiso para eliminar esta nota")
	}

	err = DeleteNote(note.ID)

	if err != nil {
		return err
	}
	return nil

}
