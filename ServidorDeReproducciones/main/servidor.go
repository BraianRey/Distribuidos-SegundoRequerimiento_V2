package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"servidorDeReproducciones/controladores"
	"servidorDeReproducciones/repositorio"
)

func main() {
	fmt.Println("===============================================")
	fmt.Println("  SERVIDOR DE REPRODUCCIONES - Puerto 3000")
	fmt.Println("===============================================")

	// Cargar reproducciones existentes (repository)
	if err := repositorio.Load(); err != nil {
		fmt.Println("⚠️  Error cargando reproducciones:", err)
	}

	// Configurar rutas
	http.HandleFunc("/reproducciones", controladores.ManejadorReproducciones)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	})

	// Iniciar servidor
	fmt.Println("✅ Servidor iniciado en http://localhost:3000")
	fmt.Println("📊 Endpoints disponibles:")
	fmt.Println("   - POST /reproducciones       (Almacenar reproducción)")
	fmt.Println("   - GET  /reproducciones       (Consultar todas)")
	fmt.Println("   - GET  /reproducciones?idUsuario=X (Consultar por usuario)")
	fmt.Println("===============================================\n")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
