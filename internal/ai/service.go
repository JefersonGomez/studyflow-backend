package ai

import "fmt"

//funciones que va a poder hacer el usuario con ia
//funcion para generar un resumen
func SummarizeNote(content string) (string, error) {
	prompt := fmt.Sprintf("Resume el siguiente texto de forma clara y concisa en español:\n\n%s", content)
	return Generate(prompt)

}

//funcion para generar preguntas de estudio

func GenerateQuestions(content string) (string, error) {
	prompt := fmt.Sprintf("Genera 10 preguntas de estudio en español basadas en el siguiente texto:\n\n%s", content)
	return Generate(prompt)

}

func GenerateStudyPlan(content string, days int) (string, error) {
	prompt := fmt.Sprintf("Crea un plan de estudio para %d días en español basado en el siguiente contenido:\n\n%s", days, content)
	return Generate(prompt)
}

//funcion para analizar PDF y extraer eventos

func AnalyzePDF(content string) (string, error) {
	prompt := fmt.Sprintf(`Analiza el siguiente programa de curso y extrae todos los eventos importantes 
    (exámenes, quizzes, proyectos, entregas) con sus fechas. 
    Responde en formato JSON con esta estructura:
    {"eventos": [{"titulo": "nombre", "fecha": "YYYY-MM-DD", "tipo": "examen|quiz|proyecto"}]}
    
    Texto del programa:\n\n%s`, content)
	return Generate(prompt)
}
