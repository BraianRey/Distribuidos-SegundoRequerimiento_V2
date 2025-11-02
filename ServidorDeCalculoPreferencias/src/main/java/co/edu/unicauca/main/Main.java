package co.edu.unicauca.main;

import java.util.ResourceBundle.Control;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosIml;
import co.edu.unicauca.configuracion.lector.LectorPropiedadesConfig;
import co.edu.unicauca.configuracion.servicios.ServidorDeObjetos;

public class Main {
    public static void main(String[] args) {

        int puertoNS = Integer.parseInt(LectorPropiedadesConfig.get("ns.port"));
        String direccionIPNS=LectorPropiedadesConfig.get("ns.host");

        //paso 1: arrancar o crear el ns
        ServidorDeObjetos.arrancarNS(direccionIPNS, puertoNS);
       
        //paso 2: crear el objeto remoto
        ControladorPreferenciasUsuariosIml objControladorPreferencias = ServidorDeObjetos.crearObjetoRemoto();
       
        //paso 3: registrar el objeto remoto en el ns
        String identificadorObjetoRemoto = "objControladorPreferenciasUsuarios";
        ServidorDeObjetos.registrarObjetoRemoto(objControladorPreferencias, direccionIPNS, puertoNS, identificadorObjetoRemoto);

    }
}

