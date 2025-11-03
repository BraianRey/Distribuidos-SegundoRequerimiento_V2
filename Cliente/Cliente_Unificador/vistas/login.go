package vistas

import (
	"bufio"
	"encoding/json"
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

	usuario, err := controller.ValidarUsuario(user, pass)
	if err != nil {
		fmt.Println("Error al verificar usuario:", err)
		return false
	}

	if usuario != nil {
		// Guardar ID del usuario en config.json
		config := make(map[string]interface{})
		configPath := "config.json"
		data, err := os.ReadFile(configPath)
		if err == nil {
			json.Unmarshal(data, &config)
		}
		config["user_id"] = usuario.ID
		newData, _ := json.MarshalIndent(config, "", "    ")
		os.WriteFile(configPath, newData, 0644)

		fmt.Println("Inicio de sesión exitoso (servidor JSON simulado)")
		return true
	}

	fmt.Println("Usuario o contraseña incorrectos")
	return false
}
