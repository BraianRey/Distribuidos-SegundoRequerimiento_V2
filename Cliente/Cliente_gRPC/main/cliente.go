package main

import (
	"fmt"
	"log"
	"strings"

	"cliente.local/grpc-cliente/conexion"
)

func main() {
	log.Println("=== Cliente gRPC ===")
	// Ejecutar en modo standalone: pedir al usuario el userID (Para probar solo gRPC)
	var userID string
	fmt.Print("Ingrese userID (por ejemplo '1' o 'user1'): ")
	fmt.Scanln(&userID)
	if strings.TrimSpace(userID) == "" {
		log.Println("No se proporcionó userID. Saliendo.")
		return
	}
	conexion.RunClienteGRPC(strings.TrimSpace(userID))
}
