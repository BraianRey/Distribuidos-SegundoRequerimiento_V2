package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Reproduccion representa una reproducción de canción por un usuario
type Reproduccion struct {
	ID        int       `json:"id"`
	IDUsuario int       `json:"idUsuario"`
	IDCancion int       `json:"idCancion"`
	Titulo    string    `json:"titulo"`
	Artista   string    `json:"artista"`
	Genero    string    `json:"genero"`
	Idioma    string    `json:"idioma"`
	FechaHora time.Time `json:"fechaHora"`
}

// Repositorio de reproducciones
var reproducciones []Reproduccion
var archivoJSON = "Reproducciones.json"

func main() {
	fmt.Println("===============================================")
	fmt.Println("  SERVIDOR DE REPRODUCCIONES - Puerto 3000")
	fmt.Println("===============================================")

	// Cargar reproducciones existentes
	cargarReproducciones()

	// Configurar rutas
	http.HandleFunc("/reproducciones", manejadorReproducciones)
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

// Manejador principal de reproducciones
func manejadorReproducciones(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "POST":
		almacenarReproduccion(w, r)
	case "GET":
		consultarReproducciones(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// POST /reproducciones - Almacenar nueva reproducción
func almacenarReproduccion(w http.ResponseWriter, r *http.Request) {
	var nuevaReproduccion Reproduccion

	// Leer body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer datos", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parsear JSON
	err = json.Unmarshal(body, &nuevaReproduccion)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Generar ID
	nuevaReproduccion.ID = len(reproducciones) + 1
	nuevaReproduccion.FechaHora = time.Now()

	// Agregar a la lista
	reproducciones = append(reproducciones, nuevaReproduccion)

	// Guardar en archivo
	err = guardarReproducciones()
	if err != nil {
		http.Error(w, "Error al guardar", http.StatusInternalServerError)
		return
	}

	// Log
	fmt.Printf("📝 Nueva reproducción almacenada:\n")
	fmt.Printf("   Usuario: %d | Canción: %d - %s | Género: %s\n",
		nuevaReproduccion.IDUsuario,
		nuevaReproduccion.IDCancion,
		nuevaReproduccion.Titulo,
		nuevaReproduccion.Genero)

	// Responder
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaReproduccion)
}

// GET /reproducciones?idUsuario=X - Consultar reproducciones
func consultarReproducciones(w http.ResponseWriter, r *http.Request) {
	idUsuarioStr := r.URL.Query().Get("idUsuario")

	// Si no se especifica usuario, devolver todas
	if idUsuarioStr == "" {
		fmt.Printf("📊 Consultando todas las reproducciones (%d)\n", len(reproducciones))
		json.NewEncoder(w).Encode(reproducciones)
		return
	}

	// Parsear ID de usuario
	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	// Filtrar por usuario
	var reproduccionesUsuario []Reproduccion
	for _, rep := range reproducciones {
		if rep.IDUsuario == idUsuario {
			reproduccionesUsuario = append(reproduccionesUsuario, rep)
		}
	}

	fmt.Printf("📊 Consultando reproducciones del usuario %d: %d encontradas\n",
		idUsuario, len(reproduccionesUsuario))

	// Responder
	json.NewEncoder(w).Encode(reproduccionesUsuario)
}

// Cargar reproducciones desde archivo JSON
func cargarReproducciones() {
	data, err := os.ReadFile(archivoJSON)
	if err != nil {
		fmt.Println("⚠️  No se encontró archivo de reproducciones, creando nuevo...")
		reproducciones = []Reproduccion{}
		guardarReproducciones()
		return
	}

	err = json.Unmarshal(data, &reproducciones)
	if err != nil {
		fmt.Println("❌ Error al cargar reproducciones:", err)
		reproducciones = []Reproduccion{}
		return
	}

	fmt.Printf("✅ Cargadas %d reproducciones desde %s\n", len(reproducciones), archivoJSON)
}

// Guardar reproducciones en archivo JSON
func guardarReproducciones() error {
	data, err := json.MarshalIndent(reproducciones, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(archivoJSON, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
