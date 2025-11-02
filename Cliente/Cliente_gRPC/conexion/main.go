package conexion

import (
	"context"
	"log"
	"time"

	menu "cliente.local/grpc-cliente/vistas"
	"google.golang.org/grpc"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

func RunClienteGRPC() {
	// Conexión al servidor de streaming
	connStream, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	// Asegurar el cierre de la conexión al finalizar main
	defer connStream.Close()
	// Crear cliente gRPC para el servicio de streaming
	clientStream := pbStream.NewAudioServiceClient(connStream)

	// Conexión al servidor de canciones
	connCancion, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	// Asegurar el cierre de la conexión al finalizar main
	defer connCancion.Close()
	// Crear cliente gRPC para el servicio de canciones
	clientCancion := pbCancion.NewCancionesServiceClient(connCancion)

	// Context con timeout para las operaciones de la aplicación
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	// Asegurar la cancelación del context al finalizar main
	defer cancel()

	// Mostrar menú principal

	menu.MostrarMenuPrincipal(clientStream, clientCancion, ctx)
}
