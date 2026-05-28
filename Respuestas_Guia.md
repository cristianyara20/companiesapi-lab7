# Respuestas a la Guía de Actividad Práctica — Parte I y Parte II

Este documento contiene las respuestas a las preguntas obligatorias de la guía (Parte I: Pasos 1, 2 y 9 — Parte II: Sección 14 tabla comparativa ampliada y prompts 7–12), enfocadas en la tecnología **Go**, usando el framework **Gin** y el ORM **GORM**.

---

## 1. Tabla Comparativa (ASP.NET Core vs Go)
*(Correspondiente al Paso 1 de la guía)*

| Concepto en ASP.NET Core | Equivalente en la tecnología asignada (Go + GORM + Gin) |
| :--- | :--- |
| **Controller** | Handlers o Controladores de Gin (ej. `CompaniaController`) |
| **Entity** | `structs` de Go (ej. `type Compania struct`) |
| **DbContext** | Puntero a la base de datos de Gorm (`*gorm.DB`) |
| **DbSet** | Modelos pasados por referencia (ej. `db.Model(&Compania{})`) |
| **Migration** | Migraciones automáticas desde código: `db.AutoMigrate(...)` |
| **Fluent API** | Struct Tags (etiquetas) en los modelos (ej. `gorm:"primaryKey"`) |
| **Repository** | Interfaces de Go y structs concretos (ej. `CompaniaRepositoryImpl`) |
| **Unit of Work** | Implementación manual coordinando `db.Begin()`, `Commit()` y `Rollback()` |
| **Service Layer** | Structs en la capa Application (ej. `CompaniaService`) |
| **Dependency Injection** | Inyección manual usando constructores (ej. `NewCompaniaService(...)`) |
| **appsettings.json** | Archivo `.env` cargado con la librería `godotenv` |
| **Program.cs / Startup.cs** | Archivo `main.go` (función `main()`) |
| **Middleware** | Funciones middleware de Gin (ej. `gin.HandlerFunc`) |
| **Logging** | Librerías externas como `go.uber.org/zap` |

---

## 2. Preguntas sobre el ORM (Investigación Paso 2)

**¿Cuál ORM van a usar?**
GORM (gorm.io/gorm).

**¿Por qué ese ORM es adecuado?**
Es el ORM más popular, estable y recomendado para Go. Posee una sintaxis declarativa limpia (Developer Friendly), soporta asociaciones complejas (1:N, N:N), auto-migraciones, pre-carga inteligente (Eager Loading) y un manejo robusto de transacciones.

**¿Cómo se define una entidad?**
Se define como un `struct` nativo de Go. Las propiedades de la tabla (como llaves o tipos) se definen utilizando *struct tags* al lado de cada campo. Ejemplo: `Id uint `gorm:"primaryKey"``.

**¿Cómo se configura una relación uno a muchos?**
En la entidad "padre" (Compañía) se agrega un campo tipo *slice* (arreglo) de la entidad "hijo" (ej. `Empleados []Empleado`). En la entidad "hijo" se declara la llave foránea obligatoria (ej. `CompaniaId uint`).

**¿Cómo se hacen migraciones?**
Se llama a la función `AutoMigrate(&Compania{}, &Empleado{})` del objeto `db`. GORM se encarga de crear las tablas automáticamente o actualizar las columnas si no existen al iniciar la aplicación.

**¿Cómo se insertan datos iniciales?**
Creando los structs con los datos llenos y usando `db.Create(&companias)`. Si son varios, se pasa un *slice* de las entidades y GORM hace la inserción múltiple en bloque (Bulk Insert) optimizando la consulta.

**¿Cómo se realizan consultas básicas?**
Usando métodos encadenados en el objeto de base de datos. Ejemplo: `db.Find(&usuarios)` (traer todos), `db.First(&usuario, id)` (traer por ID) o `db.Where("cargo = ?", "Dev").Find(&usuarios)` (buscar con condición).

