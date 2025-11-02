package co.edu.unicauca.configuracion.lector;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

public class LectorPropiedadesConfig {
    private static Properties props = new Properties();

    static {
        try (InputStream input = LectorPropiedadesConfig.class.getClassLoader()
                                                     .getResourceAsStream("application.properties")) {
            if (input == null) {
                System.out.println("No se encontró el archivo application.properties");
            } else {
                props.load(input);
            }
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public static String get(String key) {
        return props.getProperty(key);
    }
}
