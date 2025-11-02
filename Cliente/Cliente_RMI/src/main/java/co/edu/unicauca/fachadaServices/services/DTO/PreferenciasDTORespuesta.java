package co.edu.unicauca.fachadaServices.services.DTO;

import java.io.Serializable;
import java.util.List;

import lombok.Data;

@Data
public class PreferenciasDTORespuesta implements Serializable{
    private int idUsuario;
    private List<PreferenciaGeneroDTORespuesta> preferenciasGeneros;
    private List<PreferenciaArtistaDTORespuesta> preferenciasArtistas;
}
