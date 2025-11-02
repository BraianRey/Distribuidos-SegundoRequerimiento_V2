package modelos

// Usuario representa un usuario para autenticación
type Usuario struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
