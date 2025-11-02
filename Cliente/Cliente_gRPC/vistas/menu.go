package vistas

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	comp "cliente.local/grpc-cliente/componentes"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
	pbStream "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// Despliega el menú principal y maneja la navegación entre submenús
func MostrarMenuPrincipal(clientStream pbStream.AudioServiceClient, clientCancion pbCancion.CancionesServiceClient, ctx context.Context) {
	// lector de entrada estándar
	reader := bufio.NewReader(os.Stdin)
	// menú principal delega en submenus
	for {
		fmt.Print("\n1) Ver géneros\n0) Salir\nSeleccione una opción: ")
		opcion, _ := reader.ReadString('\n')
		opcion = strings.TrimSpace(opcion)
		switch opcion {
		case "0":
			fmt.Println("Saliendo...")
			return
		case "1":
			mostrarMenuGeneros(clientStream, clientCancion, ctx, reader)
		default:
			// mantiene el loop
			fmt.Println("Opción inválida.")
		}
	}
}

// Maneja la selección de géneros y delega la selección de canción
func mostrarMenuGeneros(clientStream pbStream.AudioServiceClient, clientCancion pbCancion.CancionesServiceClient, ctx context.Context, reader *bufio.Reader) {
	for {
		idGenero, ok, err := comp.SeleccionarGenero(clientCancion, ctx, reader)
		if err != nil {
			fmt.Println("Error al obtener géneros:", err)
			// Volver al menu principal
			return
		}
		if !ok {
			// Usuario eligió volver o no hay géneros -> volver al menu principal
			return
		}

		// Seleccionar canción para el género elegido
		_, detalle, ok, err := comp.SeleccionarCancion(clientCancion, ctx, idGenero, reader)
		if err != nil {
			fmt.Println("Error al obtener canción o detalles:", err)
			// Mantenerse en el submenu de géneros
			continue
		}
		if !ok {
			// Volver a la lista de géneros
			continue
		}

		// Mostrar detalles y opciones para la canción seleccionada
		mostrarDetallesCancion(clientStream, detalle, ctx, reader)
		// Después de mostrar detalles, volvemos a la lista de géneros
	}
}

// Imprime detalles y ofrece opciones (reproducir/volver)
func mostrarDetallesCancion(clientStream pbStream.AudioServiceClient, detalle *pbCancion.Cancion, ctx context.Context, reader *bufio.Reader) {
	fmt.Printf("\nDetalles de la canción:\nTítulo: %s\nArtista: %s\nAlbum: %s\nAño: %d\nDuración: %s\nGénero: %s\n",
		detalle.Titulo, detalle.Artista, detalle.Album, detalle.Anio, detalle.Duracion, detalle.Genero.Nombre)

	for {
		fmt.Print("\n1) Reproducir\n0) Volver\nSeleccione una opción: ")
		opc, _ := reader.ReadString('\n')
		opc = strings.TrimSpace(opc)
		switch opc {
		case "0":
			// Volver a la lista de canciones/géneros
			return
		case "1":
			if err := comp.ReproducirCancion(clientStream, detalle, ctx, reader); err != nil {
				log.Println("Error durante la reproducción:", err)
			}
			// Después de reproducir volvemos al menú de detalles (o se podría volver a géneros)
			return
		default:
			// Mantiene el loop en este submenu
			fmt.Println("Opción inválida.")
		}
	}
}
