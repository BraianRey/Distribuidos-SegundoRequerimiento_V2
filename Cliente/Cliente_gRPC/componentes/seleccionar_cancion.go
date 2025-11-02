package componentes

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

// SeleccionarCancion lista las canciones de un género y permite seleccionar una.
// Retorna idCancion (int32), detalle de la canción, ok (true si seleccionó), error.
func SeleccionarCancion(clientCancion pbCancion.CancionesServiceClient, ctx context.Context, idGenero int32, reader *bufio.Reader) (int32, *pbCancion.Cancion, bool, error) {
	respCanciones, err := clientCancion.ListarCancionesPorGenero(ctx, &pbCancion.GeneroId{Id: idGenero})
	if err != nil {
		return 0, nil, false, err
	}
	if len(respCanciones.Canciones) == 0 {
		fmt.Println("No hay canciones en este género.")
		return 0, nil, false, nil
	}
	fmt.Println("\nCanciones disponibles:")
	for i, c := range respCanciones.Canciones {
		fmt.Printf("%d) %s - %s\n", i+1, c.Titulo, c.Artista)
	}
	fmt.Print("Seleccione la canción por número (o 0 para volver): ")
	cancStr, _ := reader.ReadString('\n')
	cancStr = strings.TrimSpace(cancStr)
	if cancStr == "0" {
		return 0, nil, false, nil
	}
	cancIdx, err := strconv.Atoi(cancStr)
	if err != nil || cancIdx < 1 || cancIdx > len(respCanciones.Canciones) {
		fmt.Println("Opción inválida.")
		return 0, nil, false, nil
	}
	idCancion := respCanciones.Canciones[cancIdx-1].Id
	respDetalle, err := clientCancion.ObtenerDetallesCancion(ctx, &pbCancion.CancionId{Id: idCancion})
	if err != nil || respDetalle == nil {
		return 0, nil, false, err
	}
	return idCancion, respDetalle, true, nil
}
