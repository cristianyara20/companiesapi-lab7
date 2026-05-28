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

	// ── 2. Seed de Usuario Administrador (Módulo 5) ──────────────────────
	var userCount int64
	db.Model(&entities.Usuario{}).Count(&userCount)
	if userCount == 0 {
		log.Println("🌱 Insertando usuario administrador inicial...")
		hash, err := bcrypt.GenerateFromPassword([]byte("Admin123*"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ Error generando hash de contraseña para admin: %v", err)
		}
		admin := entities.Usuario{
			Nombre:         "Administrador Global",
			Correo:         "admin@companies.com",
			ContrasenaHash: string(hash),
			Rol:            "ADMIN",
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("❌ Error insertando admin inicial: %v", err)
		}
		log.Println("✅ Usuario administrador semilla insertado: admin@companies.com / Admin123*")
	} else {
		log.Println("⚠️  Seed de usuarios omitido: ya existen usuarios")
	}
}
