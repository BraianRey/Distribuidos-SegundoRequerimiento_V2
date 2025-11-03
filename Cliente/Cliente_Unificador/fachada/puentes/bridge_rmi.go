package puentes

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ModuloRMI struct{}

// Nombre retorna el nombre del módulo
func (m *ModuloRMI) Nombre() string {
	return "Ver Preferencias Musicales (RMI)"
}

// Iniciar inicia el módulo RMI
func (m *ModuloRMI) Iniciar(ctx context.Context, userID string) error {
	if userID == "" {
		fmt.Println("Debe iniciar sesión primero")
		return fmt.Errorf("usuario no autenticado")
	}

	// Evitar usar JAR: preferimos ejecutar la clase ClienteMain desde target/classes.
	classesPath := filepath.Join("..", "Cliente_RMI", "target", "classes")
	clienteMainClass := filepath.Join(classesPath, "co", "edu", "unicauca", "main", "ClienteMain.class")
	var cmd *exec.Cmd

	// Si ClienteMain.class no existe, compilar el módulo Java (mvn compile) para generar clases
	if _, err := os.Stat(clienteMainClass); err != nil {
		fmt.Println("ClienteMain.class no encontrado, intentando compilar Cliente_RMI (mvn compile)...")
		mvnCmd := exec.CommandContext(ctx, "mvn", "-f", filepath.Join("..", "Cliente_RMI", "pom.xml"), "compile", "-DskipTests")
		mvnCmd.Stdout = os.Stdout
		mvnCmd.Stderr = os.Stderr
		mvnCmd.Stdin = os.Stdin
		if err := mvnCmd.Run(); err != nil {
			fmt.Println("Compilación con Maven falló o no está disponible:", err)
			return fmt.Errorf("no se pudo compilar Cliente_RMI: %w", err)
		}
	}

	// Buscar ClienteMain.class dentro de target para determinar classpath correcto
	var clienteMainFound string
	_ = filepath.WalkDir(filepath.Join("..", "Cliente_RMI", "target"), func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "ClienteMain.class") {
			clienteMainFound = p
			return filepath.SkipDir
		}
		return nil
	})

	if clienteMainFound != "" {
		// Usar siempre target/classes como classpath
		classesRoot := filepath.Join("..", "Cliente_RMI", "target", "classes")
		fmt.Println("Usando classpath:", classesRoot)
		cmd = exec.CommandContext(ctx, "java", "-cp", classesRoot, "co.edu.unicauca.main.ClienteMain", "--user-id", userID)
	} else {
		fmt.Println("ClienteMain no disponible en target; no se ejecutará el Cliente RMI")
		return fmt.Errorf("ClienteMain no encontrado en target")
	}
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
