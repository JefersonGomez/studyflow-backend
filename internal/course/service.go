package course

import "errors"

func CreateCourseService(userID, name, description, color string) (*Course, error) {

	newCouser := &Course{

		UserID:      userID,
		Name:        name,
		Description: description,
		Color:       color,
	}

	err := CreateCourse(newCouser)
	if err != nil {
		return nil, err
	}

	return newCouser, nil

}

func GetUserCoursesByService(userID string) ([]Course, error) {

	courses, err := GetCoursesByUserID(userID)

	if err != nil {
		return nil, err
	}
	return courses, nil

}

func UpdateCourseService(courseID, userID, name, description, color string) (*Course, error) {
	curso, err := GetCourseByID(courseID)
	if err != nil {
		return nil, err
	}

	if curso.UserID != userID {
		return nil, errors.New("no tienes permiso para editar este curso")
	}

	curso.Name = name
	curso.Description = description
	curso.Color = color

	err = UpdateCourse(curso)
	if err != nil {
		return nil, err
	}

	return curso, nil
}

func DeleteCourseService(courseID, userID string) error {
	curso, err := GetCourseByID(courseID)
	if err != nil {
		return err
	}

	if curso.UserID != userID {
		return errors.New("no tienes permiso para editar este curso")
	}

	err = DeleteCourse(curso.ID)
	if err != nil {
		return err
	}

	return nil
}
