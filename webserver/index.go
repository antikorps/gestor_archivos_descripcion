package webserver

import (
	"net/http"
)

type IndexData struct {
	UsoPassword bool
	Password    string
}

/* renderizado de la página principal */
func (webServer *WebServer) index(w http.ResponseWriter, r *http.Request) {

	var indexData IndexData
	if webServer.Password != "" {
		indexData.Password = webServer.Password
		indexData.UsoPassword = true
	}

	webServer.IndexTemplate.Execute(w, indexData)
}
