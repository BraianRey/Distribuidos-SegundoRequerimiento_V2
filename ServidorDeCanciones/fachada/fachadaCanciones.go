package fachada

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"servidor.local/grpc-servidorcanciones/capaFachadaServices"
	dtos "servidor.local/grpc-servidorcanciones/capaFachadaServices/DTOs"
	componenteconexioncola "servidor.local/grpc-servidorcanciones/componenteConexionCola"
	"servidor.local/grpc-servidorcanciones/modelos"
)

// FachadaCanciones proporciona una interfaz simplificada para acceder a la lógica de negocio relacionada con las canciones y géneros
type FachadaCanciones struct {
	Generos      []modelos.Genero  // lista de géneros disponibles
	Canciones    []modelos.Cancion // lista de canciones disponibles
	filePath     string
	mu           sync.Mutex
	conexionCola *componenteconexioncola.RabbitPublisher
}

// Intenta cargar canciones desde canciones.json en el directorio del servicio.
// Si el archivo no existe, usa datos embebidos y crea el archivo.
func NuevaFachadaCanciones() *FachadaCanciones {
	f := &FachadaCanciones{filePath: "canciones.json"}

	// intentar conectar a RabbitMQ para envío de notificaciones (opcional)
	if pub, err := componenteconexioncola.NewRabbitPublisher(); err == nil {
		f.conexionCola = pub
		log.Printf("Conectado a RabbitMQ para notificaciones")
	} else {
		log.Printf("No se pudo conectar a RabbitMQ (opcional): %v", err)
	}

	// intentar leer géneros desde generos.json
	if gens, err := capaFachadaServices.LeerGeneros(); err == nil && len(gens) > 0 {
		f.Generos = gens
	}

	// intentar leer canciones
	if data, err := os.ReadFile(f.filePath); err == nil && len(data) > 0 {
		var canciones []modelos.Cancion
		if err := json.Unmarshal(data, &canciones); err == nil {
			f.Canciones = canciones
			// si no cargamos generos antes, construir mapa de generos únicos y persistir
			if len(f.Generos) == 0 {
				genMap := map[int]modelos.Genero{}
				for _, c := range canciones {
					genMap[c.Genero.ID] = c.Genero
				}
				var gens []modelos.Genero
				for _, g := range genMap {
					gens = append(gens, g)
				}
				sort.Slice(gens, func(i, j int) bool { return gens[i].ID < gens[j].ID })
				f.Generos = gens
				// persistir generos
				_ = capaFachadaServices.EscribirGeneros(f.Generos)
			}
			log.Printf("Cargadas %d canciones desde %s", len(canciones), f.filePath)
			return f
		}
		log.Printf("Error parseando %s: %v", f.filePath, err)
	}

	// fallback: datos embebidos
	generos := []modelos.Genero{{ID: 1, Nombre: "Rock"}, {ID: 2, Nombre: "Pop"}, {ID: 3, Nombre: "Clásica"}}
	canciones := []modelos.Cancion{
		{ID: 1, Titulo: "Lamento Boliviano", Artista: "Enanitos Verdes", Album: "Contrarreloj", Anio: 1994, Duracion: "4:20", Genero: generos[0], Idioma: "Español"},
		{ID: 2, Titulo: "De Música Ligera", Artista: "Soda Stereo", Album: "Canción Animal", Anio: 1990, Duracion: "3:30", Genero: generos[0], Idioma: "Español"},
		{ID: 3, Titulo: "La Flaca", Artista: "Jarabe de Palo", Album: "La Flaca", Anio: 1996, Duracion: "4:00", Genero: generos[0], Idioma: "Español"},
		{ID: 4, Titulo: "Hey Ya!", Artista: "OutKast", Album: "Speakerboxxx/The Love Below", Anio: 2003, Duracion: "3:55", Genero: generos[1], Idioma: "Inglés"},
		{ID: 5, Titulo: "Umbrella", Artista: "Rihanna", Album: "Good Girl Gone Bad", Anio: 2007, Duracion: "4:36", Genero: generos[1], Idioma: "Inglés"},
		{ID: 6, Titulo: "Sunflower", Artista: "Post Malone & Swae Lee", Album: "Spiderman Into the Spider-Verse", Anio: 2018, Duracion: "2:38", Genero: generos[1], Idioma: "Inglés"},
		{ID: 7, Titulo: "Moonlight Sonata", Artista: "Ludwig van Beethoven", Album: "Piano Sonata No. 14", Anio: 1791, Duracion: "14:59", Genero: generos[2], Idioma: "Instrumental"},
		{ID: 8, Titulo: "Fur Elise", Artista: "Ludwig van Beethoven", Album: "Bagatelle No. 25", Anio: 1792, Duracion: "5:06", Genero: generos[2], Idioma: "Instrumental"},
		{ID: 9, Titulo: "Hungarian Rhapsody No2", Artista: "Franz Liszt", Album: "Hungarian Rhapsodies", Anio: 1847, Duracion: "10:31", Genero: generos[2], Idioma: "Instrumental"},
	}
	f.Generos = generos
	f.Canciones = canciones

	// persistir datos iniciales
	if out, err := json.MarshalIndent(f.Canciones, "", "  "); err == nil {
		_ = os.WriteFile(f.filePath, out, 0644)
	}
	// persistir generos por separado
	_ = capaFachadaServices.EscribirGeneros(f.Generos)
	return f
}

