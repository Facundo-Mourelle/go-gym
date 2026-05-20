# Go Gym Tracker 🏋️‍♂️

Aplicación *full-stack* (Móvil/Web) para registrar tu progreso en el gimnasio con alto nivel de detalle. Diseñada para calcular volumen efectivo, esfuerzo (RPE/RIR) y perfiles de resistencia mecánica.

## 🚀 Tecnologías

* **Backend:** Go 1.25, Arquitectura Hexagonal, `net/http` estándar (mux Go 1.22+)
* **Frontend:** React 19, TypeScript, Vite, Tailwind CSS, Zustand (Estado global), Recharts (Gráficos), Lucide React (Íconos)
* **Base de Datos:** PostgreSQL 15 (Ejecutado localmente vía Docker)
* **Autenticación:** JWT (golang-jwt/v5)

---

## 🛠️ Setup Rápido (Scripts)

El proyecto incluye scripts para arrancar todo el stack:

```bash
# Iniciar todo (DB + Backend + Frontend)
./start.sh

# Solo backend (para desarrollo)
./start_backend.sh

# Solo frontend (para desarrollo)
./start_frontend.sh

# Ejecutar tests al iniciar
./start.sh --tests

# Detener todos los servicios
./stop.sh
```

### `start.sh` (recomendado)
Levanta PostgreSQL con Docker, espera a que esté lista, inicia el backend y el frontend automáticamente.

```
=== Services running ===
DB:      localhost:5432
Backend: localhost:8080
Frontend: http://localhost:5173
```

### `stop.sh`
Detiene backend, frontend y baja el contenedor de PostgreSQL.

---

## 🛠️ Setup Manual (paso a paso)

