package capacontroladores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	capafachadaservices "servidor.local/grpc-servidorstream/capaFachadaServices"
	pb "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// ControladorServidor implementa el servicio de streaming de audio
type ControladorServidor struct {
	pb.UnimplementedAudioServiceServer
}

// Reproduccion representa una reproducción para enviar al servidor
type Reproduccion struct {
	IDUsuario int    `json:"idUsuario"`
	IDCancion int    `json:"idCancion"`
	Titulo    string `json:"titulo"`
	Artista   string `json:"artista"`
	Genero    string `json:"genero"`
	Idioma    string `json:"idioma"`
}

// Implementación del procedimiento remoto
func (s *ControladorServidor) EnviarCancionMedianteStream(req *pb.PeticionDTO, stream pb.AudioService_EnviarCancionMedianteStreamServer) error {
	// Registrar reproducción antes de iniciar el streaming
	err := registrarReproduccion(req)
	if err != nil {
		log.Printf("⚠️ Error registrando reproducción: %v", err)
		// No fallamos el streaming si falla el registro
	}

	// Iniciar streaming de audio
	return capafachadaservices.StreamAudioFile(
		req.Titulo,
		// función para enviar fragmento al cliente
		func(data []byte) error {
			return stream.Send(&pb.FragmentoCancion{Data: data})
		})
}

// registrarReproduccion envía la información de la reproducción al ServidorDeReproducciones
func registrarReproduccion(req *pb.PeticionDTO) error {
	log.Printf("📝 Registrando reproducción: Usuario=%d, Canción=%d (%s)",
		req.IdUsuario, req.Id, req.Titulo)

	// Crear estructura de reproducción
	reproduccion := Reproduccion{
		IDUsuario: int(req.IdUsuario),
		IDCancion: int(req.Id),
		Titulo:    req.Titulo,
		Artista:   req.Artista,
		Genero:    req.Genero,
		Idioma:    req.Idioma,
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(reproduccion)
	if err != nil {
		return fmt.Errorf("error serializando reproducción: %w", err)
	}

	// Crear petición HTTP POST
	url := "http://localhost:3000/reproducciones"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando petición HTTP: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Enviar petición con timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error enviando petición: %w", err)
	}
	defer resp.Body.Close()

	// Verificar respuesta
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("servidor respondió con código %d", resp.StatusCode)
	}

	log.Printf("✅ Reproducción registrada exitosamente en ServidorDeReproducciones")
	return nil
}
