package event

import "github.com/JefersonGomez/studyflow-backend/pkg/database"

func CreateEvent(event *Event) error {
	result := database.DB.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetEventsByUserID(userID string) ([]Event, error) {
	var events []Event
	result := database.DB.Where("user_id = ?", userID).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func GetEventsByCourseID(courseID string) ([]Event, error) {
	var events []Event
	result := database.DB.Where("course_id = ?", courseID).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func GetEventByID(eventID string) (*Event, error) {
	var event Event
	result := database.DB.Where("id = ?", eventID).First(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

func UpdateEvent(event *Event) error {
	result := database.DB.Save(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteEvent(eventID string) error {
	result := database.DB.Delete(&Event{}, "id = ?", eventID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
