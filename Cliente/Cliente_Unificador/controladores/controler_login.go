package controladores

import (
	"cliente.local/unificador/modelos"
	"cliente.local/unificador/servicios"
)

// Maneja la lógica de login utilizando el ServicioLogin
type LoginController struct {
	servicio *servicios.ServicioLogin
}

// Constructor del controlador de login
func NuevoLoginController(servicio *servicios.ServicioLogin) *LoginController {
	return &LoginController{servicio: servicio}
}

// Valida las credenciales del usuario y retorna el usuario si es válido
func (c *LoginController) ValidarUsuario(usuario, password string) (*modelos.Usuario, error) {
	return c.servicio.VerificarCredenciales(usuario, password)
}
