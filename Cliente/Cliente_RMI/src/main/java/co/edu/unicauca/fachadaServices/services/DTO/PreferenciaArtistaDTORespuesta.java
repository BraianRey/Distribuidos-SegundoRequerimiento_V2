package co.edu.unicauca.fachadaServices.services.DTO;

import java.io.Serializable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * DTO para transferir información de preferencias por artista
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaArtistaDTORespuesta implements Serializable {
    private static final long serialVersionUID = 1L;

    private String nombreArtista;
    private Integer numeroPreferencias;
}