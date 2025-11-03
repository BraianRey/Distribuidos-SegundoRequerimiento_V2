package api

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"servidor.local/grpc-servidorcanciones/capaFachadaServices"
	"servidor.local/grpc-servidorcanciones/fachada"
	"servidor.local/grpc-servidorcanciones/modelos"
)

// Extraer keys de maps en multipart forms
func keysFromMap(m map[string][]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func keysFromFileHeaderMap(m map[string][]*multipart.FileHeader) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Construye el http.HandlerFunc para la ruta /canciones
func cancionesHandler(f *fachada.FachadaCanciones) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetCanciones(w, f)
		case http.MethodPost:
			handlePostCancion(w, r, f)
		default:
			http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		}
	}
}

// Construye el http.HandlerFunc para la ruta /generos
func generosHandler(f *fachada.FachadaCanciones) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			gens := f.ListarGeneros()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gens)
		case http.MethodPost:
			var g modelos.Genero
			if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
				http.Error(w, "json inválido", http.StatusBadRequest)
				return
			}
			if g.Nombre == "" {
				http.Error(w, "nombre requerido", http.StatusBadRequest)
				return
			}
			ng, err := f.CrearGenero(g.Nombre)
			if err != nil {
				http.Error(w, "error creando genero", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(ng)
		default:
			http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		}
	}
}

// Verifica si el título ya existe
func isTitleDuplicate(f *fachada.FachadaCanciones, title string) bool {
	if strings.TrimSpace(title) == "" {
		return false
	}
	for _, ex := range f.ObtenerTodasLasCanciones() {
		if strings.EqualFold(strings.TrimSpace(ex.Titulo), strings.TrimSpace(title)) {
			return true
		}
	}
	return false
}

// Responde con la lista de canciones en el formato requerido
func handleGetCanciones(w http.ResponseWriter, f *fachada.FachadaCanciones) {
	canciones := f.ObtenerTodasLasCanciones()
	type CancionDTO struct {
		Id       int    `json:"id"`
		Titulo   string `json:"titulo"`
		Artista  string `json:"artista"`
		Genero   string `json:"genero"`
		Idioma   string `json:"idioma"`
		Duracion string `json:"duracion"`
	}
	var cancionesDTO []CancionDTO
	for _, c := range canciones {
		cancionesDTO = append(cancionesDTO, CancionDTO{
			Id:       c.ID,
			Titulo:   c.Titulo,
			Artista:  c.Artista,
			Genero:   c.Genero.Nombre,
			Idioma:   c.Idioma,
			Duracion: c.Duracion,
		})
	}
	log.Printf("GET /canciones - devolviendo %d canciones", len(cancionesDTO))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cancionesDTO)
}

// Maneja la creación de una canción a partir de multipart/form-data o JSON
func handlePostCancion(w http.ResponseWriter, r *http.Request, f *fachada.FachadaCanciones) {
	defer r.Body.Close()
	log.Printf("/canciones - nueva petición: metodo=%s url=%s content-type=%s content-length=%d", r.Method, r.URL.String(), r.Header.Get("Content-Type"), r.ContentLength)

	contentType := r.Header.Get("Content-Type")
	var song modelos.Cancion

	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Printf("/canciones - ParseMultipartForm error: %v", err)
			http.Error(w, "error parseando multipart: "+err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("/canciones - multipart parsed: form.Value keys=%v form.File keys=%v", keysFromMap(r.MultipartForm.Value), keysFromFileHeaderMap(r.MultipartForm.File))
		form := r.MultipartForm
		if vals := form.Value["Titulo"]; len(vals) > 0 {
			song.Titulo = vals[0]
		}
		if vals := form.Value["Artista"]; len(vals) > 0 {
			song.Artista = vals[0]
		}
		if vals := form.Value["Album"]; len(vals) > 0 {
			song.Album = vals[0]
		}
		if vals := form.Value["Anio"]; len(vals) > 0 {
			if ai, err := strconv.Atoi(vals[0]); err == nil {
				song.Anio = ai
			}
		}

		// Verificar duplicado por título antes de guardar
		if isTitleDuplicate(f, song.Titulo) {
			log.Printf("/canciones - título ya registrado: %s", song.Titulo)
			http.Error(w, "título ya registrado", http.StatusConflict)
			return
		}

		if vals := form.Value["Duracion"]; len(vals) > 0 {
			song.Duracion = vals[0]
		}
		if vals := form.Value["GeneroID"]; len(vals) > 0 {
			if gi, err := strconv.Atoi(vals[0]); err == nil {
				song.Genero.ID = gi
			}
		} else if vals := form.Value["GeneroNombre"]; len(vals) > 0 {
			song.Genero.Nombre = vals[0]
		}
		if vals := form.Value["Idioma"]; len(vals) > 0 {
			song.Idioma = vals[0]
		}

		files := form.File["file"]
		if len(files) > 0 {
			fh := files[0]
			log.Printf("/canciones - archivo recibido: filename=%s size=%d header=%v", fh.Filename, fh.Size, fh.Header)
			fhr, err := fh.Open()
			if err != nil {
				log.Printf("/canciones - error abriendo archivo: %v", err)
				http.Error(w, "error abriendo archivo: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer fhr.Close()

			dest, err := capaFachadaServices.GuardarArchivoDesdeReader(fhr, song.Titulo)
			if err != nil {
				log.Printf("/canciones - error guardando archivo: %v", err)
				http.Error(w, "error guardando archivo: "+err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("/canciones - archivo guardado en: %s", dest)
			if song.Duracion == "" {
				if d, derr := capaFachadaServices.ObtenerDuracionFFprobe(dest); derr == nil && d != "" {
					song.Duracion = d
					log.Printf("/canciones - duracion extraida: %s", song.Duracion)
				} else if derr != nil {
					log.Printf("/canciones - fallo extrayendo duracion con ffprobe: %v", derr)
				}
			}
		}
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error leyendo body", http.StatusBadRequest)
			return
		}
		var req struct {
			modelos.Cancion
			ArchivoOrigen string `json:"archivo_origen"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "json inválido", http.StatusBadRequest)
			return
		}
		song = req.Cancion
		if req.ArchivoOrigen != "" {
			if isTitleDuplicate(f, req.Titulo) {
				log.Printf("/canciones - título ya registrado (json request): %s", req.Titulo)
				http.Error(w, "título ya registrado", http.StatusConflict)
				return
			}
			dest, err := capaFachadaServices.GuardarArchivoCancion(req.ArchivoOrigen, req.Titulo)
			if err != nil {
				http.Error(w, "error copiando archivo de audio: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if song.Duracion == "" {
				if d, derr := capaFachadaServices.ObtenerDuracionFFprobe(dest); derr == nil && d != "" {
					song.Duracion = d
				}
			}
		}
	}

	log.Printf("/canciones - registrando canción: %+v", song)
	nc, err := f.RegistrarCancion(song)
	if err != nil {
		log.Printf("/canciones - error guardando canción: %v", err)
		http.Error(w, "error guardando canción", http.StatusInternalServerError)
		return
	}
	log.Printf("/canciones - canción registrada OK: %+v", nc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(nc)
}
