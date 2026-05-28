package tests

import (
	"companies-api/application/dtos"
	"companies-api/application/services"
	"companies-api/domain/entities"
	"companies-api/infrastructure/database"
	unitofwork "companies-api/infrastructure/unit_of_work"
	"context"
	"testing"

	"go.uber.org/zap"
	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB crea una base de datos SQLite en memoria para pruebas
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Error abriendo SQLite en memoria: %v", err)
	}

	// Migrar las tablas
	database.Migrate(db)
	return db
}

// setupServices crea los servicios con UoW conectado a la DB de prueba
func setupServices(db *gorm.DB) (*services.CompaniaService, *services.EmpleadoService, *services.AuthService) {
	log, _ := zap.NewDevelopment()
	uow := unitofwork.NewUnitOfWork(db)
	return services.NewCompaniaService(uow, log),
		services.NewEmpleadoService(uow, log),
		services.NewAuthService(uow, log)
}

// ─── Tests Unitarios de CompaniaService ──────────────────────────────────────

func TestCompaniaService_Create_Success(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	dto := dtos.CreateCompaniaDTO{
		Nombre:    "Test Corp",
		Direccion: "Calle 123",
		Telefono:  "3001234567",
	}
	c, err := compSvc.Create(ctx, dto)
	if err != nil {
		t.Fatalf("Error creando compañía: %v", err)
	}
	if c.ID == 0 {
		t.Error("ID de compañía debería ser > 0")
	}
	if c.Nombre != "Test Corp" {
		t.Errorf("Nombre esperado 'Test Corp', obtenido '%s'", c.Nombre)
	}
}

func TestCompaniaService_Create_ValidationFail(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	// Nombre vacío → debe fallar validación
	dto := dtos.CreateCompaniaDTO{
		Nombre:    "",
		Direccion: "Calle 123",
		Telefono:  "3001234567",
	}
	_, err := compSvc.Create(ctx, dto)
	if err == nil {
		t.Error("Se esperaba error de validación, pero no se obtuvo ninguno")
	}
}

func TestCompaniaService_GetById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	_, err := compSvc.GetById(ctx, 999)
	if err == nil {
		t.Error("Se esperaba error 'compañía no encontrada'")
	}
}

func TestCompaniaService_Delete(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	// Crear primero
	c, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Para eliminar", Direccion: "Dir", Telefono: "1234567",
	})

	err := compSvc.Delete(ctx, c.ID)
	if err != nil {
		t.Fatalf("Error eliminando compañía: %v", err)
	}

	// Verificar que ya no existe
	_, err = compSvc.GetById(ctx, c.ID)
	if err == nil {
		t.Error("Compañía debería haberse eliminado")
	}
}

// ─── Tests Unitarios de EmpleadoService ──────────────────────────────────────

func TestEmpleadoService_Create_Success(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	// Crear compañía para asociar
	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Empresa Test", Direccion: "Dir 1", Telefono: "1234567",
	})

	dto := dtos.CreateEmpleadoDTO{
		Nombre:     "Juan",
		Apellido:   "Pérez",
		Correo:     "juan@test.com",
		Cargo:      "Dev",
		Salario:    3000000,
		CompaniaID: comp.ID,
	}
	emp, err := empSvc.Create(ctx, dto)
	if err != nil {
		t.Fatalf("Error creando empleado: %v", err)
	}
	if emp.ID == 0 {
		t.Error("ID de empleado debería ser > 0")
	}
}

func TestEmpleadoService_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Corp", Direccion: "Dir", Telefono: "1234567",
	})

	dto := dtos.CreateEmpleadoDTO{
		Nombre: "Ana", Apellido: "Gómez", Correo: "ana@test.com",
		Cargo: "QA", Salario: 2500000, CompaniaID: comp.ID,
	}
	_, _ = empSvc.Create(ctx, dto) // Primera vez OK

	// Segunda vez con mismo correo → debe fallar
	dto2 := dtos.CreateEmpleadoDTO{
		Nombre: "Otra Ana", Apellido: "López", Correo: "ana@test.com",
		Cargo: "Dev", Salario: 3000000, CompaniaID: comp.ID,
	}
	_, err := empSvc.Create(ctx, dto2)
	if err == nil {
		t.Error("Se esperaba error de correo duplicado")
	}
}

