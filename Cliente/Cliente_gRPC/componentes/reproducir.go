package componentes

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	util "cliente.local/grpc-cliente/utilidades"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// Maneja el flujo de reproducción: solicita el stream, lanza goroutines y controla la parada.
func ReproducirCancion(clientStream pbStream.AudioServiceClient, detalle *pbCancion.Cancion, ctx context.Context, reader *bufio.Reader) error {
	stream, err := clientStream.EnviarCancionMedianteStream(ctx, &pbStream.PeticionDTO{Titulo: detalle.Titulo})
	if err != nil {
		return err
	}
	readerPipe, writerPipe := io.Pipe()
	canalStop := make(chan struct{})
	canalSincronizacion := make(chan struct{})

	// Goroutine para reproducir
	go util.DecodificarReproducir(readerPipe, canalStop, canalSincronizacion)

	// Goroutine para recibir datos de la canción (responsabilidad separada en package stream)
	go RecibirCancion(stream, writerPipe, canalStop, canalSincronizacion)

	// Menú de reproducción
	for {
		fmt.Print("\nReproduciendo...\n0) Salir\nSeleccione una opción: ")
		subOpc, _ := reader.ReadString('\n')
		subOpc = strings.TrimSpace(subOpc)
		if subOpc == "0" {
			// detener reproducción
			close(canalStop)
			<-canalSincronizacion // esperar a que todo termine
			fmt.Println("Reproducción detenida.")
			break
		}
	}
	return nil
}
