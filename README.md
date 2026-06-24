# StudyFlow Backend

API REST para la plataforma StudyFlow, una herramienta de gestión académica con IA local.

## Tech Stack

- **Go** + **Gin** — Framework HTTP
- **PostgreSQL** — Base de datos
- **Ollama** — IA local (Qwen3 8B)
- **Swagger** — Documentación de la API
- **JWT** — Autenticación
- **Google OAuth** — Login

## Estructura
cmd/server/        → Punto de entrada

internal/          → Módulos del sistema (auth, courses, tasks, notes...)

pkg/               → Utilidades compartidas (database, middleware, response)

docs/              → Swagger autogenerado

storage/           → PDFs subidos

## Correr en local

```bash
# 1. Clonar el repo
git clone https://github.com/JefersonGomez/studyflow-backend.git

# 2. Copiar variables de entorno
cp .env.example .env

# 3. Correr
go run cmd/server/main.go
```

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | /api/v1/health | Health check |

Documentación completa: `http://localhost:8080/swagger/index.html`

## Progreso

- [x] Setup inicial del proyecto
- [x] Conexión a PostgreSQL con GORM
- [x] Modelos y migraciones automáticas
- [x] Autenticación con Google OAuth
- [x] Generación y validación de JWT
- [x] Middleware de rutas protegidas
- [x] CRUD de materias (courses)
- [x] CRUD de tareas (tasks)
- [x] CRUD de notas (notes)
- [x] CRUD de eventos (events)
- [x] CRUD de pizarra (whiteboard)
- [x] Subida y extracción de texto de PDFs
- [x] Integración con Ollama IA local