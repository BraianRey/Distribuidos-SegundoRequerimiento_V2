package vistas

import (
	"context"
	"fmt"

	interfaces "cliente.local/unificador/fachada"
)

// Despliega el menú principal y maneja la selección de módulos
func MostrarMenuPrincipal(ctx context.Context, modulos []interfaces.IModulo, userID string) {
	for {
		fmt.Println("\n=== MENÚ PRINCIPAL ===")
		for i, m := range modulos {
			fmt.Printf("%d) %s\n", i+1, m.Nombre())
		}
		fmt.Println("0) Salir")
		fmt.Print("Seleccione una opción: ")

		var opcion int
		fmt.Scanln(&opcion)

		if opcion == 0 {
			fmt.Println("Saliendo...")
			return
		}
		if opcion < 1 || opcion > len(modulos) {
			fmt.Println("Opción inválida.")
			continue
		}

		modulos[opcion-1].Iniciar(ctx, userID)
	}
}
