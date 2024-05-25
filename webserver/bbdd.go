package webserver

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

/* Retorna la conexión tras el create table if not exists necesario*/
func crearConexionSqlite(ruta string) *sql.DB {

	baseDatos, baseDatosError := sql.Open("sqlite", ruta)
	if baseDatosError != nil {
		log.Fatalln(baseDatosError)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS "registros" (
		"id"	INTEGER NOT NULL,
		"descripcion"	TEXT NOT NULL,
		"nombre"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
	);
	`
	// Valorar usar FTS5 y MATCH

	_, createTableError := baseDatos.Exec(createTable)
	if createTableError != nil {
		log.Fatalln("createTableError", createTableError)
	}

	return baseDatos
}