func TestEmpleadoService_Create_CompaniaNoExiste(t *testing.T) {
	db := setupTestDB(t)
	_, empSvc, _ := setupServices(db)
	ctx := context.Background()

	dto := dtos.CreateEmpleadoDTO{
		Nombre: "Test", Apellido: "Test", Correo: "test@test.com",
		Cargo: "Dev", Salario: 1000000, CompaniaID: 999,
	}
	_, err := empSvc.Create(ctx, dto)
	if err == nil {
		t.Error("Se esperaba error 'compañía no existe'")
	}
}

func TestEmpleadoService_Create_SalarioNegativo(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Corp", Direccion: "Dir", Telefono: "1234567",
	})

	dto := dtos.CreateEmpleadoDTO{
		Nombre: "Test", Apellido: "Test", Correo: "neg@test.com",
		Cargo: "Dev", Salario: -500, CompaniaID: comp.ID,
	}
	_, err := empSvc.Create(ctx, dto)
	if err == nil {
		t.Error("Se esperaba error de validación por salario negativo")
	}
}

func TestEmpleadoService_PatchPartial(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Corp", Direccion: "Dir", Telefono: "1234567",
	})

	emp, _ := empSvc.Create(ctx, dtos.CreateEmpleadoDTO{
		Nombre: "Original", Apellido: "Test", Correo: "patch@test.com",
		Cargo: "Dev", Salario: 2000000, CompaniaID: comp.ID,
	})

	nuevoNombre := "Modificado"
	updated, err := empSvc.PatchPartial(ctx, emp.ID, dtos.UpdateEmpleadoDTO{
		Nombre: &nuevoNombre,
	})
	if err != nil {
		t.Fatalf("Error en PatchPartial: %v", err)
	}
	if updated.Nombre != "Modificado" {
		t.Errorf("Nombre esperado 'Modificado', obtenido '%s'", updated.Nombre)
	}
	// El apellido no cambió
	if updated.Apellido != "Test" {
		t.Errorf("Apellido no debió cambiar, obtenido '%s'", updated.Apellido)
	}
}

func TestEmpleadoService_GetPaged(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Corp", Direccion: "Dir", Telefono: "1234567",
	})

	// Insertar 5 empleados
	for i := 0; i < 5; i++ {
		empSvc.Create(ctx, dtos.CreateEmpleadoDTO{
			Nombre: "Emp", Apellido: "Test",
			Correo: "emp" + string(rune('a'+i)) + "@test.com",
			Cargo: "Dev", Salario: 1000000, CompaniaID: comp.ID,
		})
	}

	// Pedir página 1, tamaño 3
	list, total, err := empSvc.GetPaged(ctx, 1, 3, "id", "asc", "")
	if err != nil {
		t.Fatalf("Error en GetPaged: %v", err)
	}
	if total != 5 {
		t.Errorf("Total esperado 5, obtenido %d", total)
	}
	if len(list) != 3 {
		t.Errorf("Tamaño de página esperado 3, obtenido %d", len(list))
	}
}

// ─── Tests Unitarios de AuthService ──────────────────────────────────────────

