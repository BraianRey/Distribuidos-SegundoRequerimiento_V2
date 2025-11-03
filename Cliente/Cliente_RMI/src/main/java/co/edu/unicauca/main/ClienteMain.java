package co.edu.unicauca.main;

import co.edu.unicauca.vista.Menu;

public class ClienteMain {
    public static void main(String[] args) {
        String userId = null;
        
        // Buscar el parámetro --user-id
        for (int i = 0; i < args.length; i++) {
            if (args[i].equals("--user-id") && i + 1 < args.length) {
                userId = args[i + 1];
                break;
            }
        }

        // Si no se proporcionó el user-id, usar uno por defecto
        if (userId == null) {
            System.out.println("⚠️ No se proporcionó un ID de usuario, usando modo interactivo");
        }

        Menu menu = new Menu(userId);
        menu.ejecutarMenuPrincipal();
    }
}