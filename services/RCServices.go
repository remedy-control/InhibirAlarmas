package services

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"VentanasCRQ/Utilerias"

	"VentanasCRQ/models"
)

// modificacion
func CambiarStatusSitio(psID string, psStatus string, psMensaje string, psQualifi string) bool {
	config := models.ConfigRC{
		Server:     "http://localhost:8080/Remedy/servicios/RMDUpdate?",
		Sistema:    "CRQVENTANAS",
		Formulario: "Site-EP",
		Columnas:   "%277%27=%27" + psStatus + "%27" + psQualifi + psMensaje + "",
		ID:         psID,
	}
	logFile, err := os.OpenFile("/home/remedy/VentanasCR/CambiarStatus.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Cambiar estatus", config.Server, config.Sistema, config.Formulario, config.Columnas, config.ID)
	return config.ActualizacionRC()
}

// Actualización 30/11/2023
func StatusCRQ(psID string, cambio string, fechaI string, fechaF string) bool {
	config := models.ConfigRC{
		Server:     "http://localhost:8080/Remedy/servicios/RMDUpdate?",
		Sistema:    "CRQVENTANAS",
		Formulario: "Site-EP",
		Columnas:   "%27536878365%27=%27" + cambio + "%27%20%27536878366%27=%27" + fechaI + "%27%20%27536878367%27=%27" + fechaF + "%27",
		ID:         psID,
	}
	logFile, err := os.OpenFile("/home/remedy/VentanasCR/StatusCRQ.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Cambiar estatus", config.Server, config.Sistema, config.Formulario, config.Columnas, config.ID)
	return config.ActualizacionRC()
}

func BuscarSitio(psNemonico string) (string, string, string) {
	config := models.ConfigRC{
		Server:      "http://localhost:8080/Remedy/servicios/RMDSelect?",
		Sistema:     "CRQVENTANAS",
		Formulario:  "Site-EP",
		Columnas:    "1%207%20536870915",
		Condiciones: "%27536870925%27=%27" + strings.Replace(psNemonico, " ", "%20", -1) + "%27%20%27730000001%27!=%27(PACKET%20CORE,VOLTE,SVA)%27%20%27536870974%27!=%27(AXE,TSP,BSC,UMG,UGC,MGCF,BSP,RNC,MGW,CGP,IGW,CSCF,ATS,MRFP,MEDIAX)%27",
		//	Filtro:      "%27536870925%27=%27" + strings.Replace(psNemonico, " ", "%20", -1) + "%27%20%27730000001%27=%27(GMS,IMS)%27%20%27536870974%27=%27(MGW,CGP,IGW)%27",
	}

	voRes, err := config.ConsultaRC()
	if err != "" {
		return "", err, ""
	}
	if len(voRes) == 0 {
		return "", "", ""
	}
	vsRes := voRes[0]["1"]
	vsStatus := voRes[0]["7"]
	vsStatusReason := voRes[0]["536870915"]
	return vsStatus, vsRes, vsStatusReason
	/*------------------------------------------------*/
}

func ConsultarCR(psIdCrq string) ([]map[string]string, string) {
	config := models.ConfigRC{
		Server:      "http://localhost:8080/Remedy/servicios/RMDSelect?",
		Sistema:     "CRQVENTANAS",
		Formulario:  "CHG:Infrastructure%20Change",
		Columnas:    "1000000350%201000000362",
		Condiciones: "%271000000182%27=%27" + psIdCrq + "%27",
	}
	return config.ConsultaRC()
}

func SitioRelacionado(psIdCrq string) ([]map[string]string, string) {
	config := models.ConfigRC{
		Server:      "http://localhost:8080/Remedy/servicios/RMDSelect?",
		Sistema:     "CRQVENTANAS",
		Formulario:  "CHG:Associations",
		Columnas:    "1000000206",
		Condiciones: "%271000000205%27=%27" + psIdCrq + "%27",
	}
	return config.ConsultaRC()
}

