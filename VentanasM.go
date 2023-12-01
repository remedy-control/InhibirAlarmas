package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"VentanasCRQ/Utilerias"
	"VentanasCRQ/services"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//type ErrorRC string

func main() {
	e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	getDatos := func(c echo.Context) error {
		ventanas := []Utilerias.Ventanas{}
		v := new(Utilerias.Ventanas)
		if err := c.Bind(v); err != nil {
			log.Println(err)
		}
		ventana := Utilerias.Ventanas{}
		vsFechaf := services.ParseFecha(v.FechaF)
		vfFechaActual := time.Now()
		diff := vfFechaActual.Sub(vsFechaf)
		if diff < 0 {
			sitios, err := services.SitioRelacionado(v.Cambio)
			if err != "" {
				return c.JSON(http.StatusConflict, ventanas)
			}
			for _, sitio := range sitios {
				if services.ValidarCRQ(sitio["1000000206"], v.FechaI, v.FechaF) {
					ventana = Utilerias.Ventanas{
						Cambio: v.Cambio,
						Sitio:  sitio["1000000206"],
						FechaI: v.FechaI,
						FechaF: v.FechaF,
						Ciclo:  "0",
					}
					ventanas = append(ventanas, ventana)
				} else {
					log.Println("El Sitio " + sitio["1000000206"] + " del cambio " + v.Cambio + " no se encuentra en operando, tiene tecnología PACKET CORE o ya se encuentra un registro activo")
				}

			}

			if len(ventanas) > 0 {
				Utilerias.CrearArchivo("/home/remedy/VentanasCR/Pendientes/"+ventana.Cambio+".csv", ventanas)
				//Utilerias.CrearArchivo("/home/remedy/VentanasCR/Pendientes/"+ventana.Cambio+".csv", ventanas)
			} else {
				log.Println("ERROR: Ningun Sitio del cambio " + v.Cambio + " cumple con las reglas del sistema")
			}

		}

		return c.JSON(http.StatusOK, ventanas)
	}

	cambiarStatus := func(c echo.Context) error {
		//Leer parametros del JSON
		v := new(Utilerias.Ventanas)
		if err := c.Bind(v); err != nil {
			log.Println("Error en el proceso /change")
			log.Println(err)
		}
		vsPathLocal := "/home/remedy/VentanasCR/Procesados/"
		vsPathDestino := "/home/remedy/VentanasCR/Finalizados/"
		log.Println("Cambio cancelado o terminado antes de su fecha fin: " + v.Cambio)
		services.CambioEstadoPrematuro(v.Cambio, vsPathLocal, vsPathDestino)

		return c.String(http.StatusOK, "Proceso terminado")
	}

	procesoVentanas := func(c echo.Context) error {
		//vsPathLocal := "c:/Logs/VentanasCR/Pendientes/"
		//vsPathDestino := "c:/Logs/VentanasCR/Procesados/"
		//vsPathFinal := "c:/Logs/VentanasCR/Finalizados/"
		vsPathPend := "/home/remedy/VentanasCR/Pendientes/"
		vsPathProc := "/home/remedy/VentanasCR/Procesados/"
		vsPathFinal := "/home/remedy/VentanasCR/Finalizados/"
		vsPathElim := "/home/remedy/VentanasCR/Eliminar/"
		voArchivos, _ := ioutil.ReadDir(vsPathPend + ".")
		log.Println("Se procesan los archivos en /Pendientes")
		services.ProcesarArchivos(voArchivos, vsPathPend, vsPathProc, vsPathElim, "3")
		voArchivos, _ = ioutil.ReadDir(vsPathProc + ".")
		log.Println("Se procesan los archivos en /Procesados")
		services.ProcesarArchivos(voArchivos, vsPathProc, vsPathFinal, vsPathElim, "4")
		return c.NoContent(http.StatusOK)
	}

	segundaValidacion := func(c echo.Context) error {
		vsPathFinal := "/home/remedy/VentanasCR/Finalizados/"
		vsPathElim := "/home/remedy/VentanasCR/Eliminar/"
		voArchivos, _ := ioutil.ReadDir(vsPathFinal + ".")
		log.Println("Se procesan los archivos en /Finalizados")
		services.ProcesarArchivos(voArchivos, vsPathFinal, vsPathElim, vsPathElim, "4")
		return c.NoContent(http.StatusOK)
	}

	cancelarProceso := func(c echo.Context) error {
		v := new(Utilerias.Ventanas)
		if err := c.Bind(v); err != nil {
			log.Println("Error en proceso /cancell:")
			log.Println(err)
		}
		//vsPathLocal := "c:/Logs/VentanasCR/Pendientes/"
		//vsPathDestino := "c:/Logs/VentanasCR/Procesados/"
		vsPathLocal := "/home/remedy/VentanasCR/Pendientes/"
		vsPathDestino := "/home/remedy/VentanasCR/Eliminar/"
		voArchivo := v.Cambio + ".csv"
		log.Println("Proceso cancelado antes de iniciar el cambio: " + voArchivo)
		services.CancelarProceso(voArchivo, vsPathLocal, vsPathDestino)

		return c.String(http.StatusOK, "El CRQ ha sido cancelado")
	}

	//*Proceso para generar el archivo del CRQ
	e.POST("/ventana", getDatos)
	//*Proceso cuando un CRQ es cancelado antes de su fecha inicio.
	e.POST("/cancell", cancelarProceso)
	//*Proceso cuando un CRQ es cancelado o termiando antes de su fecha fin.
	e.POST("/change", cambiarStatus)
	//*Proceso de Remedy Control (Pendientes/Procesados)
	e.GET("/proceso", procesoVentanas)
	//*Proceso de Remedy Control (Finalizados)
	e.GET("/validacion", segundaValidacion)
	e.Start(":8020")

}
