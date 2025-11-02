package co.edu.unicauca.fachadaServices.DTO;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class CancionDTOEntrada {
    private Integer id;
    private String titulo;
    private String artista;
    private String genero;
}

