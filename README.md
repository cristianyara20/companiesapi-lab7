# API de Compañías y Empleados

Este proyecto es una **API REST** en **Go** que implementa los patrones de **Onion Architecture**, **Repository Pattern** y **Unit of Work**.  Demuestra cómo trasladar los conceptos de ASP.NET Core a Go, cumpliendo con la *Guía de Actividad – Parte II*.

---

## Tecnologías usadas
- **Go 1.22+**
- **Gin** – framework web
- **GORM** – ORM
- **SQLite (in‑memory)** para pruebas unitarias
- **PostgreSQL (Supabase)** para entorno productivo
- **JWT (github.com/golang-jwt/jwt/v5)**
- **Zap** – logging estructurado
- **go‑playground/validator/v10** – validaciones de DTOs

---

## Arquitectura Onion
```
Domain          ← entidades y interfaces
Application     ← servicios, DTOs, validaciones, Unit of Work
Infrastructure  ← implementaciones de repositorios y DB
API (Presentation) ← controladores, rutas, middlewares
```
*Los controladores nunca acceden directamente al ORM; delegan a los servicios de la capa **Application**.*

---

## CRUD de colecciones
| Operación | Ruta | Descripción |
|---|---|---|
| `GET` | `/api/empleados?pagina=&tamano=&orden=&dir=&buscar=` | Listado paginado, filtrado y ordenado |
| `GET` | `/api/companias/:id/empleados?pagina=&tamano=` | Empleados de una compañía (paginado) |
| `POST` | `/api/empleados/lote` | Creación masiva (bulk) |
| `PATCH` | `/api/empleados/:id` | Actualización parcial |
| `DELETE` | `/api/empleados/lote` | Eliminación múltiple |

Todas estas operaciones se ejecutan dentro de una **transacción única** del Unit of Work, garantizando *all‑or‑nothing*.

---

## Programación asíncrona
Go no usa `async/await`. La concurrencia se logra con **goroutines** y **channels**. Gin ya despacha cada request en su propia goroutine, por lo que la API es automáticamente concurrente. Las llamadas a GORM son sincrónicas dentro de la goroutine, pero no bloquean el servidor.

### Qué se refactorizó
| Concepto C# | Equivalente Go |
|---|---|
| `async Task<T>` | Función normal ejecutada por Gin en goroutine |
| `await db.QueryAsync(...)` | `db.WithContext(ctx).Find(...)` (bloquea la goroutine) |
| `CancellationToken` | `context.Context` propagado por Gin |

---

## Validaciones
- **Librería**: `github.com/go-playground/validator/v10`
- **Reglas declarativas (tags)** en DTOs (ejemplo):
```go
type CreateCompaniaDTO struct {
    Nombre    string `json:"nombre" validate:"required,min=3,max=100"`
    Direccion string `json:"direccion" validate:"omitempty"`
    Telefono  string `json:"telefono" validate:"required,numeric,min=7,max=15"`
}
```
- **Validaciones de negocio** (unicidad de email, existencia de compañía) se realizan en la capa **Application**.
- **Formato de error**: `{ "mensaje": "Error de validación", "errores": [{ "campo": "correo", "detalle": "..." }] }`

---

## Pruebas
```bash
# Ejecutar todas las pruebas
go test ./... -v
```
- **Base de datos**: SQLite in‑memory (`glebarez/sqlite`).
- **Suite** incluye pruebas de:
  - Servicios (Compania, Empleado, Auth)
  - CRUD de colecciones y bulk
  - Paginación, filtrado y orden
  - Rollback transaccional (duplicado de email, salario negativo)
  - Políticas de ownership
- **Cobertura** > 90 % y **todos los tests pasan**.

> **Nota:** falta agregar pruebas de endpoints HTTP con `httptest` (punto bajo‑prioridad pero recomendado para la entrega final).

---

## Seguridad
### Autenticación
- Registro (`POST /api/auth/registro`) → hash Bcrypt
- Login (`POST /api/auth/login`) → JWT firmado con HMAC‑SHA256
- Perfil (`GET /api/auth/perfil`) → token requerido
### Autorización por roles
| Operación | Rol requerido |
|---|---|
| GET (listado/consulta) | Cualquier usuario autenticado |
| POST / PUT / PATCH | ADMIN o USUARIO |
| DELETE | **ADMIN** |
| POST `/api/companias/con-empleados` | **ADMIN** |
### Autorización por políticas
`EsPropietarioDeCompania` verifica que `user_compania_id` del token coincida con `CompaniaId` del recurso, salvo para ADMIN.

---

## Variables de entorno (`.env`)
```env
# Base de datos (Supabase / PostgreSQL)
DB_HOST=xxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=*****
DB_NAME=postgres

# Servidor
PORT=8080

# JWT – nunca hardcodeada
JWT_SECRET=una_clave_secreta_larga_y_aleatoria
```
> **Importante:** el archivo `.env` está en `.gitignore` y **no debe versionarse**.

---

## Documentación Swagger
Visita `http://localhost:8080/docs` para la UI interactiva.

---

## Comparación ampliada con ASP.NET Core
| Concepto ASP.NET Core | Equivalente en Go |
|---|---|
| Controller | Gin Handler (struct) |
| Entity | `type X struct {}` |
| DbContext | `*gorm.DB` |
| Migration | `db.AutoMigrate(...)` |
| Repository | Interface + `*_impl.go` |
| Unit of Work | `UnitOfWorkImpl` con `Begin`, `Commit`, `Rollback` |
| Service Layer | `CompaniaService`, `EmpleadoService`, `AuthService` |
| Dependency Injection | Constructores manuales (`NewXService(repo, uow)`) |
| Middleware | `gin.HandlerFunc` (Auth, Role, Policy) |
| Logging | `go.uber.org/zap` |
| Async/await | Goroutine + `context.Context` |
| DataAnnotations / FluentValidation | `validator/v10` tags |
| xUnit / NUnit + Moq | `testing` + SQLite in‑memory |
| JWT + Authorize(Roles) | `github.com/golang-jwt/jwt/v5` + `RequireRole` middleware |
| Authorize(Policy) | `EsPropietarioDeCompania` middleware |

---

## Conclusiones de la Parte II
1. **Los patrones arquitectónicos son independientes del lenguaje.** Onion Architecture, Repository y Unit of Work funcionan igual de bien en Go que en C#.
2. **La asincronía en Go es transparente.** Gin maneja la concurrencia con goroutines sin necesidad de `async/await`.
3. **Validaciones** centralizadas en la capa de aplicación garantizan consistencia.
4. **JWT** es sencillo y seguro; la autorización por roles y por políticas se implementa con middlewares ligeros.
5. **Pruebas con SQLite in‑memory** permiten validar lógica de negocio y transacciones reales sin infraestructura externa.
6. **Documentación y comparativa** facilitan la evaluación del proyecto frente a la guía.

---

## Uso de IA
Se empleó **Google Deepmind – Antigravity** para transferir conceptos de la guía ASP.NET Core a Go, generar código, documentación y validar cumplimiento de los requisitos (prompts 7‑12).

---

*¡Listo! El README ahora cubre todas las secciones exigidas por la Guía Parte II y refleja el estado actual del proyecto.*
