package servicios

import (
	"cliente.local/unificador/modelos"
	repositorios "cliente.local/unificador/repositorio"
)

// Maneja la lógica de negocio (autenticación).
type ServicioLogin struct {
	repo *repositorios.RepoLoginJSON
}

// Crea una nueva instancia de ServicioLogin
func NuevoServicioLogin(repo *repositorios.RepoLoginJSON) *ServicioLogin {
	return &ServicioLogin{repo: repo}
}

// Verifica las credenciales del usuario y retorna el usuario si es válido
func (s *ServicioLogin) VerificarCredenciales(usuario, password string) (*modelos.Usuario, error) {
	usuarioEncontrado, err := s.repo.BuscarUsuario(usuario, password)
	if err != nil {
		return nil, err
	}
	return usuarioEncontrado, nil
}
