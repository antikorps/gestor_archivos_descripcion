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

type EliminarData struct {
	Exito bool
}

/* Lógica para el borrado de un archivo y renderizado de la respuesta */
func (webServer *WebServer) eliminar(w http.ResponseWriter, r *http.Request) {
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

	filas, filasError := webServer.SQLite.Conexion.Query("SELECT nombre FROM registros WHERE id = ?", id)
	if filasError != nil {
		log.Println("ERROR: en la operación de eliminado ha fallado el select nombre con ID", id, filasError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	var nombreArchivo string
	for filas.Next() {
		scanError := filas.Scan(&nombreArchivo)
		if scanError != nil {
			log.Println("ERROR: en la operación de eliminado ha fallado el scaneo del nombre con id", id, scanError)
			http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
			return
		}
	}
	filas.Close()

	resultado, resultadoError := webServer.SQLite.Conexion.Exec("DELETE FROM registros WHERE id = ?", id)
	if resultadoError != nil {
		log.Println("ERROR: en la operación de eliminado ha fallado el DELETE con id", id, resultadoError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}

	var eliminarData EliminarData
	filasAfectadas, _ := resultado.RowsAffected()
	if filasAfectadas > 0 {
		eliminarData.Exito = true
	}

	webServer.SQLite.Mutex.Unlock() // MUTEX UNLOCK

	rutaArchivo := filepath.Join(webServer.RutaArchivos, fmt.Sprintf("%d_%v", id, nombreArchivo))
	eliminarError := os.Remove(rutaArchivo)
	if eliminarError != nil {
		log.Println("ERROR: en la operación de eliminado ha fallado la eliminación del archivo", rutaArchivo, eliminarError)
	}

	webServer.EliminarTemplate.Execute(w, eliminarData)
}
