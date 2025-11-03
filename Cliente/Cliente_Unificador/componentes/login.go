package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login simula autenticación contra un "servidor" JSON
func Login(username, password string) error {
	client := http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get("http://localhost:3333/users")
	if err != nil {
		return fmt.Errorf("error al conectar con servidor simulado: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Users []User `json:"users"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	for _, u := range data.Users {
		if strings.EqualFold(u.Username, username) && u.Password == password {
			return nil // autenticación exitosa
		}
	}
	return errors.New("credenciales inválidas")
}
