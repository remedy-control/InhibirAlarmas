package models

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ConfigRC struct {
	Server      string
	Sistema     string
	Formulario  string
	Columnas    string
	Condiciones string
	Orden       string
	ID          string
}

func (config *ConfigRC) ActualizacionRC() bool {
	var vsUrl bytes.Buffer
	var vsError string
	vsUrl.WriteString(config.Server)
	vsUrl.WriteString("cSistema=")
	vsUrl.WriteString(config.Sistema)
	vsUrl.WriteString("&cForma=")
	vsUrl.WriteString(config.Formulario)
	vsUrl.WriteString("&cID=")
	vsUrl.WriteString(config.ID)
	vsUrl.WriteString("&cColumnas=")
	vsUrl.WriteString(config.Columnas)
	if config.Condiciones != "" {
		vsUrl.WriteString("&cCondiciones=")
		vsUrl.WriteString(config.Condiciones)
	}
	//log.Println("URL Cambio estado: " + vsUrl.String())
	respuesta, err := http.Get(vsUrl.String())
	if err != nil {
		vsError = err.Error()
	}
	defer respuesta.Body.Close()
	contenido, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		vsError = err.Error()
	}
	for _, linea := range strings.Split(string(contenido), "\n") {
		if strings.Contains(linea, "<UPDATED") && strings.Contains(linea, "<UPDATED") {
			log.Println("El sitio con ID: " + config.ID + " fue actualizado correctamente")
			return true
		} else if strings.Contains(linea, "<ERROR") {
			indexI := strings.Index(linea, ">")
			indexF := strings.LastIndex(linea, "<")
			vsError = linea[indexI+1 : indexF]
			log.Println("ERROR AL ACTUALIZAR: " + strings.Replace(vsError, "\n", " ", -1))
			if strings.Contains(vsError, "Cambio de estado no permitido") {
				return true
			} else {
				return false
			}
		} else {
		}
	}
	return false
}

func (config *ConfigRC) ConsultaRC() ([]map[string]string, string) {
	var vsUrl bytes.Buffer
	var vsError string
	vsUrl.WriteString(config.Server)
	vsUrl.WriteString("cSistema=")
	vsUrl.WriteString(config.Sistema)
	vsUrl.WriteString("&cForma=")
	vsUrl.WriteString(config.Formulario)
	vsUrl.WriteString("&cColumnas=")
	vsUrl.WriteString(config.Columnas)
	vsUrl.WriteString("&cCondiciones=")
	vsUrl.WriteString(config.Condiciones)
	respuesta, err := http.Get(vsUrl.String())
	if err != nil {
		vsError = "Error en peticion http;: " + err.Error()
		return nil, vsError
	}
	defer respuesta.Body.Close()
	contenido, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		vsError = err.Error()
	}
	var campo map[string]string
	var res []map[string]string
	var id string
	for _, linea := range strings.Split(string(contenido), "\n") {
		if strings.Contains(linea, "<Entry>") {
			campo = make(map[string]string)
		}
		if strings.Contains(linea, "<id>") {
			indexI := strings.Index(linea, ">")
			indexF := strings.LastIndex(linea, "<")
			id = linea[indexI+1 : indexF]
		}
		if strings.Contains(linea, "<value>") {
			indexI := strings.Index(linea, ">")
			indexF := strings.LastIndex(linea, "<")
			value := linea[indexI+1 : indexF]
			campo[id] = value

		}
		if strings.Contains(linea, "</Entry>") {
			res = append(res, campo)
		}
	}
	return res, ""
}

func (config *ConfigRC) InsertRC() (string, string) {
	contenido, err := config.respuestaRC()
	var vsError string
	if err != nil {
		return "", err.Error()
	}
	for _, linea := range strings.Split(string(contenido), "\n") {
		if strings.Contains(linea, "<NUEVO") && strings.Contains(linea, "<NUEVO") {
			indexI := strings.Index(linea, ">")
			indexF := strings.LastIndex(linea, "<")
			vsID := linea[indexI+1 : indexF]
			return vsID, ""
		} else if strings.Contains(linea, "<ERROR") {
			indexI := strings.Index(linea, ">")
			indexF := strings.LastIndex(linea, "<")
			vsError = linea[indexI+1 : indexF]
		}
	}
	return "", vsError
}

func (config *ConfigRC) respuestaRC() ([]byte, error) {

	var vsUrl bytes.Buffer
	vsUrl.WriteString(config.Server)
	vsUrl.WriteString("cSistema=")
	vsUrl.WriteString(config.Sistema)
	vsUrl.WriteString("&cForma=")
	vsUrl.WriteString(config.Formulario)
	if config.ID != "" {
		vsUrl.WriteString("&cID=")
		vsUrl.WriteString(config.ID)
	}
	vsUrl.WriteString("&cColumnas=")
	vsUrl.WriteString(config.Columnas)
	if config.Condiciones != "" {
		vsUrl.WriteString("&cCondiciones=")
		vsUrl.WriteString(config.Condiciones)
	}
	respuesta, err := http.Get(vsUrl.String())
	if err != nil {
		return nil, err
	}
	defer respuesta.Body.Close()
	contenido, err := ioutil.ReadAll(respuesta.Body)
	if err != nil {
		return nil, err
	}
	return contenido, nil
}
