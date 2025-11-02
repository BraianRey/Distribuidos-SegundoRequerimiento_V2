package controladores

import "cliente.local/unificador/servicios"

// Maneja la lógica de login utilizando el ServicioLogin
type LoginController struct {
	servicio *servicios.ServicioLogin
}

// Constructor del controlador de login
func NuevoLoginController(servicio *servicios.ServicioLogin) *LoginController {
	return &LoginController{servicio: servicio}
}

// Valida las credenciales del usuario
func (c *LoginController) ValidarUsuario(usuario, password string) (bool, error) {
	return c.servicio.VerificarCredenciales(usuario, password)
}
