package database

import (
	"fmt"
	"log"
	"os"

	"companies-api/domain/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establece la conexión con Supabase PostgreSQL
// Lee las credenciales del archivo .env
func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require TimeZone=America/Bogota",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Error conectando a PostgreSQL: %v", err)
	}

	log.Println("✅ Conexión a Supabase PostgreSQL exitosa")
	return db
}

// Migrate crea o actualiza las tablas automáticamente
// Equivale a "dotnet ef database update" en EF Core
func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entities.Compania{},
		&entities.Empleado{},
	)
	if err != nil {
		log.Fatalf("❌ Error en migraciones: %v", err)
	}
	log.Println("✅ Migraciones aplicadas: tablas listas")
}