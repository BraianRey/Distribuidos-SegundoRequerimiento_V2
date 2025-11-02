package main

import (
	"log"

	"cliente.local/grpc-cliente/conexion"
)

func main() {
	log.Println("=== Cliente gRPC ===")
	conexion.RunClienteGRPC()
}
