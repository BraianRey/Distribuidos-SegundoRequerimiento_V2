package  co.edu.unicauca.vista;

import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.services.FachadaGestorUsuariosIml;
import co.edu.unicauca.utilidades.UtilidadesConsola;

import java.rmi.RemoteException;

public class Menu {
    
   
    private final FachadaGestorUsuariosIml objFachada;
    
    public Menu(FachadaGestorUsuariosIml objFachada)
    {
        this.objFachada=objFachada;
    }

    public void ejecutarMenuPrincipal()
    {
        int opcion = 0;
        do
        {
                System.out.println("==Menu==");
                System.out.println("1. Consultar preferencias del Usuario");
                System.out.println("2. Salir");

                opcion = UtilidadesConsola.leerEntero();

            switch(opcion)
            {
                case 1: opcion1(); break;               
                case 2: System.out.println("Gracias por usar el sistema"); break;
                default: System.out.println("Opcion no valida");
            }

        }while(opcion != 2);
    }

    private void opcion1()
    {
        System.out.println("Ingrese el id del usuario a consultar sus preferencias:");
        Integer idUsuario = UtilidadesConsola.leerEntero();
        try {
            //invocación del método remoto a través de la fachada
            PreferenciasDTORespuesta respuesta = this.objFachada.getReferencias(idUsuario);
            System.out.println("==Preferencias del usuario==");
            System.out.println("Generos");
            respuesta.getPreferenciasGeneros().forEach(genero -> {
                System.out.println(genero.getNombreGenero() + " Cantidad de veces escuchado: " + genero.getNumeroPreferencias());
            });

            System.out.println("\nArtistas");
            respuesta.getPreferenciasArtistas().forEach(artista -> {
                System.out.println(artista.getNombreArtista() + " Cantidad de veces escuchado: " + artista.getNumeroPreferencias());
            });
        } catch (RemoteException e) {
            System.out.println("Error al consultar las preferencias: " + e.getMessage());
        }
    }


}
