package webserver

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Registro struct {
	Id          int
	Descripcion string
	Nombre      string
}

type BuscarData struct {
	Coincidencias bool
	Password      string
	Registros     []Registro
}

func (webServer *WebServer) buscar(w http.ResponseWriter, r *http.Request) {

	busqueda := strings.TrimSpace(r.URL.Query().Get("busqueda"))
	if busqueda == "" {
		http.Error(w, "Error procesando la solicitud: el campo búsqueda no puede estar vacío", http.StatusBadRequest)
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

	var registros []Registro
	busquedaLike := "%" + busqueda + "%"
	filas, filasError := webServer.SQLite.Conexion.Query("SELECT id, descripcion, nombre FROM registros WHERE descripcion LIKE ?", busquedaLike)
	if filasError != nil {
		log.Println("ERROR: ha fallado el select id, descripcion, nombre con la búsqueda", busqueda, filasError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	for filas.Next() {
		var registro Registro
		scanError := filas.Scan(&registro.Id, &registro.Descripcion, &registro.Nombre)
		if scanError != nil {
			log.Println("ERROR: ha fallado el escaneo tras el select id, descripcion, nombre con la búsqueda", busqueda, scanError)
			http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
			return
		}
		registros = append(registros, registro)
	}
	filas.Close()
	webServer.SQLite.Mutex.Unlock() // MUTEX UNLOCK

	var buscarData BuscarData
	buscarData.Registros = registros
	buscarData.Password = url.QueryEscape(password)
	if len(registros) > 0 {
		buscarData.Coincidencias = true
	}

	webServer.BuscarTemplate.Execute(w, buscarData)
}
