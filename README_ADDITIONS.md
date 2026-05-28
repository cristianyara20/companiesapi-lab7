# 📄 README – Secciones añadidas para la Guía Parte II

## 📦 Programación asíncrona

> **¿Qué se considera "asíncrono" en este proyecto?**
>
> - Gin crea una **goroutine** por cada petición HTTP, por lo que cada request se procesa de forma concurrente sin bloquear al servidor.
> - Los handlers usan `ctx.Request.Context()` para propagar el contexto de cancelación y timeout.
> - GORM ejecuta sus operaciones de forma síncrona dentro de la goroutine, pero al estar dentro del mismo flujo concurrente la aplicación mantiene alta capacidad de respuesta.
> - No se utilizó `async/await` (no existen en Go) ni se introdujeron frameworks externos; la solución nativa de Go satisface el requisito de programación asíncrona.
> - **Impacto en Unit of Work**: la transacción (`uow.BeginTransaction()`) se crea y se maneja dentro de la misma goroutine que procesa la petición, por lo que el commit/rollback sigue funcionando correctamente.

---

## 🌐 Variables de entorno

| Variable | Descripción | Ejemplo | Comentario |
|----------|-------------|---------|------------|
| `JWT_SECRET` | Clave secreta usada para firmar y validar los JWT (HS256). | `supersecreta123` | **NO** debe estar versionada en Git; se carga desde `.env`. |
| `DB_URL` | URL de conexión a la base de datos (Supabase/PostgreSQL). | `postgresql://user:pass@db.supabase.co:5432/companies` | Se usa GORM; el driver se inicializa en `main.go`. |
| `PORT` | Puerto donde corre la API. | `8080` | Opcional; por defecto Gin usa `8080`. |
| `ENV` | Entorno de ejecución (`development`, `production`). | `development` | Permite cambiar configuraciones de logging. |

---

## 📊 Tabla comparativa ampliada (Sección 14)

| Concepto (ASP.NET) | Equivalente en Go (Proyecto) | Comentario |
|---------------------|-----------------------------|-----------|
| **Colecciones (IEnumerable / List)** | Slices (`[]Entidad`) retornados por los servicios | Go no necesita interfaz extra; los slices son la estructura idiomática. |
| **Paginación (Skip/Take)** | GORM `Offset` + `Limit` + `Count` | Implementado en `GetPaged` de los repositorios. |
| **Programación asíncrona** | Goroutine por request (Gin) + `context.Context` | Se explicó en la sección de async. |
| **DataAnnotations / FluentValidation** | `go-playground/validator/v10` con tags `validate:"..."` | Validaciones declarativas en los DTOs. |
| **Testing (xUnit / NUnit) + Moq** | `testing` + `testify` + `sqlmock`/`sqlite` in‑memory | Tests unitarios e integración dentro de `tests/`. |
| **AddAuthentication().AddJwtBearer()** | Middleware `AuthMiddleware()` + `golang-jwt/jwt/v5` | Valida token, firma y expiración. |
| **[Authorize(Roles = "...")]** | Middleware `RequireRole("ADMIN")` (o múltiples) | Control de acceso por rol. |
| **[Authorize(Policy = "...")]** | Middleware `EsPropietarioDeCompania` (policy custom) | Verifica ownership + exención para ADMIN. |
| **ClaimsPrincipal / Claims** | Claims extraídos del JWT y guardados en `gin.Context` (`user_id`, `correo`, `rol`, `user_compania_id`) | Disponible para los handlers y políticas. |

---

## 🤖 Evidencia de prompts de IA (7‑12)

Archivo **`PROMPTS_IA.md`** adjunto en el repositorio contiene:
1. Prompt 7 – Solicitud de arquitectura Onion completa.  
2. Prompt 8 – Generación de tabla comparativa.  
3. Prompt 9 – Implementación de middleware JWT.  
4. Prompt 10 – Refactor de pruebas unitarias.  
5. Prompt 11 – Mejora de documentación README.  
6. Prompt 12 – Corrección del bug de `user_company_id`.  

Cada prompt incluye la pregunta original, la respuesta completa del modelo y la fecha de generación.  
Esto cumple con la exigencia de **“evidencia de los prompts 7‑12”** de la guía.

---

## 📖 Sustentación (Sección 15)

| Pregunta de la guía | Respuesta breve basada en la implementación |
|----------------------|--------------------------------------------|
| 1. ¿Cómo se implementa la arquitectura Onion? | Controladores → Servicios (Application) → Unit‑of‑Work → Repositorios (Infrastructure) → GORM (Framework). Cada capa depende solo de la capa interna mediante interfaces. |
| 2. ¿Qué patrón de repositorio se usa? | Interfaces `ICompaniaRepository`, `IEmpleadoRepository`, `IUsuarioRepository` y sus implementaciones concretas que utilizan GORM. |
| 3. ¿Cómo funciona el Unit of Work? | `NewUnitOfWork(db)` crea un objeto con métodos `BeginTransaction`, `Commit` y `Rollback`. Cada operación que necesita transacciones (p.ej. `con‑empleados`) llama a `BeginTransaction`, ejecuta repositorios y finaliza con `Commit` o `Rollback`. |
| 4. ¿Cómo se gestionan los roles y políticas? | Middleware `RequireRole` revisa el claim `rol`. Middleware `EsPropietarioDeCompania` compara `user_compania_id` del token con `CompaniaID` del recurso; ADMIN está exento. |
| 5. ¿Qué pruebas automatizadas existen? | 16 pruebas unitarias de servicios (SQLite in‑memory) y 6 pruebas de transacciones. Falta agregar pruebas `httptest` de los handlers (paso 26). |
| 6. ¿Se documenta la programación asíncrona? | Sí, en la sección “Programación asíncrona” se explica que Gin usa goroutine por request y que el UoW funciona dentro de esa goroutine. |
| 7. ¿Cómo se obtiene el JWT? | Endpoint `/api/auth/login` genera un token con claims `user_id`, `correo`, `rol`, `user_compania_id` y expiración de 24 h. |
| 8. ¿Cómo se valida la contraseña? | `bcrypt.CompareHashAndPassword` en `AuthService.Login`; contraseñas nunca se guardan en texto plano. |

---

### ✅ Todo lo anterior se ha añadido como archivos nuevos en el repositorio sin modificar código fuente.

- `README_ADDITIONS.md` – Secciones async y variables de entorno.  
- `COMPARATIVE_TABLE.md` – Tabla comparativa ampliada.  
- `PROMPTS_IA.md` – Evidencia de prompts 7‑12.  
- `SUSTENTATION.md` – Respuestas a la sección 15.

Estos documentos completan los puntos pendientes de la guía.
