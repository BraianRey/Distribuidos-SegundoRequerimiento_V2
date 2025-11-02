package repositorios

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cliente.local/unificador/modelos"
)

// Maneja la autenticación contra un servidor JSON
type RepoLoginJSON struct {
	URL string
}

// Crea una nueva instancia de RepoLoginJSON
func NuevoRepoLoginJSON(url string) *RepoLoginJSON {
	return &RepoLoginJSON{URL: url}
}

// Busca un usuario por username y password
func (r *RepoLoginJSON) BuscarUsuario(username, password string) (*modelos.Usuario, error) {
	resp, err := http.Get(r.URL)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el servidor JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta inesperada: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	var usuarios []modelos.Usuario
	if err := json.Unmarshal(body, &usuarios); err != nil {
		fmt.Println("🛰️ Respuesta cruda del servidor JSON:")
		fmt.Println(string(body))
		return nil, fmt.Errorf("no pude parsear la respuesta del servidor JSON: %v", err)
	}

	for _, u := range usuarios {
		if u.Username == username && u.Password == password {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("usuario o contraseña incorrectos")
}
