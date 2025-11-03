package co.edu.unicauca.fachadaServices.services.DTO;

import java.io.Serializable;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * DTO para transferir información de preferencias por género
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaGeneroDTORespuesta implements Serializable {
    private static final long serialVersionUID = 1L;

    private String nombreGenero;
    private Integer numeroPreferencias;
}