package co.edu.unicauca.fachadaServices.DTO;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class PreferenciaGeneroDTORespuesta implements Serializable{
   private String nombreGenero;
   private Integer numeroPreferencias;
}


