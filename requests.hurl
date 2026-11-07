# 1. Crear un nuevo registro de álbum (POST /albums)
POST http://localhost:8080/albums
Content-Type: application/json
{
  "id_album" : 1
  "id_usuario": "user1",
  "titulo": "Album 1",
  "artista": "Artista 1",
  "descripcion": "Descripción ejemplo",
  "genero": "Rock",
  "duracion_min": 34,
  "rating": 4
}
HTTP 201
[Asserts]
jsonpath "$.id_regitro" exists
jsonpath "$.id_album" == 1
jsonpath "$.id_usuario" == "user1"
jsonpath "$.titulo" == "Album 1"
jsonpath "$.artista" == "Artista 1"
jsonpath "$.descripcion" == "Descripción ejemplo"
jsonpath "$.genero" == "Rock"
jsonpath "$.duracion_min" == 34
jsonpath "$.rating" == 4


# 2. Listar todos los registros de álbumes (GET /albums)
GET http://localhost:8080/albums
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].id_album" == 1
jsonpath "$[0].id_usuario" == "user1"
jsonpath "$[0].titulo" == "Album 1"
jsonpath "$[0].artista" == "Artista 1"
jsonpath "$[0].descripcion" == "Descripción ejemplo"
jsonpath "$[0].genero" == "Rock"
jsonpath "$[0].duracion_min" == 34
jsonpath "$[0].rating" == 4


# 3. Obtener el registro de álbum por ID (GET /albums/1)
GET http://localhost:8080/albums/1
HTTP 200
[Asserts]
jsonpath "$.id_registro" == 1
jsonpath "$.id_album" == 1
jsonpath "$.id_usuario" == "user1"
jsonpath "$.titulo" == "Album 1"
jsonpath "$.artista" == "Artista 1"
jsonpath "$.descripcion" == "Descripción ejemplo"
jsonpath "$.genero" == "Rock"
jsonpath "$.duracion_min" == 34
jsonpath "$.rating" == 4


# 4. Actualizar el registro de álbum (PUT /albums/1)
PUT http://localhost:8080/albums/1
Content-Type: application/json
{
  "id_usuario": "user3",
  "id_album": 5,
  "titulo": "Album A",
  "artista": "Artista 2",
  "descripcion": "Descripción ejemplo 2",
  "genero": "Rock",
  "duracion_min": 32,
  "rating": 3
}
HTTP 200
[Asserts]
jsonpath "$.id_registro" == 1
jsonpath "$.id_album" == 5
jsonpath "$.id_usuario" == "user3"
jsonpath "$.titulo" == "Album A"
jsonpath "$.artista" == "Artista 2"
jsonpath "$.descripcion" == "Descripción ejemplo 2"
jsonpath "$.genero" == "Rock"
jsonpath "$.duracion_min" == 32
jsonpath "$.rating" == 3


# 5. Borrar el registro de álbum (DELETE /albums/1)
DELETE http://localhost:8080/albums/1
HTTP 204
[Asserts]
status == 204


# 6. Verificar que la lista volvió a quedar vacía (GET /albums)
GET http://localhost:8080/albums
HTTP 200
[Asserts]
jsonpath "$" == null