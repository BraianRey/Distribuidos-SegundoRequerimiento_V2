package controladores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"ServidorDeReproducciones/dto"
	"ServidorDeReproducciones/repositorio"
	"ServidorDeReproducciones/utils"
)

// ManejadorReproducciones maneja GET/POST para /reproducciones
func ManejadorReproducciones(w http.ResponseWriter, r *http.Request) {
	utils.SetJSONHeaders(w)

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

func almacenarReproduccion(w http.ResponseWriter, r *http.Request) {
	var nueva dto.Reproduccion

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer datos", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &nueva); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	nueva.FechaHora = time.Now()

	agregada, err := repositorio.Add(nueva)
	if err != nil {
		http.Error(w, "Error al guardar", http.StatusInternalServerError)
		return
	}

	fmt.Printf("📝 Nueva reproducción almacenada:\n")
	fmt.Printf("   Usuario: %d | Canción: %d - %s | Género: %s\n",
		agregada.IDUsuario, agregada.IDCancion, agregada.Titulo, agregada.Genero)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agregada)
}

func consultarReproducciones(w http.ResponseWriter, r *http.Request) {
	idUsuarioStr := r.URL.Query().Get("idUsuario")

	if idUsuarioStr == "" {
		all := repositorio.GetAll()
		fmt.Printf("📊 Consultando todas las reproducciones (%d)\n", len(all))
		json.NewEncoder(w).Encode(all)
		return
	}

	idUsuario, err := strconv.Atoi(idUsuarioStr)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	res := repositorio.GetByUser(idUsuario)
	fmt.Printf("📊 Consultando reproducciones del usuario %d: %d encontradas\n", idUsuario, len(res))
	json.NewEncoder(w).Encode(res)
}
