package database

import (
	"companies-api/domain/entities"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed inserta datos iniciales si las tablas están vacías
func Seed(db *gorm.DB) {
	// ── 1. Seed de Compañías y Empleados ──────────────────────────────────
	var count int64
	db.Model(&entities.Compania{}).Count(&count)
	if count == 0 {
		log.Println("🌱 Insertando compañías iniciales...")

		// 3 compañías
		companias := []entities.Compania{
			{Nombre: "Tech Solutions S.A.S", Direccion: "Calle 45 # 10-20, Bogotá", Telefono: "3001234567"},
			{Nombre: "Innovatech Ltda", Direccion: "Carrera 15 # 88-64, Medellín", Telefono: "3109876543"},
			{Nombre: "DataCorp S.A", Direccion: "Av. El Dorado # 68B-31, Cali", Telefono: "3204567890"},
		}
		if err := db.Create(&companias).Error; err != nil {
			log.Fatalf("❌ Error insertando compañías: %v", err)
		}

		// 10 empleados distribuidos entre las 3 compañías
		empleados := []entities.Empleado{
			{Nombre: "Ana", Apellido: "Gómez", Correo: "ana.gomez@tech.com", Cargo: "Desarrolladora Backend", Salario: 3500000, CompaniaID: companias[0].ID},
			{Nombre: "Carlos", Apellido: "Rojas", Correo: "carlos.rojas@tech.com", Cargo: "QA Tester", Salario: 2800000, CompaniaID: companias[0].ID},
			{Nombre: "María", Apellido: "López", Correo: "maria.lopez@tech.com", Cargo: "DevOps Engineer", Salario: 4000000, CompaniaID: companias[0].ID},
			{Nombre: "Pedro", Apellido: "Martínez", Correo: "pedro.martinez@innovatech.com", Cargo: "Arquitecto de Software", Salario: 5000000, CompaniaID: companias[1].ID},
			{Nombre: "Laura", Apellido: "Vargas", Correo: "laura.vargas@innovatech.com", Cargo: "Analista de Sistemas", Salario: 3200000, CompaniaID: companias[1].ID},
			{Nombre: "Diego", Apellido: "Hernández", Correo: "diego.hernandez@innovatech.com", Cargo: "Desarrollador Frontend", Salario: 3800000, CompaniaID: companias[1].ID},
			{Nombre: "Sofía", Apellido: "Jiménez", Correo: "sofia.jimenez@innovatech.com", Cargo: "Scrum Master", Salario: 3600000, CompaniaID: companias[1].ID},
			{Nombre: "Andrés", Apellido: "Castillo", Correo: "andres.castillo@datacorp.com", Cargo: "DBA Senior", Salario: 4200000, CompaniaID: companias[2].ID},
			{Nombre: "Valentina", Apellido: "Moreno", Correo: "valentina.moreno@datacorp.com", Cargo: "Data Engineer", Salario: 4500000, CompaniaID: companias[2].ID},
			{Nombre: "Julián", Apellido: "Torres", Correo: "julian.torres@datacorp.com", Cargo: "ML Engineer", Salario: 5500000, CompaniaID: companias[2].ID},
		}
		if err := db.Create(&empleados).Error; err != nil {
			log.Fatalf("❌ Error insertando empleados: %v", err)
		}
		log.Println("✅ Seed de compañías y empleados completo")
	} else {
		log.Println("⚠️  Seed de compañías/empleados omitido: ya existen datos")
	}

	// ── 2. Seed de Usuarios (Administradores y de Negocio) ────────────────
	log.Println("🌱 Verificando usuarios semilla...")
	
	// Helper para insertar un usuario si no existe
	insertarUsuarioSiNoExiste := func(nombre, correo, contrasena, rol string, companiaID *uint) {
		var u entities.Usuario
		err := db.Where("correo = ?", correo).First(&u).Error
		if err != nil && (err == gorm.ErrRecordNotFound || err.Error() == "record not found") {
			log.Printf("👤 Creando usuario semilla: %s (%s)...", nombre, correo)
			hash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("❌ Error generando hash para %s: %v", correo, err)
			}
			
			nuevoUsuario := entities.Usuario{
				Nombre:         nombre,
				Correo:         correo,
				ContrasenaHash: string(hash),
				Rol:            rol,
				CompaniaID:     companiaID,
			}
			
			if err := db.Create(&nuevoUsuario).Error; err != nil {
				log.Fatalf("❌ Error al insertar usuario %s: %v", correo, err)
			}
			log.Printf("✅ Usuario %s creado exitosamente", correo)
		} else if err != nil {
			log.Printf("❌ Error consultando existencia de %s: %v", correo, err)
		} else {
			log.Printf("ℹ️  Usuario %s ya existe en la base de datos", correo)
		}
	}

	// 1. Administrador Global
	insertarUsuarioSiNoExiste("Administrador Global", "admin@companies.com", "Admin123*", "ADMIN", nil)

	// 2. Administrador Medellín (Innovatech Ltda - ID: 2)
	idMedellin := uint(2)
	insertarUsuarioSiNoExiste("Administrador Medellín", "admin@innovatech.com", "Admin123*", "ADMIN", &idMedellin)

	// 3. Usuario Bogotá (Tech Solutions S.A.S - ID: 1)
	idBogota := uint(1)
	insertarUsuarioSiNoExiste("Usuario Bogotá", "usuario@tech.com", "User123*", "USUARIO", &idBogota)
}