func TestAuthService_Register_Success(t *testing.T) {
	db := setupTestDB(t)
	_, _, authSvc := setupServices(db)
	ctx := context.Background()

	dto := dtos.RegisterDTO{
		Nombre:     "Admin Test",
		Correo:     "admin@test.com",
		Contrasena: "password123",
		Rol:        "ADMIN",
	}
	user, err := authSvc.Register(ctx, dto)
	if err != nil {
		t.Fatalf("Error registrando usuario: %v", err)
	}
	if user.ID == 0 {
		t.Error("ID de usuario debería ser > 0")
	}
	if user.Rol != "ADMIN" {
		t.Errorf("Rol esperado 'ADMIN', obtenido '%s'", user.Rol)
	}
	// La contraseña hash debe existir (pero NO es la contraseña en texto plano)
	if user.ContrasenaHash == "password123" {
		t.Error("La contraseña NO debería guardarse en texto plano")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	_, _, authSvc := setupServices(db)
	ctx := context.Background()

	dto := dtos.RegisterDTO{
		Nombre: "User1", Correo: "dup@test.com", Contrasena: "pass123", Rol: "USUARIO",
	}
	authSvc.Register(ctx, dto)

	// Segundo registro con mismo correo
	dto2 := dtos.RegisterDTO{
		Nombre: "User2", Correo: "dup@test.com", Contrasena: "pass456", Rol: "USUARIO",
	}
	_, err := authSvc.Register(ctx, dto2)
	if err == nil {
		t.Error("Se esperaba error de correo duplicado")
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	db := setupTestDB(t)
	_, _, authSvc := setupServices(db)
	ctx := context.Background()

	// Registrar
	authSvc.Register(ctx, dtos.RegisterDTO{
		Nombre: "Login Test", Correo: "login@test.com", Contrasena: "mipass123", Rol: "ADMIN",
	})

	// Login
	token, user, err := authSvc.Login(ctx, dtos.LoginDTO{
		Correo: "login@test.com", Contrasena: "mipass123",
	})
	if err != nil {
		t.Fatalf("Error en login: %v", err)
	}
	if token == "" {
		t.Error("Token JWT no debería estar vacío")
	}
	if user.Correo != "login@test.com" {
		t.Errorf("Correo esperado 'login@test.com', obtenido '%s'", user.Correo)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	_, _, authSvc := setupServices(db)
	ctx := context.Background()

	authSvc.Register(ctx, dtos.RegisterDTO{
		Nombre: "User", Correo: "wrong@test.com", Contrasena: "correcta123", Rol: "USUARIO",
	})

	_, _, err := authSvc.Login(ctx, dtos.LoginDTO{
		Correo: "wrong@test.com", Contrasena: "incorrecta",
	})
	if err == nil {
		t.Error("Se esperaba error de credenciales inválidas")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	_, _, authSvc := setupServices(db)
	ctx := context.Background()

	_, _, err := authSvc.Login(ctx, dtos.LoginDTO{
		Correo: "noexiste@test.com", Contrasena: "loquesea",
	})
	if err == nil {
		t.Error("Se esperaba error de usuario no encontrado")
	}
}

// ─── Test Transaccional (Caso Integración) ──────────────────────────────────

func TestCompaniaService_CreateConEmpleados_Success(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	dto := dtos.CreateCompaniaConEmpleadosDTO{
		Nombre:    "Transaccional Corp",
		Direccion: "Calle 1",
		Telefono:  "1234567",
		Empleados: []dtos.CreateEmpleadoDTO{
			// Note: CompaniaID is typically injected by the service during the transaction, 
			// but since it's required by validation, we need to bypass it or set a dummy value that passes validation (>0).
			// Let's set it to a dummy value 999 since the service will override it anyway with the actual created company ID.
			{Nombre: "Emp1", Apellido: "A", Correo: "e1@t.com", Cargo: "Dev", Salario: 1000000, CompaniaID: 999},
			{Nombre: "Emp2", Apellido: "B", Correo: "e2@t.com", Cargo: "QA", Salario: 2000000, CompaniaID: 999},
		},
	}
	comp, err := compSvc.CreateConEmpleados(ctx, dto)
	if err != nil {
		t.Fatalf("Error en CreateConEmpleados: %v", err)
	}
	if comp.ID == 0 {
		t.Error("ID de compañía debería ser > 0")
	}
}

func TestCompaniaService_CreateConEmpleados_Rollback_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	// Pre-insertar un empleado con correo existente
	preComp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Pre Corp", Direccion: "Dir", Telefono: "1234567",
	})
	empSvc.Create(ctx, dtos.CreateEmpleadoDTO{
		Nombre: "Existente", Apellido: "X", Correo: "duplicado@t.com",
		Cargo: "Dev", Salario: 1000000, CompaniaID: preComp.ID,
	})

	// Intentar crear compañía con un empleado que tiene correo duplicado
	dto := dtos.CreateCompaniaConEmpleadosDTO{
		Nombre:    "Rollback Corp",
		Direccion: "Calle 2",
		Telefono:  "7654321",
		Empleados: []dtos.CreateEmpleadoDTO{
			{Nombre: "OK", Apellido: "A", Correo: "ok@t.com", Cargo: "Dev", Salario: 1000000, CompaniaID: 999},
			{Nombre: "Dup", Apellido: "B", Correo: "duplicado@t.com", Cargo: "QA", Salario: 2000000, CompaniaID: 999},
		},
	}
	_, err := compSvc.CreateConEmpleados(ctx, dto)
	if err == nil {
		t.Error("Se esperaba error por correo duplicado — ROLLBACK")
	}

	// Verificar que la compañía "Rollback Corp" NO se creó (rollback exitoso)
	var count int64
	db.Model(&entities.Compania{}).Where("nombre = ?", "Rollback Corp").Count(&count)
	if count != 0 {
		t.Error("La compañía 'Rollback Corp' NO debió crearse (rollback falló)")
	}
}

func TestCompaniaService_CreateConEmpleados_Rollback_SalarioNegativo(t *testing.T) {
	db := setupTestDB(t)
	compSvc, _, _ := setupServices(db)
	ctx := context.Background()

	dto := dtos.CreateCompaniaConEmpleadosDTO{
		Nombre:    "Salario Neg Corp",
		Direccion: "Dir",
		Telefono:  "1234567",
		Empleados: []dtos.CreateEmpleadoDTO{
			{Nombre: "OK", Apellido: "A", Correo: "salok@t.com", Cargo: "Dev", Salario: 1000000, CompaniaID: 999},
			{Nombre: "Bad", Apellido: "B", Correo: "salbad@t.com", Cargo: "QA", Salario: -500, CompaniaID: 999},
		},
	}
	_, err := compSvc.CreateConEmpleados(ctx, dto)
	if err == nil {
		t.Error("Se esperaba error de validación por salario negativo — ROLLBACK")
	}

	// Verificar rollback
	var count int64
	db.Model(&entities.Compania{}).Where("nombre = ?", "Salario Neg Corp").Count(&count)
	if count != 0 {
		t.Error("La compañía NO debió crearse (rollback falló)")
	}
}

// ─── Test de Bulk Operations ────────────────────────────────────────────────

func TestEmpleadoService_CreateRange_Success(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Bulk Corp", Direccion: "Dir", Telefono: "1234567",
	})

	list := []dtos.CreateEmpleadoDTO{
		{Nombre: "B1", Apellido: "A", Correo: "b1@test.com", Cargo: "Dev", Salario: 1000000, CompaniaID: comp.ID},
		{Nombre: "B2", Apellido: "B", Correo: "b2@test.com", Cargo: "QA", Salario: 2000000, CompaniaID: comp.ID},
		{Nombre: "B3", Apellido: "C", Correo: "b3@test.com", Cargo: "PM", Salario: 3000000, CompaniaID: comp.ID},
	}

	result, err := empSvc.CreateRange(ctx, list)
	if err != nil {
		t.Fatalf("Error en CreateRange: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Se esperaban 3 empleados creados, obtenidos %d", len(result))
	}
}

