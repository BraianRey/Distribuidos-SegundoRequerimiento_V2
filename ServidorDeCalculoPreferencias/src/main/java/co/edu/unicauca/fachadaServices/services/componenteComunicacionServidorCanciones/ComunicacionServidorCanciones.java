package co.edu.unicauca.fachadaServices.services.componenteComunicacionServidorCanciones;

import co.edu.unicauca.fachadaServices.DTO.CancionDTOEntrada;
import feign.Feign;
import feign.jackson.JacksonDecoder;

import java.util.ArrayList;
import java.util.List;

public class ComunicacionServidorCanciones {
    private static final String BASE_URL = "http://localhost:3000/";
    private final CancionesRemoteClient client;

    public ComunicacionServidorCanciones(){
        this.client = Feign.builder()
            .decoder(new JacksonDecoder())
            .target(CancionesRemoteClient.class, BASE_URL);
    }

    public List<CancionDTOEntrada> obtenerCancionesRemotas(){
        try{
            List<CancionDTOEntrada> canciones = client.obtenerCanciones();
            return canciones != null ? canciones : new ArrayList<>();
        } catch (Exception e) {
            return new ArrayList<>();
        }
    }
}

