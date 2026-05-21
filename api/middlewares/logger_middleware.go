package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware registra cada petición HTTP que llega a la API
// Equivale al middleware de logging en ASP.NET Core
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		// Log al RECIBIR la petición
		logger.Info("📥 Request recibido",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", ctx.ClientIP()),
		)

		ctx.Next() // procesar el request (pasar al siguiente handler)

		// Log al RESPONDER
		logger.Info("📤 Response enviado",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("tiempo", time.Since(start)),
		)
	}
}

// ErrorHandlerMiddleware captura errores inesperados (panics)
// Equivale al UseExceptionHandler en ASP.NET Core
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("❌ Error inesperado",
					zap.Any("error", r),
					zap.String("path", ctx.Request.URL.Path),
				)
				ctx.JSON(500, gin.H{"error": "Error interno del servidor"})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