### 1. Requisitos Previos
* Instalar [Go](https://golang.org/doc/install) (v1.22+)
* Instalar [Node.js](https://nodejs.org/es/) y `npm`
* Instalar [Docker](https://www.docker.com/) y Docker Compose

### 2. Configurar variables de entorno

Crear archivo `.env` en la raíz del proyecto:

```env
DATABASE_URL="postgres://gym:gympassword@localhost:5432/gymtracker?sslmode=disable"
JWT_SECRET="your-secret-key-change-in-production"
PORT=8080
```

| Variable | Default | Descripción |
|---|---|---|
| `DATABASE_URL` | `""` | Conexión a PostgreSQL. Si está vacío usa repositorios en memoria. |
| `JWT_SECRET` | `"your-secret-key-change-in-production"` | Secreto para firmar tokens JWT |
| `PORT` | `8080` | Puerto del servidor HTTP |

### 3. Base de Datos (PostgreSQL)

```bash
docker-compose up -d
```

*Las migraciones de esquemas (`internal/repository/postgres/migrations`) se pueden aplicar para inicializar las tablas.*

### 4. Iniciar el Backend (Go API)

```bash
go run cmd/api/main.go
```

*El servidor corre en el puerto `8080`. Si no hay conexión a PostgreSQL, hace fallback automático a repositorios en memoria.*

### 5. Iniciar el Frontend (React/Vite)

```bash
cd gym-tracker-frontend
npm install
npm run dev
```

*Visita `http://localhost:5173` en tu navegador.*

---

## 📡 API Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `GET` | `/health` | ❌ | Health check |
| **Auth** | | | |
| `POST` | `/api/v1/auth/register` | ❌ | Registrar usuario |
| `POST` | `/api/v1/auth/login` | ❌ | Iniciar sesión (retorna JWT) |
| `GET` | `/api/v1/auth/me` | ✅ | Obtener usuario actual |
| **Ejercicios** | | | |
| `GET` | `/api/v1/exercises/templates` | ❌ | Listar plantillas de ejercicios |
| `GET` | `/api/v1/exercises` | ✅ | Listar ejercicios del usuario |
| `GET` | `/api/v1/exercises/{id}` | ❌ | Obtener ejercicio por ID |
| `POST` | `/api/v1/exercises/custom` | ✅ | Crear ejercicio personalizado |
| `POST` | `/api/v1/exercises/from-template` | ✅ | Crear desde plantilla |
| `PUT` | `/api/v1/exercises/{id}` | ✅ | Actualizar ejercicio |
| `DELETE` | `/api/v1/exercises/{id}` | ✅ | Eliminar ejercicio |
| `GET` | `/api/v1/patterns` | ❌ | Listar patrones de movimiento |
| **Rutinas** | | | |
| `POST` | `/api/v1/routines` | ✅ | Crear rutina |
| `GET` | `/api/v1/routines` | ✅ | Listar rutinas |
| `GET` | `/api/v1/routines/{id}` | ✅ | Obtener rutina |
| `PUT` | `/api/v1/routines/{id}` | ✅ | Actualizar rutina |
| `DELETE` | `/api/v1/routines/{id}` | ✅ | Eliminar rutina |
| `POST` | `/api/v1/routines/{id}/exercises` | ✅ | Agregar ejercicio a rutina |
| `DELETE` | `/api/v1/routines/{id}/exercises/{exerciseId}` | ✅ | Quitar ejercicio de rutina |
| **Sesiones** | | | |
| `POST` | `/api/v1/sessions` | ✅ | Iniciar sesión de entrenamiento |
| `GET` | `/api/v1/sessions` | ✅ | Listar sesiones |
| `GET` | `/api/v1/sessions/{id}` | ✅ | Obtener sesión con detalles |
| `POST` | `/api/v1/sessions/{id}/sets` | ✅ | Registrar serie en sesión |
| `PUT` | `/api/v1/sessions/{id}/sets/{setId}` | ✅ | Actualizar serie |
| `DELETE` | `/api/v1/sessions/{id}/sets/{setId}` | ✅ | Eliminar serie |
| `POST` | `/api/v1/sessions/{id}/complete` | ✅ | Completar sesión |
| **Equipamiento** | | | |
| `GET` | `/api/v1/equipment` | ❌ | Listar equipamiento disponible |
| **Métricas** | | | |
| `GET` | `/api/v1/metrics/dashboard` | ✅ | Dashboard semanal |
| `GET` | `/api/v1/metrics/progress/{exerciseId}` | ✅ | Progreso histórico de un ejercicio |

---

## ✨ Características Principales

* **Dashboard Dinámico:** Estadísticas semanales calculadas en tiempo real (sesiones, volumen total, series) e historial rápido de "Actividad Reciente".
* **Navegación Global:** Barra de navegación inferior diseñada para *Mobile-first*, facilitando el salto entre Inicio, Historial, Ejercicios y Progreso.
* **Sesión Activa UX:**
  * Temporizador de descanso (Rest Timer) automático de 90s al completar un set.
  * Botones rápidos (+/-) para ajustar peso y repeticiones.
  * Slider interactivo para registrar el Esfuerzo Percibido (RPE/RIR).
  * Contexto histórico: Vista rápida de qué hiciste la última vez en ese mismo ejercicio.
* **Diccionario de Ejercicios:** Explorador de ejercicios con filtros en tiempo real y etiquetas visuales de equipamiento y grupos musculares.
* **Métricas y Progreso:** Gráficos de tendencias con `recharts` que visualizan tu Progreso del 1RM (Fórmula de Epley) a lo largo del tiempo.
* **Autenticación JWT:** Los usuarios pueden registrarse e iniciar sesión. Los datos personales (ejercicios, rutinas, sesiones) están protegidos por usuario.

---

## 🏗️ Arquitectura y Dominio (Backend)

Análisis de Arquitectura en forma "Bottom-Up" implementada bajo principios de código limpio y Arquitectura Hexagonal:

```
internal/
├── api/           # Handlers HTTP, router, middleware
│   ├── handler/   # Auth, Exercise, Session, Workout, Routine, Metrics, Equipment
│   └── middleware/ # JWT, logging
├── config/        # Configuración (JWT, puerto, DB)
├── domain/        # Modelo de dominio puro
│   ├── calculator/ # Cálculos de volumen, 1RM, etc.
│   └── resistance/ # Perfiles de resistencia mecánica
├── repository/    # Interfaces de repositorio
│   ├── memory/    # Implementación en memoria (fallback)
│   └── postgres/  # Implementación con PostgreSQL
├── service/       # Lógica de negocio (casos de uso)
└── storage/       # Almacenamiento de archivos
```

### Muscle Groups
Contiene tipos para referirse a los principales grupos musculares de los que vamos a mantener registro.

### Equipment & Resistance Profiles
Contiene tipos de equipamiento a utilizar para registrar de manera más acertada los ejercicios. Se dividen principalmente en:
- **Pesos Libres:** El peso es directo, el perfil de resistencia depende directamente del brazo de momento.
- **Máquinas/Poleas:** El peso depende del sistema empleado. El peso efectivo se calcula mediante *perfiles de resistencia* definidos en `internal/domain/resistance` evaluando ventajas mecánicas y pérdida por fricción.

Se incluye la marca del equipamiento para distinguir posibles variaciones sutiles de tensión entre distintas máquinas.

### Movement Patterns
Contiene los patrones de movimiento que realizan los grupos musculares declarados, especialmente los que contribuyen a la hipertrofia.

### Exercise
Para construir un ejercicio necesitamos: *Equipamiento*, *Patrones de movimiento Principales* y *Secundarios*.
La categorización depende de la contribución del patrón de movimiento y el rango (ROM) empleado.

> [!Note]
> Esta asignación no es derivada automáticamente, queda a decisión del usuario (vía Builder) para permitir flexibilidad total dependiendo del equipo utilizado.

### Workout (Routines)
Registro de plantillas/planes de entrenamiento que agrupan ejercicios ordenados para repetir en el tiempo.

### Session & Performed Sets
Registro vivo de la sesión en curso. A la hora de registrar una serie se captura:
- Repeticiones realizadas y en reserva (RIR/RPE).
- Número de serie dentro de la sesión.
- Peso bruto (sin procesar) y Equipamiento.

A partir del peso sin procesar, se calcula automáticamente el **Peso Efectivo** utilizando el perfil de resistencia del módulo `resistance`.

---

## 📁 Estructura del Proyecto

```
go-gym/
├── cmd/api/main.go          # Entry point del backend
├── internal/                 # Lógica de negocio (hexagonal)
├── gym-tracker-frontend/     # Frontend React + Vite
├── scripts/                  # Utilidades
├── docker-compose.yml        # PostgreSQL container
├── start.sh                  # Inicia todo el stack
├── stop.sh                   # Detiene todo el stack
├── start_backend.sh          # Inicia solo backend
├── start_frontend.sh         # Inicia solo frontend
├── .env                      # Variables de entorno (local)
└── profile.env               # Perfil JWT de ejemplo para testing
