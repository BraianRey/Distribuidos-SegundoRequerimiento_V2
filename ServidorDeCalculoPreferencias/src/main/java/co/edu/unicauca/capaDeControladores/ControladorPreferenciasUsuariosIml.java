
package co.edu.unicauca.capaDeControladores;

import java.rmi.RemoteException;
import java.rmi.server.UnicastRemoteObject;
import co.edu.unicauca.fachadaServices.DTO.PreferenciasDTORespuesta;
import co.edu.unicauca.fachadaServices.services.IPreferenciasService;


public class ControladorPreferenciasUsuariosIml extends UnicastRemoteObject  implements ControladorPreferenciasUsuariosInt{

    private IPreferenciasService servicioFachadaPreferencias;

    public ControladorPreferenciasUsuariosIml(IPreferenciasService servicioFachadaPreferencias) throws RemoteException {
        super();
        this.servicioFachadaPreferencias = servicioFachadaPreferencias;
    }

    @Override
    public PreferenciasDTORespuesta getReferencias(Integer id) throws RemoteException {
        return this.servicioFachadaPreferencias.getReferencias(id);
    }
}