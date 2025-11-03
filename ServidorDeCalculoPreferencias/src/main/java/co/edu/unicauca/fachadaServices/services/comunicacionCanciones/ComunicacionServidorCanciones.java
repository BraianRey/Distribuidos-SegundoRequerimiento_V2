package co.edu.unicauca.fachadaServices.services.comunicacionCanciones;

import co.edu.unicauca.fachadaServices.DTO.CancionDTOEntrada;
import feign.Feign;
import feign.Logger;
import feign.jackson.JacksonDecoder;

import java.util.ArrayList;
import java.util.List;

public class ComunicacionServidorCanciones {
    private static final String BASE_URL = "http://localhost:8080";
    private final CancionesRemoteClient client;

    public ComunicacionServidorCanciones(){
        System.out.println("🔧 Inicializando cliente Feign para: " + BASE_URL);

        this.client = Feign.builder()
                .decoder(new JacksonDecoder())
                .logger(new Logger.JavaLogger().appendToFile("feign.log"))
                .logLevel(Logger.Level.FULL)
                .target(CancionesRemoteClient.class, BASE_URL);

        System.out.println("✅ Cliente Feign inicializado");
    }

    public List<CancionDTOEntrada> obtenerCancionesRemotas(){
        System.out.println("\n" + "=".repeat(50));
        System.out.println("🔗 CONSULTANDO SERVIDOR DE CANCIONES");
        System.out.println("   URL: " + BASE_URL + "/canciones");
        System.out.println("=".repeat(50));

        try{
            System.out.println("📤 Enviando petición HTTP GET...");

            List<CancionDTOEntrada> canciones = client.obtenerCanciones();

            System.out.println("📥 Respuesta recibida");

            if (canciones == null) {
                System.err.println("❌ Respuesta nula del servidor");
                return new ArrayList<>();
            }

            System.out.println("✅ Canciones deserializadas: " + canciones.size());

            if (!canciones.isEmpty()) {
                CancionDTOEntrada primera = canciones.get(0);
                System.out.println("📋 Ejemplo: " + primera);
            }

            System.out.println("=".repeat(50) + "\n");

            return canciones;

        } catch (feign.RetryableException e) {
            System.err.println("\n❌ ERROR DE CONEXIÓN:");
            System.err.println("   No se pudo conectar a " + BASE_URL);
            System.err.println("   ¿Está el servidor de canciones ejecutándose?");
            System.err.println("   Mensaje: " + e.getMessage());
            return new ArrayList<>();

        } catch (feign.codec.DecodeException e) {
            System.err.println("\n❌ ERROR DE DESERIALIZACIÓN:");
            System.err.println("   No se pudo convertir JSON a objetos");
            System.err.println("   Mensaje: " + e.getMessage());
            e.printStackTrace();
            return new ArrayList<>();

        } catch (Exception e) {
            System.err.println("\n❌ ERROR INESPERADO:");
            System.err.println("   Tipo: " + e.getClass().getName());
            System.err.println("   Mensaje: " + e.getMessage());
            e.printStackTrace();
            return new ArrayList<>();
        }
    }
}