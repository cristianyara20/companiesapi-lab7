package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorDetail representa el detalle de un campo que falló la validación
type ValidationErrorDetail struct {
	Campo   string `json:"campo"`
	Detalle string `json:"detalle"`
}

// ValidationError representa el error general de validación
type ValidationError struct {
	Mensaje string                  `json:"mensaje"`
	Errores []ValidationErrorDetail `json:"errores"`
}

func (e *ValidationError) Error() string {
	return e.Mensaje
}

var validate = validator.New()

func init() {
	// Registrar nombre del tag JSON para los errores en lugar del nombre del struct
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct valida cualquier struct y devuelve un ValidationError si hay fallos
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var details []ValidationErrorDetail
	for _, fe := range validationErrors {
		details = append(details, ValidationErrorDetail{
			Campo:   fe.Field(),
			Detalle: getErrorMessage(fe),
		})
	}

	return &ValidationError{
		Mensaje: "Error de validación",
		Errores: details,
	}
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "El campo es obligatorio"
	case "email":
		return "Formato de correo electrónico inválido"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Debe tener al menos %s caracteres", fe.Param())
		}
		return fmt.Sprintf("El valor mínimo es %s", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("No debe exceder %s caracteres", fe.Param())
		}
		return fmt.Sprintf("El valor máximo es %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Debe ser mayor que %s", fe.Param())
	case "numeric":
		return "Debe contener únicamente dígitos numéricos"
	case "oneof":
		return fmt.Sprintf("Debe ser uno de los siguientes valores: %s", fe.Param())
	default:
		return fmt.Sprintf("Fallo en la regla de validación: %s", fe.Tag())
	}
}
