package main

import (
	"context"
	"time"

	"cliente.local/unificador/controladores"
	interfaces "cliente.local/unificador/fachada"
	"cliente.local/unificador/fachada/puentes"
	repositorios "cliente.local/unificador/repositorio"
	"cliente.local/unificador/servicios"
	"cliente.local/unificador/vistas"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	apiURL := "http://localhost:3333/users"

	// Inyección de dependencias
	repo := repositorios.NuevoRepoLoginJSON(apiURL)
	service := servicios.NuevoServicioLogin(repo)
	controller := controladores.NuevoLoginController(service)

	// Mostrar login y obtener userID
	userID, success := vistas.MostrarLogin(controller)
	if !success {
		return
	}

	// Configurar módulos
	modulos := []interfaces.IModulo{
		&puentes.ModuloGRPC{},
		&puentes.ModuloRMI{},
	}

	vistas.MostrarMenuPrincipal(ctx, modulos, userID)
}
