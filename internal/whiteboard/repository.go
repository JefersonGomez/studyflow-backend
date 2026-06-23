package whiteboard

import "github.com/JefersonGomez/studyflow-backend/pkg/database"

func CreateWhiteboard(whiteboard *Whiteboard) error {

	result := database.DB.Create(whiteboard)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func GetWhiteboardsByCourseID(courseID string) ([]Whiteboard, error) {

	var whiteboards []Whiteboard

	result := database.DB.Where("course_id = ?", courseID).Find(&whiteboards)

	if result.Error != nil {
		return nil, result.Error
	}
	return whiteboards, nil

}

func GetWhiteboardByID(whiteboardID string) (*Whiteboard, error) {

	var whiteboard Whiteboard

	result := database.DB.Where("id = ?", whiteboardID).First(&whiteboard)

	if result.Error != nil {
		return nil, result.Error
	}
	return &whiteboard, nil

}

func UpdateWhiteboard(whiteboard *Whiteboard) error {

	result := database.DB.Save(whiteboard)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func DeleteWhiteboard(whiteboardID string) error {

	result := database.DB.Delete(&Whiteboard{}, "id =?", whiteboardID)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