**¿Cómo maneja el ORM las transacciones?**
De dos formas: 
1. *Manualmente:* Iniciando la sesión transaccional con `tx := db.Begin()`, haciendo operaciones con esa sesión `tx` y finalizando explícitamente con `tx.Commit()` o `tx.Rollback()`.
2. *Automáticamente:* Utilizando un bloque anónimo `db.Transaction(func(tx *gorm.DB) error { ... })`, donde el commit o rollback se infiere automáticamente si la función interna devuelve un error.

**¿El ORM implementa Unit of Work internamente?**
Sí, implícitamente cuando se utiliza el método de transacciones. Sin embargo, para aplicar rigurosamente la Clean Architecture a nivel de negocio y aislar los repositorios, la mejor práctica es encapsularlo en una interfaz Unit of Work propia administrada por la capa Infrastructure.

---

## 3. Evidencia Paso 6: Crear Migraciones
*(Correspondiente al Paso 6 de la guía)*

* **Comando para crear migración:**
  En Go con GORM no se utiliza una herramienta CLI externa para generar archivos estáticos de migración. La migración se define directamente en código declarando las entidades en la función `AutoMigrate` dentro del archivo `infrastructure/database/connection.go`:
  ```go
  db.AutoMigrate(&entities.Compania{}, &entities.Empleado{})
  ```

* **Comando para aplicar migración:**
  Se ejecuta automáticamente al iniciar la aplicación desde la consola:
  ```bash
  go run main.go
  ```
  Esto levanta la conexión y ejecuta la función `database.Migrate(db)`.

* **Tablas generadas:**
  * `companias`
  * `empleados`

* **Relación entre Compañía y Empleado:**
  Relación de **Uno a Muchos (1:N)**. Una compañía tiene muchos empleados, y un empleado pertenece obligatoriamente a una única compañía a través de la llave foránea `CompaniaID`.
  * En `Compania`: `Empleados []Empleado gorm:"foreignKey:CompaniaID"`
  * En `Empleado`: `CompaniaID uint gorm:"not null;index"`

---

## 4. Preguntas obligatorias sobre Unit of Work (Paso 9)

**1. ¿Qué es Unit of Work?**
Es un patrón de diseño arquitectónico que agrupa un conjunto de operaciones en la base de datos (crear, actualizar, borrar) dentro de una única transacción. Asegura que todas las operaciones se guarden juntas o se cancelen juntas, manteniendo los datos siempre consistentes y atómicos.

**2. ¿Qué problema resuelve?**
Previene la pérdida de integridad referencial o estados corruptos (datos guardados parcialmente). Por ejemplo, si se crea una empresa pero falla la creación de sus empleados, evita que la empresa se guarde vacía y sola, garantizando que el sistema sea confiable frente a fallos de red, caídas o reglas de validación.

**3. ¿Qué relación tiene con Repository Pattern?**
El Unit of Work actúa como el "director de orquesta" de los Repositorios. Varios repositorios comparten el mismo objeto de transacción proveído por el Unit of Work, asegurando que todos operen coordinadamente sobre el mismo contexto y la misma conexión transaccional temporal.

**4. ¿El ORM seleccionado ya implementa Unit of Work internamente?**
Sí, GORM maneja sus propias transacciones internas y agrupa consultas vinculadas. No obstante, implementar el patrón Unit of Work de manera explícita (como se hizo en el proyecto) permite que la capa de Dominio y la capa de Aplicación no dependan en absoluto de GORM.

**5. ¿Qué objeto representa la unidad de trabajo en Go?**
El puntero o sesión de transacción a la base de datos de GORM (`tx *gorm.DB`), el cual se inicializa al llamar a `db.Begin()`.

**6. ¿Dónde se ubica Unit of Work dentro de Onion Architecture?**
Su **interfaz abstracta** (los contratos de lo que debe hacer: Commit, Rollback, Begin) se ubica en el núcleo de la arquitectura (`Domain`). Su **implementación técnica** y real (donde se toca directamente GORM y la BD) se ubica en la capa más externa: `Infrastructure`.

**7. ¿Los repositorios llaman directamente a Save, Commit o Flush?**
¡No! Los repositorios solo ejecutan operaciones CRUD pasivas en la memoria transaccional (`Create`, `Update`). Quien decide en qué momento exacto confirmar todos esos cambios y llamar a `Commit()` es el Unit of Work, orquestado desde la capa de Servicios (`Application`).

