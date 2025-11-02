package co.edu.unicauca.configuracion.servicios;

import java.rmi.Naming;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosInt;

public class ClienteDeObjetos {
    
    public static ControladorPreferenciasUsuariosInt obtenerObjetoRemoto(String direccionIPNS, int puertoNS, String identificadorObjetoRemoto) 
    {
        String URLRegistro;
        URLRegistro  = "rmi://" + direccionIPNS + ":" + puertoNS + "/"+identificadorObjetoRemoto;
        try
        {
            return (ControladorPreferenciasUsuariosInt) Naming.lookup(URLRegistro);
        }
        catch (Exception e)
        {
            System.out.println("Excepcion en obtencion del objeto remoto"+ e);
            return null;
        }
    }   
}


