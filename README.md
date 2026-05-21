# API de Compañías y Empleados

Este proyecto es una API REST funcional construida en **Go (Golang)** aplicando los patrones de diseño **Onion Architecture, Repository Pattern y Unit of Work**. Su propósito es demostrar la transferencia de conceptos arquitectónicos aprendidos en ASP.NET Core a un ecosistema tecnológico alternativo.

## Tecnología usada
- **Lenguaje:** Go (Golang)
- **Framework Web:** Gin (github.com/gin-gonic/gin)
- **Base de Datos:** PostgreSQL (Alojada en Supabase)

## ORM usado
**GORM** (gorm.io/gorm). Es el ORM más popular y recomendado en el ecosistema Go. Se eligió porque ofrece una API amigable, facilita la creación de relaciones complejas (uno a muchos), posee un mecanismo nativo y seguro de transacciones, y permite realizar migraciones automáticas (`AutoMigrate`) de forma sencilla.

## Arquitectura aplicada
Se aplica la **Onion Architecture** (Arquitectura de Cebolla) garantizando una estricta separación de responsabilidades:
- **Core / Domain:** No depende de nada.
- **Application:** Depende solo de Domain.
- **Infrastructure / API:** Dependen de Application y Domain.

## Estructura del proyecto
- `/domain`: Núcleo del sistema. Contiene entidades (modelos de datos) y las interfaces de los repositorios y del Unit of Work.
- `/application`: Servicios de negocio y DTOs (Data Transfer Objects). Orquesta la lógica y el Unit of Work.
- `/infrastructure`: Implementaciones concretas. Contiene la conexión a base de datos, migraciones, el seeding de datos iniciales y los repositorios reales.
- `/api`: Controladores, middlewares y definición de rutas REST usando Gin.
- `main.go`: Punto de entrada de la aplicación e inyección de dependencias.

## Entidades
- **Compañía:** `Id`, `Nombre`, `Direccion`, `Telefono`, `FechaCreacion`.
- **Empleado:** `Id`, `Nombre`, `Apellido`, `Correo`, `Cargo`, `Salario`, `CompaniaId`.

## Relación entre entidades
**Uno a Muchos (1:N)**. 
Una compañía puede tener múltiples empleados asociados, mientras que un empleado pertenece a una sola compañía obligatoriamente. Se modela mediante la llave foránea `CompaniaId` dentro del struct de Empleado.

## Repository Pattern
Se crearon interfaces abstractas en el Dominio (`CompaniaRepository` y `EmpleadoRepository`) con operaciones como `GetAll`, `GetById`, `Create`, `Update` y `Delete`. Las capas superiores (Servicios) operan sobre estas interfaces sin tener conocimiento de las consultas SQL ni del ORM subyacente.

## Unit of Work

### ¿Qué es Unit of Work?
Es un patrón de diseño que agrupa un conjunto de operaciones que modifican la base de datos dentro de una única transacción. Si todas las operaciones son exitosas, los datos se confirman (Commit). Si ocurre algún error en medio del proceso, toda la operación se descarta (Rollback) para mantener la integridad de la base de datos.

### ¿Cómo se implementó en esta tecnología?
En Go con GORM, se creó la estructura `UnitOfWorkImpl`. Expone un método `BeginTransaction()` que solicita una transacción a GORM (`db.Begin()`). Cuando este método es llamado, todos los repositorios creados comparten esa misma sesión en curso (`tx`), asegurando que todos operen sobre la misma transacción.

### ¿Cómo se manejan las transacciones?
El manejador principal de la transacción es el Unit of Work. A nivel de negocio, la capa de Servicios llama a `BeginTransaction()`, realiza sus operaciones, y finaliza explícitamente llamando a `Commit()` o `Rollback()`.

### ¿Cómo se hace commit?
Se invoca `uow.Commit()`, el cual ejecuta internamente `tx.Commit().Error` en GORM y aplica todos los cambios físicamente en la BD simultáneamente.

### ¿Cómo se hace rollback?
Se invoca `uow.Rollback()`, que llama a `tx.Rollback().Error` en GORM y descarta todos los cambios realizados en memoria.

