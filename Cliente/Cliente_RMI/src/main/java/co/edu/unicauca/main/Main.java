package co.edu.unicauca.main;

import co.edu.unicauca.vista.Menu;

/**
 * Clase principal del Cliente RMI
 * Este cliente se conecta al ServidorDeCalculoPreferencias mediante RMI
 * para consultar las preferencias musicales de los usuarios
 */
public class Main {
    
    public static void main(String[] args) {
        System.out.println("===============================================");
        System.out.println("   CLIENTE RMI - SISTEMA DE STREAMING");
        System.out.println("===============================================\n");
        // Buscar parametro --user-id si existe
        String userId = null;
        for (int i = 0; i < args.length; i++) {
            if ("--user-id".equals(args[i]) && i + 1 < args.length) {
                userId = args[i + 1];
                break;
            }
        }

        // Iniciar el menú del cliente pasando el userId (puede ser null)
        Menu menu = new Menu(userId);
        menu.ejecutarMenuPrincipal();
    }
}
