package puentes

import (
	"context"
)

type ModuloRMI struct{}

// Nombre retorna el nombre del módulo
func (m *ModuloRMI) Nombre() string {
	return "Cliente RMI"
}

// Iniciar inicia el módulo RMI
func (m *ModuloRMI) Iniciar(ctx context.Context) error {
	// IMPLEMENTAR LÓGICA DEL CLIENTE RMI AQUÍ
	return nil
}

// Cerrar cierra el módulo RMI
func (m *ModuloRMI) Cerrar() error {
	return nil
}
