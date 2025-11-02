package componentes

import (
	"fmt"
	"io"
	"log"
	"sync"

	pb "servidor.local/grpc-servidorstream/serviciosStreaming"
)

// Recibe los fragmentos del stream y los escribe al pipe.
func RecibirCancion(stream pb.AudioService_EnviarCancionMedianteStreamClient, writer *io.PipeWriter, canalStop <-chan struct{}, canalSincronizacion chan struct{}) {
	noFragmento := 0
	var once sync.Once // para cerrar writer solo 1 vez
	// función para cerrar writer de forma segura
	closeWriter := func() {
		once.Do(func() { _ = writer.Close() })
	}

	for {
		select {
		case <-canalStop:
			fmt.Println("\nRecibirCancion: stop recibido, cerrando writer.")
			closeWriter()
			<-canalSincronizacion
			fmt.Println("✓ Reproducción finalizada.")
			return
		default:
			fragmento, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Canción recibida completa.")
				closeWriter()
				<-canalSincronizacion
				fmt.Println("✓ Reproducción finalizada.")
				return
			}
			if err != nil {
				log.Printf("Error recibiendo chunk (stream.Recv): %T - %v", err, err)
				closeWriter()
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada por error en recv.")
				return
			}
			noFragmento++
			fmt.Printf("\n Fragmento #%d recibido (%d bytes) reproduciendo ...", noFragmento, len(fragmento.Data))
			if _, err := writer.Write(fragmento.Data); err != nil {
				log.Printf("Error escribiendo en pipe: %v", err)
				closeWriter()
				<-canalSincronizacion
				fmt.Println("Reproducción finalizada.")
				return
			}
		}
	}
}
