package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"servidor.local/grpc-servidorcanciones/api"
	"servidor.local/grpc-servidorcanciones/controladores"
	"servidor.local/grpc-servidorcanciones/fachada"
	pb "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

const addr = ":9000"

func main() {
	log.Println("===============================================")
	log.Println("   SERVIDOR DE CANCIONES - INICIANDO")
	log.Println("===============================================")

	// Inicializar la fachada
	log.Println("📚 Inicializando catálogo de canciones...")
	fachadaCanciones := fachada.NuevaFachadaCanciones()
	log.Printf("✅ Catálogo cargado: %d canciones disponibles\n", len(fachadaCanciones.Canciones))

	// Iniciar servidor HTTP REST en goroutine
	log.Println("🌐 Iniciando servidor HTTP REST en puerto 8080...")
	go func() {
		api.StartHTTP(fachadaCanciones, ":8080")
	}()

	// Dar tiempo para que el servidor HTTP inicie
	time.Sleep(500 * time.Millisecond)
	log.Println("✅ Servidor HTTP REST listo")

	// Iniciar el servidor gRPC
	log.Printf("🔌 Iniciando servidor gRPC en puerto %s...\n", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌ No se pudo escuchar en %s: %v", addr, err)
	}

	// Crear servidor gRPC y registrar el servicio
	grpcServer := grpc.NewServer()
	pb.RegisterCancionesServiceServer(grpcServer, controladores.NuevoCancionesController(fachadaCanciones))

	log.Println("===============================================")
	log.Println("✅ SERVIDOR LISTO")
	log.Println("   - HTTP REST: http://localhost:8080")
	log.Println("   - gRPC:      localhost:9000")
	log.Println("===============================================\n")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("❌ Error al servir gRPC: %v", err)
	}
}
