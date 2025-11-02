package co.edu.unicauca.main;

import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosIml;
import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosInt;
import co.edu.unicauca.fachadaServices.services.IPreferenciasService;
import co.edu.unicauca.fachadaServices.services.PreferenciasServiceImpl;

public class Main {

    public static void main(String[] args) {
        try {
            System.out.println("===============================================");
            System.out.println("  SERVIDOR DE CALCULO DE PREFERENCIAS - RMI");
            System.out.println("===============================================\n");

            // Crear el registro RMI en el puerto 1099
            System.out.println("📡 Creando registro RMI en puerto 2020...");
            Registry registro = LocateRegistry.createRegistry(2020);
            System.out.println("✅ Registro RMI creado exitosamente");

            // Crear la implementación del servicio de preferencias
            System.out.println("\n🔧 Inicializando servicio de preferencias...");
            IPreferenciasService objPreferenciasService = new PreferenciasServiceImpl();

            // Crear el controlador (ya se exporta automáticamente si extiende UnicastRemoteObject)
            System.out.println("📤 Creando objeto remoto...");
            ControladorPreferenciasUsuariosInt objRemoto = new ControladorPreferenciasUsuariosIml(objPreferenciasService);

            // Registrar el objeto en el registro RMI (SIN exportObject)
            String nombreObjeto = "ObjetoRemotoPreferencias";
            registro.rebind(nombreObjeto, objRemoto);

            System.out.println("✅ Objeto remoto registrado: " + nombreObjeto);
            System.out.println("\n" + "=".repeat(47));
            System.out.println("🚀 SERVIDOR LISTO PARA RECIBIR PETICIONES");
            System.out.println("=".repeat(47));
            System.out.println("📍 Puerto: 2020");
            System.out.println("🔑 Nombre: " + nombreObjeto);
            System.out.println("\n💡 Presiona Ctrl+C para detener el servidor\n");

        } catch (Exception e) {
            System.err.println("\n❌ ERROR al iniciar el servidor RMI:");
            System.err.println("   " + e.getMessage());

            if (e.getMessage() != null && e.getMessage().contains("Address already in use")) {
                System.err.println("\n💡 El puerto 1099 ya está en uso.");
                System.err.println("   Cierra cualquier otro servidor RMI o cambia el puerto.");
            }

            e.printStackTrace();
        }
    }
}