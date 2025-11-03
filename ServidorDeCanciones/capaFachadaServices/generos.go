package capaFachadaServices

import (
	"encoding/json"
	"os"

	"servidor.local/grpc-servidorcanciones/modelos"
)

var generosFile = "generos.json"

// Lee el archivo generos.json si existe
func LeerGeneros() ([]modelos.Genero, error) {
	if data, err := os.ReadFile(generosFile); err == nil && len(data) > 0 {
		var gens []modelos.Genero
		if err := json.Unmarshal(data, &gens); err != nil {
			return nil, err
		}
		return gens, nil
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

// Escribe la lista de géneros en generos.json
func EscribirGeneros(gens []modelos.Genero) error {
	out, err := json.MarshalIndent(gens, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(generosFile, out, 0644)
}
