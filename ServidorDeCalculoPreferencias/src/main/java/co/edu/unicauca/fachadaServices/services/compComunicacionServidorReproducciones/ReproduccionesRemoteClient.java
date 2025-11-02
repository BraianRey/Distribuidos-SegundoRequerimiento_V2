package co.edu.unicauca.fachadaServices.services.compComunicacionServidorReproducciones;

import co.edu.unicauca.fachadaServices.DTO.ReproduccionesDTOEntrada;
import feign.Headers;
import feign.Param;
import feign.RequestLine;

import java.util.List;

public interface ReproduccionesRemoteClient {
    @RequestLine("GET /Reproducciones?idUsuario={idUsuario}")
    @Headers("Accept: application/json")
    List<ReproduccionesDTOEntrada> obtenerReproducciones(@Param("idUsuario") Integer idUsuario);
}


