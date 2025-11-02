package co.edu.unicauca.configuracion.servicios;

import java.rmi.RemoteException;
import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;
import java.rmi.server.UnicastRemoteObject;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosIml;
import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosInt;
import co.edu.unicauca.fachadaServices.services.IPreferenciasService;
import co.edu.unicauca.fachadaServices.services.PreferenciasServiceImpl;

public class ServidorDeObjetos {

    public static void main(String args[]) throws RemoteException {
        try {
            System.out.println("===============================================");
            System.out.println("  SERVIDOR DE CALCULO DE PREFERENCIAS - RMI");
            System.out.println("===============================================\n");

            // Crear el registro RMI en el puerto 1099
            Registry registro = LocateRegistry.createRegistry(2020);
            System.out.println("✅ Registro RMI creado en puerto 2020");

            // Crear la implementación del servicio
            IPreferenciasService objPreferenciasService = new PreferenciasServiceImpl();

            // Crear el controlador
            ControladorPreferenciasUsuariosInt objRemoto = new ControladorPreferenciasUsuariosIml(objPreferenciasService);

            // Exportar el objeto remoto
            ControladorPreferenciasUsuariosInt skeleton = (ControladorPreferenciasUsuariosInt)
                    UnicastRemoteObject.exportObject(objRemoto, 0);

            // Registrar el objeto con el nombre correcto
            registro.rebind("ObjetoRemotoPreferencias", skeleton);

            System.out.println("✅ Objeto remoto registrado: ObjetoRemotoPreferencias");
            System.out.println("✅ Servidor listo para recibir peticiones");
            System.out.println("===============================================\n");

        } catch (RemoteException e) {
            System.err.println("❌ Error al iniciar servidor RMI:");
            System.err.println("   " + e.getMessage());
            e.printStackTrace();
        }
    }
}