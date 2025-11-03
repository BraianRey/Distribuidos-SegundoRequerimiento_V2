package co.edu.unicauca.fachadaServices.DTO;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
@JsonIgnoreProperties(ignoreUnknown = true)
public class CancionDTOEntrada {
    private Integer id;
    private String titulo;
    private String artista;
    private String genero;
    private String idioma;
    private String duracion;

    @Override
    public String toString() {
        return "Cancion[id=" + id + ", titulo=" + titulo + ", artista=" + artista +
                ", genero=" + genero + ", idioma=" + idioma + "]";
    }
}