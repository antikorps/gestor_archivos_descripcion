package webserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type SubirData struct {
	Nombre      string
	Descripcion string
}

/* Lógica para la súbida de un archivo y renderizado de la respuesta */
func (webServer *WebServer) subir(w http.ResponseWriter, r *http.Request) {

	/* Limitar la lectura del cuerpo a los megas de la configuración

	MaxBytesReader cierra la conexión si se supera, así que devolverá un mensaje del tipo:

	La conexión ha sido reiniciada
	La conexión al servidor fue reiniciada mientras la página se cargaba.

	No es la forma más "usable" de gestionar esta excepción, pero sí la más cómoda...

	La otra opción sería leer el cuerpo de forma progresiva y si se sobrepasa el límite sin haber llegado al io.EOF es demasiado tamaño

	*/
	r.Body = http.MaxBytesReader(w, r.Body, webServer.Megas*1024*1024)

	descripcion := r.FormValue("descripcion")
	if descripcion == "" {
		http.Error(w, "Error procesando la solicitud: falta la descripción del archivo o el archivo enviado es demasido grande", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if password != webServer.Password {
		http.Error(w, "Error procesando la solicitud: usuario no autorizado", http.StatusForbidden)
		return
	}

	archivo, archivoHeader, archivoError := r.FormFile("archivo")
	if archivoError != nil {
		http.Error(w, "Error procesando la solicitud: falta el campo archivo o es demasido grande", http.StatusBadRequest)
		return
	}
	defer archivo.Close()

	// Carpeta archivos
	carpetaError := os.MkdirAll("archivos", 0777)
	if carpetaError != nil && !os.IsExist(carpetaError) {
		log.Println("ERROR: en el servidor manejando la carpeta archivos", carpetaError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}

	nombreArchivo := normalizar(string(archivoHeader.Filename))

	// Guardar a base de datos
	webServer.SQLite.Mutex.Lock() // MUTEX LOCK
	insertRegistro := "INSERT INTO registros (descripcion, nombre) VALUES (?,?);"

	insertResultado, insertRegistroError := webServer.SQLite.Conexion.Exec(insertRegistro, descripcion, nombreArchivo)
	if insertRegistroError != nil {
		log.Println("ERROR: ha fallado el insert del registro (descripcion, nombre) en la base de datos", insertRegistroError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	idRegistro, idRegistroError := insertResultado.LastInsertId()
	if idRegistroError != nil {
		log.Println("ERROR: no se ha podido obtener el id del registro tras la inserción en la base de datos", idRegistroError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}

	webServer.SQLite.Mutex.Unlock() // MUTEX UNLOCK

	// Crear fichero
	rutaArchivo := filepath.Join("archivos", fmt.Sprintf("%d_%v", idRegistro, nombreArchivo))
	fichero, ficheroError := os.Create(rutaArchivo)
	if ficheroError != nil {
		log.Println("ERROR: no se ha podido crear el fichero", rutaArchivo, ficheroError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}
	defer fichero.Close()

	// Copiar a fichero
	_, copiaError := io.Copy(fichero, archivo)
	if copiaError != nil {
		log.Println("ERROR: no se ha podido copiar el fichero", rutaArchivo, copiaError)
		http.Error(w, "Error en el servidor: por favor, inténtalo de nuevo más tarde", http.StatusInternalServerError)
		return
	}

	subirData := SubirData{
		Nombre:      nombreArchivo,
		Descripcion: descripcion,
	}

	webServer.SubirTemplate.Execute(w, subirData)
}