// ✅ NUEVO: Obtener todas las canciones
func (f *FachadaCanciones) ObtenerTodasLasCanciones() []modelos.Cancion {
	f.mu.Lock()
	defer f.mu.Unlock()
	log.Printf("Solicitadas todas las canciones, se tienen %d canciones", len(f.Canciones))
	return f.Canciones
}

// ListarGeneros devuelve la lista de géneros disponibles
func (f *FachadaCanciones) ListarGeneros() []modelos.Genero {
	log.Printf("Se a solicitado la lista de géneros, se tienen %d géneros", len(f.Generos))
	return f.Generos
}

// CrearGenero crea un nuevo genero si no existe (comparación por nombre normalizada)
func (f *FachadaCanciones) CrearGenero(nombre string) (modelos.Genero, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if nombre == "" {
		return modelos.Genero{}, nil
	}
	// normalizar nombre (trim + lower)
	norm := strings.TrimSpace(strings.ToLower(nombre))
	for _, g := range f.Generos {
		if strings.ToLower(strings.TrimSpace(g.Nombre)) == norm {
			return g, nil
		}
	}
	maxG := 0
	for _, g := range f.Generos {
		if g.ID > maxG {
			maxG = g.ID
		}
	}
	ng := modelos.Genero{ID: maxG + 1, Nombre: nombre}
	f.Generos = append(f.Generos, ng)
	sort.Slice(f.Generos, func(i, j int) bool { return f.Generos[i].ID < f.Generos[j].ID })
	// persistir generos
	if err := capaFachadaServices.EscribirGeneros(f.Generos); err != nil {
		return ng, err
	}
	return ng, nil
}

// ListarCancionesPorGenero devuelve las canciones de un género específico
func (f *FachadaCanciones) ListarCancionesPorGenero(id int32) []modelos.Cancion {
	var canciones []modelos.Cancion
	for _, c := range f.Canciones {
		if int32(c.Genero.ID) == id {
			canciones = append(canciones, c)
		}
	}
	log.Printf("Se a solicitado la lista de canciones del género %d, se tienen %d canciones", id, len(canciones))
	return canciones
}

// ObtenerDetallesCancion devuelve los detalles de una canción específica
func (f *FachadaCanciones) ObtenerDetallesCancion(id int32) *modelos.Cancion {
	for _, c := range f.Canciones {
		if int32(c.ID) == id {
			log.Printf("Se a solicitado los detalles de la canción '%s'", c.Titulo)
			return &c
		}
	}
	log.Printf("No se encontraron detalles para la canción con ID %d", id)
	return nil
}

