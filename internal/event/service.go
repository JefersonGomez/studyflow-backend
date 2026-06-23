package event

import (
	"errors"
	"time"
)

func CreateEventService(userID, courseID, title, description, eventType string, startDate time.Time, endDate *time.Time) (*Event, error) {
	newEvent := &Event{
		UserID:      userID,
		Title:       title,
		Description: description,
		Type:        eventType,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if courseID != "" {
		newEvent.CourseID = &courseID
	}

	err := CreateEvent(newEvent)
	if err != nil {
		return nil, err
	}
	return newEvent, nil
}

func GetUserEventsService(userID string) ([]Event, error) {
	events, err := GetEventsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func GetCourseEventsService(courseID string) ([]Event, error) {
	events, err := GetEventsByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func UpdateEventService(eventID, userID, title, description, eventType string, startDate time.Time, endDate *time.Time) (*Event, error) {
	event, err := GetEventByID(eventID)
	if err != nil {
		return nil, err
	}

	if event.UserID != userID {
		return nil, errors.New("no tienes permiso para editar este evento")
	}

	event.Title = title
	event.Description = description
	event.Type = eventType
	event.StartDate = startDate
	event.EndDate = endDate

	err = UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func DeleteEventService(eventID, userID string) error {
	event, err := GetEventByID(eventID)
	if err != nil {
		return err
	}

	if event.UserID != userID {
		return errors.New("no tienes permiso para eliminar este evento")
	}

	err = DeleteEvent(event.ID)
	if err != nil {
		return err
	}
	return nil
}
