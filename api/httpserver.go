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

// helpers para debug: extraer keys de maps en multipart forms
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

// StartHTTP arranca el servidor HTTP para administrar canciones y géneros.
func StartHTTP(f *fachada.FachadaCanciones, listenAddr string) {
	http.HandleFunc("/canciones", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			defer r.Body.Close()

			log.Printf("/canciones - nueva petición: metodo=%s url=%s content-type=%s content-length=%d", r.Method, r.URL.String(), r.Header.Get("Content-Type"), r.ContentLength)

			contentType := r.Header.Get("Content-Type")
			var song modelos.Cancion
			if strings.HasPrefix(contentType, "multipart/form-data") {
				log.Printf("/canciones - procesando multipart/form-data")
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
			return
		}
		http.Error(w, "ruta no encontrada", http.StatusNotFound)
	})

	http.HandleFunc("/generos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			gens := f.ListarGeneros()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(gens)
			return
		}
		if r.Method == http.MethodPost {
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
			return
		}
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	})

	log.Printf("Servidor HTTP para canciones escuchando en %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
