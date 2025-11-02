package vistas

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cliente.local/unificador/controladores"
)

// Mestra la interfaz de login en consola
func MostrarLogin(controller *controladores.LoginController) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Usuario: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("Contraseña: ")
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)

	ok, err := controller.ValidarUsuario(user, pass)
	if err != nil {
		fmt.Println("Error al verificar usuario:", err)
		return false
	}

	if ok {
		fmt.Println("Inicio de sesión exitoso (servidor JSON simulado)")
		return true
	}

	fmt.Println("Usuario o contraseña incorrectos")
	return false
}
