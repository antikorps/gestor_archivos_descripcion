package webserver

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sync"
)

type WebServer struct {
	SQLite           *SQLite
	Router           *http.ServeMux
	RutaArchivos     string
	Megas            int64
	Password         string
	IndexTemplate    *template.Template
	SubirTemplate    *template.Template
	BuscarTemplate   *template.Template
	EliminarTemplate *template.Template
}

type SQLite struct {
	Mutex    sync.Mutex
	Conexion *sql.DB
}

type AppConfiguracion struct {
	RutaBaseDatos string
	RutaArchivos  string
	Password      string
	Megas         int64
	Index         string
	Subir         string
	Buscar        string
	Eliminar      string
}

/* Lógica de la aplicación: prepara la conexión a base de datos SQLite y el manejador del servidor web*/
func CrearApp(configuracion AppConfiguracion) WebServer {
	var webServer WebServer
	webServer.RutaArchivos = configuracion.RutaArchivos
	webServer.Password = configuracion.Password
	webServer.Megas = configuracion.Megas

	// Index
	plantillaIndex, plantillaIndexError := template.New("index").Parse(configuracion.Index)
	if plantillaIndexError != nil {
		log.Fatalln("ERROR FATAL: ha fallado el parseo de index.html", plantillaIndexError)
	}
	webServer.IndexTemplate = plantillaIndex

	// Subir
	plantillaSubir, plantillaSubirError := template.New("subir").Parse(configuracion.Subir)
	if plantillaSubirError != nil {
		log.Fatalln("ERROR FATAL: ha fallado el parseo de subir.html", plantillaSubirError)
	}
	webServer.SubirTemplate = plantillaSubir

	// Buscar
	plantillaBuscar, plantillaBuscarError := template.New("buscar").Parse(configuracion.Buscar)
	if plantillaBuscarError != nil {
		log.Fatalln("ERROR FATAL: ha fallado el parse de buscar.html", plantillaBuscarError)
	}
	webServer.BuscarTemplate = plantillaBuscar

	// Eliminar
	plantillaEliminar, plantillaEliminarError := template.New("eliminar").Parse(configuracion.Eliminar)
	if plantillaEliminarError != nil {
		log.Fatalln("ERROR FATAL: ha fallado el parse de eliminar.html", plantillaBuscarError)
	}
	webServer.EliminarTemplate = plantillaEliminar

	// SQLite con Mutex por si hay concurrencia de escritura
	webServer.SQLite = &SQLite{
		Conexion: crearConexionSqlite(configuracion.RutaBaseDatos),
	}

	webServer.crearRouter()
	webServer.manejadorRouter()

	return webServer
}
