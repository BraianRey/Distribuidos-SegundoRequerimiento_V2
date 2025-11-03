package co.edu.unicauca.configuracion.lector;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

/**
 * Lector de propiedades de configuración desde application.properties
 */
public class LectorPropiedadesConfig {
    private static Properties props = new Properties();

    static {
        try (InputStream input = LectorPropiedadesConfig.class.getClassLoader()
                .getResourceAsStream("Cliente/Cliente_RMI/src/main/resources/application.properties")) {
            if (input == null) {
                System.out.println("⚠️  No se encontró el archivo application.properties");
            } else {
                props.load(input);
            }
        } catch (IOException e) {
            System.out.println("❌ Error al cargar application.properties: " + e.getMessage());
        }
    }

    /**
     * Obtener el valor de una propiedad
     */
    public static String get(String key) {
        return props.getProperty(key);
    }

    /**
     * Obtener la IP del servidor RMI
     * @return IP del servidor (por defecto: localhost)
     */
    public String obtenerIPServidor() {
        return props.getProperty("servidor.ip", "localhost");
    }

    /**
     * Obtener el puerto del servidor RMI
     * @return Puerto del servidor (por defecto: 1099)
     */
    public int obtenerPuertoServidor() {
        String puerto = props.getProperty("servidor.puerto", "2020");
        try {
            return Integer.parseInt(puerto);
        } catch (NumberFormatException e) {
            System.out.println("⚠️  Puerto inválido en configuración, usando 1099 por defecto");
            return 1010;
        }
    }
}