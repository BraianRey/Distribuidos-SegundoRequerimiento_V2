package fachada

import (
	"servidorDeReproducciones/dto"
	"servidorDeReproducciones/repositorio"
)

// Aplica reglas de negocio (si hubiera) y delega al repositorio
func RegistrarReproduccion(r dto.Reproduccion) (dto.Reproduccion, error) {
	// Aquí se podrían añadir validaciones u otros efectos (ej: publicar eventos)
	return repositorio.Add(r)
}

// Devuelve todas las reproducciones
func ObtenerTodas() []dto.Reproduccion {
	return repositorio.GetAll()
}

// Devuelve reproducciones filtradas por usuario
func ObtenerPorUsuario(idUsuario int) []dto.Reproduccion {
	return repositorio.GetByUser(idUsuario)
}
