package utilidades

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// package-level sync.Once para inicializar speaker solo 1 vez
var speakerInitOnce sync.Once
var speakerSampleRate beep.SampleRate // se usa para resampling si es necesario

// Decodifica desde el reader y reproduce hasta EOF o hasta que llegue stop.
// Usa un sync.Once local para cerrar canalSincronizacion solo 1 vez y evita panics por double close.
// Realiza resampling si el sampleRate de la canción no coincide con el usado en speaker.
func DecodificarReproducir(reader io.Reader, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	// decodificar el stream de audio
	streamer, format, err := mp3.Decode(io.NopCloser(reader))
	if err != nil {
		log.Printf("error decodificando MP3: %v", err)
		safeClose(canalSincronizacion) // cerrar canalSincronizacion para evitar deadlock en RecibirCancion
		return
	}

	// inicializar speaker una sola vez con sampleRate fijo (44100 Hz)
	speakerInitOnce.Do(func() {
		speakerSampleRate = beep.SampleRate(44100)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic inicializando speaker: %v", r)
			}
		}()
		if err := speaker.Init(speakerSampleRate, speakerSampleRate.N(time.Second/2)); err != nil {
			log.Printf("inicializando speaker: %v", err)
		}
	})

	var streamerToPlay beep.Streamer = streamer
	// Si el sampleRate del audio no coincide con el del speaker, hacer resampling
	if format.SampleRate != speakerSampleRate {
		streamerToPlay = beep.Resample(4, format.SampleRate, speakerSampleRate, streamer)
	}

	// Callback que se dispara al terminar la reproducción natural
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamerToPlay, beep.Callback(func() {
		close(done)
	})))

	// Asegurar que canalSincronizacion solo se cierre una vez
	var once sync.Once
	closeSync := func() { once.Do(func() { safeClose(canalSincronizacion) }) }

	// Goroutine que escucha stop o finalización natural
	go func() {
		select {
		// Si llega stop, limpiar y cerrar streamer
		case <-canalStop:
			log.Println("DecodificarReproducir: stop recibido, deteniendo reproducción.")
			speaker.Clear()      // Limpiar cola de reproducción
			_ = streamer.Close() // Cerrar streamer para liberar recursos
			closeSync()
		// Si termina naturalmente, cerrar canalSincronizacion
		case <-done:
			log.Println("DecodificarReproducir: reproducción terminó naturalmente.")
			closeSync() // Cerrar canalSincronizacion
		}
	}()
}

// safeClose evita panics por cerrar canales múltiples veces
func safeClose(ch chan struct{}) {
	if ch == nil {
		return
	}
	// Se usa recover para evitar panic si el canal ya está cerrado
	defer func() { _ = recover() }()
	close(ch)
}
