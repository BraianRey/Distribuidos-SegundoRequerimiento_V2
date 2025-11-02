package puentes

import (
	"context"

	clientgrpc "cliente.local/grpc-cliente/conexion"
)

type ModuloGRPC struct{}

// Nombre retorna el nombre del módulo
func (m *ModuloGRPC) Nombre() string { return "Cliente gRPC" }

// Iniciar inicia el módulo gRPC
func (m *ModuloGRPC) Iniciar(ctx context.Context) error {
	clientgrpc.RunClienteGRPC()
	return nil
}