func TestEmpleadoService_DeleteRange_Success(t *testing.T) {
	db := setupTestDB(t)
	compSvc, empSvc, _ := setupServices(db)
	ctx := context.Background()

	comp, _ := compSvc.Create(ctx, dtos.CreateCompaniaDTO{
		Nombre: "Del Corp", Direccion: "Dir", Telefono: "1234567",
	})

	emp1, _ := empSvc.Create(ctx, dtos.CreateEmpleadoDTO{
		Nombre: "D1", Apellido: "A", Correo: "d1@test.com", Cargo: "Dev", Salario: 1000000, CompaniaID: comp.ID,
	})
	emp2, _ := empSvc.Create(ctx, dtos.CreateEmpleadoDTO{
		Nombre: "D2", Apellido: "B", Correo: "d2@test.com", Cargo: "QA", Salario: 2000000, CompaniaID: comp.ID,
	})

	err := empSvc.DeleteRange(ctx, []uint{emp1.ID, emp2.ID})
	if err != nil {
		t.Fatalf("Error en DeleteRange: %v", err)
	}

	// Verificar que ya no existen
	_, err1 := empSvc.GetById(ctx, emp1.ID)
	_, err2 := empSvc.GetById(ctx, emp2.ID)
	if err1 == nil || err2 == nil {
		t.Error("Los empleados deberían haberse eliminado")
	}
}
