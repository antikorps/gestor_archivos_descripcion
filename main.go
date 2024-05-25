package main

import (
	_ "embed"
	"flag"
	"fmt"
	"gestor_archivos_por_descripcion/webserver"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed plantillas_html/index.html
var indexEmbed string

//go:embed plantillas_html/subir.html
var subirEmbed string

//go:embed plantillas_html/buscar.html
var buscarEmbed string

//go:embed plantillas_html/eliminar.html
var eliminarEmbed string

type Argumentos struct {
	Puerto   int64
	Megas    int64
	Password string
}

func main() {

	var argumentos Argumentos

	flag.Int64Var(&argumentos.Puerto, "puerto", 8000, "puerto para la aplicación web")
	flag.Int64Var(&argumentos.Megas, "megas", 20, "tamaño máximo en megas de los archivos que se pueden subir")
	flag.StringVar(&argumentos.Password, "password", "", "password en caso de que sea necesario autentificación para operaciones de lectura/escritura")

	flag.Parse()

	ejecutableRuta, ejecutableRutaError := os.Executable()
	if ejecutableRutaError != nil {
		log.Fatalln("ERROR FATAL: no ha podido recuperarse la ruta del ejecutable", ejecutableRutaError)
	}
	rutaDirectorio := filepath.Dir(ejecutableRuta)
	rutaBaseDatos := filepath.Join(rutaDirectorio, "bbdd.sqlite")
	rutaArchivos := filepath.Join(rutaDirectorio, "archivos")

	configuracion := webserver.AppConfiguracion{
		RutaBaseDatos: rutaBaseDatos,
		RutaArchivos:  rutaArchivos,
		Password:      argumentos.Password,
		Megas:         argumentos.Megas,
		Index:         indexEmbed,
		Subir:         subirEmbed,
		Buscar:        buscarEmbed,
		Eliminar:      eliminarEmbed,
	}

	webserver := webserver.CrearApp(configuracion)
	log.Println("INFO: aplicación web iniciada en puerto", argumentos.Puerto)

	http.ListenAndServe(":"+fmt.Sprint(argumentos.Puerto), webserver.Router)
}
