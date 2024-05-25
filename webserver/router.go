package webserver

import "net/http"

func (webServer *WebServer) crearRouter() {
	webServer.Router = http.NewServeMux()
}

/* manejador (G: get, P: post) de las siguientes rutas: G /, P /subir, G /buscar, G /descargar, G /eliminar*/
func (webServer *WebServer) manejadorRouter() {
	webServer.Router.HandleFunc("GET /", webServer.index)
	webServer.Router.HandleFunc("POST /subir", webServer.subir)
	webServer.Router.HandleFunc("GET /buscar", webServer.buscar)
	webServer.Router.HandleFunc("GET /descargar", webServer.descargar)
	webServer.Router.HandleFunc("GET /eliminar", webServer.eliminar)
}
