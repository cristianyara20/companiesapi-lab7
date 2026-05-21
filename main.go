package main

import (
	"companies-api/api/routes"
	"companies-api/infrastructure/database"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 1. Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Archivo .env no encontrado, usando variables del sistema")
	}

	// 2. Configurar Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("🚀 Iniciando API de Compañías y Empleados")
	sugar.Info("📦 Stack: Go · Gin · GORM · PostgreSQL (Supabase)")

	// 3. Conexión a Base de Datos
	db := database.Connect()

	// 4. Ejecutar migraciones (AutoMigrate crea/actualiza tablas)
	database.Migrate(db)

	// 5. Insertar datos iniciales si la BD está vacía
	database.Seed(db)

	// 6. Configurar router (rutas + middlewares + DI)
	router := routes.Setup(db, logger)

	// 7. Iniciar el servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("✅ Servidor listo",
		zap.String("url", "http://localhost:"+port),
		zap.String("endpoints", "http://localhost:"+port+"/api/companias"),
	)

	if err := router.Run(":" + port); err != nil {
		logger.Fatal("❌ Error al iniciar el servidor", zap.Error(err))
	}
}
