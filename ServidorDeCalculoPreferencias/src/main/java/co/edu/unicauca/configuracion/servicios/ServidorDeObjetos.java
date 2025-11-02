package co.edu.unicauca.configuracion.servicios;

import java.net.MalformedURLException;
import java.rmi.Naming;
import java.rmi.RemoteException;
import java.rmi.registry.LocateRegistry;
import java.rmi.registry.Registry;

import co.edu.unicauca.capaDeControladores.ControladorPreferenciasUsuariosIml;
import co.edu.unicauca.fachadaServices.services.PreferenciasServiceImpl;

public class ServidorDeObjetos
{
   
    public static void arrancarNS(String direccionIPNS, int puertoNS) 
    {       
        try
        {

            Registry registro = LocateRegistry.getRegistry(direccionIPNS,puertoNS);  
            String[] vectorObjetosRemotos=registro.list();
            System.out.println("El rmiRegistry se ha obtenido y se encuentra escuchando en el puerto: " + puertoNS); 
           
            System.out.println("Objetos remotos registrados en el NS:");
            for (String idObjetosRemotos : vectorObjetosRemotos) {
                System.out.println("Id del objeto remoto : "+idObjetosRemotos );
            }
            
        }
        catch(RemoteException e)
        {
                System.out.println("El rmiRegistry no se localizó en el puerto: " + puertoNS);
                crearNS(puertoNS);
        }
        
    }

    private static void crearNS(int puertoNS)
    {
        try
        {
            Registry registro = LocateRegistry.createRegistry(puertoNS);
            System.out.println("El rmiRegistry  se ha creado en el puerto: " + puertoNS+" "+registro.toString());
        
        }
        catch(RemoteException e)
        {
                System.out.println("El rmiRegistry no se logró crear en el puerto: " + puertoNS);
        }
    }

    public static void registrarObjetoRemoto(ControladorPreferenciasUsuariosIml objetoRemoto, String direccionIPNS, int puertoNS, String identificadorObjetoRemoto)
	{
            String UrlRegistro = "rmi://"+direccionIPNS+":"+puertoNS+"/"+identificadorObjetoRemoto;
            try
            {
                Naming.rebind(UrlRegistro, objetoRemoto);
                System.out.println("Se realizó el registro del objeto remoto en el ns ubicado en la dirección: " +direccionIPNS+" y "+ "puerto "+puertoNS);
            } catch (RemoteException e)
            {
                System.out.println("Error en el registro del objeto remoto");
                e.printStackTrace();
            } catch (MalformedURLException e)
            {
                System.out.println("Error url inválida");                  
                e.printStackTrace();
            }
	}	

    public static ControladorPreferenciasUsuariosIml crearObjetoRemoto() 
    {
        ControladorPreferenciasUsuariosIml objControladorPreferencias=null;
        try {
            objControladorPreferencias = new ControladorPreferenciasUsuariosIml(new PreferenciasServiceImpl());
        } catch (RemoteException e) {            
            e.printStackTrace();
        }
        return objControladorPreferencias;   
    }   
}