func CrearAlarma(psSitio string, psCambio string) (string, string) {
	vsCondicionesDefault := "%271000000000%27=%27INHIBIR%20ALARMA-SITIO%20ATORADO%20DESPUES%20DE%20TRES%20INTENTOS%27%20%271000000151%27=%27El%20sitio%20" + psSitio + "%20del%20cambio%20" + psCambio + "%20no%20se%20ha%20podido%20regresar%20a%20OPERANDO%20despues%20de%20tres%20intentos%20por%20errores%20en%20el%20servicio%20de%20Remedy%20Control,%20favor%20de%20regresarlo%20manualmente%27%20%271000000163%27=%274000%27%20%271000000162%27=%274000%27%20%271000000099%27=%273%27%20%271000000215%27=%275000%27%20%27536890044%27=%270%27%20%271000000063%27=%27FALLA%27%20%271000000064%27=%27DISPONIBILIDAD%27%20%271000000065%27=%27SIN%20SERVICIO%27%20%27536890003%27=%271%27%20%271000000251%27=%27TELCEL%27%20%271000000014%27=%27TRANSFORMACION%20DIGITAL%27%20%271000000217%27=%27CORP-TD-DIGITALIZACION%20DE%20SERVICIOS%27"
	config := models.ConfigRC{
		Server:     "http://localhost:8080/Remedy/servicios/RMDInsert?",
		Sistema:    "CRQVENTANAS",
		Formulario: "HPD:Help%20Desk",
		Columnas:   "%27303497400%27=%27" + psSitio + "%27%20" + vsCondicionesDefault,
	}
	return config.InsertRC()
}

func ParseFecha(psFecha string) time.Time {

	parse := strings.Split(psFecha, " ")
	//fmt.Printf("%q\n", parse)
	fecha := strings.Split(parse[0], "/")
	//fmt.Printf("%q\n", fecha)
	viDia, _ := strconv.Atoi(fecha[0])
	viMes, _ := strconv.Atoi(fecha[1])
	viAño, _ := strconv.Atoi(fecha[2])
	hora := strings.Split(parse[1], ":")
	//fmt.Printf("%q\n", hora)
	viHora, _ := strconv.Atoi(hora[0])
	viMin, _ := strconv.Atoi(hora[1])
	//viSeg, _ := strconv.Atoi(hora[2])

	return time.Date(viAño, time.Month(viMes), viDia, viHora, viMin, 0, 0, time.Local)

	//return parse
}

func CancelarProceso(poArchivo string, psPathlocal string, psPathDestino string) {
	voRegistros := Utilerias.LeerArchivo(psPathlocal + poArchivo)
	Utilerias.CrearArchivo(psPathDestino+poArchivo, voRegistros)
	os.Remove(psPathlocal + poArchivo)
}

func CambioEstadoPrematuro(poCambio string, psPathProc string, psPathFinal string) {
	poArchivo := poCambio + ".csv"
	voRegistros := Utilerias.LeerArchivo(psPathProc + poArchivo)
	//Validar archivo en Procesados
	if voRegistros == nil {
		log.Println("Archivo " + psPathProc + poCambio + ".csv no encontrado")
	} else {
		sitios, vsErr := SitioRelacionado(poCambio)
		if vsErr != "" {
			log.Println("/change. Error al consultar sitios: " + vsErr)
		}
		for _, sitio := range sitios {
			_, idSitio, _ := BuscarSitio(sitio["1000000206"])
			vsMensaje := "%20%27536878362%27=%27Actualizado%20a%20OPERANDO%20por%20fin%20del%20cambio%20" + poCambio + "%27"
			vsQualifi := "%20%27536870915%27=%27-%27"

			CambiarStatusSitio(idSitio, "4", vsMensaje, vsQualifi)
		}
		Utilerias.CrearArchivo(psPathFinal+poArchivo, voRegistros)
		os.Remove(psPathProc + poArchivo)
	}
}