**8. ¿Cómo se revierte una operación cuando ocurre un error?**
Se ejecuta la instrucción `uow.Rollback()`. Esto envía una instrucción a la base de datos que descarta y borra inmediatamente todos los comandos SQL de inserción o actualización que estuvieran temporalmente en caché o memoria transaccional.

**9. ¿Cómo se garantiza que varias operaciones se guarden como una sola unidad?**
Asegurándose de que, al llamar a `BeginTransaction()`, el Unit of Work asigne la **misma sesión transaccional temporal** a todos los repositorios involucrados. Así, cualquier query ejecutada por cualquier repositorio viaja exactamente por el mismo túnel hacia la base de datos.

**10. ¿Qué ventajas tiene usar Unit of Work en una API empresarial?**
- Reduce drásticamente los cuellos de botella al abrir y cerrar conexiones innecesarias a la BD.
- Previene inconsistencias graves en datos relacionales o financieros.
- Facilita el rollback atómico en lógicas de negocio extensas o de múltiples pasos.
- Permite que el código sea testeable (Mock Unit Of Work) sin tocar bases de datos reales.

---

## 5. Tabla Comparativa Ampliada — Parte II (Sección 14)

*(Conceptos nuevos introducidos en la Parte II)*

| Concepto en ASP.NET Core | Equivalente en la tecnología asignada (Go + GORM + Gin) |
| :--- | :--- |
| **Endpoints de colección** (`IEnumerable<T>` / `List<T>`) | Slices de Go (`[]Entity`) devueltos en un envelope JSON con metadatos: `{datos, pagina, tamano, total, totalPaginas}` |
| **Paginación** (`Skip(n).Take(m)` en LINQ) | `query.Offset((pagina-1)*tamano).Limit(tamano)` en GORM + `Count(&total)` previo |
| **Ordenamiento dinámico** (`OrderBy(campo)`) | `query.Order(ordenCol + " " + ordenDir)` — columna sanitizada contra un mapa permitido para prevenir SQL Injection |
| **`async` / `await` + `Task<T>`** | No existe equivalente directo. Gin despacha cada petición en una **goroutine** automáticamente. El código es sincrónico dentro de la goroutine pero concurrente a nivel de servidor. |
| **`CancellationToken`** | `context.Context` propagado en todas las llamadas: `r.db.WithContext(ctx).Find(...)` |
| **`DataAnnotations`** (`[Required]`, `[EmailAddress]`, `[Range]`) | Struct tags del paquete `validator`: `` `validate:"required,email,gt=0"` `` |
| **`FluentValidation`** (reglas de negocio con mensajes) | Validación manual en la capa de Servicios con consultas al repositorio + tipo `ValidationError` personalizado |
| **`xUnit` / `NUnit`** | Paquete `testing` de la stdlib de Go (`func TestXxx(t *testing.T)`) |
| **`Moq`** (mocks de repositorios) | Base de datos SQLite en memoria (`glebarez/sqlite` + `:memory:`) usando los repositorios reales — no se necesitan mocks |
| **`AddAuthentication().AddJwtBearer()`** | `github.com/golang-jwt/jwt/v5` + middleware `AuthMiddleware()` registrado en Gin con `router.Use(...)` |
| **`[Authorize(Roles = "ADMIN")]`** | Middleware `RequireRole("ADMIN")` aplicado a grupos de rutas: `comp.Use(authGuard, adminRoleGuard)` |
| **`[Authorize(Policy = "EsPropietario")]`** | Middleware `EsPropietarioDeCompania(empleadoService)` — función `gin.HandlerFunc` que evalúa claims del token contra el recurso solicitado |
| **`IAuthorizationHandler`** / **`IAuthorizationRequirement`** | No existe concepto formal. El handler es el propio middleware: recibe el contexto de Gin, extrae claims, consulta el repositorio y decide permitir o denegar. |
| **`ClaimsPrincipal`** / **`HttpContext.User.Claims`** | Claims guardados en el contexto de Gin por `AuthMiddleware()`: `ctx.Set("user_id", ...)`, recuperados con `ctx.Get("user_id")` |
| **`PasswordHasher<T>`** (ASP.NET Identity) | `golang.org/x/crypto/bcrypt` — `bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)` y `bcrypt.CompareHashAndPassword(hash, pass)` |

