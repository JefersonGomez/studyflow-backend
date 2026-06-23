package studyfile

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var texto strings.Builder
	totalPages := r.NumPage()

	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		content, err := page.GetPlainText(nil)
		if err == nil {
			texto.WriteString(content)
		}
	}

	return texto.String(), nil
}

func UploadStudyFileService(userID, courseID, filename, storagePath string) (*Studyfile, error) {
	// extraer texto del PDF ya guardado en disco
	parsedText, err := ExtractTextFromPDF(storagePath)
	if err != nil {
		parsedText = "" // si falla la extracción, guardamos vacío pero no fallamos
	}

	newFile := &Studyfile{
		UserID:      userID,
		CourseID:    &courseID,
		FileName:    filename,
		StoragePath: storagePath,
		ParsedText:  parsedText,
		AiProcessed: false,
	}

	err = CreateStudyFile(newFile)
	if err != nil {
		return nil, err
	}

	return newFile, nil
}

func GetCourseFilesService(courseID string) ([]Studyfile, error) {
	files, err := GetStudyFilesByCourseID(courseID)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func DeleteStudyFileService(fileID, userID string) error {
	file, err := GetStudyFileByID(fileID)
	if err != nil {
		return err
	}

	if file.UserID != userID {
		return errors.New("no tienes permiso para eliminar este archivo")
	}

	// eliminar el archivo físico del disco
	if err := os.Remove(file.StoragePath); err != nil {
		return fmt.Errorf("error al eliminar el archivo del disco: %w", err)
	}

	return DeleteStudyFile(file.ID)
}
