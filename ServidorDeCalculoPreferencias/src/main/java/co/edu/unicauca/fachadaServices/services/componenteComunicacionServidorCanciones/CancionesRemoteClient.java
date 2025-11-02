package co.edu.unicauca.fachadaServices.services.componenteComunicacionServidorCanciones;

import co.edu.unicauca.fachadaServices.DTO.CancionDTOEntrada;
import feign.Headers;
import feign.RequestLine;

import java.util.List;

public interface CancionesRemoteClient {
    @RequestLine("GET /canciones")
    @Headers("Accept: aplication/json")
    List<CancionDTOEntrada> obtenerCanciones();
}