---

## 6. Evidencia de Prompts de IA — Parte II (Prompts 7 al 12)

### Prompt 7 — CRUD de colecciones

> "Tengo una API REST en Go/Gin con GORM que ya hace CRUD de entidades individuales con Onion Architecture y Unit of Work. Explícame cómo agregar operaciones sobre colecciones: creación masiva (bulk insert), actualización parcial con PATCH, eliminación múltiple, y cómo añadir paginación, filtrado y ordenamiento al endpoint de listado, respetando las capas y sin que el repositorio confirme la transacción."

**Resultado aplicado:** `CreateRange`, `DeleteRange`, `GetPaged`, `PatchPartial` en `empleado_repository_impl.go`. El servicio orquesta `BeginTransaction()` + `Commit()` / `Rollback()`.

---

### Prompt 8 — Programación asíncrona

> "En Go/Gin, ¿el manejo de peticiones y el acceso a datos con GORM es sincrónico o asíncrono? Si soporta asincronía, muéstrame cómo refactorizar un servicio y un repositorio de sincrónico a asíncrono (goroutines o el mecanismo correspondiente) y qué precauciones tomar con la sesión y la transacción del Unit of Work."

**Resultado aplicado:** Se determinó que Gin despacha cada petición en una goroutine automáticamente. GORM es sincrónico por llamada pero no bloquea el servidor. Se propaga `context.Context` en todas las capas (`WithContext(ctx)`). No se requiere refactorización adicional.

---

### Prompt 9 — Validaciones

> "¿Cuál es la librería o mecanismo de validación recomendado para Go/Gin? Muéstrame cómo validar los DTOs de Compañía y Empleado (campos obligatorios, longitud, formato de correo, salario positivo) y cómo devolver el error de validación con el código HTTP adecuado, manteniendo la validación en la capa de aplicación."

**Resultado aplicado:** `github.com/go-playground/validator/v10` con struct tags. Centralizador `ValidateStruct()` en `application/validation/validator.go`. Formato de error unificado `{mensaje, errores: [{campo, detalle}]}`.

---

### Prompt 10 — Pruebas

> "¿Cuál es el framework de pruebas más usado en Go? Muéstrame cómo escribir pruebas unitarias para un servicio y pruebas de integración para los endpoints de Compañía y Empleado, incluyendo una prueba que verifique el rollback del endpoint transaccional cuando falla la creación de un empleado."

**Resultado aplicado:** Paquete `testing` de la stdlib + `glebarez/sqlite` para BD in-memory. 16+ pruebas en `tests/services_test.go`. Pruebas de rollback transaccional con verificación directa en BD.

---

### Prompt 11 — JWT por roles

> "Explícame cómo implementar autenticación con JWT en Go/Gin: entidad Usuario, registro, login que devuelve un token y cómo proteger endpoints exigiendo un rol específico (p. ej. solo ADMIN puede eliminar). Indica la librería recomendada y dónde ubicar cada parte en Onion Architecture."

**Resultado aplicado:** `github.com/golang-jwt/jwt/v5` + `golang.org/x/crypto/bcrypt`. Entidad `Usuario` en Domain, `AuthService` en Application, repositorio en Infrastructure, `AuthController` + `AuthMiddleware()` + `RequireRole()` en API/Presentation.

---

### Prompt 12 — JWT por políticas

> "¿Cuál es la diferencia entre autorización por roles y por políticas (claims/policies) en Go/Gin? Muéstrame cómo implementar una política basada en reglas o claims (ej.: un usuario solo puede editar empleados de su propia compañía) y compárala con las policies de ASP.NET Core ([Authorize(Policy=...)] con requirements y handlers)."

**Resultado aplicado:** Middleware `EsPropietarioDeCompania(empleadoService)` en `auth_middleware.go`. Evalúa el claim `user_compania_id` del token contra `empleado.CompaniaID` del recurso. ADMIN está exento. Equivalente manual del `IAuthorizationHandler` de ASP.NET Core sin necesitar Casbin.
