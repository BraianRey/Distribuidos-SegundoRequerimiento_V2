package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"servidor.local/grpc-servidorcanciones/api"
	"servidor.local/grpc-servidorcanciones/controladores"
	"servidor.local/grpc-servidorcanciones/fachada"
	pb "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

const addr = ":9000"

func main() {
	// Inicializar la fachada y el servidor gRPC
	fachadaCanciones := fachada.NuevaFachadaCanciones()

	// Iniciar servidor HTTP para registrar canciones via Postman (API)
	go func() {
		api.StartHTTP(fachadaCanciones, ":8080")
	}()

	// Iniciar el servidor gRPC
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("No se pudo escuchar en %s: %v", addr, err)
	}
	// Crear servidor gRPC y registrar el servicio de canciones
	grpcServer := grpc.NewServer()
	// Registrar el servicio de canciones con su controlador
	pb.RegisterCancionesServiceServer(grpcServer, controladores.NuevoCancionesController(fachadaCanciones))

	log.Printf("ServidorDeCanciones escuchando en %s (gRPC)", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
