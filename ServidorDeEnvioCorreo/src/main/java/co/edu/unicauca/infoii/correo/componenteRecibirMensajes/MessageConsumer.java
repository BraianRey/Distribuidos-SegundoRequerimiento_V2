package co.edu.unicauca.infoii.correo.componenteRecibirMensajes;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Service;

import co.edu.unicauca.infoii.correo.DTOs.CancionAlmacenarDTOInput;

// Servicio que consume mensajes de la cola de RabbitMQ.
@Service
public class MessageConsumer {

    @RabbitListener(queues = "notificaciones_canciones")
    public void receiveMessage(CancionAlmacenarDTOInput objClienteCreado) {
        System.out.println("-----------------Datos de la canción recibidos-----------------");
        System.out.println("-----------------Enviando correo electrónico-----------------");



        System.out.println("-----------------Correo enviado al cliente con los siguientes datos:-----------------");
        System.out.println("Título:   " + objClienteCreado.getTitulo());
        System.out.println("Artista:  " + objClienteCreado.getArtista());
        System.out.println("Género:   " + objClienteCreado.getGenero());
        System.out.println("Álbum:    " + objClienteCreado.getAlbum());
        System.out.println("Año:      " + objClienteCreado.getAnio());
        System.out.println("Duración: " + objClienteCreado.getDuracion());
        System.out.println("Fecha de Almacenamiento: " + objClienteCreado.getFecha_registro());
        System.out.println("Frase motivadora: " + "Disfruta de las pequeñas cosas");
        System.out.println("----------------------------------------------------");
    }
}