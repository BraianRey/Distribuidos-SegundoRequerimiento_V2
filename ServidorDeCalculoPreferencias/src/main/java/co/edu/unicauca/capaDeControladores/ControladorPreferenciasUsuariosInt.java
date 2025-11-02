package co.edu.unicauca.capaDeControladores;

import java.rmi.Remote;
import java.rmi.RemoteException;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;

public interface ControladorPreferenciasUsuariosInt extends Remote{
    public PreferenciasDTORespuesta getReferencias(Integer id) throws RemoteException;
}