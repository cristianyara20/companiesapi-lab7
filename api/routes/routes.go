package routes

import (
	"companies-api/api/controllers"
	"companies-api/api/middlewares"
	"companies-api/application/services"
	unitofwork "companies-api/infrastructure/unit_of_work"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup configura el router con todas las rutas y la inyección de dependencias
// Este archivo es el equivalente a Program.cs en ASP.NET Core:
//   - Registra middlewares
//   - Registra servicios (DI manual)
//   - Define las rutas
func Setup(db *gorm.DB, logger *zap.Logger) *gin.Engine {
	router := gin.New()

	// ── Middlewares globales ──────────────────────────────────────────────
	router.Use(gin.Recovery())                             // manejo de panics
	router.Use(middlewares.LoggerMiddleware(logger))       // logging de requests
	router.Use(middlewares.ErrorHandlerMiddleware(logger)) // manejo de errores

	// ── Inyección de dependencias (manual, de abajo hacia arriba) ─────────
	// Flujo: Controller → Service → UnitOfWork → Repository → GORM → DB
	uow := unitofwork.NewUnitOfWork(db)
	companiaService := services.NewCompaniaService(uow, logger)
	empleadoService := services.NewEmpleadoService(uow, logger)
	companiaCtrl := controllers.NewCompaniaController(companiaService, logger)
	empleadoCtrl := controllers.NewEmpleadoController(empleadoService, logger)

	// ── Rutas de la API ───────────────────────────────────────────────────
	api := router.Group("/api")
	{
		// Rutas de compañías
		comp := api.Group("/companias")
		comp.GET("", companiaCtrl.GetAll)                            // Listar todas
		comp.GET("/:id", companiaCtrl.GetById)                       // Buscar por ID
		comp.POST("", companiaCtrl.Create)                           // Crear
		comp.PUT("/:id", companiaCtrl.Update)                        // Actualizar
		comp.DELETE("/:id", companiaCtrl.Delete)                     // Eliminar
		comp.GET("/:id/empleados", companiaCtrl.GetEmpleados)        // Empleados de una compañía
		comp.POST("/con-empleados", companiaCtrl.CreateConEmpleados) // Transaccional

		// Rutas de empleados
		emp := api.Group("/empleados")
		emp.GET("", empleadoCtrl.GetAll)        // Listar todos
		emp.GET("/:id", empleadoCtrl.GetById)   // Buscar por ID
		emp.POST("", empleadoCtrl.Create)       // Crear
		emp.PUT("/:id", empleadoCtrl.Update)    // Actualizar
		emp.DELETE("/:id", empleadoCtrl.Delete) // Eliminar
	}
	// ── Rutas Swagger UI (Documentación Interactiva) ──────────────────────
	router.StaticFile("/swagger.yaml", "./swagger.yaml")
	router.GET("/docs", func(c *gin.Context) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>API Docs - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/swagger.yaml',
      dom_id: '#swagger-ui',
    });
  };
</script>
</body>
</html>`
		c.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	// ── Ruta principal (Raíz) ─────────────────────────────────────────────
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API REST con Go lista",
			"docs":    "Visita /docs para ver la documentación de Swagger",
		})
	})

	return router
}
