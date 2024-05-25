# Gestor de archivos por descripción
Esta aplicación web permite gestionar (subir, descargar, eliminar) una colección de archivos de forma que estén online y sean disponibles desde cualquier dispositivo. 

## Ejecución 
Descargar la versión correspondiente al sistema operativo/arquitectura desde Releases y ejecutar

En GNU/Linux igual es necesario otorgar permisos específicos de ejecución:
```bash
sudo chmod +x gestor_archivos_por_descripcion
```
### Opciones
```bash
Usage of ./gestor_archivos_por_descripcion:
  -password string
        password en caso de que sea necesario autentificación para operaciones de lectura/escritura
  -megas int
        tamaño máximo en megas de los archivos que se pueden subir (default 20)
  -puerto int
        puerto para la aplicación web (default 8000)
```

Mediante **--pasword** se puede requerir que cualquier operación de lectura y escritura esté protegida por esta contraseña.

Por defecto, cualquier solicitud con un archivo mayor de 20 megas se rechazará, si fuera necesario ampliar este número puede hacerse con **--megas**.

Se puede asignar un puerto específico mediante **--puerto**, de lo contrario intentará iniciarse en el 8000.

## Funcionamiento
La aplicación web creará en el mismo directorio en el que se encuentra el ejecutable un archivo llamado **bbdd.sqlite** y una carpeta con el nombre **archivos**.

En la carpeta "archivos" se guardarán todos los archivos que se vayan subiendo. El nombre será una normalización del nombre del fichero enviado con un prefijo numérico que corresponde al identificador de la base de datos. De esta forma, es posible tener distintos archivos con el mismo nombre y con una descripción diferente (o no, depende del usuario).

El archivo "bbdd.sqlite" es una base de datos SQLite3 con una estructura muy básica:
```sql
CREATE TABLE IF NOT EXISTS "registros" (
		"id"	INTEGER NOT NULL,
		"descripcion"	TEXT NOT NULL,
		"nombre"	TEXT NOT NULL,
		PRIMARY KEY("id" AUTOINCREMENT)
);
```
Las búsquedas se hacen mediante un LIKE:
```sql
SELECT * FROM registros WHERE descripcion LIKE "%XXX%"
```
donde XXX es el cadena introducida en el campo búsqueda del formulario.

### Tamaño demasiado grande
Cuando la solciitud del archivo sea mayor que el tamaño aceptado, devolverá un cierre de conexión. Depende del navegador, pero será algo del tipo:
```
La conexión ha sido reiniciada
La conexión al servidor fue reiniciada mientras la página se cargaba.
```
No es la forma más recomendable de gestionar este error, pero sí la más cómoda. No afecta al funcionamiento de la aplicación, simplemente es un cierre quizá "abrupto", pero seguro.


### Automatización
Se ha intentado simplificar los endpoints y los parámetros de URL (o query strings) por si es necesaria la automatización. 
### Subir
```
curl -v -F descripcion="{descripción}" -F password="{contraseña}" -F archivo=@{ruta}  http://localhost:8000/subir
```
donde {descripción} es la descripción del archivo, {ruta} es la ruta al archivo y {contraseña} es el password (con el valor codificado o "URL encoded"). Si la aplicación no lo usa puede dejarse el valor vacío o directamente omitirse

### Buscar
```
GET:

/buscar?busqueda={ejemplo}&password={contraseña}
```
donde {ejemplo} es la cadena que se quiere buscar y {contraseña} es el password (con el valor codificado o "URL encoded"). Si la aplicación no lo usa puede dejarse el valor vacío o directamente omitirse
### Descargar
```
GET:

/descargar?id={identificador}&password={contraseña}
```
donde {identificador} es el id asignado en la base de datos y {contraseña} es el password (con el valor codificado o "URL encoded"). Si la aplicación no lo usa puede dejarse el valor vacío o directamente omitirse.
### Eliminar
```
GET:

/eliminar?id={identificador}&password={contraseña}
```
donde {identificador} es el id asignado en la base de datos y {contraseña} es la contraseña (con el valor codificado o "URL encoded"). Si la aplicación no lo usa puede dejarse el valor vacío o directamente omitirse.

## Capturas
![Captura](https://i.imgur.com/GhUdzne.png)

![Captura](https://i.imgur.com/xhKDpgI.png)

![Captura](https://i.imgur.com/mcIvpbE.png)

![Captura](https://i.imgur.com/lyEpMOC.png)