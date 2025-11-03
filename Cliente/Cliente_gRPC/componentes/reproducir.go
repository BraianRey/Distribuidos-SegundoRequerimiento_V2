package componentes

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	util "cliente.local/grpc-cliente/utilidades"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// Maneja el flujo de reproducción: solicita el stream, lanza goroutines y controla la parada.
func ReproducirCancion(clientStream pbStream.AudioServiceClient, detalle *pbCancion.Cancion, ctx context.Context, reader *bufio.Reader, providedUserID string) error {
	// Determinar ID de usuario: requerimos que el usuario sea proporcionado por el Unificador
	if providedUserID == "" {
		return fmt.Errorf("debe proporcionar userID (ejecutado en modo standalone no permitido sin userID)")
	}
	var userIDInt int
	// Intentar convertir directamente
	if uid, err := strconv.Atoi(providedUserID); err == nil {
		userIDInt = uid
	} else {
		re := regexp.MustCompile("\\d+")
		found := re.FindString(providedUserID)
		if found == "" {
			return fmt.Errorf("error al convertir ID de usuario proporcionado: %v", err)
		}
		uid2, err2 := strconv.Atoi(found)
		if err2 != nil {
			return fmt.Errorf("error al convertir ID de usuario extraído: %v", err2)
		}
		userIDInt = uid2
	}
	fmt.Printf("   detalle.Idioma: '%s'\n", detalle.Idioma) // ← CRÍTICO

	// Crear petición con TODOS los campos necesarios
	stream, err := clientStream.EnviarCancionMedianteStream(ctx, &pbStream.PeticionDTO{
		Id:        detalle.Id,
		Titulo:    detalle.Titulo,
		Artista:   detalle.Artista,
		Album:     detalle.Album,
		Anio:      detalle.Anio,
		Duracion:  detalle.Duracion,
		Genero:    detalle.Genero.Nombre,
		Idioma:    detalle.Idioma,
		IdUsuario: int32(userIDInt),
	})
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
