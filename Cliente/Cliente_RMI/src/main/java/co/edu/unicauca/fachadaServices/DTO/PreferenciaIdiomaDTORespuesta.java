package co.edu.unicauca.fachadaServices.DTO;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * DTO para transferir información de preferencias por idioma
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaIdiomaDTORespuesta implements Serializable {
    private static final long serialVersionUID = 1L;

    private String nombreIdioma;
    private Integer numeroPreferencias;
}