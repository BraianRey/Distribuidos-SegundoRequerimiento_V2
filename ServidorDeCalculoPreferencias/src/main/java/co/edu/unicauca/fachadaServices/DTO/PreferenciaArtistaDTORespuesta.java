package co.edu.unicauca.fachadaServices.DTO;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaArtistaDTORespuesta implements Serializable{
   private String nombreArtista;
   private Integer numeroPreferencias;
}


