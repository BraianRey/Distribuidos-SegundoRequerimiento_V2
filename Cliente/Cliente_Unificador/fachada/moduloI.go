package interfaces

import "context"

// IModulo define el comportamiento común de todos los módulos del sistema.
type IModulo interface {
	Nombre() string
	Iniciar(ctx context.Context, userID string) error
}
