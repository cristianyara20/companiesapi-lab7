package routes

import (
	"companies-api/api/controllers"
	"companies-api/api/middlewares"
	"companies-api/application/services"
	unitofwork "companies-api/infrastructure/unit_of_work"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Setup configura el router con todas las rutas y la inyección de dependencias
func Setup(db *gorm.DB, logger *zap.Logger) *gin.Engine {
	router := gin.New()

	// ── Middlewares globales ──────────────────────────────────────────────
	router.Use(gin.Recovery())                             // manejo de panics
	router.Use(middlewares.LoggerMiddleware(logger))       // logging de requests
	router.Use(middlewares.ErrorHandlerMiddleware(logger)) // manejo de errores
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ── Inyección de dependencias (manual, de abajo hacia arriba) ─────────
	uow := unitofwork.NewUnitOfWork(db)
	companiaService := services.NewCompaniaService(uow, logger)
	empleadoService := services.NewEmpleadoService(uow, logger)
	authService := services.NewAuthService(uow, logger)

	companiaCtrl := controllers.NewCompaniaController(companiaService, logger)
	empleadoCtrl := controllers.NewEmpleadoController(empleadoService, logger)
	authCtrl := controllers.NewAuthController(authService, logger)

	// Middlewares reutilizables
	authGuard := middlewares.AuthMiddleware()
	adminRoleGuard := middlewares.RequireRole("ADMIN")
	anyRoleGuard := middlewares.RequireRole("ADMIN", "USUARIO")
	ownershipGuard := middlewares.EsPropietarioDeCompania(empleadoService)

	// ── Rutas de la API ───────────────────────────────────────────────────
	api := router.Group("/api")
	{
		// Rutas públicas de Autenticación
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/registro", authCtrl.Registro)
			authGroup.POST("/login", authCtrl.Login)
			authGroup.GET("/perfil", authGuard, authCtrl.Perfil) // requiere token
		}

		// Rutas protegidas de compañías
		comp := api.Group("/companias")
		comp.Use(authGuard)
		{
			comp.GET("", anyRoleGuard, companiaCtrl.GetAll)                            // Listar todas (Cualquier usuario autenticado)
			comp.POST("", anyRoleGuard, companiaCtrl.Create)                           // Crear
			comp.POST("/con-empleados", adminRoleGuard, companiaCtrl.CreateConEmpleados) // Transaccional (Solo ADMIN)
			comp.GET("/:id", anyRoleGuard, companiaCtrl.GetById)                       // Buscar por ID
			comp.PUT("/:id", anyRoleGuard, companiaCtrl.Update)                        // Actualizar
			comp.DELETE("/:id", adminRoleGuard, companiaCtrl.Delete)                   // Eliminar (Solo ADMIN)
			comp.GET("/:id/empleados", anyRoleGuard, companiaCtrl.GetEmpleados)        // Empleados de una compañía
		}

		// Rutas protegidas de empleados
		emp := api.Group("/empleados")
		emp.Use(authGuard)
		{
			emp.GET("", anyRoleGuard, empleadoCtrl.GetAll)                                // Listar todos (paged, filtered, sorted)
			emp.POST("", anyRoleGuard, empleadoCtrl.Create)                               // Crear
			emp.POST("/lote", anyRoleGuard, empleadoCtrl.CreateRange)                     // Creación masiva (Bulk)
			emp.DELETE("/lote", adminRoleGuard, empleadoCtrl.DeleteRange)                 // Eliminación múltiple (Solo ADMIN)
			emp.GET("/:id", anyRoleGuard, empleadoCtrl.GetById)                           // Buscar por ID
			emp.PUT("/:id", anyRoleGuard, ownershipGuard, empleadoCtrl.Update)            // Reemplazo completo (Propietario / ADMIN)
			emp.PATCH("/:id", anyRoleGuard, ownershipGuard, empleadoCtrl.Patch)           // Actualización parcial (Propietario / ADMIN)
			emp.DELETE("/:id", anyRoleGuard, ownershipGuard, empleadoCtrl.Delete)         // Eliminar individual (Propietario / ADMIN)
		}
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
      persistAuthorization: true,
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
