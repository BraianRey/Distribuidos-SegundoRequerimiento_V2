package co.edu.unicauca.fachadaServices.DTO;

import java.io.Serializable;
import java.util.List;

import lombok.Data;

@Data
public class PreferenciasDTORespuesta implements Serializable{
    private static final long serialVersionUID = 1L;

    private int idUsuario;
    private List<PreferenciaGeneroDTORespuesta> preferenciasGenero;
    private List<PreferenciaArtistaDTORespuesta> preferenciasArtista;
    private List<PreferenciaIdiomaDTORespuesta>  preferenciasIdioma;
}
