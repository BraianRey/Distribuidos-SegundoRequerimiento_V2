package repositorio

import (
	"encoding/json"
	"fmt"
	"os"
	"servidorDeReproducciones/dto"
	"sync"
)

var (
	reproducciones []dto.Reproduccion
	archivoJSON    = "Reproducciones.json"
	mu             sync.Mutex
)

// Load carga las reproducciones desde el archivo JSON (si existe)
func Load() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(archivoJSON)
	if err != nil {
		// archivo no existe -> inicializar vacío
		reproducciones = []dto.Reproduccion{}
		return saveInternal() // Usar función interna sin lock
	}

	err = json.Unmarshal(data, &reproducciones)
	if err != nil {
		reproducciones = []dto.Reproduccion{}
		return err
	}

	fmt.Printf("✅ Cargadas %d reproducciones desde %s\n", len(reproducciones), archivoJSON)
	return nil
}

// Save persiste las reproducciones en disco (con lock)
func Save() error {
	mu.Lock()
	defer mu.Unlock()
	return saveInternal()
}

// saveInternal persiste sin bloquear (asume que el mutex ya está bloqueado)
func saveInternal() error {
	data, err := json.MarshalIndent(reproducciones, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(archivoJSON, data, 0644)
}

// Add agrega una nueva reproduccion y la persiste
func Add(r dto.Reproduccion) (dto.Reproduccion, error) {
	mu.Lock()
	defer mu.Unlock()

	r.ID = len(reproducciones) + 1
	reproducciones = append(reproducciones, r)

	if err := saveInternal(); err != nil { // Usar función interna
		return dto.Reproduccion{}, err
	}
	return r, nil
}

// GetAll devuelve todas las reproducciones (copia)
func GetAll() []dto.Reproduccion {
	mu.Lock()
	defer mu.Unlock()
	copySlice := make([]dto.Reproduccion, len(reproducciones))
	copy(copySlice, reproducciones)
	return copySlice
}

// GetByUser devuelve reproducciones filtradas por idUsuario
func GetByUser(idUsuario int) []dto.Reproduccion {
	mu.Lock()
	defer mu.Unlock()
	var res []dto.Reproduccion
	for _, rep := range reproducciones {
		if rep.IDUsuario == idUsuario {
			res = append(res, rep)
		}
	}
	return res
}
