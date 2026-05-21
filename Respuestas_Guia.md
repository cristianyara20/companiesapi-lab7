# Respuestas a la Guía de Actividad Práctica

Este documento contiene las respuestas a las preguntas obligatorias de la guía (Pasos 1, 2 y 9), enfocadas en la tecnología **Go**, usando el framework **Gin** y el ORM **GORM**.

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

## 3. Preguntas obligatorias sobre Unit of Work (Paso 9)

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
