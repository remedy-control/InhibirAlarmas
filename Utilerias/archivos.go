package Utilerias

import (
	"encoding/csv"
	"log"
	"os"
)

type Ventanas struct {
	Cambio string `json:"cambio"`
	Sitio  string `json:"sitio"`
	FechaI string `json:"fechaI"`
	FechaF string `json:"fechaF"`
	Ciclo  string `json:"ciclo"`
	Error  string `json:"error"`
}

func LeerArchivo(psRuta string) []Ventanas {
	if f, err := os.Open(psRuta); err == nil {
		defer f.Close()
		datos := []Ventanas{}
		if lineas, err := csv.NewReader(f).ReadAll(); err == nil {
			for _, linea := range lineas {
				if len(linea) < 5{
					ventana := Ventanas{
						Cambio: linea[0],
						Sitio:  linea[1],
						FechaI: linea[2],
						FechaF: linea[3],
						Ciclo: "0",
					}
					datos = append(datos, ventana)
				}else{
					ventana := Ventanas{
						Cambio: linea[0],
						Sitio:  linea[1],
						FechaI: linea[2],
						FechaF: linea[3],
						Ciclo: linea[4],
					}
					datos = append(datos, ventana)
				}
				
			}
			return datos
		}
	}

	return nil
}

func RemoveVentana(ventanas []Ventanas, eliminar Ventanas) []Ventanas{
	for i, ventana := range ventanas {
		if ventana.Cambio == eliminar.Cambio {
			return append(ventanas[:i], ventanas[i+1:]...)
		}
	}
	return ventanas
}

func CrearArchivo(psRuta string, psDatos []Ventanas) {
	if f, err := os.Create(psRuta); err == nil {
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		fila := []string{}
		for _, dato := range psDatos {
			fila = []string{dato.Cambio, dato.Sitio, dato.FechaI, dato.FechaF, dato.Ciclo}
			if err := w.Write(fila); err != nil {
				log.Println("Error al crear el archivo", err)
			}
		}

		log.Println("Archivo " + psRuta + " creado exitosamente")

	}

}
