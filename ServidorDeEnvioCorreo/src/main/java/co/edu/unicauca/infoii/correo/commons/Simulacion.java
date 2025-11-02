package co.edu.unicauca.infoii.correo.commons;

// Clase para simular procesos largos en la consola.
public class Simulacion {
    public static void simular(int tiempoTotal, String mensaje) {
        System.out.print(mensaje + " ");
        System.out.flush(); 
        int pasos = 20; 
        int delay = tiempoTotal / pasos; 
    
        for (int i = 0; i < pasos; i++) {
            try {
                Thread.sleep(delay);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();                
                return;
            }
    
            System.out.print("# ");
            System.out.flush(); 
        }    
        System.out.println("Finalizado");
        System.out.flush();
    }

    // Muestra el nombre del hilo actual en la consola.
    public static void mostrarHiloActual() {
        System.out.println("Hilo actual: " + Thread.currentThread().getName());
        System.out.flush(); 
    }
}
