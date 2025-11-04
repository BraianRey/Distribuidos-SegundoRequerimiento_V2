package dto

import "time"

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
