package api

import (
	"log"
	"net/http"

	"servidor.local/grpc-servidorcanciones/fachada"
)

// StartHTTP arranca el servidor HTTP para administrar canciones y géneros.
func StartHTTP(f *fachada.FachadaCanciones, listenAddr string) {
	http.HandleFunc("/canciones", cancionesHandler(f))
	http.HandleFunc("/generos", generosHandler(f))

	log.Printf("Servidor HTTP para canciones escuchando en %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
