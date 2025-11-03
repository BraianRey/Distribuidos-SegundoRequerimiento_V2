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
        
        // Iniciar el menú del cliente
        Menu menu = new Menu();
        menu.ejecutarMenuPrincipal();
    }
}
