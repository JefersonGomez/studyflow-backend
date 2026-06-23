package whiteboard

import (
	"errors"

	"gorm.io/datatypes"
)

func CreateWhiteboardService(userID, courseID, title string, elements datatypes.JSON) (*Whiteboard, error) {
	newWhiteboard := &Whiteboard{
		UserID:   userID,
		Title:    title,
		Elements: elements,
	}

	if courseID != "" {
		newWhiteboard.CourseID = &courseID
	}

	err := CreateWhiteboard(newWhiteboard)
	if err != nil {
		return nil, err
	}
	return newWhiteboard, nil
}

func GetCourseWhiteboardsService(courseID string) ([]Whiteboard, error) {
	whiteboards, err := GetWhiteboardsByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return whiteboards, nil
}

func UpdateWhiteboardService(whiteboardID, userID, title string, elements datatypes.JSON) (*Whiteboard, error) {
	whiteboard, err := GetWhiteboardByID(whiteboardID)
	if err != nil {
		return nil, err
	}

	if whiteboard.UserID != userID {
		return nil, errors.New("no tienes permiso para editar esta pizarra")
	}

	whiteboard.Title = title
	whiteboard.Elements = elements

	err = UpdateWhiteboard(whiteboard)
	if err != nil {
		return nil, err
	}
	return whiteboard, nil
}

func DeleteWhiteboardService(whiteboardID, userID string) error {
	whiteboard, err := GetWhiteboardByID(whiteboardID)
	if err != nil {
		return err
	}

	if whiteboard.UserID != userID {
		return errors.New("no tienes permiso para eliminar esta pizarra")
	}

	err = DeleteWhiteboard(whiteboard.ID)
	if err != nil {
		return err
	}
	return nil
}