## Tabla comparativa con ASP.NET Core
| Concepto en ASP.NET Core | Equivalente en Go (con GORM y Gin) |
| ------------------------ | ------------------------ |
| Controller | Gin Handlers (Controladores como structs) |
| Entity | Structs de Go (ej. `type Compania struct`) |
| DbContext | El objeto conexión de base de datos (`*gorm.DB`) |
| DbSet | Modelos invocados desde la conexión (ej. `db.Model(&Entity{})`) |
| Migration | Migraciones en código mediante `db.AutoMigrate(&Entities...)` |
| Fluent API | Etiquetas o Struct Tags (ej. `gorm:"primaryKey"`) |
| Repository | Interfaces y structs inyectados (ej. `CompaniaRepositoryImpl`) |
| Unit of Work | Struct personalizado coordinando `db.Begin()`, `Commit()` y `Rollback()` |
| Service Layer | Structs en la capa de Application (`CompaniaService`) |
| Dependency Injection | Inyección manual usando funciones constructoras (`NewService(...)`) |
| appsettings.json | Archivo `.env` + paquete `godotenv` |
| Program.cs / Startup.cs | Archivo `main.go` |
| Middleware | Funciones middleware de Gin (`gin.HandlerFunc`) |
| Logging | Paquete externo `go.uber.org/zap` |

## Endpoints

**Compañías**
- `GET /api/companias` (Listar)
- `GET /api/companias/:id` (Obtener)
- `POST /api/companias` (Crear)
- `PUT /api/companias/:id` (Actualizar)
- `DELETE /api/companias/:id` (Eliminar)
- `GET /api/companias/:id/empleados` (Obtener empleados por compañía)

**Empleados**
- `GET /api/empleados` (Listar)
- `GET /api/empleados/:id` (Obtener)
- `POST /api/empleados` (Crear)
- `PUT /api/empleados/:id` (Actualizar)
- `DELETE /api/empleados/:id` (Eliminar)

## Endpoint transaccional
`POST /api/companias/con-empleados`
Se encarga de crear de manera íntegra y atómica una compañía junto con una lista de empleados. Utiliza el patrón **Unit of Work** para asegurar que, si existe cualquier error durante la creación (ej. correo duplicado de un empleado), no se guardará ni la compañía, ni el resto de empleados.

## Instalación y Configuración de base de datos
1. Tener Go instalado (v1.20+).
2. Clonar/descargar el repositorio en tu máquina.
3. Crear un archivo llamado `.env` en la raíz del proyecto con la siguiente estructura:
   ```env
   DB_HOST=db.tuproyecto.supabase.co
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=TuPasswordSeguro
   DB_NAME=postgres
   PORT=8080
   ```
4. Ejecutar el comando `go mod tidy` para instalar dependencias (Gin, Gorm, Zap, etc.).

## Migraciones y Ejecución del proyecto
1. En la raíz del proyecto abre tu terminal y ejecuta:
   ```bash
   go run main.go
   ```
2. GORM se encargará automáticamente de crear y aplicar las migraciones a PostgreSQL (`AutoMigrate`). 
3. El archivo `seed.go` validará si la BD está vacía y llenará las tablas con **3 compañías y 10 empleados** de demostración.
4. El servidor se ejecutará en `http://localhost:8080`.

## Pruebas con Swagger/Postman
El servidor está expuesto en HTTP. Puedes usar herramientas como Postman o Thunder Client e importar las colecciones apuntando a `http://localhost:8080/api/...`. El formato del `body` en los `POST` debe ser JSON puro.

## Logging
El proyecto implementa un registro de eventos asíncrono y de alto rendimiento utilizando `go.uber.org/zap`. Todas las operaciones CRUD críticas y los inicios de transacciones, éxitos y fallos registran su estado en la consola para facilitar su rastreo.

## Uso de IA
Se utilizaron herramientas de IA (Google Deepmind) para transferir y equiparar conceptos abstractos de C# / ASP.NET a una estructura de Go idiomática, optimizar la creación de Struct Tags, guiar la correcta implementación del patrón transaccional, generar scripts automatizados para correcciones de base de datos y elaborar la documentación del sistema.

## Conclusiones
Se demostró de forma contundente que los principios de Clean Architecture y patrones de persistencia (Repository y Unit of Work) no están atados a un lenguaje específico como C# o Java. En Go, aunque la inyección de dependencias es manual y los objetos de base de datos se manejan explícitamente mediante punteros, la arquitectura resultante es igual de modular, altamente escalable y mucho más liviana.
