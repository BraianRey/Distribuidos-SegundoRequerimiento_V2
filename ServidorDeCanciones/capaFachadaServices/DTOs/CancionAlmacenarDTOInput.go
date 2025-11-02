package dtos

// Representa los datos necesarios para almacenar una nueva canción
type CancionAlmacenarDTOInput struct {
	Titulo        string `json:"titulo"`
	Artista       string `json:"artista"`
	Album         string `json:"album"`
	Anio          int    `json:"anio"`
	Duracion      string `json:"duracion"`
	Genero        string `json:"genero"`
	FechaRegistro string `json:"fecha_registro"`
}
