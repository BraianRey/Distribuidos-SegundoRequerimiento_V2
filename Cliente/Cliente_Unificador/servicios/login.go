package servicios

import repositorios "cliente.local/unificador/repositorio"

// Maneja la lógica de negocio (autenticación).
type ServicioLogin struct {
	repo *repositorios.RepoLoginJSON
}

// Crea una nueva instancia de ServicioLogin
func NuevoServicioLogin(repo *repositorios.RepoLoginJSON) *ServicioLogin {
	return &ServicioLogin{repo: repo}
}

// Verifica las credenciales del usuario
func (s *ServicioLogin) VerificarCredenciales(usuario, password string) (bool, error) {
	usuarioEncontrado, err := s.repo.BuscarUsuario(usuario, password)
	if err != nil {
		return false, err
	}
	return usuarioEncontrado != nil, nil
}
