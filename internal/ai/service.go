package ai

import (
	"fmt"
	"strings"
)

func SummarizeNote(content string) (string, error) {
	if len(content) > 3000 {
		content = content[:3000]
	}
	prompt := fmt.Sprintf("Resume el siguiente texto de forma clara y concisa en español:\n\n%s", content)
	return Generate(prompt)
}

func GenerateQuestions(content string) (string, error) {
	if len(content) > 3000 {
		content = content[:3000]
	}
	prompt := fmt.Sprintf(`Genera exactamente 10 preguntas de estudio en español basadas en este texto.
Responde ÚNICAMENTE con un array JSON válido. Sin texto adicional, sin markdown, sin bloques de código.
Usa este formato exacto:
[{"question":"¿Pregunta?","answer":"Respuesta correcta"},{"question":"¿Pregunta?","answer":"Respuesta correcta"}]

Texto: %s`, content)
	return Generate(prompt)
}

func GenerateStudyPlan(content string, days int) (string, error) {
	if len(content) > 3000 {
		content = content[:3000]
	}
	prompt := fmt.Sprintf("Crea un plan de estudio para %d días en español basado en el siguiente contenido:\n\n%s", days, content)
	return Generate(prompt)
}

func AnalyzePDF(content string) (string, error) {
	// Buscar la palabra "cronograma" (case insensitive)
	lowerContent := strings.ToLower(content)
	idx := strings.Index(lowerContent, "cronograma")

	if idx != -1 {
		// Extraer desde donde aparece "cronograma"
		content = content[idx:]
		// Limitar a 4000 caracteres desde ahí
		if len(content) > 4000 {
			content = content[:4000]
		}
	} else {
		// Si no encuentra "cronograma", usar las últimas 4000 chars como antes
		if len(content) > 4000 {
			content = content[len(content)-4000:]
		}
	}

	prompt := fmt.Sprintf(`Eres un asistente académico. Analiza el siguiente texto de un programa de curso universitario.
    Busca una tabla de CRONOGRAMA con columnas de sesión, fecha y actividades.
    Extrae ÚNICAMENTE las actividades evaluadas: Quiz, Examen, Tarea, Seguimiento, Avance, Investigación, Propuesta, Proyecto.
    USA SOLO las fechas literales del texto, NO inventes fechas.
    Convierte fechas al formato YYYY-MM-DD usando el año 2026.
    Responde SOLO en JSON sin texto adicional ni explicaciones:
    {"eventos": [{"titulo": "nombre actividad", "fecha": "YYYY-MM-DD", "tipo": "examen|quiz|proyecto|tarea"}]}
    
    Texto:\n\n%s`, content)
	return Generate(prompt)
}
