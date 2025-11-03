package puentes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ModuloRMI struct{}

// Nombre retorna el nombre del módulo
func (m *ModuloRMI) Nombre() string {
	return "Ver Preferencias Musicales (RMI)"
}

// Iniciar inicia el módulo RMI
func (m *ModuloRMI) Iniciar(ctx context.Context) error {
	// Leer ID del usuario del archivo de configuración
	config := make(map[string]interface{})
	configPath := "config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error al leer configuración:", err)
		return err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Error al parsear configuración:", err)
		return err
	}

	userID, ok := config["user_id"].(string)
	if !ok || userID == "" {
		fmt.Println("Debe iniciar sesión primero")
		return fmt.Errorf("usuario no autenticado")
	}

	// Ejecutar el cliente RMI con el ID del usuario
	cmd := exec.CommandContext(ctx, "java", "-jar",
		filepath.Join("Cliente_RMI", "target", "Cliente_RMI-1.0-SNAPSHOT.jar"),
		"--user-id", userID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("Error al ejecutar Cliente RMI:", err)
		return err
	}
	return nil
}

// Cerrar cierra el módulo RMI
func (m *ModuloRMI) Cerrar() error {
	return nil
}