func ProcesarArchivos(poArchivos []os.FileInfo, psPathlocal string, psPathDestino string, psPathEliminar string, psStatus string) {
	for _, voArchivo := range poArchivos {
		//fmt.Println(voArchivo.Name())
		logFile, err := os.OpenFile("/home/remedy/VentanasCR/ProcesarArchivos.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Panic(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
		log.SetFlags(log.Lshortfile | log.LstdFlags)

		voRegistros := Utilerias.LeerArchivo(psPathlocal + voArchivo.Name())
		voRegistrosValidos := []Utilerias.Ventanas{}
		var vfFecha time.Time
		var vfFechaF string
		var vfFechaIni string
		var vsMensaje string
		var vsQualifi string
		var validateChange bool = true
		var validateAux bool
		var validaIntento bool = false
		var vsCambio string
		var vsCiclo string
		var validateProcess bool = false
		for _, voRegistro := range voRegistros {
			vsCambio = voRegistro.Cambio
			if psStatus == "3" {
				vfFecha = ParseFecha(voRegistro.FechaI)
				vfFechaF = ParseFecha(voRegistro.FechaF).String()
				vfFechaIni = ParseFecha(voRegistro.FechaI).String()
				log.Println("Procesar Archivos: ", vfFechaF, vfFechaIni, reflect.TypeOf(vfFechaF), voRegistro, voRegistro)

				vsMensaje = "%20%27536878362%27=%27Actualizado%20a%20OPERANDO%20SIN%20RADIAR%20por%20inicio%20del%20cambio%20" + voRegistro.Cambio + "%27"
				vsQualifi = "%20%27536870915%27=%27POR%20REGLA%20AUTOMATICA%20PARA%20INHIBIR%20INCIDENTES%27"

				vfFechaActual := time.Now()
				diff := vfFechaActual.Sub(vfFecha)
				if diff >= 0 {
					validateProcess = true
					os.Remove(psPathlocal + voArchivo.Name())
					vsStatus, idSitio, _ := BuscarSitio(strings.Replace(voRegistro.Sitio, " ", "%20", -1))
					if vsStatus == "4" {
						log.Println("Realizando cambio del sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio)
						validateAux = CambiarStatusSitio(idSitio, psStatus, vsMensaje, vsQualifi)
						log.Println("Procesar Archivos2: ", validateAux)
						StatusCRQ(idSitio, vsCambio, vfFechaIni, vfFechaF)

						if !validateAux {
							log.Println("El sitio " + voRegistro.Sitio + " no se actualizó correctamente")
							validateChange = false
						}
						voRegistrosValidos = append(voRegistrosValidos, voRegistro)
					} else {
						log.Println("El sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio + " ya no se encuentra OPERANDO, se excenta del proceso")
						validateChange = false
					}
				} else {
					voRegistrosValidos = append(voRegistrosValidos, voRegistro)
				}
			} else if psStatus == "4" {
				vfFecha = ParseFecha(voRegistro.FechaF)
				vsCiclo = voRegistro.Ciclo
				vsMensaje = "%20%27536878362%27=%27Actualizado%20a%20OPERANDO%20por%20fin%20del%20cambio%20" + voRegistro.Cambio + "%27"
				vsQualifi = "%20%27536870915%27=%27-%27"
				vfFechaIni = "-"
				vfFechaF = "-"
				vsCambio = "-"
				vfFechaActual := time.Now()
				diff := vfFechaActual.Sub(vfFecha)
				if diff >= 0 {
					validateProcess = true
					os.Remove(psPathlocal + voArchivo.Name())
					vsStatus, idSitio, vsStatusReason := BuscarSitio(strings.Replace(voRegistro.Sitio, " ", "%20", -1))
					log.Println("El sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio + " se encuentra en status " + vsStatus + " y tiene motivo de estado " + vsStatusReason)
					if vsStatus == "3" && vsStatusReason == "POR REGLA AUTOMATICA PARA INHIBIR INCIDENTES" {
						log.Println("Realizando cambio del sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio)
						validateAux = CambiarStatusSitio(idSitio, psStatus, vsMensaje, vsQualifi)
						StatusCRQ(idSitio, vsCambio, vfFechaIni, vfFechaF)
						if !validateAux {
							ciclo, _ := strconv.Atoi(vsCiclo)
							ciclo++
							vsCiclo := strconv.Itoa(ciclo)
							if ciclo > 3 {
								log.Println("El sitio " + voRegistro.Sitio + " del cambio: " + vsCambio + " no se actualizó correctamente en mas de tres intentos se genera alarma.")
								inc, err := CrearAlarma(voRegistro.Sitio, vsCambio)
								if err != "" {
									log.Println("Se genera el incidente " + inc + " para regresar manualmente el sitio a OPERANDO")
									ventana := Utilerias.Ventanas{
										Cambio: voRegistro.Cambio,
										Sitio:  voRegistro.Sitio,
										FechaI: voRegistro.FechaI,
										FechaF: voRegistro.FechaF,
										Ciclo:  "0",
									}
									voRegistrosValidos = append(voRegistrosValidos, ventana)
								} else {
									log.Println("ERROR. No se logro generar la alarma para el sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio + ", favor de revisar")
									voRegistrosValidos = append(voRegistrosValidos, voRegistro)
								}
								validaIntento = true
							} else {
								log.Println("El sitio " + voRegistro.Sitio + " del cambio: " + vsCambio + " no se actualizó correctamente. Intento " + vsCiclo)
								ventana := Utilerias.Ventanas{
									Cambio: voRegistro.Cambio,
									Sitio:  voRegistro.Sitio,
									FechaI: voRegistro.FechaI,
									FechaF: voRegistro.FechaF,
									Ciclo:  vsCiclo,
								}
								voRegistrosValidos = append(voRegistrosValidos, ventana)
								validateChange = false
							}
						} else {
							voRegistrosValidos = append(voRegistrosValidos, voRegistro)
						}
					} else {
						log.Println("El sitio: " + voRegistro.Sitio + " del cambio: " + vsCambio + " fue movido manualemente durante la ejecución del cambio, se excenta del proceso")
						validateChange = true
					}
				} else {
					voRegistrosValidos = append(voRegistrosValidos, voRegistro)
				}
			}
		}
		if len(voRegistrosValidos) > 0 {
			if validateProcess {
				if psStatus == "3" {
					if validateChange {
						log.Println("Todos los sitios del cambio: " + vsCambio + " se procesaron correctamente")
					} else {
						log.Println("Uno o mas sitios del cambio: " + vsCambio + " no se procesaron correctamente")
					}
					Utilerias.CrearArchivo(psPathDestino+voArchivo.Name(), voRegistrosValidos)
				} else if psStatus == "4" {
					if validateChange {
						if validaIntento {
							log.Println("Uno o mas sitios del cambio: " + vsCambio + " no pudo regresar a OPERANDO y se generó una alarma. El registro sera descartado")
							Utilerias.CrearArchivo(psPathEliminar+voArchivo.Name(), voRegistrosValidos)
						} else {
							log.Println("Todos los sitios del cambio: " + vsCambio + " se procesaron correctamente")
							Utilerias.CrearArchivo(psPathDestino+voArchivo.Name(), voRegistrosValidos)
						}
					} else {
						log.Println("Uno o mas sitios del cambio: " + vsCambio + " no se procesaron correctamente. Rollback")
						Utilerias.CrearArchivo(psPathlocal+voArchivo.Name(), voRegistrosValidos)
					}
				}
			}
		} else {
			log.Println("ERROR: Ningun Sitio del cambio " + vsCambio + " sigue cumpliendo con las reglas del sistema")
			Utilerias.CrearArchivo(psPathEliminar+voArchivo.Name(), voRegistros)
		}
		logFile.Close()
	}
}

func ValidarCRQ(psSitio string, psFechaI string, psFechaF string) bool {
	vsPathLocal := "/home/remedy/VentanasCR/Pendientes/"
	vsPathLocalProc := "/home/remedy/VentanasCR/Procesados/"
	//vsPathLocal := "c:/Logs/VentanasCR/Pendientes/"
	//vsPathLocalProc := "c:/Logs/VentanasCR/Procesados/"
	voArchivos, _ := ioutil.ReadDir(vsPathLocal + ".")
	voArchivosProce, _ := ioutil.ReadDir(vsPathLocalProc + ".")
	vsStatus, _, _ := BuscarSitio(strings.Replace(psSitio, " ", "%20", -1))
	//log.Println("ValidarCRQ", "\""+vsStatus+"\"")
	if vsStatus == "4" {
		if len(voArchivos) > 0 {
			for _, voArchivo := range voArchivos {
				//fmt.Println(voArchivo.Name())
				voRegistros := Utilerias.LeerArchivo(vsPathLocal + voArchivo.Name())
				if validarCarpeta(voRegistros, psSitio, psFechaI, psFechaF) {
					return false
				}
			}
			for _, voArchivo := range voArchivosProce {
				//fmt.Println(voArchivo.Name())
				voRegistros := Utilerias.LeerArchivo(vsPathLocalProc + voArchivo.Name())
				if validarCarpeta(voRegistros, psSitio, psFechaI, psFechaF) {
					return false
				}
			}
		}
	} else {
		return false
	}

	return true
}

func validarCarpeta(psRegistros []Utilerias.Ventanas, psSitio string, psFechaI string, psFechaF string) bool {
	vtFechaI := ParseFecha(psFechaI)
	vtFechaF := ParseFecha(psFechaF)

	for _, voRegistro := range psRegistros {
		vtDifI := vtFechaI.Sub(ParseFecha(voRegistro.FechaF))
		vtDifF := vtFechaF.Sub(ParseFecha(voRegistro.FechaI))
		if voRegistro.Sitio == psSitio && vtDifI < 0 && vtDifF >= 0 {
			return true
		}
	}
	return false
}

func LeerActulizar(psPath string) (err error) {
	voSitiosVentana := Utilerias.LeerArchivo(psPath)
	if voSitiosVentana == nil {
		return errors.New("Archivo" + psPath + "no encontrado")
	}
	for _, voSitioVentana := range voSitiosVentana {
		_, idSitio, _ := BuscarSitio(voSitioVentana.Sitio)
		vsMensaje := "Se%20actualizo%20estatus%20del%20sitio%20a%20OPERANDO%20por%20la%20finalización%20de%20las%20actividades%20del%20cambio%20" + voSitioVentana.Cambio
		err := CambiarStatusSitio(idSitio, "4", vsMensaje, "%20%27536870915%27=%27-%27")
		if !err {
			//CrearAlarma(voSitioVentana.Sitio)
			return errors.New("Se genero un incidentes para el sitio" + voSitioVentana.Sitio + " del cambio " + voSitioVentana.Cambio)
		}
	}
	os.Remove(psPath)
	return nil
}

/*func fncActualizar(poVentanas Utilerias.Ventanas) error {
	_, idSitio, _ := BuscarSitio(poVentanas.Sitio)
	vsMensaje := "Se%20actualizo%20estatus%20del%20sitio%20a%20OPERANDO%20por%20la%20finalización%20de%20las%20actividades%20del%20cambio%20" + poVentanas.Cambio
	err := CambiarStatusSitio(idSitio, "4", vsMensaje, "%20%27536870915%27=%27-%27")
	if !err {
		//CrearAlarma(poVentanas.Sitio)
		return errors.New("Se genero un incidentes para el sitio" + poVentanas.Sitio + " del cambio " + poVentanas.Cambio)
	}
	return nil
}*/
