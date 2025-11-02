package componentes

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/types/known/emptypb"
	pbCancion "servidor.local/grpc-servidorcanciones/serviciosCanciones"
)

// SeleccionarGenero lista géneros y permite al usuario seleccionar uno.
// Retorna idGenero (int32), ok (true si se seleccionó una opción válida), error.
func SeleccionarGenero(clientCancion pbCancion.CancionesServiceClient, ctx context.Context, reader *bufio.Reader) (int32, bool, error) {
	respGeneros, err := clientCancion.ListarGeneros(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, false, err
	}
	if len(respGeneros.Generos) == 0 {
		fmt.Println("No hay géneros disponibles.")
		return 0, false, nil
	}
	fmt.Println("\nGéneros disponibles:")
	for i, g := range respGeneros.Generos {
		fmt.Printf("%d) %s\n", i+1, g.Nombre)
	}
	fmt.Print("Seleccione el género por número (o 0 para volver): ")
	genStr, _ := reader.ReadString('\n')
	genStr = strings.TrimSpace(genStr)
	if genStr == "0" {
		return 0, false, nil
	}
	genIdx, err := strconv.Atoi(genStr)
	if err != nil || genIdx < 1 || genIdx > len(respGeneros.Generos) {
		fmt.Println("Opción inválida.")
		return 0, false, nil
	}
	return respGeneros.Generos[genIdx-1].Id, true, nil
}
