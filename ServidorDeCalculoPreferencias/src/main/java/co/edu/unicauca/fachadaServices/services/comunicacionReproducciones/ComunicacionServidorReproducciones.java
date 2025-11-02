package co.edu.unicauca.fachadaServices.services.comunicacionReproducciones;

import co.edu.unicauca.fachadaServices.DTO.ReproduccionesDTOEntrada;
import feign.Feign;
import feign.jackson.JacksonDecoder;

import java.util.ArrayList;
import java.util.List;

public class ComunicacionServidorReproducciones {
    private static final String BASE_URL = "http://localhost:3000";
    private final ReproduccionesRemoteClient client;

    public ComunicacionServidorReproducciones() {
        this.client = Feign.builder()
            .decoder(new JacksonDecoder())
            .target(ReproduccionesRemoteClient.class, BASE_URL);
    }

    public List<ReproduccionesDTOEntrada> obtenerReproduccionesRemotas(Integer idUsuario){
        try{
            List<ReproduccionesDTOEntrada> reproducciones = client.obtenerReproducciones(idUsuario);
            return reproducciones != null ? reproducciones : new ArrayList<>();
        } catch (Exception e){
            return new ArrayList<>();
        }
    }
}


