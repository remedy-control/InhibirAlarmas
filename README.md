># CRVentanas

## Descripción del proyecto

Este proyecto tiene como finalidad evitar que se generen incidentes de tipo alarmas cuando uno o mas sitios se encuentran en mantenimiento programado a través de un CRQ.

- Toma en cuenta las fechas programadas de ejecución del cambio los sitios pasan a estado `OPERANDO SIN RADIAR`, `POR REGLA AUTOMATICA PARA INHIBIR INCIDENTES` de manera temporal. 
- Al término de la fecha de programa de fin cambian los sitios a estado `OPERANDO`.

## Estado del proyecto

Este proyecto actualmente se encuentra en uso y está desplegado en los 3 ambientes (desarrollo, QA y producción MX 166) únicamente en un servidor por ambiente.

## Requerimientos

Se recomienda encarecidamente respetar estos puntos para poder usar la aplicación de forma correcta:

-   Tener instalado GO

## Ejecutar

Divido que el programa se ejecuta por factores externos no se puede ejecutar de manera local.

## Despliegue en el servidor

Para desplegar la aplicación, se requiere colocar el archivo `VentanasCRQ` en el servidor de aplicaciones, siguiendo estos pasos:

1.  Accede al servidor de aplicaciones.
2.  Accede a la terminal y buscar numero del servio de InhibirAlarmas con el comando `ps -fea|grep InhibirAlarmas`.
3.  Encontrado el servico lo borramos con el comando `kill -9 (numero del servicio)`
4.  Navega hasta la ubicación `/home/remedy/VentanasCR` del servidor.
5.  En la anterior ruta mencionada, coloca el arhcivo `VentanasCRQ`, ubicado en la ruta `/InhibirAlarmas/VentanasCRQ` del proyecto ya ejecutado.
6.  Una vez que hayas colocado el archivo en el servidor, modificar el nombre o borrar el archico `InhibirAlarmas`.
7.  Al archivo `VentanasCRQ` cambiarle el nombre a `InhibirAlarmas` y darle los permisos `755`.
8.  En la terminal alzar el proyecto con el comando `nohup /home/remedy/VentanasCR/InhibirAlarmas &`.
9.  Verificar que esta arriba con el comando `ps -fea|grep InhibirAlarmas`. 

## Importante

Cabe recalcar que siendo este programa en lenguaje GO la forma de ejecutar el proyecto es:

1.  En su editor de código (Visual studo code), cambiar las configuraciones de go con los comandos:
    -  $env:GOOS="linux"  
    -  $Env:GOARCH="amd64"
2. Ya podemos ejecutarlo con el comando `go build`, ya con esto tendremos el archivo `VentanasCRQ`

