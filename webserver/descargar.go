package webserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/* Lógica para la descarga de un archivo y renderizado de la respuesta */
func (webServer *WebServer) descargar(w http.ResponseWriter, r *http.Request) {
	idQuery := strings.TrimSpace(r.URL.Query().Get("id"))
	if idQuery == "" {
		http.Error(w, "Error procesando la solicitud: el campo id no puede estar vacío", http.StatusBadRequest)
		return
	}

	id, idError := strconv.Atoi(idQuery)
	if idError != nil {
		http.Error(w, "Error procesando la solicitud: el campo id tiene que ser numérico", http.StatusBadRequest)
		return
	}

	passwordBruto := strings.TrimSpace(r.URL.Query().Get("password"))
	password, passwordError := url.QueryUnescape(passwordBruto)
	if passwordBruto != "" && passwordError != nil {
		http.Error(w, "Error procesando la solicitud: asegúrate que el password es correcto y está escapado", http.StatusForbidden)
		return
	}
	if password != webServer.Password {
		http.Error(w, "Error procesando la solicitud: usuario no autorizado", http.StatusForbidden)
		return
	}

	webServer.SQLite.Mutex.Lock() // MUTEX LOCK

	filas, filasError := webServer.SQLite.Conexion.Query("SELECT nombre FROM registros WHERE id = ? LIMIT 1", id)
	if filasError != nil {
		log.Println("ERROR: en el select id, nombre con el ID", id, filasError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	var nombreArchivo string
	for filas.Next() {
		scanError := filas.Scan(&nombreArchivo)
		if scanError != nil {
			log.Println("ERROR: ha fallado el escaneo del id, nombre con id", id, scanError)
			http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
			return
		}
	}
	filas.Close()

	webServer.SQLite.Mutex.Unlock() // MUTEX UNLOCK

	if nombreArchivo == "" {
		http.Error(w, "Error en la solicitud: el campo id no tiene coincidencias en la base de datos", http.StatusBadRequest)
		return
	}

	rutaArchivo := filepath.Join(webServer.RutaArchivos, fmt.Sprintf("%d_%v", id, nombreArchivo))

	archivo, archivoError := os.Open(rutaArchivo)
	if archivoError != nil {
		log.Println("ERROR al abrir el archivo", rutaArchivo, archivoError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	defer archivo.Close()

	archivoInfo, archivoInfoError := os.Stat(rutaArchivo)
	if archivoInfoError != nil {
		log.Println("ERROR al obtener la información del archivo", archivoInfo, archivoInfoError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+nombreArchivo)
	w.Header().Set("Content-Length", strconv.FormatInt(archivoInfo.Size(), 10))

	http.ServeFile(w, r, rutaArchivo)

}
