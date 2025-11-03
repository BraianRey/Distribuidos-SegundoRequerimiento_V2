package puentes

import (
	"context"

	clientgrpc "cliente.local/grpc-cliente/conexion"
)

type ModuloGRPC struct{}

// Nombre retorna el nombre del módulo
func (m *ModuloGRPC) Nombre() string { return "Cliente gRPC" }

// Iniciar inicia el módulo gRPC
func (m *ModuloGRPC) Iniciar(ctx context.Context, userID string) error {
	// Ejecutar cliente gRPC pasando el userID de la sesión actual
	clientgrpc.RunClienteGRPC(userID)
	return nil
}