// Añade una canción nueva (genera ID automáticamente si es 0) y persiste en canciones.json
func (f *FachadaCanciones) RegistrarCancion(c modelos.Cancion) (modelos.Cancion, error) {
	if c.Genero.ID == 0 && c.Genero.Nombre != "" {
		ng, err := f.CrearGenero(c.Genero.Nombre)
		if err != nil {
			return c, err
		}
		c.Genero = ng
	} else if c.Genero.ID != 0 {
		// copiar generos actuales para buscarlos
		f.mu.Lock()
		gensCopy := make([]modelos.Genero, len(f.Generos))
		copy(gensCopy, f.Generos)
		f.mu.Unlock()

		found := false
		for _, g := range gensCopy {
			if g.ID == c.Genero.ID {
				c.Genero = g
				found = true
				break
			}
		}
		if !found {
			// intentar leer generos desde archivo
			if gens, err := capaFachadaServices.LeerGeneros(); err == nil && len(gens) > 0 {
				for _, g := range gens {
					if g.ID == c.Genero.ID {
						c.Genero = g
						found = true
						break
					}
				}
			}
		}
		// si aun no se encuentra, ignorar ID y crear por nombre si se proporciona
		if !found && c.Genero.Nombre != "" {
			ng, err := f.CrearGenero(c.Genero.Nombre)
			if err != nil {
				return c, err
			}
			c.Genero = ng
		}
	}

	// Modificar la lista de canciones
	f.mu.Lock()
	// verificar duplicado por título (case-insensitive)
	for _, ex := range f.Canciones {
		if strings.EqualFold(strings.TrimSpace(ex.Titulo), strings.TrimSpace(c.Titulo)) {
			f.mu.Unlock()
			return modelos.Cancion{}, fmt.Errorf("título ya registrado: %s", c.Titulo)
		}
	}
	defer f.mu.Unlock()

	// Asignar ID
	maxID := 0
	for _, ex := range f.Canciones {
		if ex.ID > maxID {
			maxID = ex.ID
		}
	}
	if c.ID == 0 {
		c.ID = maxID + 1
	} else {
		// Evitar IDs duplicados: si el ID ya existe, reasignar
		for _, ex := range f.Canciones {
			if ex.ID == c.ID {
				c.ID = maxID + 1
				break
			}
		}
	}

	f.Canciones = append(f.Canciones, c)

	// Escribir canciones a archivo
	if out, err := json.MarshalIndent(f.Canciones, "", "  "); err == nil {
		if werr := os.WriteFile(f.filePath, out, 0644); werr != nil {
			log.Printf("Error guardando canciones: %v", werr)
			return c, werr
		}
	} else {
		log.Printf("Error serializando canciones: %v", err)
		return c, err
	}

	log.Printf("Canción registrada: %s (ID %d)", c.Titulo, c.ID)

	// Enviar notificación a RabbitMQ (siempre que la conexión exista)
	if f.conexionCola != nil {
		// Fecha de registro (ahora)
		fecha := time.Now().UTC().Format(time.RFC3339)

		// Usar DTO para la representación del mensaje (sin ID, con nombre de genero)
		dto := dtos.CancionAlmacenarDTOInput{
			Titulo:        c.Titulo,
			Artista:       c.Artista,
			Album:         c.Album,
			Anio:          c.Anio,
			Duracion:      c.Duracion,
			Genero:        c.Genero.Nombre,
			FechaRegistro: fecha,
		}

		notif := componenteconexioncola.NotificacionCancion{
			Titulo:        dto.Titulo,
			Artista:       dto.Artista,
			Album:         dto.Album,
			Anio:          dto.Anio,
			Duracion:      dto.Duracion,
			Genero:        dto.Genero,
			FechaRegistro: dto.FechaRegistro,
			Mensaje:       fmt.Sprintf("Canción registrada con ID %d", c.ID),
		}

		if perr := f.conexionCola.PublicarNotificacion(notif); perr != nil {
			log.Printf("Error publicando notificación a RabbitMQ: %v", perr)
		} else {
			log.Printf("Notificación enviada a RabbitMQ para la canción ID %d", c.ID)
		}
	}

	return c, nil
}
