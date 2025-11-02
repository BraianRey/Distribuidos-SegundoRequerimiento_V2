package co.edu.unicauca.fachadaServices.services.calculadorPreferencias;

import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;

import co.edu.unicauca.fachadaServices.DTO.CancionDTOEntrada;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaArtistaDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaGeneroDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciaIdiomaDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.DTO.ReproduccionesDTOEntrada;

public class CalculadorPreferencias {

    public PreferenciasDTORespuesta calcular(Integer idUsuario,
                                             List<CancionDTOEntrada> canciones,
                                             List<ReproduccionesDTOEntrada> reproducciones) {

        System.out.println("\n🔢 INICIANDO CÁLCULO DE PREFERENCIAS");
        System.out.println("   Canciones disponibles: " + canciones.size());
        System.out.println("   Reproducciones a procesar: " + reproducciones.size());

        // Crear mapa de canciones por ID
        Map<Integer, CancionDTOEntrada> mapaCanciones = canciones.stream()
                .filter(Objects::nonNull)
                .filter(c -> c.getId() != null)
                .collect(Collectors.toMap(CancionDTOEntrada::getId, c -> c, (a,b) -> a));

        System.out.println("   Mapa de canciones creado: " + mapaCanciones.size() + " entradas");

        // Contadores
        Map<String, Integer> contadorGeneros = new HashMap<>();
        Map<String, Integer> contadorArtistas = new HashMap<>();
        Map<String, Integer> contadorIdiomas = new HashMap<>();  // ← AGREGAR

        int reproduccionesProcesadas = 0;
        int reproduccionesIgnoradas = 0;

        for (ReproduccionesDTOEntrada r : reproducciones) {
            Integer idCancion = r.getIdCancion();

            if (idCancion == null) {
                System.out.println("   ⚠️ Reproducción sin ID de canción, ignorando");
                reproduccionesIgnoradas++;
                continue;
            }

            CancionDTOEntrada c = mapaCanciones.get(idCancion);
            if (c == null) {
                System.out.println("   ⚠️ Canción ID " + idCancion + " no encontrada en el catálogo");
                reproduccionesIgnoradas++;
                continue;
            }

            String genero = c.getGenero() == null ? "Desconocido" : c.getGenero();
            String artista = c.getArtista() == null ? "Desconocido" : c.getArtista();
            String idioma = c.getIdioma() == null ? "Desconocido" : c.getIdioma();  // ← AGREGAR

            contadorGeneros.put(genero, contadorGeneros.getOrDefault(genero, 0) + 1);
            contadorArtistas.put(artista, contadorArtistas.getOrDefault(artista, 0) + 1);
            contadorIdiomas.put(idioma, contadorIdiomas.getOrDefault(idioma, 0) + 1);  // ← AGREGAR

            reproduccionesProcesadas++;
        }

        System.out.println("   ✅ Procesadas: " + reproduccionesProcesadas);
        System.out.println("   ⚠️ Ignoradas: " + reproduccionesIgnoradas);
        System.out.println("   📊 Géneros únicos: " + contadorGeneros.size());
        System.out.println("   👤 Artistas únicos: " + contadorArtistas.size());
        System.out.println("   🌍 Idiomas únicos: " + contadorIdiomas.size());

        // Convertir a DTOs y ordenar por preferencias (descendente)
        List<PreferenciaGeneroDTORespuesta> preferenciasGeneros = contadorGeneros.entrySet().stream()
                .map(e -> new PreferenciaGeneroDTORespuesta(e.getKey(), e.getValue()))
                .sorted(Comparator.comparingInt(PreferenciaGeneroDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaGeneroDTORespuesta::getNombreGenero))
                .collect(Collectors.toList());

        List<PreferenciaArtistaDTORespuesta> preferenciasArtistas = contadorArtistas.entrySet().stream()
                .map(e -> new PreferenciaArtistaDTORespuesta(e.getKey(), e.getValue()))
                .sorted(Comparator.comparingInt(PreferenciaArtistaDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaArtistaDTORespuesta::getNombreArtista))
                .collect(Collectors.toList());

        // ← AGREGAR: Calcular preferencias por idioma
        List<PreferenciaIdiomaDTORespuesta> preferenciasIdiomas = contadorIdiomas.entrySet().stream()
                .map(e -> new PreferenciaIdiomaDTORespuesta(e.getKey(), e.getValue()))
                .sorted(Comparator.comparingInt(PreferenciaIdiomaDTORespuesta::getNumeroPreferencias).reversed()
                        .thenComparing(PreferenciaIdiomaDTORespuesta::getNombreIdioma))
                .collect(Collectors.toList());

        // Crear respuesta
        PreferenciasDTORespuesta respuesta = new PreferenciasDTORespuesta();
        respuesta.setIdUsuario(idUsuario);
        respuesta.setPreferenciasGenero(preferenciasGeneros);
        respuesta.setPreferenciasArtista(preferenciasArtistas);
        respuesta.setPreferenciasIdioma(preferenciasIdiomas);  // ← AGREGAR

        System.out.println("✅ CÁLCULO COMPLETADO\n");

        return respuesta;
    }
}